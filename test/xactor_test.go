package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aronfan/plat.mini/xactor"
)

func Test_AgentManualCleanup(t *testing.T) {
	k1 := "123456"
	k2 := "654321"

	opt := &xactor.AgentManagerOption{
		Duration: 1 * time.Second,
		Cleanup:  make(chan any, 100),
	}
	am := xactor.NewAgentManagerWithOption(opt)
	{
		fnCall := func(v *xactor.Call) { v.Response(v.Req) }
		fnDone := func() { fmt.Println("Done.") }
		opt1 := &xactor.AgentOption{
			Key:      k1,
			CallFn:   fnCall,
			DoneFn:   fnDone,
			Duration: 1 * time.Second,
			Cleanup:  am.GetCleanup(),
		}
		ag1 := xactor.NewAgentWithOption(opt1)
		opt1.Key = k2
		ag2 := xactor.NewAgentWithOption(opt1)
		if ok, _ := am.Add(k1, ag1); ok {
			ag1.Start()
		}
		if ok, _ := am.Add(k2, ag2); ok {
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

		ok, atdel := am.Add(k1, ag1)
		fmt.Println("at delete:", atdel)
		fmt.Println("add:", ok)

		req := "flush"
		resp, _ := ag1.Call(req)
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
		Cleanup:  make(chan any, 100),
	}
	am := xactor.NewAgentManagerWithOption(opt)
	am.Start()
	{
		fnCall := func(v *xactor.Call) { v.Response(v.Req) }
		fnDone := func() { fmt.Println("Done.") }
		opt1 := &xactor.AgentOption{
			Key:      k1,
			CallFn:   fnCall,
			DoneFn:   fnDone,
			Duration: 1 * time.Second,
			Cleanup:  am.GetCleanup(),
		}
		ag1 := xactor.NewAgentWithOption(opt1)
		opt1.Key = k2
		ag2 := xactor.NewAgentWithOption(opt1)
		if ok, _ := am.Add(k1, ag1); ok {
			ag1.Start()
		}
		if ok, _ := am.Add(k2, ag2); ok {
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

		ok, atdel := am.Add(k1, ag1)
		fmt.Println("at delete:", atdel)
		fmt.Println("add:", ok)

		req := "flush"
		resp, _ := ag1.Call(req)
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
