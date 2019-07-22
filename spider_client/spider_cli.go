package spider_client

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// 发送一个请求需要构造的结构体
type Req struct {
	httpReq  *http.Request
	resultCh chan *Result
}

// 请求结果
type Result struct {
	Resp     *http.Response // 原生内容
	Err      error          // 错误信息
	Body     string         // Body
	UseProxy Proxy          // 使用的代理
}

const directProxy = "direct"

type Proxy string // 代理, 格式: http://127.0.0.1:8080/

type consumerMessage struct {
	proxy Proxy
	data  string
}

type proxyStatus struct {
	notifyConsumerCh chan string // 调度器协程向消费协程发消息的通道
	complainedTimes  int         // 被投诉过的次数
	consumerNum      int         // 消费者数量
}

type Client struct {
	//httpClient http.Client // 标准库的http客户端
	reqCh                          chan *Req              // 请求通道
	addProxyCh                     chan Proxy             // 代理通道, 有可用的代理时,发送到这个通道, 可以增加请求消费者, 从而提高并发
	complainProxyCh                chan Proxy             // 投诉代理通道
	consumerEveryProxy             int                    // 每个代理启动几个消费者
	assignedProxyMap               map[Proxy]*proxyStatus // 已经使用的消费者
	assignedProxyMapLock           sync.RWMutex           // map操作锁
	needDirect                     bool                   // 是否需要直连
	consumerMessageCh              chan consumerMessage   // 消费协程向调度器协程发消息的通道
	minIntervalPerProxyEveryTowReq time.Duration          // 2次请求最小间隔
}

func (c *Client) GetProxyCh() chan Proxy {
	return c.addProxyCh
}

func (c *Client) GetUsingProxy() []Proxy {
	list := make([]Proxy, 0)
	c.assignedProxyMapLock.RLock()
	for p := range c.assignedProxyMap {
		list = append(list, p)
	}
	c.assignedProxyMapLock.RUnlock()
	return list
}

func (c *Client) readProxyMap(key Proxy) (*proxyStatus, bool) {
	c.assignedProxyMapLock.RLock()
	p, ok := c.assignedProxyMap[key]
	c.assignedProxyMapLock.RUnlock()
	return p, ok
}

func (c *Client) writeProxyMap(key Proxy, v *proxyStatus) {
	c.assignedProxyMapLock.Lock()
	c.assignedProxyMap[key] = v
	c.assignedProxyMapLock.Unlock()
}

func (c *Client) delProxyMap(key Proxy) {
	c.assignedProxyMapLock.Lock()
	delete(c.assignedProxyMap, key)
	c.assignedProxyMapLock.Unlock()
}

func New(consumerEveryProxy, reqChSize int, minIntervalPerProxyEveryTowReq time.Duration, needDirect bool) *Client { // 创建一个新的客户端
	client := &Client{
		reqCh:                          make(chan *Req, reqChSize),
		addProxyCh:                     make(chan Proxy, 2),
		complainProxyCh:                make(chan Proxy, 2),
		consumerEveryProxy:             consumerEveryProxy,
		assignedProxyMap:               make(map[Proxy]*proxyStatus),
		needDirect:                     needDirect,
		consumerMessageCh:              make(chan consumerMessage, 10),
		minIntervalPerProxyEveryTowReq: minIntervalPerProxyEveryTowReq,
	}
	go client.startScheduler()
	return client
}

// 只调用一次
func (c *Client) startScheduler() {
	go func() {
		if c.needDirect {
			c.addProxyCh <- directProxy
		}
	}()
	for {
		select {
		case p := <-c.addProxyCh: // 有可用的代理了
			//log.Println(p)
			// 这个代理已经在使用了不进行分配
			if _, ok := c.readProxyMap(p); ok || len(c.assignedProxyMap) > 500 {
				continue
			}
			go c.proxyLeader(p)
		case p := <-c.complainProxyCh: // 有代理被投诉
			if s, ok := c.readProxyMap(p); ok {
				s.complainedTimes++                              // 增加投诉次数
				if s.complainedTimes >= c.consumerEveryProxy/4 { // 被投诉超过5次, 直接终止使用这个代理
					go func() { s.notifyConsumerCh <- "exit" }() // leader 可能在blocked中, 开协程发送指令
					c.delProxyMap(p)
					log.Println("可用代理", len(c.assignedProxyMap))
				} else { // 冻结代理
					go func() { c.assignedProxyMap[p].notifyConsumerCh <- "block" }()
				}
			}
		case m := <-c.consumerMessageCh:
			go c.onConsumerMessage(m) // 开协程处理消息
		}
	}
}

// 代理老大
func (c *Client) proxyLeader(p Proxy) {
	// 创建和协程通信的通道
	notifyConsumerCh := make(chan string)
	status := &proxyStatus{
		notifyConsumerCh: notifyConsumerCh,
		consumerNum:      c.consumerEveryProxy,
	}
	c.writeProxyMap(p, status)
	// 创建请求通道, 转发client请求通道的请求
	reqCh := make(chan *Req)
	msgChs := make([]chan string, c.consumerEveryProxy)
	for i := range msgChs {
		msgChs[i] = make(chan string, 2)
	}
	// 启动消费者
	for i := 0; i < status.consumerNum; i++ {
		go consumer(p, msgChs[i], reqCh)
	}
	// 请求转发
	exit := false
	go func() {
		for !exit {
			reqCh <- <-c.reqCh
			if c.minIntervalPerProxyEveryTowReq > 0 {
				//log.Println(c.minIntervalPerProxyEveryTowReq)
				time.Sleep(c.minIntervalPerProxyEveryTowReq)
			}
		}
	}()
	// 消息转发
	go func() {
		for {
			msg := <-status.notifyConsumerCh
			if msg == "block" { // 收到一次冻结命令, 就杀死一个消费者(有些代理可能有并发限制, 降低并发提高成功率)
				//log.Println("block", p)
				time.Sleep(time.Second)
				if len(msgChs) > 0 {
					msgChs[len(msgChs)-1] <- "exit"
					msgChs = msgChs[0 : len(msgChs)-1]
					status.consumerNum--
				}
			} else {
				for i := 0; i < status.consumerNum; i++ {
					msgChs[i] <- msg
				}
				if msg == "exit" {
					//log.Println("exit", p)
					exit = true
					break
				}
			}
		}
	}()
}

// 收到消费者消息的回调
func (c *Client) onConsumerMessage(m consumerMessage) {

}

// 请求消费者
func consumer(p Proxy, fromScheduler chan string, reqCh chan *Req) { // 请求消费者
	proxyFn := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(string(p))
	}
	if p == directProxy {
		proxyFn = nil
	}
	httpClient := http.Client{
		Transport:     &http.Transport{Proxy: proxyFn},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       15 * time.Second,
	}
DONE: // 轮训请求通道
	for {
		select {
		case req := <-reqCh:
			//log.Println("req", p)
			doReq(httpClient, req, p) // 这里不能开协程, 保证一个消费者同时只处理一个请求
		case msg := <-fromScheduler:
			if msg != "" {
				switch msg {
				case "exit":
					break DONE
				}
			}
		}
	}
	//log.Println("consumer exit", p)
}

// 执行请求
func doReq(client http.Client, req *Req, p Proxy) {
	result := new(Result)               // 构造结果
	resp, err := client.Do(req.httpReq) // 执行请求
	result.Resp = resp
	result.UseProxy = p
	if err != nil { // 请求出错返回错误
		//log.Println(err)
		result.Err = err
		req.resultCh <- result
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		result.Err = err
		req.resultCh <- result
		return
	}
	result.Body = string(data) // 请求成功, 发送结果
	req.resultCh <- result
}

func (c *Client) Get(url string, retry int) (result *Result, err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	result, err = c.Do(request, retry)
	return
}

func (c *Client) Do(r *http.Request, retry int) (result *Result, err error) {
	req := new(Req)
	req.httpReq = r
	req.resultCh = make(chan *Result, 1)
	for retry++; retry > 0; retry-- {
		c.reqCh <- req          // 添加请求到队列
		result = <-req.resultCh // 等待请求被执行
		err = result.Err
		if err != nil { // 如果请求失败了, 则投诉这个代理
			//log.Println(err)
			c.ComplainProxy(result.UseProxy)
		}
		if err == nil && result.Resp.StatusCode == 200 {
			return
		}
	}
	return
}

// 投诉代理
func (c *Client) ComplainProxy(proxy Proxy) {
	//return
	go func() {
		if proxy == directProxy {
			return
		}
		c.complainProxyCh <- proxy
	}()
}
