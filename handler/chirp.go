package handler

import (
	"chirpy/core/config"
	"chirpy/dto"
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
)

func NewChirp(apiCfg *config.ApiConfig) *Chirp {
	return &Chirp{
		ApiCfg: apiCfg,
	}
}

type Chirp struct {
	ApiCfg *config.ApiConfig
}

func (c *Chirp) Create(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(token, c.ApiCfg.Secret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirp, err := validateChirp(userID, r)
	if err != nil {
		fmt.Printf("error validating chirp: %s\n", err)
		errDto := dto.Response{
			Error: err.Error(),
		}
		errorResponse, _ := json.Marshal(errDto)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errorResponse)
	} else {
		createChirp, err := c.ApiCfg.Db.CreateChirp(r.Context(), database.CreateChirpParams{
			ID:        uuid.New(),
			UserID:    userID,
			Body:      chirp.Body,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			errDto := dto.Response{
				Error: err.Error(),
			}
			errorResponse, _ := json.Marshal(errDto)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write(errorResponse)
			return
		}

		w.WriteHeader(http.StatusCreated)
		bodySlice, err := json.Marshal(createChirp)
		if err != nil {
			errDto := dto.Response{
				Error: err.Error(),
			}
			errorResponse, _ := json.Marshal(errDto)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write(errorResponse)
			return
		}
		_, err = w.Write(bodySlice)
	}
}

func (c *Chirp) GetAll(w http.ResponseWriter, r *http.Request) {
	chirps, err := c.ApiCfg.Db.GetChirps(r.Context())
	if err != nil {
		errDto := dto.Response{
			Error: err.Error(),
		}
		errorResponse, _ := json.Marshal(errDto)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write(errorResponse)
		return
	}

	responseBytes, err := json.Marshal(chirps)
	if err != nil {
		errDto := dto.Response{
			Error: err.Error(),
		}
		errorResponse, _ := json.Marshal(errDto)
		_, err = w.Write(errorResponse)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
}

func (c *Chirp) GetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")
	if chirpId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chirp, err := c.ApiCfg.Db.GetChirpById(r.Context(), uuid.MustParse(chirpId))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	responseBytes, err := json.Marshal(chirp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
}

func (c *Chirp) Delete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(token, c.ApiCfg.Secret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirpId := r.PathValue("chirpID")
	if chirpId == "" {
		w.WriteHeader(http.StatusBadRequest)
		errDto := dto.Response{
			Error: "chirp id cannot be empty",
		}
		errorResponse, _ := json.Marshal(errDto)
		_, err = w.Write(errorResponse)
		return
	}

	chirp, err := c.ApiCfg.Db.GetChirpById(r.Context(), uuid.MustParse(chirpId))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if userID != chirp.UserID {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = c.ApiCfg.Db.DeleteChirp(r.Context(), chirp.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validateChirp(userID uuid.UUID, r *http.Request) (*dto.ChirpRequest, error) {
	invalidWords := []string{"kerfuffle", "sharbert", "fornax"}
	replacement := "****"

	decoder := json.NewDecoder(r.Body)
	var chirp dto.ChirpRequest
	err := decoder.Decode(&chirp)
	if err != nil {
		return nil, fmt.Errorf("error decoding request: %s\n", err)
	}

	isValid, err := func() (bool, error) {
		if len(chirp.Body) > 140 {
			return false, fmt.Errorf("chirp is too long")
		}
		if userID == uuid.Nil {
			return false, fmt.Errorf("user id cannot be empty")
		}
		return true, nil
	}()

	if isValid {
		result := chirp.Body
		for _, word := range invalidWords {
			re := regexp.MustCompile("(?i)" + regexp.QuoteMeta(word))
			result = re.ReplaceAllString(result, replacement)
		}

		return &chirp, nil
	} else {
		return nil, err
	}
}
