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
		apiCfg: apiCfg,
	}
}

type Admin struct {
	apiCfg *config.ApiConfig
}

// Metrics
// @Summary Show metrics
// @Description Show metrics of total fileserver hits
// @Tags Admin handler
// @Produce text/html
// @Success 200 {string} string "HTML response"
// @Router /admin/metrics [get]
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
		`, admin.apiCfg.FileserverHits.Load())))
	if err != nil {
		panic(err)
	}
}

// Reset
// @Summary Reset database
// @Description Reset database if on dev profile and clear all users, also reset fileserver hits
// @Tags Admin handler
// @Produce json
// @Success 204 {object} nil "No Content"
// @Failure 500 {object} dto.Response
// @Router /admin/reset [post]
func (admin *Admin) Reset(w http.ResponseWriter, r *http.Request) {
	if admin.apiCfg.Platform == "dev" {
		err := admin.apiCfg.Db.ClearUsers(r.Context())
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
		admin.apiCfg.FileserverHits.Store(0)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
