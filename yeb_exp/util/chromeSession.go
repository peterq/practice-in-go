package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
)

type ChromeSession struct {
	finish                context.CancelFunc
	executor              cdp.Executor
	cdp                   *chromedp.CDP
	cookie                string
	ctx                   context.Context
	user                  string
	pass                  string
	th                    *chromedp.TargetHandler
	needCapture           bool
	err                   error
	captureInfo           *TJson
	captureFrameNode      *cdp.Node
	bigImg                *[]byte
	smallImg              *[]byte
	captchaFrameContextId runtime.ExecutionContextID
	uaCh                  chan string
	mobileCh              <-chan string
}

func (ss *ChromeSession) handleEvent(evt chromedp.Evt) {
	//log.Println(fmt.Sprintf("%#v %#v", evt.Msg, evt.Evt))
	//log.Println(evt.Msg.Method)
	switch e := evt.Evt.(type) {
	case *network.EventResponseReceived:
		ss.handReceiveResponse(e)
	case *network.EventRequestWillBeSent:
		ss.handRequestWillBeSent(e)
	case *runtime.EventConsoleAPICalled:
		ss.handConsoleAPICalled(e)
	case *page.EventFrameNavigated:
		ss.handleFrameNavigated(e)
	case *runtime.EventExecutionContextCreated:
		ss.handleContextCreated(e)
	}
}

// js 执行上下文被创建
func (ss *ChromeSession) handleContextCreated(e *runtime.EventExecutionContextCreated) {
	aux := TJson{}
	json.Unmarshal(e.Context.AuxData, &aux)
	b := true
	if aux["isDefault"].(bool) {
		chromedp.Evaluate(StrReplace(GetAssetStr("/inject.js"), map[string]interface{}{
			"[contextId]": e.Context.ID,
			"[frameId]":   aux["frameId"],
			"[u]":         ss.user,
			"[p]":         ss.pass,
		}), &b, func(params *runtime.EvaluateParams) *runtime.EvaluateParams {
			return params.
				WithContextID(e.Context.ID).
				WithIncludeCommandLineAPI(true)
		}).Do(ss.ctx, ss.th)
	}
}

// iframe url 确定
func (ss *ChromeSession) handleFrameNavigated(evt *page.EventFrameNavigated) {
}

func (ss *ChromeSession) getNodeBox(id cdp.NodeID) (x, y int64, box *dom.BoxModel, err error) {
	box, err = dom.GetBoxModel().WithNodeID(id).Do(ss.ctx, ss.th)
	if err != nil {
		return
	}
	c := len(box.Content)
	for i := 0; i < c; i += 2 {
		x += int64(box.Content[i])
		y += int64(box.Content[i+1])
	}
	x /= int64(c / 2)
	y /= int64(c / 2)
	return
}

// 搜索指定节点
func matchBackendNodeId(node *cdp.Node, id cdp.BackendNodeID) *cdp.Node {
	if node.BackendNodeID == id {
		return node
	}
	if node.ChildNodeCount > 0 {
		if node.ContentDocument != nil {
			if match := matchBackendNodeId(node.ContentDocument, id); match != nil {
				return match
			}
		}
		for _, n := range node.Children {
			match := matchBackendNodeId(n, id)
			if match != nil {
				return match
			}
		}
	}
	return nil
}

// http 消息拦截
func (ss *ChromeSession) handReceiveResponse(evt *network.EventResponseReceived) {
	if strings.Contains(evt.Response.URL, "https://promoprod.alipay.com/campaign/lotteryWithLogonInfo.json") {
		log.Println(ss.getResponseStr(evt.RequestID))
		time.Sleep(1 * time.Second)
		//log.Println(<-ss.mobileCh)
		//var b bool
		//ss.cdp.Run(ss.ctx, chromedp.Evaluate(`location.reload()`, &b))
		//ss.cdp.Run(ss.ctx, chromedp.Evaluate(`notifyChromedp('ua.ok', {ua: window.json_ua})`, &b))
	}
}

func (ss *ChromeSession) getNodeInFrame(selector string) (cdp.NodeID, error) {
	if ss.captureFrameNode == nil {
		ro, _, e := (&runtime.EvaluateParams{
			Expression:    `document.querySelector('#tcaptcha_iframe')`,
			ReturnByValue: false,
		}).
			Do(ss.ctx, ss.th)
		LogErr(e)
		iframe, e := dom.DescribeNode().WithObjectID(ro.ObjectID).Do(ss.ctx, ss.th)
		LogErr(e)
		tree, e := dom.GetDocument().WithDepth(-1).WithPierce(true).Do(ss.ctx, ss.th)
		LogErr(e)
		// 搜索到iframe node
		ss.captureFrameNode = matchBackendNodeId(tree, iframe.BackendNodeID).ContentDocument.Children[1]
	}
	return dom.QuerySelector(ss.captureFrameNode.NodeID, selector).Do(ss.ctx, ss.th)
}

// 请求即将发送
func (ss *ChromeSession) handRequestWillBeSent(evt *network.EventRequestWillBeSent) {
}

// console被调用, 用来通信
func (ss *ChromeSession) handConsoleAPICalled(evt *runtime.EventConsoleAPICalled) {
	if evt.Type != "debug" || len(evt.Args) != 1 || evt.Args[0].Type != "string" {
		return
	}
	str := ""
	json.Unmarshal(evt.Args[0].Value, &str)
	//log.Println(str, strings.Index(str, "__notify__"))
	if strings.Index(str, "__notify__") == 0 {
		str = strings.Replace(str, "__notify__", "", 1)
		js := new(TJson)
		err := js.UnmarshalJSON([]byte(str))
		if err == nil {
			ss.handleNotify(js)
		}
	}
}

// 页面传递消息回调
func (ss *ChromeSession) handleNotify(js *TJson) {
	msg := *js
	if msg["type"].(string) == "ua.ok" {
		//ss.uaCh <- msg["data"].(map[string]interface {})["ua"].(string)
		m := <-ss.mobileCh
		log.Println(m)
		return
		var b bool
		time.Sleep(3 * time.Second)
		ss.cdp.Run(ss.ctx, chromedp.Evaluate(`document.querySelector('#ant-render-id-pages_outside_components_mobile_mobile > div > input').value=''`, &b))
		ss.cdp.Run(ss.ctx, chromedp.SendKeys("#ant-render-id-pages_outside_components_mobile_mobile > div > input", m))
		return
		time.Sleep(time.Second)
		ss.cdp.Run(ss.ctx, chromedp.Evaluate(fmt.Sprintf(`doSubmit('%s')`, m), &b))
	}
}

// 获取响应文本
func (ss *ChromeSession) getResponseStr(id network.RequestID) string {
	bin, err := network.GetResponseBody(id).Do(ss.ctx, ss.executor)
	if err != nil {
		return ""
	}
	return string(bin)
}

// 获取响应二进制数据
func (ss *ChromeSession) getResponse(id network.RequestID) []byte {
	bin, err := network.GetResponseBody(id).Do(ss.ctx, ss.executor)
	if err != nil {
		return []byte{}
	}
	return bin
}
