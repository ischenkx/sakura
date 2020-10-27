package main

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"net/http"
)

var logger = logrus.WithField("source", "registrator")

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		logger.Errorf("failed to create env client:", err)
		return
	}

	http.HandleFunc("/port", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		logger.Info("Registering node "+id)
		info, err := cli.ContainerInspect(context.Background(), id)
		if err != nil {
			logger.Info("err:", err)
			w.WriteHeader(500)
			return
		}
		ports := info.NetworkSettings.Ports
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