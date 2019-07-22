package util

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log"
)

func GetUa(ctx context.Context, link string, mobileCh <-chan string) (uaCh chan string) {
	uaCh = make(chan string)
	var err error
	// create context
	ctxt, cancel := context.WithCancel(ctx)

	// 启动 chrome
	c, err := chromedp.New(ctxt, chromedp.WithErrorf(func(s string, i ...interface{}) {}))
	if err != nil {
		log.Fatal(err)
	}
	ss := ChromeSession{
		finish:   cancel,
		executor: nil,
		cdp:      c,
		cookie:   "",
		ctx:      ctxt,
		uaCh:     uaCh,
		mobileCh: mobileCh,
	}

	// 事件监听
	c.Run(ctxt, network.Enable())
	err = c.Run(ctxt, chromedp.ActionFunc(func(c context.Context, executor cdp.Executor) error {
		go func() {
			ss.executor = executor
			th := executor.(*chromedp.TargetHandler)
			ss.th = th
		Receive:
			for evt := range th.EvtCh {
				go ss.handleEvent(evt)
				select {
				case <-ctxt.Done():
					break Receive
				default:
				}
			}
		}()
		return nil
	}))
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		defer cancel()
		// 分享页面
		err = c.Run(ctxt, chromedp.Navigate(link))
		if err != nil {
			log.Println(err)
			return
		}
		<-ctxt.Done()

		// 关闭chrome
		err1 := c.Shutdown(ctxt)
		if err1 != nil {
			log.Println(err1)
			return
		}

		// 等待chrome关闭
		err1 = c.Wait()
		if err1 != nil {
			log.Println(err1)
		}
	}()
	return
}
