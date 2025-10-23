package config

import (
	"chirpy/internal/database"
	"net/http"
	"sync/atomic"
)

func NewApiConfig(db *database.Queries, platform, secret *string) *ApiConfig {
	return &ApiConfig{
		FileserverHits: atomic.Int32{},
		Db:             db,
		Platform:       *platform,
		Secret:         *secret,
	}
}

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             *database.Queries
	Platform       string
	Secret         string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
