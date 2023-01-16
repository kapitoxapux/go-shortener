package main

import (
	"log"
	"net/http"

	"myapp/internal/app/config"
	"myapp/internal/app/handler"
)

func main() {

	env := config.SetEnvConf("localhost:8080", "/")

	server := &http.Server{
		Addr:    env.Address,
		Handler: handler.NewRoutes(),
	}

	log.Fatal(server.ListenAndServe())
}
