package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ClientID, UserID string
}

func (c Claims) Valid() error {
	return nil
}

type Auth struct {
	secret string
}

func (auth *Auth) Register(id, userId string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{id, userId}).SignedString([]byte(auth.secret))
}

func (auth *Auth) Authorize(token string) (string, string, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(auth.secret), nil
	})
	return claims.ClientID, claims.UserID, err
}

func New(secret string) *Auth {
	return &Auth{secret}
}
