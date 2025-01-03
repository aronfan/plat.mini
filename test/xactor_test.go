package test

import (
	"fmt"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/aronfan/plat.mini/xactor"
	"go.uber.org/goleak"
)

func Test_AgentForceStop(t *testing.T) {
	nums := 100
	opt := &xactor.AgentManagerOption{
		Duration: 1 * time.Second,
		Idlechan: make(chan any, nums),
	}
	am := xactor.NewAgentManagerWithOption(opt)
	defer am.Stop()
	{
		fnCall := func(v *xactor.Call) { v.Response(0, v.Req) }
		fnDone := func() { fmt.Println("Done.") }
		for i := 1; i <= nums; i++ {
			key := fmt.Sprintf("%d", i)
			opt := &xactor.AgentOption{
				Key:      key,
				CallFn:   fnCall,
				DoneFn:   fnDone,
				Duration: 1 * time.Second,
				Idlechan: am.GetIdlechan(),
			}
			agent := xactor.NewAgentWithOption(opt)
			if err := am.Add(key, agent); err == nil {
				agent.Start()
			}
		}
	}

	fmt.Println("Len=", am.Len())
	{
		succ := 0
		fail := 0
		for i := 1; i <= nums; i++ {
			key := fmt.Sprintf("%d", i)
			val, _ := am.Val(key)
			input := fmt.Sprintf("hello, world #%d", i)
			_, output, err := val.Call(input)
			if err != nil {
				fail += 1
			} else {
				if input != output.(string) {
					fail += 1
				} else {
					succ += 1
				}
			}
		}
		fmt.Println("succ=", succ, "fail=", fail)
	}

	{
		maxRetries := 3
		for i := 1; i <= maxRetries; i++ {
			fmt.Printf("i=%d Len=%d\n", i, am.Len())
			am.StopAgents(func(agent *xactor.Agent) bool {
				_, _, err := agent.Call("flush")
				if err != nil {
					fmt.Println("key:", agent.GetKey(), "failed to flush")
					return false
				}
				agent.Stop()
				return true
			})
			len := am.Len()
			fmt.Printf("i=%d Len=%d\n", i, len)
			if len == 0 {
				break
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}

	fmt.Println("Len=", am.Len())
}

func Test_AgentManualCleanup(t *testing.T) {
	k1 := "123456"
	k2 := "654321"
	opt := &xactor.AgentManagerOption{
		Duration: 1 * time.Second,
		Idlechan: make(chan any, 100),
	}
	am := xactor.NewAgentManagerWithOption(opt)
	{
		fnCall := func(v *xactor.Call) { v.Response(0, v.Req) }
		fnDone := func() { fmt.Println("Done.") }
		opt1 := &xactor.AgentOption{
			Key:      k1,
			CallFn:   fnCall,
			DoneFn:   fnDone,
			Duration: 1 * time.Second,
			Idlechan: am.GetIdlechan(),
		}
		ag1 := xactor.NewAgentWithOption(opt1)
		opt1.Key = k2
		ag2 := xactor.NewAgentWithOption(opt1)
		if err := am.Add(k1, ag1); err == nil {
			ag1.Start()
		}
		if err := am.Add(k2, ag2); err == nil {
			ag2.Start()
		}
	}

	len := am.Len()
	if len != 2 {
		t.Fatal("not equal 2")
	}
	fmt.Println("len=", len)

	{
		val, _ := am.Val(k1)
		fmt.Println(val.Call("hello, world"))

		expires := time.Now().Add(30 * time.Second).Unix()
		ag1 := am.MarkDel(k1, expires)

		err := am.Add(k1, ag1)
		fmt.Println("at delete:", err == xactor.ErrAgentAtDel)
		fmt.Println("add:", err == nil)

		req := "flush"
		_, resp, _ := ag1.Call(req)
		if req != resp.(string) {
			t.Fatal("not equal")
		}

		// now it's safe to stop
		ag1.Stop()

		// now it's safe to delete
		am.Del(k1, ag1)
	}

	fmt.Println("len=", am.Len())

	{
		expires := time.Now().Add(30 * time.Second).Unix()
		ag2 := am.MarkDel(k2, expires)
		fmt.Println("at delete:", am.AtDel(k2))
		ag2.Stop()
		am.Del(k2, ag2)
	}

	fmt.Println("len=", am.Len())
}

func Test_AgentAutoCleanup(t *testing.T) {
	k1 := "123456"
	k2 := "654321"
	opt := &xactor.AgentManagerOption{
		Duration: 1 * time.Second,
		Idlechan: make(chan any, 100),
	}
	am := xactor.NewAgentManagerWithOption(opt)
	am.Start()
	{
		fnCall := func(v *xactor.Call) { v.Response(0, v.Req) }
		fnDone := func() { fmt.Println("Done.") }
		opt1 := &xactor.AgentOption{
			Key:      k1,
			CallFn:   fnCall,
			DoneFn:   fnDone,
			Duration: 1 * time.Second,
			Idlechan: am.GetIdlechan(),
		}
		ag1 := xactor.NewAgentWithOption(opt1)
		opt1.Key = k2
		ag2 := xactor.NewAgentWithOption(opt1)
		if err := am.Add(k1, ag1); err == nil {
			ag1.Start()
		}
		if err := am.Add(k2, ag2); err == nil {
			ag2.Start()
		}
	}

	len := am.Len()
	if len != 2 {
		t.Fatal("not equal 2")
	}
	fmt.Println("len=", len)

	{
		val, _ := am.Val(k1)
		fmt.Println(val.Call("hello, world"))

		expires := time.Now().Add(-30 * time.Second).Unix()
		ag1 := am.MarkDel(k1, expires)

		err := am.Add(k1, ag1)
		fmt.Println("at delete:", err == xactor.ErrAgentAtDel)
		fmt.Println("add:", err == nil)

		req := "flush"
		_, resp, _ := ag1.Call(req)
		if req != resp.(string) {
			t.Fatal("not equal")
		}

		// now it's safe to stop
		ag1.Stop()
		time.Sleep(2 * time.Second)
	}

	fmt.Println("len=", am.Len())

	{
		// trigger manager to cleanup k2 agent
		ag2, _ := am.Val(k2)
		ag2.SetLast(time.Now().Add(-400 * time.Second))
		time.Sleep(3 * time.Second)
	}

	fmt.Println("len=", am.Len())

	am.Stop()
}

func Test_AgentCoreDump(t *testing.T) {
	k1 := "123456"
	opt := &xactor.AgentManagerOption{
		Duration: 1 * time.Second,
		Idlechan: make(chan any, 100),
	}
	am := xactor.NewAgentManagerWithOption(opt)
	am.Start()
	defer am.Stop()
	{
		fnCall := func(v *xactor.Call) {
			defer func() {
				if err := recover(); err != nil {
					stack := string(debug.Stack())
					ss := strings.Split(stack, "\n")
					for i := 0; i < len(ss); i++ {
						str := strings.Replace(ss[i], "\t", "    ", -1)
						fmt.Println(str)
					}
				}
			}()

			var ms map[string]int
			ms["abc"] = 10 // cause coredump

			v.Response(0, v.Req)
		}
		fnDone := func() { fmt.Println("Done.") }
		opt1 := &xactor.AgentOption{
			Key:      k1,
			CallFn:   fnCall,
			DoneFn:   fnDone,
			Duration: 1 * time.Second,
			Idlechan: am.GetIdlechan(),
		}
		ag1 := xactor.NewAgentWithOption(opt1)
		if err := am.Add(k1, ag1); err == nil {
			ag1.Start()
		}
	}

	{
		ag1, _ := am.Val(k1)
		fmt.Println(ag1.Call("abc"))

		ag1.SetLast(time.Now().Add(-400 * time.Second))
		time.Sleep(10 * time.Second)
	}
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
