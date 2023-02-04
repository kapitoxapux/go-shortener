package main

import (
	"log"
	"net/http"

	"myapp/internal/app/config"
	"myapp/internal/app/handler"
)

func main() {
	config.SetConfig()

	address := config.GetConfigAddress()

	server := &http.Server{
		Addr:    address,
		Handler: handler.GzipMiddleware(handler.NewRoutes()),
	}

	log.Fatal(server.ListenAndServe())
}
