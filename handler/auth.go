package handler

import (
	"chirpy/core/config"
	"chirpy/dto"
	"chirpy/internal/auth"
	"encoding/json"
	"net/http"
)

func NewAuth(apiCfg *config.ApiConfig) *Auth {
	return &Auth{
		ApiConfig: apiCfg,
	}
}

type Auth struct {
	ApiConfig *config.ApiConfig
}

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginRequest)
	if err != nil {
		errDto := dto.Response{
			Error: err.Error(),
		}
		errorResponse, _ := json.Marshal(errDto)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errorResponse)
	}

	user, err := a.ApiConfig.Db.GetUserByEmail(r.Context(), loginRequest.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write([]byte("Incorrect email or password"))
		return
	}

	hash, err := auth.CheckPasswordHash(loginRequest.Password, user.HashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		resp, err := json.Marshal(dto.Response{
			Error: "Incorrect email or password",
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(resp)
		return
	}

	if hash {
		userResponse := dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		responseBytes, err := json.Marshal(userResponse)
		if err != nil {
			errDto := dto.Response{
				Error: err.Error(),
			}
			errorResponse, _ := json.Marshal(errDto)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write(errorResponse)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(responseBytes)
		return
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		resp, err := json.Marshal(dto.Response{
			Error: "Incorrect email or password",
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(resp)
		return
	}
}
