package xdb

import (
	"context"
	"net"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
)

var protocol string = "mysql+ssh"

func MysqlOverSsh(dsn string, cli *ssh.Client) string {
	mysql.RegisterDialContext(protocol, func(ctx context.Context, addr string) (net.Conn, error) {
		return cli.Dial("tcp", addr)
	})
	// replace "tcp" to "mysql+ssh"
	return strings.Replace(dsn, "tcp", protocol, 1)
}
