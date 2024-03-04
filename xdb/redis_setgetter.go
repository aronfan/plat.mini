package xdb

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type SetGetter[T any] interface {
	Init()
	*T
}

func Getter[T any, PT SetGetter[T]](rdb *redis.Client, key string) (*T, bool, error) {
	t := PT(new(T))
	c := rdb.HGetAll(context.Background(), key)
	val, err := c.Result()
	switch {
	case err != nil:
		return nil, false, err
	case len(val) == 0:
		t.Init()
		return t, true, nil
	default:
		if err := c.Scan(t); err != nil {
			return nil, false, err
		} else {
			return t, false, nil
		}
	}
}

func Setter[T any](key string, t *T) func(redis.Pipeliner) error {
	return func(pipe redis.Pipeliner) error {
		pipe.HSet(context.Background(), key, t)
		return nil
	}
}
