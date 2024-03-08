package test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/aronfan/plat.mini/xtoken"
	"github.com/golang-jwt/jwt/v5"
)

func TestJWT(t *testing.T) {
	key := "hello"

	expirationTime := time.Now().Add(30 * time.Second)
	j := xtoken.NewJWT([]byte(key), jwt.SigningMethodHS256)
	token, err := j.Token(map[string]any{
		"id":  int64(1),
		"iss": "server",
		"sub": "aron",
		"exp": expirationTime.Unix(),
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

	id, ok := data["id"]
	if ok {
		dumpType(id)
	}

	exp, ok := data["exp"]
	if ok {
		dumpType(exp)
	}
}

func dumpType(t any) {
	typeOf := reflect.TypeOf(t)
	fmt.Println("Type:", typeOf)

	valueOf := reflect.ValueOf(t)
	fmt.Println("Value:", valueOf)
}
