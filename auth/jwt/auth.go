package jwt

import (
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

func (auth *Auth) Register(id string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{id}).SignedString([]byte(auth.secret))
}

func (auth *Auth) Authorize(token string) (string, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(auth.secret), nil
	})
	return claims.Data, err
}

func New(secret string) *Auth {
	return &Auth{secret}
}