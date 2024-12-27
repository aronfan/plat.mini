package xdb

import (
	"time"

	"github.com/aronfan/plat.mini/xcm"
	"github.com/aronfan/plat.mini/xlog"
	"github.com/aronfan/plat.mini/xssh"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql(conf *xcm.Config) (*gorm.DB, *ssh.Client, error) {
	var err error
	var db *gorm.DB
	var cli *ssh.Client
	defer func() {
		if err != nil {
			if db != nil {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}
			if cli != nil {
				cli.Close()
			}
		}
	}()

	dsn := conf.MysqlConfig.Dsn

	if conf.MysqlConfig.Ssh.Enable {
		cli, err = xssh.SshClientWithKeyFile(conf.MysqlConfig.Ssh.Addr,
			conf.MysqlConfig.Ssh.User,
			conf.MysqlConfig.Ssh.Keyfile,
			conf.MysqlConfig.Ssh.Keypass)
		if err != nil {
			xlog.Error("SshClientWithKeyFile", zap.Error(err))
			return nil, nil, err
		}
		dsn = MysqlOverSsh(dsn, cli)
	}

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		xlog.Error("gorm.Open", zap.Error(err))
		return nil, nil, err
	}

	sdb, _ := db.DB()
	sdb.SetMaxIdleConns(10)
	sdb.SetMaxOpenConns(100)
	sdb.SetConnMaxLifetime(time.Hour)

	xlog.Info("Mysql init ok")
	return db, cli, nil
}

func InitRedis(conf *xcm.Config) (*redis.Client, *ssh.Client, error) {
	var err error
	var db *redis.Client
	var cli *ssh.Client
	defer func() {
		if err != nil {
			if db != nil {
				db.Close()
			}
			if cli != nil {
				cli.Close()
			}
		}
	}()

	if conf.RedisConfig.Ssh.Enable {
		cli, err = xssh.SshClientWithKeyFile(conf.RedisConfig.Ssh.Addr,
			conf.RedisConfig.Ssh.User,
			conf.RedisConfig.Ssh.Keyfile,
			conf.RedisConfig.Ssh.Keypass)
		if err != nil {
			xlog.Error("Redis", zap.Error(err))
			return nil, nil, err
		}
	}

	opt := NewRedisOptionsWithUrl(conf.RedisConfig.Url)
	db, err = RedisOverSsh(opt, cli)
	if err != nil {
		xlog.Error("Redis", zap.Error(err))
		return nil, nil, err
	}

	xlog.Info("Redis init ok")
	return db, cli, nil
}
