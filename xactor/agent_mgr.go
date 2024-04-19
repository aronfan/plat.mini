package xactor

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/vladopajic/go-actor/actor"
)

var (
	ErrAgentAtDel error = errors.New("agent at del")
	ErrAgentExist error = errors.New("agent exists")
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
	idlechan chan any
}

type AgentManagerOption struct {
	Duration time.Duration
	Idlechan chan any
}

func (am *AgentManager) Add(k string, agent *Agent) error {
	am.lock.Lock()
	defer am.lock.Unlock()

	if am.atDel(k) {
		return ErrAgentAtDel
	} else {
		delete(am.delete, k)
	}

	if _, ok := am.agents[k]; ok {
		return ErrAgentExist
	}

	am.agents[k] = agent
	return nil
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

func (am *AgentManager) Val(k string) (*Agent, error) {
	am.lock.RLock()
	defer am.lock.RUnlock()

	if am.atDel(k) {
		return nil, ErrAgentAtDel
	} else {
		if agent, ok := am.agents[k]; ok {
			return agent, nil
		} else {
			return nil, nil
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
		if am.timer != nil {
			am.timer.Stop()
			am.timer = nil
		}
		return actor.WorkerEnd
	case msg, ok := <-am.inMbx.ReceiveC():
		if ok {
			fmt.Println(msg)
		}
		return actor.WorkerContinue
	case key, ok := <-am.idlechan:
		if ok {
			go am.onIdle(key.(string))
		}
		return actor.WorkerContinue
	case <-am.timer.C:
		am.doCleanup()
		am.timer.Reset(am.duration)
		return actor.WorkerContinue
	}
}

func (am *AgentManager) onIdle(key string) {
	expires := time.Now().Add(30 * time.Second).Unix()
	agent := am.MarkDel(key, expires)
	if agent != nil {
		agent.Call("flush")
		agent.Stop()
		am.Del(key, agent)
	}
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

func (am *AgentManager) GetIdlechan() chan any {
	return am.idlechan
}

func (am *AgentManager) Start() {
	if am.actor == nil {
		actor := actor.New(am)
		am.actor = actor
		am.timer = time.NewTimer(am.duration)
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

func (am *AgentManager) StopAgents(cb func(agent *Agent) bool) {
	am.lock.Lock()
	defer am.lock.Unlock()

	for key, agent := range am.agents {
		ok := cb(agent)
		if ok {
			delete(am.agents, key)
			delete(am.delete, key)
		}
	}
}

func NewAgentManager() *AgentManager {
	opt := &AgentManagerOption{
		Duration: 30 * time.Second,
		Idlechan: make(chan any, 10000),
	}
	return NewAgentManagerWithOption(opt)
}

func NewAgentManagerWithOption(opt *AgentManagerOption) *AgentManager {
	return &AgentManager{
		lock:     sync.RWMutex{},
		agents:   make(map[string]*Agent),
		delete:   make(map[string]*deleteCtx),
		inMbx:    actor.NewMailbox[any](actor.OptAsChan()),
		actor:    nil,
		timer:    nil,
		duration: opt.Duration,
		idlechan: opt.Idlechan,
	}
}
