package idea

import "container/list"

type Queue struct {
	out, in         chan interface{}   // 出, 进通道
	actionCh        chan *queueAction  // 队列操作请求通道
	members         list.List          // 真实的队列
	current         interface{}        // 当前值
	receivers       list.List          // 接受者队列
	currentReceiver <-chan interface{} // 当前接受者
	init            bool               // 是否已经初始化
}

const (
	actionEnqueue = iota // 入队
	actionDequeue        // 出队
	actionLen            // 查长度
)

type queueAction struct {
	action  int
	payload interface{}
}

func (q *Queue) _init() {
	if q.init {
		return
	}
	go func() {
		for {
			if q.current != nil && q.currentReceiver != nil { // 有接受者而且队列里有值
				q.currentReceiver <- q.current
				q.current = nil
				q.currentReceiver = nil
				if q.members.Len() > 0 {
					q.current = q.members
				}
			}
		}
	}()
}

func (q *Queue) Enqueue(payload interface{}) {
	q.actionCh <- &queueAction{
		payload: payload,
		action:  actionEnqueue,
	}
}

func (q *Queue) Dequeue(payload interface{}) interface{} {
	receive := make(chan interface{})
	q.actionCh <- &queueAction{
		action:  actionDequeue,
		payload: receive,
	}
	return <-receive
}

func NewQueue() *Queue {
	q := new(Queue)
	q._init()
	return q
}
