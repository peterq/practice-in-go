package spider_client

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func RefreshPool(interval time.Duration, outputCh chan Proxy) chan Proxy {
	checked := make(map[Proxy]bool)
	c := New(2, 2, 0, true)
	checkCh := make(chan Proxy, 10)
	for i := 0; i < 100; i++ { // 开 50 个检测代理的协程
		go checkProxy(checkCh, outputCh)
	}

	go func() {
		lastClearTime := time.Now().UnixNano()
		clearInterval := int64(time.Minute)
		for {
			list := getProxyList(c)
			for _, p := range list {
				if _, ok := checked[p]; !ok { // 防止重复检测
					checkCh <- p
				}
			}
			time.Sleep(interval)
			if time.Now().UnixNano()-lastClearTime > clearInterval {
				log.Println("清空已检测代理")
				checked = make(map[Proxy]bool) // 清空它, 有些代理不稳定, 当时不能用, 现在又可以用了
				lastClearTime = time.Now().UnixNano()
			}
		}
	}()
	return checkCh
}

func checkProxy(input chan Proxy, output chan Proxy) {
	var p Proxy
	httpClient := http.Client{
		Transport: &http.Transport{Proxy: func(_ *http.Request) (*url.URL, error) {
			return url.Parse(string(p))
		}},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       5 * time.Second,
	}
	for {
		p = <-input
		codes := []string{
			"https://food.boohee.com/fb/v1/foods/fd645f00",
			"https://www.baidu.com",
			"https://www.baidu.com",
			"https://www.baidu.com",
			"https://www.baidu.com",
		}
		success := 0
		for _, code := range codes {
			request, err := http.NewRequest("GET", code, nil)
			if err != nil {
				//log.Println(err)
				continue
			}
			resp, err := httpClient.Do(request)
			if err != nil {
				//log.Println(err)
				continue
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				//log.Println(err)
				continue
			}
			if len(data) < 10 {
				//log.Println("check ok", string(data))
				continue
			}
			success++
		}
		//log.Println(success)
		if success == len(codes) {
			output <- p
		}
	}
}

func getProxyList(c *Client) []Proxy {
	u := "http://www.90api.cn/vip.php?key=584198000&sl=3000"
	//u := "http://www.90api.cn/vip.php?key=584198450&sl=3000"
	re, err := c.Get(u, 2)
	list := make([]Proxy, 0)
	if err != nil {
		log.Println("代理接口调用出错", err)
		return list
	}
	lines := strings.Split(re.Body, "\r\n")
	for index := range lines {
		line := lines[index]
		//log.Println(line)
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		p := "http://" + parts[0] + ":" + parts[1] + "/"
		list = append(list, Proxy(p))
	}
	return list
}
