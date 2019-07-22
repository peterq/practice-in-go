package loader

import (
	"log"
)

type initFunc func()

var mp = map[string]initFunc{}

func StartApp(name string) {
	if fn, ok := mp[name]; ok {
		fn()
	} else {
		log.Fatal("app: " + name + " is not exist")
	}
}
