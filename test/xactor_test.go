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
		if am.Add(k1, ag1) {
			ag1.Start()
			fmt.Println("add:", am.Add(k1, ag1))
		}
		if am.Add(k2, ag2) {
			ag2.Start()
			fmt.Println("add:", am.Add(k2, ag2))
		}
	}

	len := am.Len()
	if len != 2 {
		t.Fatal("not equal 2")
	}
	fmt.Println("len=", len)

	{
		val := am.Val(k1)
		fmt.Println(val.Call("hello, world"))

		expires := time.Now().Add(30 * time.Second).Unix()
		ag1 := am.MarkDel(k1, expires)

		fmt.Println("at delete:", am.AtDel(k1))
		fmt.Println("add:", am.Add(k1, ag1))

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
