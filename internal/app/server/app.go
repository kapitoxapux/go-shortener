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

func RemoveWorkers(storage *service.Service, inputCh chan *service.Shorter) {
	shorters := make([]string, 0, 10)

	workersCount := 10
	workerChs := make([]chan *service.Shorter, 0, workersCount)
	fanOutChs := service.FanOut(inputCh, workersCount)
	for _, fanOutCh := range fanOutChs {
		workerCh := make(chan *service.Shorter)
		service.NewWorker(fanOutCh, workerCh)
		workerChs = append(workerChs, workerCh)
	}
	for id := range service.FanIn(workerChs...) {
		shorters = append(shorters, id)
		storage.Storage.RemoveShorts(shorters)
	}
}

func (a *App) Run() error {
	route := chi.NewRouter()
	address := config.GetConfigAddress()
	registerHTTPEndpoints(route, *a.service, *a.channel)
	a.httpServer = &http.Server{
		Addr:    address,
		Handler: handler.GzipMiddleware(route),
	}

	go RemoveWorkers(a.service, a.channel.InputChannel)

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
