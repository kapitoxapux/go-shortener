package main

import (
	"log"
	"net/http"

	"myapp/pkg/handler"
)

func main() {

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: handler.NewRoutes(),
	}

	log.Fatal(server.ListenAndServe())
}
