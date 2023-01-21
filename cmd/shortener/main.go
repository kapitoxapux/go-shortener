package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"myapp/internal/app/config"
	"myapp/internal/app/handler"
)

func main() {

	addr := flag.String("a", "localhost:8080", "SERVER_ADDRESS")
	base := flag.String("b", "http://localhost:8080", "BASE_URL")
	path := flag.String("f", config.GetStoragePath(), "FILE_STORAGE_PATH")

	flag.Parse()

	serverAdress := os.Getenv("SERVER_ADDRESS")
	if serverAdress == "" {
		serverAdress = *addr
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = *base
	}

	env := config.SetEnvConf(serverAdress, baseURL, *path)

	server := &http.Server{
		Addr:    env.Address,
		Handler: handler.NewRoutes(),
	}

	log.Fatal(server.ListenAndServe())
}
