package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"strings"
)

type Output struct {
	IP, Token string
}

func main() {

	loadBalancerAddr := "localhost:6969"
	apiAddr := "localhost:7878"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		id := uuid.New().String()

		res, err := http.Post(fmt.Sprintf("http://%s/auth", apiAddr), "text/html", strings.NewReader(id))

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(401)
			return
		}

		if res.StatusCode == 500 {
			fmt.Println("(api) status code: 500")
			w.WriteHeader(401)
			return
		}

		token, err := ioutil.ReadAll(res.Body)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(401)
			return
		}

		res, err = http.Post(fmt.Sprintf("http://%s/address", loadBalancerAddr), "text/html", strings.NewReader(id))

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(401)
			return
		}

		if res.StatusCode == 500 {
			fmt.Println("(lb) status code: 500")
			w.WriteHeader(401)
			return
		}

		addr, err := ioutil.ReadAll(res.Body)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(401)
			return
		}

		json.NewEncoder(w).Encode(Output{string(addr), string(token)})
	})

	fmt.Println(http.ListenAndServe("localhost:8888", nil))
}