package xdb

import (
	"context"
	"net"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/ssh"
)

func NewRedisOptions(addr, user, pass string, db int) *redis.Options {
	return &redis.Options{Addr: addr, Username: user, Password: pass, DB: db}
}

func NewRedisOptionsWithUrl(url string) *redis.Options {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil
	}
	return opt
}

func RedisOverSsh(opt *redis.Options, tunnel *ssh.Client) (*redis.Client, error) {
	if tunnel != nil {
		opt.Dialer = func(ctx context.Context, network, raddr string) (conn net.Conn, e error) {
			return tunnel.Dial(network, raddr)
		}
		opt.ReadTimeout = -2
		opt.WriteTimeout = -2
	}

	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); nil != err {
		return nil, err
	}

	return rdb, nil
}
