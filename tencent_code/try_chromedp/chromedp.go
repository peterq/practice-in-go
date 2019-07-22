package try_chromedp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func Test() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	err = c.Run(ctxt, click())
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

type tJson map[string]interface{}

func (t *tJson) MarshalJSON() ([]byte, error) {
	return json.Marshal(*t)
}
func (t *tJson) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, (*map[string]interface{})(t))
	log.Println(err)
	return err
}

func click() chromedp.Tasks {
	var str string
	return chromedp.Tasks{
		network.Enable(),
		network.SetRequestInterception([]*network.RequestPattern{
			{
				URLPattern:        "https://imgcache.qq.com/ptlogin/v4/style/42/images/loading.gif",
				ResourceType:      "",
				InterceptionStage: network.InterceptionStageRequest,
			},
		}),
		chromedp.Navigate(`https://xui.ptlogin2.qq.com/cgi-bin/xlogin?appid=756044602&style=42&s_url=http://wtlogin.qq.com/&pt_no_onekey=1&pt_no_auth=1&daid=499&wt_force_pwd=1`),
		chromedp.ActionFunc(func(c context.Context, cdp cdp.Executor) error {
			var res tJson
			err := cdp.Execute(c, runtime.CommandEvaluate, &tJson{
				"expression": "location",
			}, &res)
			th := cdp.(*chromedp.TargetHandler)
			go func() {
				for evt := range th.EvtCh {
					log.Println(fmt.Sprintf("%#v, %#v", evt.Evt, evt.Msg))
				}
			}()
			log.Println(err, res)
			return nil
		}),
		chromedp.WaitVisible("#u"),
		chromedp.Evaluate(`
document.querySelector('#u').value = "1056668021@qq.com"
document.querySelector('#p').value = "zhangyunxiao521"
        `, &str),
		chromedp.Click("#go"),
		chromedp.Sleep(1000 * time.Second),
	}
}
