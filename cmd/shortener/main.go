package main

import (
	"log"
	// "net/http"

	// "myapp/internal/app/config"
	// "myapp/internal/app/handler"
	// "myapp/internal/app/repository"
	// "myapp/internal/app/storage"
	"myapp/internal/app/server"
)

func main() {

	app := server.NewApp()
	if err := app.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}

}
