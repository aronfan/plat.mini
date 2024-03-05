package xssh

import (
	"os"
	"os/user"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
)

func SshClient(addr, user, pass string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	return ssh.Dial("tcp", addr, config)
}

func SshClientWithKeyFile(addr, user, keyfile, keypass string) (*ssh.Client, error) {
	realfile, err := expandTilde(keyfile)
	if nil != err {
		return nil, err
	}

	keyBytes, err := os.ReadFile(realfile)
	if nil != err {
		return nil, err
	}

	var signer ssh.Signer
	if keypass != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(keypass))
	} else {
		signer, err = ssh.ParsePrivateKey(keyBytes)
	}

	if nil != err {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	return ssh.Dial("tcp", addr, config)
}

func expandTilde(path string) (string, error) {
	if len(path) > 0 && path[0] == '~' {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		path = filepath.Join(usr.HomeDir, path[1:])
	}
	return path, nil
}
