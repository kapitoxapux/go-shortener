package main

import (
	"context"
	"log"
	"myapp/internal/app/server"
	"os/signal"
	"syscall"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := server.NewApp()
	if err := app.Run(ctx); err != nil {
		log.Fatalf("%s", err.Error())
	}

}
