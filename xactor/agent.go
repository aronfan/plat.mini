package xactor

import (
	"time"

	"github.com/vladopajic/go-actor/actor"
)

var (
	agentIdleTimeout = 300 * time.Second
)

type Agent struct {
	key    string
	last   time.Time
	actor  actor.Actor
	report bool

	// message
	inMbx  actor.MailboxReceiver[any]
	outMbx actor.MailboxSender[any]
	fnCall func(*Call)
	fnDone func()

	// timer
	fnTimer  func()
	timer    *time.Timer
	duration time.Duration
	idlechan chan any
}

type AgentOption struct {
	Key      string
	CallFn   func(*Call)
	DoneFn   func()
	TimerFn  func()
	Duration time.Duration
	Idlechan chan any
}

func (agent *Agent) DoWork(c actor.Context) actor.WorkerStatus {
	select {
	case <-c.Done():
		if agent.timer != nil {
			agent.timer.Stop()
			agent.timer = nil
		}
		if agent.fnDone != nil {
			agent.fnDone()
		}
		return actor.WorkerEnd
	case msg, ok := <-agent.inMbx.ReceiveC():
		if ok {
			switch t := msg.(type) {
			case Responser:
				agent.SetLast(time.Now())
				if v, ok := t.(*Call); ok {
					if agent.fnCall != nil {
						agent.fnCall(v)
					} else {
						v.Response(0, nil)
					}
				}
			}
		}
		return actor.WorkerContinue
	case <-agent.timer.C:
		if agent.fnTimer != nil {
			agent.fnTimer()
		}
		agent.timer.Reset(agent.duration)
		if !agent.report && (time.Since(agent.last) >= agentIdleTimeout) {
			select {
			case agent.idlechan <- agent.key:
				agent.report = true
			default:
			}
		}
		return actor.WorkerContinue
	}
}

func (agent *Agent) GetKey() string {
	return agent.key
}

func (agent *Agent) Post(req any) {
	actor.Idle(actor.OptOnStart(func(c actor.Context) {
		agent.outMbx.Send(c, req)
	})).Start()
}

func (agent *Agent) Call(req any) (int, any, error) {
	call := newCall(req)
	actor.Idle(actor.OptOnStart(func(c actor.Context) {
		agent.outMbx.Send(c, call)
	})).Start()
	return call.WaitCall()
}

func (agent *Agent) SetLast(t time.Time) {
	agent.last = t
}

func (agent *Agent) Start() {
	if agent.actor == nil {
		actor := actor.New(agent)
		agent.actor = actor
		agent.timer = time.NewTimer(agent.duration)
		agent.last = time.Now()
		actor.Start()
	}
}

func (agent *Agent) Stop() {
	actor := agent.actor
	if actor != nil {
		actor.Stop()
		agent.actor = nil
	}
}

func NewAgentWithOption(opt *AgentOption) *Agent {
	mbx := actor.NewMailbox[any](actor.OptAsChan(), actor.OptCapacity(100))
	return &Agent{
		key:      opt.Key,
		last:     time.Now(),
		actor:    nil,
		report:   false,
		inMbx:    mbx,
		outMbx:   mbx,
		fnCall:   opt.CallFn,
		fnDone:   opt.DoneFn,
		fnTimer:  opt.TimerFn,
		timer:    nil,
		duration: opt.Duration,
		idlechan: opt.Idlechan,
	}
}
