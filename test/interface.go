package test

import (
	"fmt"
)

type People interface {
	Speak(string) string
}

type Stduent struct{}

func (stu *Stduent) Speak(think string) (talk string) {
	if think == "bitch" {
		talk = "You are a good boy"
	} else {
		talk = "hi"
	}
	return
}

func test() {
	var peo People = &Stduent{}
	//var peo People = Stduent{} // 编译不通过
	think := "bitch"
	fmt.Println(peo.Speak(think))
}
