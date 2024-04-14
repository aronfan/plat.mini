package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aronfan/plat.mini/xactor"
	"github.com/vladopajic/go-actor/actor"
)

func Test_Agent(t *testing.T) {
	agent := xactor.NewAgent(func(v *xactor.Call) { v.Response(v.Req) })

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
