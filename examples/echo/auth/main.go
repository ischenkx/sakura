package main

import (
	"github.com/RomanIschenko/notify/pkg/auth/jwt"
	jwt2 "github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

const JwtSecret = "javainuse-secret-key"

func main() {
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=ascii")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers","Content-Type,access-control-allow-origin, access-control-allow-headers")

		query := r.URL.Query()

		client, user := query.Get("client"), query.Get("user")

		claims := jwt.Claims{
			ClientID: client,
			UserID:  user,
		}

		token, err := jwt2.NewWithClaims(jwt2.SigningMethodHS256, claims).SignedString([]byte(JwtSecret))

		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			return
		}

		w.Write([]byte(token))
	})

	http.ListenAndServe("localhost:4646", nil)
}
