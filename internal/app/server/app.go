package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"myapp/internal/app/config"
	"myapp/internal/app/handler"
	"myapp/internal/app/service"
	"myapp/internal/app/storage"
)

type App struct {
	httpServer *http.Server
	service    *service.Service
	channel    *service.Channel
}

func NewApp() *App {
	InputCh := make(chan *service.Shorter)
	listener := service.NewListener(InputCh)

	db := GetDB()
	service := service.NewService(db)

	return &App{
		service: service,
		channel: listener,
	}
}

func GetDB() service.Storage {
	config.SetConfig()
	if status, _ := handler.ConnectionDBCheck(); status == http.StatusOK {

		return storage.NewDB()
	}
	if pathStorage := config.GetConfigPath(); pathStorage != "" {

		return storage.NewFileDB()
	}

	return storage.NewInMemDB()
}

func registerHTTPEndpoints(router *chi.Mux, service service.Service, channel service.Channel) {
	h := handler.NewHandler(service, channel)
	router.Post("/", h.SetShortAction)
	router.Get("/{`\\w+$`}", h.GetShortAction)
	router.Post("/api/shorten", h.GetJSONShortAction)
	router.Get("/api/user/urls", h.GetUserURLAction)
	router.Get("/ping", h.GetPingAction)
	router.Post("/api/shorten/batch", h.GetBatchAction)
	router.Delete("/api/user/urls", h.RemoveBatchAction)
}

func (a *App) Run(ctx context.Context) error {
	route := chi.NewRouter()
	address := config.GetConfigAddress()
	registerHTTPEndpoints(route, *a.service, *a.channel)

	a.httpServer = &http.Server{
		Addr:    address,
		Handler: handler.CustomMiddleware(route),
	}

	go service.RemoveWorkers(a.service, a.channel.InputChannel)

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}

	}()

	<-ctx.Done()

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	quit := make(chan struct{}, 1)

	go func() {
		time.Sleep(3 * time.Second)
		quit <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("server shutdown: %w", ctx.Err())
	case <-quit:
		log.Println("finished")
	}

	return nil
}
