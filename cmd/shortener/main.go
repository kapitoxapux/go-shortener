package main

import (
	"log"
	"net/http"

	"myapp/internal/app/handler"
)

func main() {

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: handler.NewRoutes(),
	}

	log.Fatal(server.ListenAndServe())
}
