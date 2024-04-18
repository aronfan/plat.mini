package xactor

import (
	"time"

	"github.com/vladopajic/go-actor/actor"
)

type Agent struct {
	key     string
	last    time.Time
	actor   actor.Actor
	cleanup chan any
	report  bool

	// message
	inMbx  actor.MailboxReceiver[any]
	outMbx actor.MailboxSender[any]
	fnCall func(*Call)
	fnDone func()

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
			case int:
			case string:
			case Responser:
				agent.SetLast(time.Now())
				if v, ok := t.(*Call); ok {
					if agent.fnCall != nil {
						agent.fnCall(v)
					} else {
						v.Response(0, nil)
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
		if !agent.report && (time.Since(agent.last) > 300*time.Second) {
			select {
			case agent.cleanup <- agent.key:
				agent.report = true
			default:
			}
		}
		return actor.WorkerContinue
	}
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
	mbx := actor.NewMailbox[any](actor.OptAsChan())
	return &Agent{
		key:      opt.Key,
		last:     time.Now(),
		cleanup:  opt.Cleanup,
		report:   false,
		actor:    nil,
		inMbx:    mbx,
		outMbx:   mbx,
		fnCall:   opt.CallFn,
		fnDone:   opt.DoneFn,
		duration: opt.Duration,
		timer:    nil,
		fnTimer:  opt.TimerFn,
	}
}
