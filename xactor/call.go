package xactor

import (
	"errors"
	"time"
)

var (
	ErrResponserIsFull  error = errors.New("responser is full")
	ErrResponserClosed  error = errors.New("responser closed")
	ErrResponserTimeout error = errors.New("responser timeout")
)

type Responser interface {
	Response(int, any) error
}

type Call struct {
	Req any
	To  time.Duration

	Code   int
	Resp   any
	RespCh chan any
}

func newCall(req any) *Call {
	return &Call{
		Req:    req,
		RespCh: make(chan any, 1),
		To:     5 * time.Second,
	}
}

func (call *Call) Response(code int, resp any) error {
	call.Code = code
	call.Resp = resp

	select {
	case call.RespCh <- code:
		return nil
	default:
		return ErrResponserIsFull
	}
}

func (call *Call) WaitCall() (int, any, error) {
	select {
	case _, ok := <-call.RespCh:
		if ok {
			return call.Code, call.Resp, nil
		} else {
			return 0, nil, ErrResponserClosed
		}
	case <-time.After(call.To):
		return 0, nil, ErrResponserTimeout
	}
}
