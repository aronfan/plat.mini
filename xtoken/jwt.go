package xtoken

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	key    []byte
	method jwt.SigningMethod
}

func NewJWT(key []byte, method jwt.SigningMethod) *JWT {
	return &JWT{
		key:    key,
		method: method,
	}
}

func (j *JWT) Token(data map[string]any) (string, error) {
	claims := make(jwt.MapClaims)
	for k, v := range data {
		claims[k] = v
	}
	t := jwt.NewWithClaims(j.method, claims)
	return t.SignedString(j.key)
}

func (j *JWT) Parse(token string) (map[string]any, error) {
	claims := make(jwt.MapClaims)
	_, err := jwt.ParseWithClaims(token, &claims, func(*jwt.Token) (interface{}, error) {
		return j.key, nil
	})
	if err != nil {
		return nil, err
	}
	data := make(map[string]any)
	for k, v := range claims {
		data[k] = v
	}
	return data, nil
}
