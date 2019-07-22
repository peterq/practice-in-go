package test

import (
	"log"
	"unsafe"
)

const leng = 88

var bigData = [leng]bool{}

func testMemmory() {
	log.Println(bigData)
	p := (*[leng]byte)(unsafe.Pointer(&bigData))
	for i := 0; i < 40; i++ {
		if i%3 != 0 {
			(*p)[i] = 1
		}
	}
	log.Println(bigData)

	select {}
}
