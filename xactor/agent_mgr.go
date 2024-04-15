package xactor

import (
	"sync"
)

type AgentManager struct {
	lock   sync.RWMutex
	agents map[string]*Agent
}

func (am *AgentManager) Add(k string, agent *Agent) bool {
	am.lock.Lock()
	defer am.lock.Unlock()

	_, ok := am.agents[k]
	if ok {
		return false
	} else {
		am.agents[k] = agent
		return true
	}
}

func (am *AgentManager) Del(k string) *Agent {
	am.lock.Lock()
	defer am.lock.Unlock()

	agent, ok := am.agents[k]
	if ok {
		delete(am.agents, k)
		return agent
	} else {
		return nil
	}
}

func (am *AgentManager) Val(k string) *Agent {
	am.lock.RLock()
	defer am.lock.RUnlock()

	agent, ok := am.agents[k]
	if ok {
		return agent
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
	return &AgentManager{lock: sync.RWMutex{}, agents: make(map[string]*Agent)}
}
