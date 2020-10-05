package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"net/http"
)

func main() {

	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
		return
		}

	http.HandleFunc("/port", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		info, err := cli.ContainerInspect(context.Background(), id)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		ports := info.NetworkSettings.Ports
		fmt.Println(ports)
		for _, x := range ports {
			w.Write([]byte(x[0].HostPort))
			break
		}


	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
	})
	http.ListenAndServe(":9090", nil)
}