package xdb

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func BeginTx(rdb *redis.Client, cb func(tx *redis.Tx) error, keys ...string) error {
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err := rdb.Watch(context.Background(), cb, keys...)
		if err == nil {
			// it's ok
			break
		} else if err == redis.TxFailedErr {
			// lock failed
			continue
		} else {
			return err
		}
	}
	return nil
}

func CommitTx(tx *redis.Tx, cb func(pipe redis.Pipeliner) error) error {
	_, err := tx.TxPipelined(context.Background(), cb)
	return err
}

func Commit(rdb *redis.Client, cb func(pipe redis.Pipeliner) error) error {
	_, err := rdb.Pipelined(context.Background(), cb)
	return err
}
