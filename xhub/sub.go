package xhub

import (
	"sync"
	"time"
)

type MessageSub struct {
	key    string
	to     time.Duration
	ch     chan any
	event  func(any, bool) error
	timer  func() error
	initWG sync.WaitGroup
	finiWG sync.WaitGroup
}

func NewMessageSub(key string, event func(any, bool) error, timer func() error) *MessageSub {
	return &MessageSub{
		key:    key,
		to:     60 * time.Second,
		ch:     make(chan any, 10),
		event:  event,
		timer:  timer,
		initWG: sync.WaitGroup{},
		finiWG: sync.WaitGroup{},
	}
}

func (sub *MessageSub) Start() {
	sub.initWG.Add(1)
	sub.finiWG.Add(1)

	go func() {
		defer func() {
			sub.finiWG.Done()
		}()

		t := time.NewTimer(sub.to)
		sub.initWG.Done()

	OuterLoop:
		for {
			select {
			case msg, ok := <-sub.ch:
				if !ok {
					sub.event(nil, true)
					break OuterLoop
				}
				err := sub.event(msg, false)
				if err != nil {
					break OuterLoop
				}
			case <-t.C:
				err := sub.timer()
				if err != nil {
					break OuterLoop
				}
				t.Reset(sub.to)
			}
		}
	}()

	sub.initWG.Wait()
}

func (sub *MessageSub) Stop() {
	close(sub.ch)
}

func (sub *MessageSub) Wait() {
	sub.finiWG.Wait()
}

func (sub *MessageSub) Pub(msg any) {
	sub.ch <- msg
}
