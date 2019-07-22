package main

//go:generate go run loader/app_map_generator.go

import (
	"flag"
	"funny/loader"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var name = flag.String("n", "leet_code", "app 名称")
	flag.Parse()
	loader.StartApp(*name)
}
