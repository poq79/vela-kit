package chanutil

import (
	"fmt"
	"sync/atomic"
)

const (
	OPEN uint32 = iota + 1
	ClOSE
)

type StateType uint32

type Queue struct {
	state uint32
	count uint64
	value chan interface{}
}

/*
	for {
		select {
			case <-tx
			case <-tx
		}
	}


*/

func NewQueue(size int) *Queue {
	if size == 0 {
		return &Queue{
			state: OPEN,
			count: 0,
			value: make(chan interface{}),
		}
	}

	return &Queue{
		state: OPEN,
		count: 0,
		value: make(chan interface{}, size),
	}
}

func (q *Queue) Open() bool {
	return atomic.LoadUint32(&q.state) == OPEN
}

func (q *Queue) Close() error {
	if atomic.LoadUint32(&q.state) != OPEN {
		return nil
	}
	atomic.StoreUint32(&q.state, ClOSE)
	close(q.value)
	return nil
}

func (q *Queue) Publish(v interface{}) error {
	if atomic.LoadUint32(&q.state) != OPEN {
		return fmt.Errorf("queue close")
	}
	q.value <- v
	return nil
}

type Subscriber interface {
	OnMessage(v interface{})
	OnClose()
}

func (q *Queue) Subscribe(sub Subscriber) {
	for v := range q.value {
		sub.OnMessage(v)
	}

	sub.OnClose()
}
