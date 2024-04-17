package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aronfan/plat.mini/xactor"
)

func Test_AgentManager(t *testing.T) {
	k1 := "123456"
	k2 := "654321"
	am := xactor.NewAgentManager()
	{
		fnCall := func(v *xactor.Call) { v.Response(v.Req) }
		fnDone := func() { fmt.Println("Done.") }
		ag1 := xactor.NewAgentWithDone(fnCall, fnDone)
		ag2 := xactor.NewAgentWithDone(fnCall, fnDone)
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

func Test_AgentTimer(t *testing.T) {
	k1 := "123456"
	k2 := "654321"
	am := xactor.NewAgentManager()
	am.SetTimer(1 * time.Second)
	am.Start()
	{
		fnCall := func(v *xactor.Call) { v.Response(v.Req) }
		fnDone := func() { fmt.Println("Done.") }
		ag1 := xactor.NewAgentWithDone(fnCall, fnDone)
		ag2 := xactor.NewAgentWithDone(fnCall, fnDone)
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
		expires := time.Now().Add(-30 * time.Second).Unix()
		ag2 := am.MarkDel(k2, expires)
		fmt.Println("at delete:", am.AtDel(k2))
		ag2.Stop()
		time.Sleep(2 * time.Second)
	}

	fmt.Println("len=", am.Len())

	am.Stop()
}
