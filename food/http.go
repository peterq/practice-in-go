package food

import (
	"funny/proxy"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	httpClient http.Client
	reqCh      chan *Req
}

type Req struct {
	httpReq  *http.Request
	resultCh chan *Result
}

type Result struct {
	resp *http.Response
	err  error
	data string
}

func IpPoolStart() {
	defer MetaCollection.Upsert(bson.M{"_id": "ip_pool"}, proxy.OkPool)
	var p map[string]*proxy.Proxy
	err := MetaCollection.FindId("ip_pool").One(&p)
	if err == nil {
		proxy.OkPool = p
		for k, v := range p {
			log.Println(k, *v)
		}
	}
	go proxy.CheckOKPool(15*time.Second, func() {
		MetaCollection.Upsert(bson.M{"_id": "ip_pool"}, proxy.OkPool)
	})
	for { // 每1/2分钟爬取一次代理
		proxy.RefreshPool(nil)
		time.Sleep(30 * time.Second)
	}
}

func NewClient(actorNum int) Client {
	client := Client{
		reqCh: make(chan *Req, actorNum),
		httpClient: http.Client{
			Transport: &http.Transport{Proxy: func(_ *http.Request) (*url.URL, error) {
				p := proxy.GetAProxy()
				if p == nil {
					return nil, nil
				}
				//根据定义Proxy func(*Request) (*url.URL, error)这里要返回url.URL
				return url.Parse("http://" + p.Host + ":" + p.Port + "/")
			}},
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       5 * time.Second,
		},
	}
	for i := 0; i < actorNum; i++ {
		go client.consumer()
	}
	return client
}

func (c *Client) consumer() { // 请求消费者
	for {
		req := <-c.reqCh                          // 读取请求
		result := new(Result)                     // 构造结果
		resp, err := c.httpClient.Do(req.httpReq) // 执行请求
		result.resp = resp
		if err != nil { // 请求出错返回错误, 接受下一个请求
			result.err = err
			req.resultCh <- result
			continue
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			result.err = err
			req.resultCh <- result
			continue
		}
		result.data = string(data) // 请求成功, 发送结果
		req.resultCh <- result
	}
}

func (c *Client) Get(url string, retry int) (data string, err error) {
	req := new(Req)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.httpReq = request
	req.resultCh = make(chan *Result, 1)
	var result *Result
	for retry++; retry > 0; retry-- {
		c.reqCh <- req          // 添加请求到队列
		result = <-req.resultCh // 等待请求被执行
		data = result.data
		err = result.err
		if err == nil && result.resp.StatusCode == 200 {
			return
		}
	}
	return
}
