package wasm

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func Init() {
	_, f, _, _ := runtime.Caller(0)
	dir := path.Dir(f)
	// 编译给js用的go程序
	//cmd := exec.Command("env")
	cmd := exec.Command("go", "build", "-o", dir+"/public/app.wasm", dir+"/go_js/main.go")
	os.Setenv("GOARCH", "wasm")
	os.Setenv("GOOS", "js")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
	cmd.Wait()

	// 启动一个web服务器
	http.ListenAndServe("0.0.0.0:8080", http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			if strings.HasSuffix(req.URL.Path, ".wasm") {
				resp.Header().Set("content-type", "application/wasm")
			}

			http.FileServer(http.Dir(dir+"/public")).ServeHTTP(resp, req)
		}))
}
