package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aronfan/plat.mini/xactor"
	"github.com/vladopajic/go-actor/actor"
)

func Test_AgentStart(t *testing.T) {
	agent := xactor.NewAgent(func(v *xactor.Call) { v.Response(v.Req) }, nil)

	a := actor.New(agent)
	a.Start()
	defer a.Stop()

	fmt.Println("o1")
	agent.Post(1)
	agent.Post("hello, world")

	time.Sleep(time.Second)

	fmt.Println("o2")
	req := "hello, world #2"
	resp, err := agent.Call(req)
	if err != nil {
		t.Fatal(err)
	}
	str, _ := resp.(string)
	if str != req {
		t.Fatal("not equal")
	}
	fmt.Println(str)
}

func Test_AgentStop(t *testing.T) {
	var i int = 0

	fnCall := func(v *xactor.Call) { v.Response(v.Req) }
	fnDone := func() { i += 1 }
	ag1 := xactor.NewAgent(fnCall, fnDone)
	ag2 := xactor.NewAgent(fnCall, fnDone)

	a1 := actor.New(ag1)
	a2 := actor.New(ag2)
	a1.Start()
	a2.Start()

	actors := make(map[int]actor.Actor)
	agents := make(map[int]*xactor.Agent)
	actors[1] = a1
	actors[2] = a2
	agents[1] = ag1
	agents[2] = ag2

	// delete 1
	{
		a3, ok := actors[1]
		if ok {
			delete(actors, 1)
			delete(agents, 1)
			a3.Stop()
		}
	}
	time.Sleep(time.Second)

	// delete 2
	{
		a4, ok := actors[2]
		if ok {
			delete(actors, 2)
			delete(agents, 2)
			a4.Stop()
		}
	}
	time.Sleep(time.Second)

	fmt.Println("actors num:", len(actors))
	fmt.Println("agents num:", len(agents))
	if i != 2 {
		t.Fatal("i != 2")
	}
}
