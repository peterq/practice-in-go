package yeb_exp

import (
	"log"
	"net/http"
)

func doCallApi(mobileCh <-chan string, uaCh <-chan *jua) {
	for mobile := range getUser() {
		jua := <-uaCh
		req, _ := http.NewRequest("GET", "https://promoprod.alipay.com/campaign/lotteryWithLogonInfo.json", nil)
		q := req.URL.Query()
		for k, v := range appConfig.InviteParam {
			q.Set(k, v)
		}
		q.Set("bindMobile", mobile)
		q.Set("json_ua", jua.ua)
		req.URL.RawQuery = q.Encode()
		res, err := apiClient.Do(req, 0)
		req.Header.Set("User-Agent", "Mozil                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         la/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.80 Mobile Safari/537.36")
		if err != nil {
			log.Println(err)
		}
		log.Println(mobile, res.Body, req.URL.String())
		//time.Sleep(3 * time.Second)
	}
}
