package main

import (
	"github.com/RomanIschenko/notify/auth/jwt"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func main()  {
	auth := jwt.New("secret-key")
	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		token, err := auth.Register(uuid.New().String(), uuid.New().String())
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		w.Write([]byte(token))
	})
	http.ListenAndServe("localhost:6060", nil)
}
