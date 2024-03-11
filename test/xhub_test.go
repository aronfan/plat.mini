package test

import (
	"testing"
	"time"

	"github.com/aronfan/plat.mini/xhub"
)

func TestXHub(t *testing.T) {
	hub := xhub.NewMessageHub()

	i := 0

	key := "xyz"
	sub1 := xhub.NewMessageSub(key, func(any, bool) error {
		i++
		return nil
	}, func() error {
		return nil
	})
	sub1.Start()
	sub2 := xhub.NewMessageSub(key, func(any, bool) error {
		i += 2
		return nil
	}, func() error {
		return nil
	})
	sub2.Start()

	hub.Add(key, sub1)
	hub.Add(key, sub2)

	hub.Pub(key, "abc")
	time.Sleep(200 * time.Millisecond)
	if i != 3 {
		t.Fatalf("i=%d", i)
	}

	t.Log("i=", i)

	hub.Stop(key)

	sub1.Wait()
	sub2.Wait()
	t.Log("sub1 & sub2 exit...")

	hub.Del(key, sub1)
	hub.Del(key, sub2)

	kn, sn := hub.Nums()
	if kn != 0 || sn != 0 {
		t.Fatalf("keynums=%d, subnums=%d", kn, sn)
	}

	t.Log("keynums=", kn, "subnums=", sn)
}
