package xcm

import (
	"sync"
	"time"
)

type TimeoutBasedMerger[T any] struct {
	To     time.Duration
	ch     chan *T
	initWG sync.WaitGroup
	finiWG sync.WaitGroup
}

func NewTimeoutBasedMerger[T any]() *TimeoutBasedMerger[T] {
	return &TimeoutBasedMerger[T]{
		To:     500 * time.Millisecond,
		ch:     make(chan *T),
		initWG: sync.WaitGroup{},
		finiWG: sync.WaitGroup{},
	}
}

func NewTimeoutBasedMergerWithTimeout[T any](timeout time.Duration) *TimeoutBasedMerger[T] {
	merger := NewTimeoutBasedMerger[T]()
	merger.To = timeout
	return merger
}

func (merger *TimeoutBasedMerger[T]) Start(cb func([]*T)) func(T) {
	merger.initWG.Add(1)
	merger.finiWG.Add(1)

	go func() {
		to := merger.To
		t := time.NewTimer(to)
		defer func() {
			t.Stop()
			merger.finiWG.Done()
		}()

		var evts []*T
		merger.initWG.Done()

	OuterLoop:
		for {
			select {
			case evt, ok := <-merger.ch:
				if !ok {
					break OuterLoop
				}
				evts = append(evts, evt)
				t.Reset(to)
			case <-t.C:
				if len(evts) > 0 {
					cb(evts)
					evts = nil
				}
			}
		}

	}()

	merger.initWG.Wait()

	return func(evt T) {
		merger.ch <- &evt
	}
}

func (merger *TimeoutBasedMerger[T]) Stop() {
	close(merger.ch)
	merger.finiWG.Wait()
}
