package xes

import (
	"funny/spider_client"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"strings"
	"time"
)

var Session *mgo.Session
var InfoCollection *mgo.Collection

// 学而思礼品兑换数据采集
func Init() {
	connectMongo()
	ch := make(chan int, 100)
	go queryTask(ch)
	for i := 120000; i < 130000; i++ {
		ch <- i
	}
}

func insertTask(ch chan *bson.M) {
	var arr []interface{}
	mustInsert := false
Task:
	for {
		ok := false
		var data *bson.M
		select {
		case data, ok = <-ch:
			if !ok {
				break Task
			}
		case <-time.After(5 * time.Second):
			mustInsert = true
		}
		if data != nil {
			arr = append(arr, *data)
		}
		if (len(arr) > 10 || mustInsert) && len(arr) > 0 {
			mustInsert = false
			if err := InfoCollection.Insert(arr...); err == nil {
				arr = arr[0:0]
			} else {
				arr = arr[0:0]
				log.Println(err)
			}
		}
	}
}

func queryTask(ch chan int) {
	client := spider_client.New(200, 200, 0, true)
	saveChan := make(chan *bson.M, 10)
	go insertTask(saveChan)
	for {
		i, ok := <-ch
		num := strconv.Itoa(i)
		num = strings.Repeat("0", 6-len(num)) + num
		if !ok {
			close(saveChan)
			break
		}
		result, err := client.Get("http://dh.gyspbj.com/history/tal_"+num, 1)
		if err != nil {
			log.Println(err)
		}
		info, err := getInfo(result.Body)
		if err != nil {
			log.Println(err)
		}
		//log.Println(num)
		if info != nil {
			log.Println(num)
			(*info)["worker_id"] = num
			(*info)["_id"] = i
			saveChan <- info
		}
	}
}

func connectMongo() {
	url := "mongodb://root:root@127.0.0.1:27017/admin"
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	Session = s

	Session.SetMode(mgo.Monotonic, true)
	InfoCollection = Session.DB("playground").C("xes_worker")
}

func getInfo(str string) (*bson.M, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		return nil, errors.Wrap(err, "文档解析错误")
	}
	if doc.Find("title").Text() != "已兑换-记录" {
		return nil, nil
	}
	result := bson.M{}
	mp := map[string]string{
		"name":    ".his_apply_info > *:nth-child(1)",
		"mobile":  ".his_apply_info > *:nth-child(2)",
		"address": ".his_apply_info > *:nth-child(3)",
		"order":   ".his_apply_info > *:nth-child(5)",
	}
	for k, v := range mp {
		str := doc.Find(v).Text()
		t := strings.Split(str, "：")
		str = t[len(t)-1]
		result[k] = str
	}
	return &result, nil
}
