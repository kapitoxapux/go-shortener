package main

import (
	"log"
	"net/http"

	"github.com/kapitoxapux/go-shortener/pkg/handler/handler"
)

func main() {

	routes := handler.EndpointsHandler()

	// mux := http.NewServeMux()
	// mux.Handle("/", http.HandlerFunc(PostAction))
	// mux.Handle("/*", http.HandlerFunc(GetAction))

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: routes,
	}
	log.Fatal(server.ListenAndServe())

}
