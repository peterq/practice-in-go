package proxy

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var pool = make(map[string]*Proxy)
var OkPool = make(map[string]*Proxy)
var PoolLock = new(sync.Mutex)

type Proxy struct {
	Host, Port, T string
	state         int // 0 未测试, 1 通过, 2 未通过
	LastUsingTime int64
}

// 从代理池中选出一个代理, 优先选取未被使用过的
func GetAProxy() *Proxy {
	return nil
	if len(OkPool) == 0 {
		return nil
	}
	var theOne *Proxy
	PoolLock.Lock()
	for _, p := range OkPool {
		if theOne == nil || theOne.LastUsingTime > p.LastUsingTime {
			theOne = p
		}
	}
	PoolLock.Unlock()
	theOne.LastUsingTime = time.Now().UnixNano()
	return theOne
}

func RefreshPool(cb func()) error {
	pool = make(map[string]*Proxy)
	checkCh := make(chan *Proxy, 10)   // 发送要检测的代理
	checkedCh := make(chan *Proxy, 10) // 已经检测完毕的代理
	available := checkAvailable(checkConnection(checkCh, checkedCh), checkedCh)
	go func() {
		for {
			p, ok := <-available
			if !ok {
				break
			}
			p.state = 1
			PoolLock.Lock()
			OkPool[ProxyKey(p)] = p
			PoolLock.Unlock()
		}
		log.Println("finised")
	}()
	for i := 1; i <= 1; i++ {
		list, err := getListFromXigua(i)
		if err != nil {
			return err
		}
		toCheck := make([]*Proxy, 0)
		for _, p := range list {
			key := ProxyKey(p)
			if _, ok := pool[key]; !ok {
				p.state = 2
				pool[key] = p
				toCheck = append(toCheck, p)
			}

		}
		done := waitNum(checkedCh, len(toCheck))
		for _, p := range toCheck {
			time.Sleep(10 * time.Millisecond)
			checkCh <- p
		}
		<-done
		log.Println("可用代理", len(OkPool))
	}
	log.Println("代理爬取完毕, 总数:", len(pool), "/", len(pool))
	return nil
}

func ProxyKey(p *Proxy) string {
	return p.Host + ":" + p.Port
}

func CheckOKPool(interval time.Duration, cb func()) {
	for {
		noUsefulNum := 0
		//
		for k, p := range OkPool {
			go func(k string, p *Proxy) {
				toDel := true
				defer func() {
					if toDel {
						//log.Println("代理失效", k)
						PoolLock.Lock()
						noUsefulNum++
						delete(OkPool, k)
						PoolLock.Unlock()
					}
				}()
				resp, err := proxyTest("http://"+ProxyKey(p)+"/", "https://food.boohee.com/fb/v1/foods/malachenpigourou")
				if err != nil {
					//log.Println(err)
					return
				}
				data, err := ioutil.ReadAll(resp.Body)
				if data != nil && err == nil {
					//log.Println(strings.Index(string(data), "星火米袋"), string(data))
					if strings.Index(string(data), "麻辣陈皮狗肉") > 0 {
						toDel = false
					}

				} else {
					log.Println(err)
				}
			}(k, p)
		}
		if cb != nil {
			go cb()
		}
		time.Sleep(interval)
	}
}

func proxyTest(proxy, u string) (resp *http.Response, err error) {
	proxyFn := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxy)
	}
	httpClient := http.Client{
		Transport:     &http.Transport{Proxy: proxyFn},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       5 * time.Second,
	}
	request, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}
	resp, err = httpClient.Do(request)
	return
}

func waitNum(checked chan *Proxy, num int) chan bool {
	done := make(chan bool)
	go func() {
		for i := 0; i < num; i++ {
			<-checked
		}
		done <- true
	}()
	return done

}

// 检测是否能连接上
func checkConnection(list chan *Proxy, checked chan *Proxy) chan *Proxy {
	ch := make(chan *Proxy, 10)
	go func() {
		for {
			//time.Sleep(200 *time.Millisecond)
			p, ok := <-list
			//log.Println(ProxyKey(p))
			if !ok {
				close(ch)
				break
			}
			//log.Println(p)
			go func(p *Proxy) {
				//log.Println(p)
				addr := p.Host + ":" + p.Port
				conn, err := net.DialTimeout("tcp", addr, time.Second*3)
				if err == nil { // 能连接上发送给下一个检测步骤
					conn.Close()
					//log.Println("连接上了", addr)
					ch <- p
				} else { // 不能连接上, 检测完毕
					//log.Println("连接不上", addr, err)
					checked <- p
				}
			}(p)
		}
	}()
	return ch
}

// 检测是否可用
func checkAvailable(ch chan *Proxy, checked chan *Proxy) chan *Proxy {
	available := make(chan *Proxy, 3)
	go func() {
		for {
			p, ok := <-ch
			if !ok {
				log.Println("not ok")
				break
			}
			go func(p *Proxy) {
				defer func() { checked <- p }() // 发送检测完毕事件
				resp, err := proxyTest("http://"+ProxyKey(p)+"/", "https://food.boohee.com/fb/v1/foods/malachenpigourou")
				if err != nil {
					//log.Println(err)
					return
				}
				data, err := ioutil.ReadAll(resp.Body)
				if data != nil && err == nil {
					str := string(data)
					//log.Println(str)
					if strings.Index(str, "麻辣陈皮狗肉") > 0 {
						available <- p
					}

				}
			}(p)
		}
	}()
	return available

}

func getListFromXici(page int) (list []*Proxy, err error) {
	client := http.Client{
		/*Transport:      &http.Transport{Proxy: func(_ *http.Request) (*url.URL, error) {
			return url.Parse("http://127.0.0.1:1081/")//根据定义Proxy func(*Request) (*url.URL, error)这里要返回url.URL
		}},*/
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       5 * time.Second,
	}
	u := "http://www.xicidaili.com/nn/" + strconv.Itoa(page)
	log.Println(u)
	request, err := http.NewRequest("GET", u, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.35 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err, "err")
		return
	} else {
		defer resp.Body.Close()
	}
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	if len(s) < 20 {
		log.Println(string(s))
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(s))
	doc.Find("#ip_list > tbody").Children().Each(func(i int, selection *goquery.Selection) {
		ip, _ := selection.Find("td:nth-child(2)").Html()
		port, _ := selection.Find("td:nth-child(3)").Html()
		t, _ := selection.Find("td:nth-child(6)").Html()
		p := &Proxy{ip, port, t, 0, 0}
		if ip != "" {
			list = append(list, p)
		}
	})
	return
}

func getListFrom66(page int) (list []*Proxy, err error) {
	each := 5
	client := http.Client{
		/*Transport:      &http.Transport{Proxy: func(_ *http.Request) (*url.URL, error) {
			return url.Parse("http://127.0.0.1:1081/")//根据定义Proxy func(*Request) (*url.URL, error)这里要返回url.URL
		}},*/
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       5 * time.Second,
	}
	for i := (page-1)*each + 1; i <= page*each; i++ {
		u := "http://www.66ip.cn/" + strconv.Itoa(i) + ".html"
		//log.Println(u)
		err = nil
		request, err := http.NewRequest("GET", u, nil)
		request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.35 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
		resp, err := client.Do(request)
		if err != nil {
			log.Println(err, "err")
			return nil, err
		} else {
			defer func(*http.Response) { resp.Body.Close() }(resp)
		}
		s, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		if len(s) < 20 {
			log.Println(string(s))
		}
		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(s))
		doc.Find("#main > div > div:nth-child(1) > table > tbody").Children().Each(func(i int, selection *goquery.Selection) {
			if i == 0 {
				return
			}
			ip, _ := selection.Find("td:nth-child(1)").Html()
			port, _ := selection.Find("td:nth-child(2)").Html()
			p := &Proxy{ip, port, "", 0, 0}
			if ip != "" {
				list = append(list, p)
			}
		})
	}
	return
}

func getListFromXigua(page int) (list []*Proxy, err error) {
	if page > 1 {
		return
	}
	each := 5
	client := http.Client{
		/*Transport:      &http.Transport{Proxy: func(_ *http.Request) (*url.URL, error) {
			return url.Parse("http://127.0.0.1:1081/")//根据定义Proxy func(*Request) (*url.URL, error)这里要返回url.URL
		}},*/
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       5 * time.Second,
	}
	for i := (page-1)*each + 1; i <= page*each; i++ {
		u := "http://api3.xiguadaili.com/ip/?tid=559006267855792&num=50000"
		u = "http://www.90api.cn/vip.php?key=584198450&sl=3000"
		//log.Println(u)
		err = nil
		request, err := http.NewRequest("GET", u, nil)
		request.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.35 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
		resp, err := client.Do(request)
		if err != nil {
			log.Println(err, "err")
			return nil, err
		} else {
			defer func(*http.Response) { resp.Body.Close() }(resp)
		}
		s, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		lines := strings.Split(string(s), "\r\n")
		for index := range lines {
			line := lines[index]
			//log.Println(line)
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				continue
			}
			p := &Proxy{parts[0], parts[1], "", 0, 0}
			//log.Println(ProxyKey(p))
			//time.Sleep(100 * time.Millisecond)
			list = append(list, p)
		}

	}
	return
}
