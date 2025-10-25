package handler

import (
	"chirpy/core/config"
	"chirpy/dto"
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

func NewAuth(apiCfg *config.ApiConfig) *Auth {
	return &Auth{
		apiConfig: apiCfg,
	}
}

type Auth struct {
	apiConfig *config.ApiConfig
}

// Login
// @Summary Login
// @Description Login the user to application and return JWT token
// @Tags Auth handler
// @Accept application/json
// @Produce application/json
// @Param loginRequest body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.UserResponse "User response"
// @Failure 400 {object} dto.Response "Bad request"
// @Failure 401 {object} dto.Response "Unauthorized"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /api/login [post]
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

	user, err := a.apiConfig.Db.GetUserByEmail(r.Context(), loginRequest.Email)
	if err != nil {
		errDto := dto.Response{
			Error: err.Error(),
		}
		errorResponse, _ := json.Marshal(errDto)
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write(errorResponse)
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
		jwt, err := auth.MakeJWT(user.ID, a.apiConfig.Secret, time.Duration(60*60)*time.Second)
		if err != nil {
			errDto := dto.Response{
				Error: err.Error(),
			}
			errorResponse, _ := json.Marshal(errDto)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write(errorResponse)
			return
		}
		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		_, err = a.apiConfig.Db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:  refreshToken,
			UserID: user.ID,
			ExpiresAt: sql.NullTime{
				Time:  time.Now().AddDate(0, 0, 60),
				Valid: true,
			},
			RevokedAt: sql.NullTime{},
			CreatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userResponse := dto.UserResponse{
			ID:           user.ID,
			Email:        user.Email,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			IsChirpyRed:  user.IsChirpyRed,
			Token:        jwt,
			RefreshToken: refreshToken,
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

// Refresh
// @Summary Refresh token
// @Description Refresh the token using our refresh token
// @Tags Auth handler
// @Accept application/json
// @Produce application/json
// @Success 200 {object} dto.Token "Token response"
// @Failure 400 {object} dto.Response "Bad request"
// @Failure 401 {object} nil "Unauthorized"
// @Failure 500 {object} dto.Response "Internal server error"
// @Router /api/refresh [post]
func (a *Auth) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	refresh, err := a.apiConfig.Db.GetTokenForRefreshToken(r.Context(), refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if refreshToken != refresh.Token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jwt, err := auth.MakeJWT(refresh.UserID, a.apiConfig.Secret, time.Duration(60*60)*time.Second)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	tokenResponse := dto.Token{Token: jwt}
	tokenRespBytes, err := json.Marshal(tokenResponse)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_, err = w.Write(tokenRespBytes)
}

// Revoke
// @Summary Revoke token
// @Description Revoke the refresh token for the passed in Bearer token
// @Tags Auth handler
// @Security BearerAuth
// @Accept application/json
// @Produce application/json
// @Success 204 "No content"
// @Failure 500 {object} nil "Internal server error"
// @Router /api/revoke [post]
func (a *Auth) Revoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = a.apiConfig.Db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
