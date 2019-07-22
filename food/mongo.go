package food

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
	"time"
)

var Session *mgo.Session
var FoodsCollection *mgo.Collection
var WordsCollection *mgo.Collection
var MetaCollection *mgo.Collection
var TaskCollection *mgo.Collection
var DetailCollection *mgo.Collection

type SJson map[string]interface{}

func (j *SJson) Get(key string) *interface{} {
	path := strings.Split(key, ".")
	last := j
	for index, k := range path {
		if index != len(path)-1 {
			v := (*last)[k].(SJson)
			last = &v
		} else {
			v := (*last)[k]
			return &v

		}
	}
	return nil
}

func saveMetaInterval() {
	for {
		time.Sleep(time.Second * 5)
		saveMeta()
	}
}

var Meta = SJson{
	"cats": []int{1, 2, 3, 4},
	"cats_progress": SJson{
		"1": 23,
		"2": 46,
		"3": SJson{},
	},
}

func saveMeta() {
	if MetaCollection != nil {
		log.Println("保存元信息", Meta)
		_, err := MetaCollection.Upsert(bson.M{"_id": "meta"}, Meta)
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println("保存meta失败")
	}
}

func loadMeta() {
	result := SJson{}
	err := MetaCollection.Find(bson.M{"_id": "meta"}).One(&result)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("读取meta", result)
	Meta = result
}

func conMongo() {
	url := "mongodb://root:root@127.0.0.1:27017/admin"
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	s.Refresh()
	Session = s

	Session.SetMode(mgo.Monotonic, true)
	FoodsCollection = Session.DB("playground").C("foods")
	WordsCollection = Session.DB("playground").C("foods_words")
	MetaCollection = Session.DB("playground").C("foods_meta")
	TaskCollection = Session.DB("playground").C("foods_task")
	DetailCollection = Session.DB("playground").C("foods_detail")
}
