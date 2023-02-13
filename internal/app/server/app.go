package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"

	"myapp/internal/app/config"
	"myapp/internal/app/handler"
	"myapp/internal/app/repository"
	"myapp/internal/app/service"
	"myapp/internal/app/storage"
)

var repo repository.Repository

type App struct {
	httpServer *http.Server
	service    *service.Service
}

func NewApp() *App {

	db := GetDB()
	service := service.NewService(*db)

	// if status, _ := storage.ConnectionDBCheck(); status == 200 {
	// 	repo = repository.NewRepository(config.GetStorageDB())
	// } else {
	// 	repo = nil
	// }

	return &App{
		service: service,
	}
}

func GetDB() *service.Storage {

	config.SetConfig()

	if status, _ := service.ConnectionDBCheck(); status == 200 {
		return storage.NewDB()
	}

	if pathStorage := config.GetConfigPath(); pathStorage != "" {
		return storage.NewFileDB()
	}

	return storage.NewInMemDB()
}

func registerHTTPEndpoints(router *chi.Mux, service service.Service) {
	h := handler.NewHandler(service)

	router.Post("/", h.SetShortAction)
	router.Get("/{`\\w+$`}", h.GetShortAction)
	router.Post("/api/shorten", h.GetJSONShortAction)
	router.Get("/api/user/urls", h.GetUserURLAction)
	router.Get("/ping", h.GetPingAction)
	router.Post("/api/shorten/batch", h.GetBatchAction)

}

func (a *App) Run() error {

	route := chi.NewRouter()

	address := config.GetConfigAddress()

	registerHTTPEndpoints(route, *a.service)

	a.httpServer = &http.Server{
		Addr:    address,
		Handler: route,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)

}
