package xactor

import (
	"time"

	"github.com/vladopajic/go-actor/actor"
)

type Agent struct {
	key     string
	last    time.Time
	cleanup chan any

	// message
	inMbx  actor.MailboxReceiver[any]
	outMbx actor.MailboxSender[any]
	fnCall func(*Call)
	fnDone func()
	actor  actor.Actor

	// timer
	duration time.Duration
	timer    *time.Timer
	fnTimer  func()
}

type AgentOption struct {
	Key      string
	CallFn   func(*Call)
	DoneFn   func()
	TimerFn  func()
	Duration time.Duration
	Cleanup  chan any
}

func (agent *Agent) DoWork(c actor.Context) actor.WorkerStatus {
	select {
	case <-c.Done():
		if agent.fnDone != nil {
			agent.fnDone()
		}
		if agent.timer != nil {
			agent.timer.Stop()
			agent.timer = nil
		}
		return actor.WorkerEnd
	case msg, ok := <-agent.inMbx.ReceiveC():
		if ok {
			switch t := msg.(type) {
			case int:
			case string:
			case Responser:
				agent.SetLast(time.Now())
				if v, ok := t.(*Call); ok {
					if agent.fnCall != nil {
						agent.fnCall(v)
					} else {
						v.Response(nil)
					}
				}
			default:
			}
		}
		return actor.WorkerContinue
	case <-agent.timer.C:
		if agent.fnTimer != nil {
			agent.fnTimer()
		}
		agent.timer.Reset(agent.duration)
		if time.Since(agent.last) > 300*time.Second {
			agent.cleanup <- agent.key
		}
		return actor.WorkerContinue
	}
}

func (agent *Agent) Post(req any) {
	actor.Idle(actor.OptOnStart(func(c actor.Context) {
		agent.outMbx.Send(c, req)
	})).Start()
}

func (agent *Agent) Call(req any) (any, error) {
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
		agent.last = time.Now()
		actor := actor.New(agent)
		agent.actor = actor
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
	mbx := actor.NewMailbox[any](actor.OptAsChan())
	return &Agent{
		key:      opt.Key,
		last:     time.Now(),
		cleanup:  opt.Cleanup,
		inMbx:    mbx,
		outMbx:   mbx,
		fnCall:   opt.CallFn,
		fnDone:   opt.DoneFn,
		actor:    nil,
		duration: opt.Duration,
		timer:    time.NewTimer(opt.Duration),
		fnTimer:  opt.TimerFn,
	}
}
