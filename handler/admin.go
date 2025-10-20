package handler

import (
	"chirpy/core/config"
	"chirpy/dto"
	"encoding/json"
	"fmt"
	"net/http"
)

func NewAdmin(apiCfg *config.ApiConfig) *Admin {
	return &Admin{
		ApiCfg: apiCfg,
	}
}

type Admin struct {
	ApiCfg *config.ApiConfig
}

func (admin *Admin) Metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprintf(`
		<html>
		  <body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		  </body>
		</html>
		`, admin.ApiCfg.FileserverHits.Load())))
	if err != nil {
		panic(err)
	}
}

func (admin *Admin) Reset(w http.ResponseWriter, r *http.Request) {
	if admin.ApiCfg.Platform == "dev" {
		err := admin.ApiCfg.Db.ClearUsers(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errDto := dto.Response{
				Error: err.Error(),
			}
			errorResponse, _ := json.Marshal(errDto)
			_, err = w.Write(errorResponse)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		admin.ApiCfg.FileserverHits.Store(0)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
