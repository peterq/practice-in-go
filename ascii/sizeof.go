package main

import (
	"log"
	"unsafe"
)

type t1 struct {
	a int
}

type t2 struct {
	_ [0]func(string, []string) bool
	_ [0]func(string, []string) bool
	_ [0]func(string, []string) bool
	a int
}

func main() {
	var t11 t1
	var t22 t2
	log.Println(unsafe.Sizeof(t11))
	log.Println(unsafe.Sizeof(t22))
}
