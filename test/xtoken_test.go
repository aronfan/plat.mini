package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aronfan/plat.mini/xtoken"
	"github.com/golang-jwt/jwt/v5"
)

func TestJWT(t *testing.T) {
	key := "hello"

	j := xtoken.NewJWT([]byte(key), jwt.SigningMethodHS256)
	token, err := j.Token(map[string]any{
		"id":  1,
		"iss": "server",
		"sub": "aron",
		"exp": time.Second * 30,
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)

	j2 := xtoken.NewJWT([]byte(key), jwt.SigningMethodHS256)
	data, err := j2.Parse(token)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(data)
}
