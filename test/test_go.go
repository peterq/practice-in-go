package test

import (
	"fmt"
	"log"
	"runtime"
	"time"
	"unsafe"
)

func TestGoroutine() {

	var a = [10]int16{}
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				a[i]++
				if a[i] == 0 {
					log.Fatal(a)
				}
				runtime.Gosched()
			}
		}(i)
	}
	time.Sleep(10000 * time.Millisecond)
	fmt.Println(a)
}

func Convert() {

	// 强制类型转换
	var a = []float64{1.2, 3.14, 5}
	var b []int
	b = ((*[1 << 20]int)(unsafe.Pointer(&a[0])))[:len(a):cap(a)]
	log.Println(b)

	c := []int{1, 2, 3} // 有3个元素的切片, len和cap都为3
	e := c[0 : 2 : cap(c)+1]
	log.Println(e)
}
