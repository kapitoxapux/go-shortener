package main

import (
	"log"
	"net/http"
	"os"

	"myapp/internal/app/config"
	"myapp/internal/app/handler"
)

func main() {

	serverAdress := os.Getenv("SERVER_ADDRESS")
	if serverAdress == "" {
		serverAdress = "localhost:8080"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "/app"
	}

	env := config.SetEnvConf(serverAdress, baseURL)

	server := &http.Server{
		Addr:    env.Address,
		Handler: handler.NewRoutes(),
	}

	log.Fatal(server.ListenAndServe())
}
