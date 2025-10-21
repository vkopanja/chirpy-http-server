// Package main is used for application startup
package main

import (
	"chirpy/core/config"
	"chirpy/handler"
	"chirpy/internal/database"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env file: %s\n", err)
	}

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(fmt.Sprintf("failed opening database: %s", err))
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("error closing database: %s\n", err)
		}
	}(db)

	queries := database.New(db)

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	apiCfg := config.NewApiConfig(queries, &platform)
	auth := handler.NewAuth(apiCfg)
	health := handler.NewHealth()
	admin := handler.NewAdmin(apiCfg)
	user := handler.NewUser(apiCfg)
	chirp := handler.NewChirp(apiCfg)

	mux.Handle("/app", http.StripPrefix("/app", apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets", apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir("./assets")))))
	mux.HandleFunc("GET /admin/metrics", admin.Metrics)
	mux.HandleFunc("POST /admin/reset", admin.Reset)

	mux.HandleFunc("GET /api/healthz", health.ServeHTTP)

	//auth
	mux.HandleFunc("POST /api/login", auth.Login)
	// users
	mux.HandleFunc("POST /api/users", user.Create)

	// chirps
	mux.HandleFunc("POST /api/chirps", chirp.Create)
	mux.HandleFunc("GET /api/chirps", chirp.GetAll)
	mux.HandleFunc("GET /api/chirps/{chirpID}", chirp.GetChirpById)

	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("failed starting server: %s", err))
	}
}
