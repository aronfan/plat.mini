package xid

import (
	"errors"
	"sync"
	"time"
)

const (
	workerBits uint8 = 10
	numberBits uint8 = 12

	workerMax int64 = -1 ^ (-1 << workerBits)
	numberMax int64 = -1 ^ (-1 << numberBits)

	// 41-bits for timestamp, can last about 68 years
	timeShift   uint8 = workerBits + numberBits
	workerShift uint8 = numberBits

	epoch int64 = 1652751577000 // 2022-05-17 09:39:37
)

type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerID  int64 // worker id
	number    int64 // increment number, 4096 max
}

func NewWorker(workerID int64) (*Worker, error) {
	if workerID < 0 || workerID > workerMax {
		return nil, errors.New("worker ID excess of quantity")
	}

	return &Worker{
		timestamp: 0,
		workerID:  workerID,
		number:    0,
	}, nil
}

func (w *Worker) Get() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now().UnixNano() / 1e6

	if w.timestamp == now {
		w.number++
		if w.number > numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}

	ID := int64((now-epoch)<<timeShift | (w.workerID << workerShift) | (w.number))
	return ID
}

func (w *Worker) Unmarshal(id int64) map[string]int64 {
	t := ((id >> timeShift) + epoch) / 1e3
	number := id & (1<<workerShift - 1)
	worker := id & (1<<timeShift - 1) >> workerShift
	return map[string]int64{
		"id":     id,
		"time":   t,
		"worker": worker,
		"number": number,
	}
}
