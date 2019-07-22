package spider_client

import (
	"log"
	"time"
)

func Init() {
	client := New(2, 10, 50*time.Millisecond, true)
	res, err := client.Get("http://baidu.com", 2)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res)
}
