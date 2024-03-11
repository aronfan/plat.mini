package xhub

import "sync"

type MessageHub struct {
	lock sync.RWMutex
	subs map[string][]*MessageSub
}

func NewMessageHub() *MessageHub {
	return &MessageHub{
		lock: sync.RWMutex{},
		subs: make(map[string][]*MessageSub),
	}
}

func (hub *MessageHub) Add(key string, sub *MessageSub) {
	if key == "" || sub == nil {
		return
	}

	hub.lock.Lock()
	defer hub.lock.Unlock()

	subList, ok := hub.subs[key]
	if ok {
		idx := -1
		for i := 0; i < len(subList); i++ {
			if subList[i] == sub {
				idx = i
				break
			}
		}
		if idx == -1 {
			subList = append(subList, sub)
			hub.subs[key] = subList
		}
	} else {
		subList = append(subList, sub)
		hub.subs[key] = subList
	}
}

func (hub *MessageHub) Del(key string, sub *MessageSub) {
	if key == "" || sub == nil {
		return
	}

	hub.lock.Lock()
	defer hub.lock.Unlock()

	subList, ok := hub.subs[key]
	if ok {
		if len(subList) == 1 && subList[0] == sub {
			delete(hub.subs, key)
		} else {
			idx := -1
			for i := 0; i < len(subList); i++ {
				if subList[i] == sub {
					idx = i
					break
				}
			}
			if idx != -1 {
				newList := append(subList[:idx], subList[idx+1:]...)
				hub.subs[key] = newList
			}
		}
	}
}

func (hub *MessageHub) Pub(key string, msg string) {
	if key == "" || msg == "" {
		return
	}

	hub.lock.Lock()
	defer hub.lock.Unlock()

	subList, ok := hub.subs[key]
	if ok {
		for i := 0; i < len(subList); i++ {
			sub := subList[i]
			sub.Pub(msg)
		}
	}
}

func (hub *MessageHub) Stop(key string) {
	if key == "" {
		return
	}

	hub.lock.Lock()
	defer hub.lock.Unlock()

	subList, ok := hub.subs[key]
	if ok {
		for i := 0; i < len(subList); i++ {
			sub := subList[i]
			sub.Stop()
		}
	}
}

func (hub *MessageHub) Nums() (int, int) {
	hub.lock.RLock()
	defer hub.lock.RUnlock()

	keyNums := 0
	subNums := 0
	for _, v := range hub.subs {
		keyNums += 1
		subNums += len(v)
	}
	return keyNums, subNums
}
