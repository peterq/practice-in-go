package main

import (
	//"time"
	"io/ioutil"
	"log"
	"net/http"
	"syscall/js"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//defer func() {recover()}()
	resp, err := http.Get("http://127.0.0.1:8080")
	handleErr(err)
	d, err := ioutil.ReadAll(resp.Body)
	handleErr(err)
	jsLog(string(d))
	jsLog("Hello world Go/wasm!")
	//js.Global().Get("document").Call("getElementById", "app").Set("innerText", time.Now().String())
}

func handleErr(err error) {
	if err != nil {
		js.Global().Get("console").Call("log", err.Error())
		panic(err)
	}
}

func jsLog(d string) {
	js.Global().Get("console").Call("log", string(d))
	log.Println(d)
}
