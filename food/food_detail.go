package food

import (
	"encoding/json"
	"funny/spider_client"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"time"
)

var spider *spider_client.Client

func getDetailLoop() {

	spider = spider_client.New(10, 100, 20*time.Millisecond, false)
	ch := spider.GetProxyCh()
	addProxyCh := spider_client.RefreshPool(time.Second, ch)
	var old bson.M
	err := MetaCollection.Find(bson.M{"_id": "proxy"}).One(&old)
	if err == nil {
		if list, ok := old["list"]; ok {
			for _, s := range list.([]interface{}) {
				go func(s string) { addProxyCh <- spider_client.Proxy(s) }(s.(string))
			}
		}
	}
	go func() {
		for {
			list := spider.GetUsingProxy()
			log.Println("现在可用代理", len(list))
			MetaCollection.Upsert(bson.M{"_id": "proxy"}, bson.M{"list": list})
			time.Sleep(time.Second * 5)
		}
	}()
	time.Sleep(time.Second * 3)
	taskCh := make(chan string, 20)
	cal := calulateSpeed()
	// 启动1000个任务接受协程
	for i := 0; i < 2000; i++ {
		go receiveDetailTask(taskCh, cal)
	}

	lastMaxId := 0
	for {
		var tasks []interface{}
		err := TaskCollection.
			Find(bson.M{"fetched": false, "_id": bson.M{"$gt": lastMaxId}}).
			Select(bson.M{"code": 1, "_id": 1}).
			Sort("_id").Limit(4000).All(&tasks)
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		if len(tasks) == 0 {
			log.Println("详情爬取完成")
			os.Exit(0)
		}
		lastMaxId = int(tasks[len(tasks)-1].(bson.M)["_id"].(int))
		for _, task := range tasks {
			taskCh <- task.(bson.M)["code"].(string)
			//log.Println(len(tasks), tasks)
		}
	}
}

func calulateSpeed() chan struct{} {
	ch := make(chan struct{})
	go func() {
		count := 0
		total := 0
		after := time.After(time.Second)
		for {
			select {
			case <-after:
				after = time.After(time.Second)
				log.Println("获取速度(条/每秒)", count, "总数", total)
				count = 0
			case <-ch:
				count++
				total++
			}
		}
	}()
	return ch
}

func receiveDetailTask(codeCh chan string, cal chan struct{}) {
	for {
		code := <-codeCh
		u := "https://food.boohee.com/fb/v1/foods/" + code
		//str, err := cli.Get(u, 2)
		res, err := spider.Get(u, 10)
		str := res.Body
		//log.Println(u, str)
		if err != nil {
			log.Println(err)
			continue
		}
		var j SJson
		err = json.Unmarshal([]byte(str), &j)
		if err != nil {
			log.Println(err)
			continue
		}
		if c, ok := j["code"]; ok && c.(string) == code {

			//log.Println(str)
			go func() {
				cal <- struct{}{}
				return
				TaskCollection.Update(bson.M{"code": code}, bson.M{"$set": bson.M{"fetched": true}})
				DetailCollection.Upsert(bson.M{"_id": j["id"]}, j)
			}()
		}
	}
}

func waitSomeDetailTask(num int, taskOkCh chan string) chan bool {
	done := make(chan bool)
	go func() {
		for i := 0; i < num; i++ {
			<-taskOkCh
		}
		done <- true
	}()
	return done
}
