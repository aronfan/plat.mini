package xactor

import (
	"fmt"

	"github.com/vladopajic/go-actor/actor"
)

type Agent struct {
	inMbx  actor.MailboxReceiver[any]
	outMbx actor.MailboxSender[any]
	fnCall func(*Call)
}

func (agent *Agent) DoWork(c actor.Context) actor.WorkerStatus {
	select {
	case <-c.Done():
		return actor.WorkerEnd
	case msg, ok := <-agent.inMbx.ReceiveC():
		if ok {
			switch t := msg.(type) {
			case int:
				fmt.Println("int:", t)
			case string:
				fmt.Println("string:", t)
			case Responser:
				fmt.Println("Responser:", t)
				if v, ok := t.(*Call); ok {
					if agent.fnCall != nil {
						agent.fnCall(v)
					} else {
						v.Response(nil)
					}
				}
			default:
				fmt.Printf("unknown type: %T\n", t)
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

func NewAgent(fn func(*Call)) *Agent {
	mbx := actor.NewMailbox[any](actor.OptAsChan())
	return &Agent{inMbx: mbx, outMbx: mbx, fnCall: fn}
}
