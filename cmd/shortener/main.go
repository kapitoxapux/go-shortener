package main

import (
	"log"
	"myapp/internal/app/server"
)

func main() {

	app := server.NewApp()
	if err := app.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}

}
