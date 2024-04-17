package xactor

import (
	"fmt"
	"sync"
	"time"

	"github.com/vladopajic/go-actor/actor"
)

type deleteCtx struct {
	expires int64
}

func (ctx *deleteCtx) NotExpired() bool {
	return ctx.expires >= time.Now().Unix()
}

func (ctx *deleteCtx) Expired() bool {
	return ctx.expires < time.Now().Unix()
}

type AgentManager struct {
	lock   sync.RWMutex
	agents map[string]*Agent
	delete map[string]*deleteCtx

	inMbx actor.MailboxReceiver[any]
	actor actor.Actor

	timer    *time.Timer
	duration time.Duration
}

func (am *AgentManager) Add(k string, agent *Agent) (bool, bool) {
	am.lock.Lock()
	defer am.lock.Unlock()

	if am.atDel(k) {
		return false, true
	} else {
		delete(am.delete, k)
	}

	if _, ok := am.agents[k]; ok {
		return false, false
	}

	am.agents[k] = agent
	return true, false
}

func (am *AgentManager) MarkDel(k string, expires int64) *Agent {
	am.lock.Lock()
	defer am.lock.Unlock()

	if am.atDel(k) {
		return nil
	} else {
		delete(am.delete, k)
	}

	if agent, ok := am.agents[k]; ok {
		am.delete[k] = &deleteCtx{expires: expires}
		return agent
	}
	return nil
}

func (am *AgentManager) atDel(k string) bool {
	_, ok := am.agents[k]
	if ok {
		if _, ok := am.delete[k]; ok {
			return true
		}
	}
	return false
}

func (am *AgentManager) AtDel(k string) bool {
	am.lock.RLock()
	defer am.lock.RUnlock()

	return am.atDel(k)
}

func (am *AgentManager) Del(k string, ag *Agent) bool {
	am.lock.Lock()
	defer am.lock.Unlock()

	if am.atDel(k) {
		if agent, ok := am.agents[k]; ok && (agent == ag) {
			delete(am.delete, k)
			delete(am.agents, k)
			return true
		}
	}
	return false
}

func (am *AgentManager) Val(k string) (*Agent, bool) {
	am.lock.RLock()
	defer am.lock.RUnlock()

	if am.atDel(k) {
		return nil, true
	} else {
		if agent, ok := am.agents[k]; ok {
			return agent, false
		} else {
			return nil, false
		}
	}
}

func (am *AgentManager) Len() int {
	am.lock.RLock()
	defer am.lock.RUnlock()

	return len(am.agents)
}

func (am *AgentManager) DoWork(c actor.Context) actor.WorkerStatus {
	select {
	case <-c.Done():
		return actor.WorkerEnd
	case msg, ok := <-am.inMbx.ReceiveC():
		if ok {
			fmt.Println(msg)
		}
		return actor.WorkerContinue
	case <-am.timer.C:
		am.doCleanup()
		am.timer.Reset(am.duration)
		return actor.WorkerContinue
	}
}

func (am *AgentManager) SetTimer(d time.Duration) {
	am.duration = d
	am.timer = time.NewTimer(am.duration)
}

func (am *AgentManager) doCleanup() {
	var keys []string

	am.lock.Lock()
	for key, val := range am.delete {
		if val.Expired() {
			keys = append(keys, key)
		}
	}
	for _, key := range keys {
		delete(am.delete, key)
		delete(am.agents, key)
	}
	am.lock.Unlock()

	for _, key := range keys {
		fmt.Println("clean:", key)
	}
}

func (am *AgentManager) Start() {
	if am.actor == nil {
		actor := actor.New(am)
		am.actor = actor
		actor.Start()
	}
}

func (am *AgentManager) Stop() {
	actor := am.actor
	if actor != nil {
		actor.Stop()
		am.actor = nil
	}
}

func NewAgentManager() *AgentManager {
	return &AgentManager{
		lock:   sync.RWMutex{},
		agents: make(map[string]*Agent),
		delete: make(map[string]*deleteCtx),
		inMbx:  actor.NewMailbox[any](actor.OptAsChan()),
		actor:  nil,
	}
}
