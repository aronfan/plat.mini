package test

import (
	"fmt"
	"testing"

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
		}
		if am.Add(k2, ag2) {
			ag2.Start()
		}
	}

	len := am.Len()
	if len != 2 {
		t.Fatal("not equal 2")
	}
	fmt.Println("len=", len)

	{
		ag1 := am.Del(k1)

		req := "hello, world"
		resp, _ := ag1.Call(req)
		if req != resp.(string) {
			t.Fatal("not equal")
		}

		// now it's safe to stop
		ag1.Stop()
	}

	fmt.Println("len=", am.Len())

	{
		ag2 := am.Del(k2)
		ag2.Stop()
	}

	fmt.Println("len=", am.Len())
}
