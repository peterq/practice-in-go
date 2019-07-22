package test

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type payload struct {
	Id_ bson.ObjectId `bson:"_id"`
	A   string        `json:"a_78"`
}

func testMongo() {
	url := "mongodb://root:root@127.0.0.1:27017/admin"
	s, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	Session := s

	Session.SetMode(mgo.Monotonic, true)
	c := Session.DB("playground").C("test_id")

	id := bson.NewObjectId()
	a := payload{A: "hello", Id_: id}
	err = c.Insert(a)
	bin, _ := bson.Marshal(a)
	log.Println(string(bin))
	log.Println(fmt.Sprintf("%#v %#v", err, id.String()))
}
