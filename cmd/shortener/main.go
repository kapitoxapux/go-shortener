package main

import (
	"log"
	"myapp/pkg/handler"
	"net/http"
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
