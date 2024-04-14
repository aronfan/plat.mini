package xactor

import (
	"github.com/vladopajic/go-actor/actor"
)

type Agent struct {
	inMbx  actor.MailboxReceiver[any]
	outMbx actor.MailboxSender[any]
	fnCall func(*Call)
	fnDone func()
}

func (agent *Agent) DoWork(c actor.Context) actor.WorkerStatus {
	select {
	case <-c.Done():
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

func NewAgent(fnCall func(*Call), fnDone func()) *Agent {
	mbx := actor.NewMailbox[any](actor.OptAsChan())
	return &Agent{inMbx: mbx, outMbx: mbx, fnCall: fnCall, fnDone: fnDone}
}