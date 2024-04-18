package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aronfan/plat.mini/xcm"
	"github.com/aronfan/plat.mini/xdb"
	"github.com/aronfan/plat.mini/xssh"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/ssh"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
func TestRawRedis(t *testing.T) {
	if err := xcm.LoadConfigFile("config.yaml"); err != nil {
		t.Error(err)
		return
	}
	if conf, err := xcm.MapToStruct[testConfig](); err != nil {
		t.Error(err)
		return
	} else {
		opt := xdb.NewRedisOptionsWithUrl(conf.RedisConfig.Url)
		rdb, err := xdb.RedisOverSsh(opt, nil)
		if err != nil {
			t.Errorf("%v", err)
			return
		}

		runRedisOp(rdb)
		runRedisTx(rdb)
		runRedisPipe(rdb)
		runRedisPubsub(rdb)
	}
}
*/

func TestSshRedis(t *testing.T) {
	if err := xcm.LoadConfigFile("config.yaml"); err != nil {
		t.Error(err)
		return
	}
	if conf, err := xcm.MapToStruct[testConfig](); err != nil {
		t.Error(err)
		return
	} else {
		var cli *ssh.Client
		opt := xdb.NewRedisOptionsWithUrl(conf.RedisConfig.Url)
		if conf.RedisConfig.Ssh.Enable {
			cli, err = xssh.SshClientWithKeyFile(conf.RedisConfig.Ssh.Addr,
				conf.RedisConfig.Ssh.User,
				conf.RedisConfig.Ssh.Keyfile,
				conf.RedisConfig.Ssh.Keypass)
			if err != nil {
				t.Errorf("%v", err)
				return
			}
			defer cli.Close()
		}
		rdb, err := xdb.RedisOverSsh(opt, cli)
		if err != nil {
			t.Errorf("%v", err)
			return
		}
		defer rdb.Close()

		runRedisOp(rdb)
		runRedisTx(rdb)
		runRedisPipe(rdb)
	}
}

func runRedisOp(rdb *redis.Client) {
	var ctx = context.Background()

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output:
	// key value
	// key2 does not exist

	rdb.Del(ctx, "key")
}

func runRedisPipe(rdb *redis.Client) {
	ctx := context.Background()

	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < 3; i++ {
			pipe.Set(ctx, fmt.Sprintf("key%d", i), i+1, 0)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	cmds, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < 3; i++ {
			pipe.Get(ctx, fmt.Sprintf("key%d", i))
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, cmd := range cmds {
		fmt.Println(cmd.(*redis.StringCmd).Val())
	}

	key := "pet_10001"
	pet, init, err := xdb.Getter[Pet](rdb, key)
	if err != nil {
		fmt.Println(err)
		return
	}
	if init {
		pet.PttID = 101
		pet.Level = 1
		pet.Name = "qianqian"
	} else {
		pet.Name = "qianqian"
	}

	xdb.Commit(rdb, xdb.Setter(key, pet))
}

func runRedisTx(rdb *redis.Client) {
	del := false
	key := "pet_10001"

	xdb.BeginTx(rdb, func(tx *redis.Tx) error {
		pet, init, err := xdb.Getter[Pet](rdb, key)
		if err != nil {
			return err
		}
		if init {
			pet.PttID = 9527
			pet.Level = 1
			pet.Name = "xiaobai"
		} else {
			if pet.LevelUp() >= 5 {
				del = true
			}
		}

		return xdb.CommitTx(tx, xdb.Setter(key, pet))
	}, key)

	pet, _, err := xdb.Getter[Pet](rdb, key)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("finally %+v\n", pet)

	if del {
		rdb.Del(context.Background(), key)
	}
}

/*
func runRedisPubsub(rdb *redis.Client) {
	ctx := context.Background()
	channel := "mychannel"
	pubsub := rdb.Subscribe(ctx, channel)
	defer pubsub.Close()

	go func() {
		time.Sleep(1 * time.Second)
		rdb.Publish(context.Background(), channel, "payload")
	}()

	ch := pubsub.Channel()
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
		break
	}
}
*/

func TestSshMysql(t *testing.T) {
	if err := xcm.LoadConfigFile("config.yaml"); err != nil {
		t.Error(err)
		return
	}
	if conf, err := xcm.MapToStruct[testConfig](); err != nil {
		t.Error(err)
		return
	} else {
		dsn := conf.MysqlConfig.Dsn

		if conf.MysqlConfig.Ssh.Enable {
			cli, err := xssh.SshClientWithKeyFile(conf.MysqlConfig.Ssh.Addr,
				conf.MysqlConfig.Ssh.User,
				conf.MysqlConfig.Ssh.Keyfile,
				conf.MysqlConfig.Ssh.Keypass)
			if err != nil {
				t.Errorf("%v", err)
				return
			}
			defer cli.Close()
			dsn = xdb.MysqlOverSsh(dsn, cli)
		}

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Errorf("%v", err)
			return
		}
		sqlDB, _ := db.DB()
		defer sqlDB.Close()

		rows, err := db.Raw("SHOW DATABASES").Rows()
		if err != nil {
			t.Errorf("%v", err)
			return
		}

		var databaseName string
		for rows.Next() {
			if err := rows.Scan(&databaseName); err != nil {
				t.Errorf("%v", err)
				return
			}
			t.Log("database:", databaseName)
		}
	}
}
