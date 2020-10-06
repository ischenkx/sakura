package jwt

import (
	"fmt"
	"github.com/RomanIschenko/notify/pubsub"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Data string
}

func (c Claims) Valid() error {
	return nil
}

type Auth struct {
	secret string
}

func (auth *Auth) Register(id pubsub.ClientID) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{string(id)}).SignedString([]byte(auth.secret))
}

func (auth *Auth) Authorize(token string) (pubsub.ClientID, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(auth.secret), nil
	})
	fmt.Println(claims, err)
	return pubsub.ClientID(claims.Data), err
}

func New(secret string) *Auth {
	return &Auth{secret}
}