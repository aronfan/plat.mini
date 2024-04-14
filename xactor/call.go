package xactor

import (
	"errors"
	"time"
)

var (
	ErrResponserIsFull  error = errors.New("responser is full")
	ErrResponserClosed  error = errors.New("responser was closed")
	ErrResponserTimeout error = errors.New("responser timeout")
)

type Responser interface {
	Response(any) error
}

type Call struct {
	Req    any
	RespCh chan any
	To     time.Duration
}

func newCall(req any) *Call {
	return &Call{
		Req:    req,
		RespCh: make(chan any, 1),
		To:     5 * time.Second,
	}
}

func (call *Call) Response(resp any) error {
	select {
	case call.RespCh <- resp:
		return nil
	default:
		return ErrResponserIsFull
	}
}

func (call *Call) WaitCall() (any, error) {
	select {
	case v, ok := <-call.RespCh:
		if ok {
			return v, nil
		} else {
			return nil, ErrResponserClosed
		}
	case <-time.After(call.To):
		return nil, ErrResponserTimeout
	}
}
