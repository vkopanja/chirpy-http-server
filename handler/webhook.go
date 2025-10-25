package handler

import (
	"chirpy/core/config"
	"chirpy/dto"
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
)

func NewWebhook(apiCfg *config.ApiConfig) *Webhook {
	return &Webhook{
		apiCfg: apiCfg,
	}
}

type Webhook struct {
	apiCfg *config.ApiConfig
}

func (wh *Webhook) CatchWebhook(w http.ResponseWriter, r *http.Request) {
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if key != wh.apiCfg.PolkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var webhook dto.UserWebhook
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&webhook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if webhook.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// check if user exists
	user, err := wh.apiCfg.Db.GetUserByID(r.Context(), webhook.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err = wh.apiCfg.Db.UpdateUserChirpyRedByID(r.Context(), database.UpdateUserChirpyRedByIDParams{
		IsChirpyRed: true,
		ID:          user.ID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
