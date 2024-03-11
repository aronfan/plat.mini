package xhub

import "sync"

type MessageSub struct {
	key    string
	cb     func(msg any) error
	ch     chan any
	initWG sync.WaitGroup
	finiWG sync.WaitGroup
}

func NewMessageSub(key string, cb func(any) error) *MessageSub {
	return &MessageSub{
		key:    key,
		cb:     cb,
		ch:     make(chan any, 10),
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

		sub.initWG.Done()

	OuterLoop:
		for {
			select {
			case msg, ok := <-sub.ch:
				if !ok {
					break OuterLoop
				}
				err := sub.cb(msg)
				if err != nil {
					break OuterLoop
				}
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
