package dlt

import (
	"container/heap"
	"sync"
	"sync/atomic"
	"time"
)

type element struct {
	Key    string
	Index  int
	Expire int64
}

// 过期队列
type expirationQueue []*element

func (q expirationQueue) Len() int {
	return len(q)
}

func (q expirationQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].Index, q[j].Index = i, j
}

func (q expirationQueue) Less(i, j int) bool {
	return q[i].Expire < q[j].Expire
}

func (q *expirationQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*q = old[0 : n-1]
	return item
}

func (q *expirationQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(*element)
	item.Index = n
	*q = append(*q, item)
}

func (q *expirationQueue) Top() *element {
	if q.Len() == 0 {
		return nil
	}
	return (*q)[0]
}

func (q *expirationQueue) Update(item *element, expire int64) {
	item.Expire = expire
	heap.Fix(q, item.Index)
}

func (q *expirationQueue) Remove(item *element) {
	heap.Remove(q, item.Index)
}

// 倒计时器
type CountdownTimer struct {
	seconds  int64
	mutex    sync.Mutex
	q        *expirationQueue
	mapper   map[string]*element
	callback func(string)
}

// 创建倒计时器
func NewCountdownTimer(seconds int64, callback func(string)) *CountdownTimer {
	q := make(expirationQueue, 0)
	timer := CountdownTimer{
		q:        &q,
		callback: callback,
		seconds:  seconds,
		mapper:   make(map[string]*element),
	}
	go timer.checkExpirationQueue()
	return &timer
}

// 更新时间
func (timer *CountdownTimer) Update(key string) {
	e := element{
		Key:    key,
		Expire: time.Now().Unix() + atomic.LoadInt64(&timer.seconds),
	}
	timer.Remove(key)

	timer.mutex.Lock()
	defer timer.mutex.Unlock()
	heap.Push(timer.q, &e)
	timer.mapper[key] = &e
}

// 删除计时
func (timer *CountdownTimer) Remove(key string) {
	timer.mutex.Lock()
	defer timer.mutex.Unlock()

	element, ok := timer.mapper[key]
	if !ok {
		return
	}
	delete(timer.mapper, key)
	timer.q.Remove(element)
}

// 设置过期时间
func (timer *CountdownTimer) SetExpiration(seconds int64) {
	if seconds > 0 {
		atomic.StoreInt64(&timer.seconds, seconds)
	}
}

// 检查过期队列
func (timer *CountdownTimer) checkExpirationQueue() {
	t := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-t.C:
			now := time.Now().Unix()
			keys := make([]string, 0)

			timer.mutex.Lock()
			for {
				if timer.q.Len() == 0 {
					break
				}

				top := timer.q.Top()
				if top == nil || top.Expire > now {
					break
				}
				timer.q.Remove(top)
				delete(timer.mapper, top.Key)
				keys = append(keys, top.Key)
			}
			timer.mutex.Unlock()

			if timer.callback != nil {
				for _, key := range keys {
					timer.callback(key)
				}
			}
		}
	}
}
