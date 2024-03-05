package xtest

import (
	"fmt"
	"testing"

	"github.com/aronfan/plat.mini/xlog"
	"go.uber.org/zap"
)

func TestInitLog(t *testing.T) {
	_ = xlog.InitLog()
	xlog.Debug("hello, world")
	xlog.Error("world", zap.Error(fmt.Errorf("out of memory")))
}
