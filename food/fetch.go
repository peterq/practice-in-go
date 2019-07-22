package food

import (
	"encoding/json"
	"fmt"
	"github.com/huichen/sego"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const sectionSize = 1000

var insertMap = make(map[int]section)

type section struct {
	sec  *[sectionSize]bool
	lock sync.Mutex
}

func noUseWord(word string) bool {
	mp := map[string]bool{
		"（": true,
		"）": true,
		"°": true,
		"，": true,
	}
	_, ok := mp[word]
	return ok
}

func hasInsert(id int) bool {
	sectionIndex := id / sectionSize
	offset := id - sectionSize*sectionIndex
	sec, ok := insertMap[sectionIndex]
	if !ok {
		insertMap[sectionIndex] = section{sec: &[sectionSize]bool{}, lock: sync.Mutex{}}
		sec = insertMap[sectionIndex]
	}
	return sec.sec[offset]
}

func markInsert(id int) {
	sectionIndex := id / sectionSize
	offset := id - sectionSize*sectionIndex
	sec, ok := insertMap[sectionIndex]
	if !ok {
		insertMap[sectionIndex] = section{sec: &[sectionSize]bool{}}
		sec = insertMap[sectionIndex]
	}
	sec.lock.Lock()
	sec.sec[offset] = true
	sec.lock.Unlock()
}

func fetchCats() {
	foodCh := make(chan SJson, 5)
	go receiveFood(foodCh)
	for cat := 1; cat <= 12; cat++ {
		for page := 1; page <= 10; page++ {
			go fetchList(cat, page, foodCh)
		}
	}
}

// 搜索关键词循环
func searchKeywordLoop() {
	taskCh := make(chan string, 10)
	taskOkCh := make(chan string, 10)
	foodCh := make(chan SJson, 10)
	go receiveFood(foodCh)
	go receiveSearchTask(taskCh, taskOkCh, foodCh)
	for {
		var words []interface{}
		err := WordsCollection.Find(bson.M{"searched": false}).Sort("-_id").Limit(100).All(&words)
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		for _, word := range words {
			taskCh <- word.(bson.M)["_id"].(string)
		}
		<-waitSearchKeywordDone(taskOkCh, len(words), func(s string) {
			WordsCollection.Upsert(bson.M{"_id": s}, bson.M{"searched": true})
		})
	}
}

// 等待一轮关键词搜索完成
func waitSearchKeywordDone(taskOkCh chan string, num int, cb func(string)) chan bool {
	done := make(chan bool)
	go func() {
		for i := 0; i < num; i++ {
			str := <-taskOkCh
			cb(str)
		}
		done <- true
	}()
	return done

}

// 搜索一个关键词
func receiveSearchTask(taskCh, taskOkCh chan string, foodCh chan SJson) {
	for {
		keyword := <-taskCh
		go func(keyword string) {
			totalPage := 10
			for page := 1; page <= totalPage; page++ {
				u := fmt.Sprintf("https://food.boohee.com/fb/v1/search?q=%s&page=%d", url.QueryEscape(keyword), page)
				str, err := cli.Get(u, 2)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Println(u, str)
				var j SJson
				err = json.Unmarshal([]byte(str), &j)
				if err != nil {
					log.Println(err)
					continue
				}
				if p, ok := j["total_pages"]; ok {
					p1 := int(p.(float64))
					if p1 < totalPage {
						totalPage = p1
					}
				}
				if items, ok := j["items"]; ok {
					for _, item := range items.([]interface{}) {
						foodCh <- SJson(item.(map[string]interface{}))
					}
				}

			}
			taskOkCh <- keyword
		}(keyword)
	}
}

// 食物条目处理器
func receiveFood(ch chan SJson) {
	// 关键词通道
	wordsCh := make(chan []string, 10)
	go receiveWords(wordsCh) // 关键词处理协程

	// 分词器
	var segmenter sego.Segmenter
	segmenter.LoadDictionary(strings.Split(os.Getenv("GOPATH"), ":")[0] + "/src/github.com/huichen/sego/data/dictionary.txt")
	segmenter.Dictionary()

	for {
		food := <-ch
		if id, ok := food["id"]; ok {
			idInt := int(id.(float64))
			if !hasInsert(idInt) {
				delete(food, "id")
				FoodsCollection.Upsert(bson.M{"_id": idInt}, food)
				markInsert(idInt)
				// 分词
				if name, ok := food["name"]; ok {
					segments := segmenter.Segment([]byte(name.(string)))
					output := sego.SegmentsToSlice(segments, false)
					wordsCh <- output
				}
			}
		}
	}
}

// 关键字处理
func receiveWords(ch chan []string) {
	dealt := make(map[string]bool)
	for {
		words := <-ch
		for _, word := range words {
			_, ok := dealt[word]
			if len(word) <= 1 || ok || noUseWord(word) {
				continue
			}
			c, err := WordsCollection.FindId(word).Count()
			if c == 0 {
				err = WordsCollection.Insert(bson.M{"_id": word, "searched": false})
				if err == nil {
				}
			}
			dealt[word] = true
		}
	}
}

func fetchList(cat, page int, foodCh chan SJson) error {
	url := fmt.Sprintf("https://food.boohee.com/fb/v1/foods?value=%d&kind=group&page=%d", cat, page)
	//log.Println(url)
	str, err := cli.Get(url, 2)
	if err != nil {
		return err
	}
	log.Println(url, str)
	var j SJson
	err = json.Unmarshal([]byte(str), &j)
	if err != nil {
		return err
	}
	if foods, ok := j["foods"]; ok {
		for _, food := range foods.([]interface{}) {
			foodCh <- SJson(food.(map[string]interface{}))
		}
	}
	return nil
}
