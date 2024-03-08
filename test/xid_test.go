package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aronfan/plat.mini/xid"
)

func TestNewWorker(t *testing.T) {
	worker, err := xid.NewWorker(1)
	if err != nil {
		t.Fatal(err)
	}

	id := worker.Get()
	fmt.Println(id)

	m := worker.Unmarshal(id)
	timeStr := time.Unix(m["time"], 0).In(time.Local).Format("2006-01-02 15:04:05")
	fmt.Println(m, timeStr)
}

func TestWorker_UnId(t *testing.T) {
	worker, err := xid.NewWorker(1)
	if err != nil {
		t.Fatal(err)
	}

	m := worker.Unmarshal(14971982361661440)
	timeStr := time.Unix(m["time"], 0).In(time.Local).Format("2006-01-02 15:04:05")
	fmt.Println(m, timeStr)
}
