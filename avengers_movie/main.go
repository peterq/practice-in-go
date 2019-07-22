package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"
)

var globalCtx context.Context
var cancelGlobalCtx context.CancelFunc
var dingUrl = "https://oapi.dingtalk.com/robot/send?access_token=40da525c975c6a7d433eaf85854e12b0a99579348ae54775a382694501b1b7f0"

//var dingUrl = "https://oapi.dingtalk.com/robot/send?access_token=d50545dca18bcaa92e507c8e27cb1b6f78139050eb40bffeef5cde4303954e47"
var resultUrl = "https://maoyan.com/cinema/15658?poi=80370653&movieId=248172"
var dingBotNotify = newDingBot(dingUrl)
var html []byte
var cookies = []*http.Cookie{
	{
		Name:  "JSESSIONID",
		Value: "23A6BC7FF315CD8A582852A5B17D3136.node1",
	},
}
var result string
var resultLock sync.Mutex

var httpClient = http.Client{
	Transport:     nil,
	CheckRedirect: nil,
	Jar:           nil,
	Timeout:       0,
}

func main() {
	log.Println("启动成功")
	globalCtx, cancelGlobalCtx = context.WithCancel(context.Background())

	httpClient.Jar, _ = cookiejar.New(&cookiejar.Options{
		PublicSuffixList: nil,
	})
	u, _ := url.Parse("http://jwglnew.hunnu.edu.cn")
	log.Println(u, cookies)
	httpClient.Jar.SetCookies(u, cookies)

	go checkResult()
	go heartBeat()
	<-globalCtx.Done()
}

func checkResult() {
	for true {
		currentResult := grab()
		resultLock.Lock()
		resultLock.Unlock()
		if result != currentResult {
			log.Println("有变化")
			result = currentResult
			sendToDing("有变化啦!!!!")
		} else {
			log.Println("没变化")
		}
		// 每分钟查一次
		time.Sleep(time.Minute)
	}
}

func sendToDing(title string) {
	resultLock.Lock()
	defer resultLock.Unlock()
	if len(result) < 5 {
		return
	}
	js := `{
     "msgtype": "markdown",
     "markdown": {"title":"", "text":"" }}`
	postData := map[string]interface{}{}
	json.Unmarshal([]byte(js), &postData)
	md := `
# %s
%s

`
	md = fmt.Sprintf(md, title, result)
	dingBotNotify.sendMarkDown(title, md)
}

func heartBeat() {
	time.Sleep(5 * time.Second)
	for true {
		nextTime := time.Now().Add(time.Hour * 2)
		sendToDing("我还活着, 下次时间:" + nextTime.Format("15:04:05"))
		// 2 小时发送一次消息, 用来确认还活着
		time.Sleep(2 * time.Hour)
	}
}

func grab() string {
	req, _ := http.NewRequest("GET", resultUrl, nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	content := string(bin)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(content))
	return doc.Find("#app > div.show-list.active > div.show-date").Text()
}

type dingBot struct {
	apiUrl string
}

type tjson map[string]interface{}

func newDingBot(url string) *dingBot {
	return &dingBot{
		apiUrl: url,
	}
}

func (bot *dingBot) sendMarkDown(title, content string) {
	postData := tjson{
		"msgtype": "markdown",
		"markdown": tjson{
			"title": title,
			"text":  "# " + title + "\n\n" + content,
		},
	}
	postRaw, _ := json.Marshal(postData)
	resp, err := http.Post(bot.apiUrl, "application/json", bytes.NewReader(postRaw))
	if err != nil {
		log.Println(err)
		return
	}
	bin, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(bin))
}
