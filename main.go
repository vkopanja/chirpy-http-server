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

	_ "chirpy/docs"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Chirpy API
// @version 1.0
// @description This is Chirpy API created for boot.dev HTTP servers course.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env file: %s\n", err)
	}

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

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

	apiCfg := config.NewApiConfig(queries, &platform, &secret, &polkaKey)
	auth := handler.NewAuth(apiCfg)
	health := handler.NewHealth()
	admin := handler.NewAdmin(apiCfg)
	user := handler.NewUser(apiCfg)
	chirp := handler.NewChirp(apiCfg)
	webhook := handler.NewWebhook(apiCfg)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	mux.Handle("/app", http.StripPrefix("/app", apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets", apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir("./assets")))))
	mux.HandleFunc("GET /admin/metrics", admin.Metrics)
	mux.HandleFunc("POST /admin/reset", admin.Reset)

	mux.HandleFunc("GET /api/healthz", health.ServeHTTP)

	//auth
	mux.HandleFunc("POST /api/login", auth.Login)
	mux.HandleFunc("POST /api/refresh", auth.Refresh)
	mux.HandleFunc("POST /api/revoke", auth.Revoke)

	// users
	mux.HandleFunc("POST /api/users", user.Create)
	mux.HandleFunc("PUT /api/users", user.Update)

	// chirps
	mux.HandleFunc("POST /api/chirps", chirp.Create)
	mux.HandleFunc("GET /api/chirps", chirp.GetAllWithFilter)
	mux.HandleFunc("GET /api/chirps/{chirpID}", chirp.GetChirpById)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", chirp.Delete)

	// webhooks
	mux.HandleFunc("POST /api/polka/webhooks", webhook.CatchWebhook)

	err = server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("failed starting server: %s", err))
	}
}
