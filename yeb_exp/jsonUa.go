package yeb_exp

import (
	"context"
	"log"
	"time"
)

type jua struct {
	ua    string
	abort func()
	ctx   context.Context
}

func (j *jua) provide(ch chan *jua) {
	time.Sleep(3 * time.Second)
LOOP:
	for i := 0; i < 5; i++ {
		ch <- j
		select {
		case <-time.After(1 * time.Second):
		case <-j.ctx.Done():
			break LOOP
		}
	}
}

func getJsonUa() <-chan *jua {
	ch := make(chan *jua, 10)
	go func() {
		for ua := range freshJsonUa() {
			ctx, cancel := context.WithCancel(appCtx)
			j := &jua{
				ua:    ua,
				abort: cancel,
				ctx:   ctx,
			}
			log.Println(ua)
			go j.provide(ch)
		}
	}()
	return ch
}

func freshJsonUa() <-chan string {
	return make(chan string)
}
