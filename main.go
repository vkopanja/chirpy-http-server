// Package main is used for application startup
package main

import (
	"chirpy/handler"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	health := handler.NewHealth()

	mux.Handle("/app", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/app/assets", http.StripPrefix("/app/assets", http.FileServer(http.Dir("./assets"))))
	mux.HandleFunc("/healthz", health.ServeHTTP)

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("failed starting server: %s", err))
	}
}
