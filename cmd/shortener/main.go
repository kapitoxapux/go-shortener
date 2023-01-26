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
	path := flag.String("f", "json.txt", "FILE_STORAGE_PATH")

	flag.Parse()

	serverAdress := os.Getenv("SERVER_ADDRESS")
	if serverAdress == "" {
		os.Setenv("SERVER_ADDRESS", *addr)
		serverAdress = *addr
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		os.Setenv("BASE_URL", *base)
		baseURL = *base
	}

	storagePath := os.Getenv("FILE_STORAGE_PATH")
	if storagePath == "" {
		os.Setenv("FILE_STORAGE_PATH", *path)
		storagePath = *path
	}

	env := config.SetEnvConf(serverAdress, baseURL, storagePath)

	server := &http.Server{
		Addr:    env.Address,
		Handler: handler.GzipMiddleware(handler.NewRoutes()),
	}

	log.Fatal(server.ListenAndServe())
}
