package xtest

import (
	"testing"
	"time"

	"github.com/aronfan/plat.mini/xcm"
)

func TestXCM(t *testing.T) {
	if err := xcm.LoadConfigFile("config.yaml"); err != nil {
		t.Error(err)
		return
	}

	config, err := xcm.MapToStruct[testConfig]()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", config)

	merger := xcm.NewTimeoutBasedMerger[int]()
	onEvent := merger.Start(func(evts []*int) { t.Log("evts length=", len(evts)) })
	onEvent(1)
	onEvent(2)
	time.Sleep(1 * time.Second) // trigger the timeout
	merger.Stop()
	t.Log("merger stopped...")
}
