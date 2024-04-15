package xactor

import (
	"sync"
	"time"
)

type deleteCtx struct {
	expires int64
}

func (ctx *deleteCtx) NotExpired() bool {
	return ctx.expires >= time.Now().Unix()
}

type AgentManager struct {
	lock   sync.RWMutex
	agents map[string]*Agent
	delete map[string]*deleteCtx
}

func (am *AgentManager) Add(k string, agent *Agent) bool {
	am.lock.Lock()
	defer am.lock.Unlock()

	_, ok := am.agents[k]
	if ok {
		return false
	} else {
		ctx, ok := am.delete[k]
		if ok {
			if ctx.NotExpired() {
				return false
			} else {
				delete(am.delete, k)
			}
		}
		am.agents[k] = agent
		return true
	}
}

func (am *AgentManager) MarkDel(k string, expires int64) *Agent {
	am.lock.Lock()
	defer am.lock.Unlock()

	agent, ok := am.agents[k]
	if ok {
		ctx, ok := am.delete[k]
		if ok {
			if ctx.expires < expires {
				ctx.expires = expires
			}
		} else {
			am.delete[k] = &deleteCtx{expires: expires}
		}
		return agent
	} else {
		delete(am.delete, k)
		return nil
	}
}

func (am *AgentManager) AtDel(k string) bool {
	am.lock.RLock()
	defer am.lock.RUnlock()

	_, ok := am.agents[k]
	if ok {
		ctx, ok := am.delete[k]
		if ok {
			return ctx.NotExpired()
		} else {
			return false
		}
	} else {
		return false
	}
}

func (am *AgentManager) Del(k string, ag *Agent) bool {
	am.lock.Lock()
	defer am.lock.Unlock()

	agent, ok := am.agents[k]
	if ok {
		if agent == ag {
			delete(am.delete, k)
			delete(am.agents, k)
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (am *AgentManager) Val(k string) *Agent {
	am.lock.RLock()
	defer am.lock.RUnlock()

	agent, ok := am.agents[k]
	if ok {
		_, ok := am.delete[k]
		if ok {
			return nil
		} else {
			return agent
		}
	} else {
		return nil
	}
}

func (am *AgentManager) Len() int {
	am.lock.RLock()
	defer am.lock.RUnlock()

	return len(am.agents)
}

func NewAgentManager() *AgentManager {
	return &AgentManager{lock: sync.RWMutex{},
		agents: make(map[string]*Agent), delete: make(map[string]*deleteCtx)}
}
