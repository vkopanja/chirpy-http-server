package handler

import (
	"chirpy/core/config"
	"chirpy/dto"
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func NewUser(apiCfg *config.ApiConfig) *User {
	return &User{
		ApiCfg: apiCfg,
	}
}

type User struct {
	ApiCfg *config.ApiConfig
}

func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var createUser dto.CreateOrUpdateUser
	err := decoder.Decode(&createUser)
	if err != nil {
		fmt.Printf("error decoding user request: %s\n", err)
		errDto := dto.Response{
			Error: err.Error(),
		}
		errorResponse, _ := json.Marshal(errDto)
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(errorResponse)
	}

	hash, err := auth.HashPassword(createUser.Password)
	if err != nil {
		fmt.Printf("error hashing password: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		return
	}

	user, err := u.ApiCfg.Db.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New(),
		Email:          createUser.Email,
		HashedPassword: hash,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		fmt.Printf("error creating user: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		errDto := dto.Response{
			Error: err.Error(),
		}
		errorResponse, _ := json.Marshal(errDto)
		_, err = w.Write(errorResponse)
		return
	}

	userDto, err := json.Marshal(dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
	if err != nil {
		fmt.Printf("error creating user: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(userDto)
}

func (u *User) Update(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(token, u.ApiCfg.Secret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var updateUser dto.CreateOrUpdateUser
	err = decoder.Decode(&updateUser)
	if err != nil {
		fmt.Printf("error decoding user request: %s\n", err)
		errDto := dto.Response{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		errorResponse, _ := json.Marshal(errDto)
		_, err = w.Write(errorResponse)
	}

	user, err := u.ApiCfg.Db.GetUserByID(r.Context(), userID)
	if err != nil {
		fmt.Printf("error getting user: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.ID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	hashedPassword, err := auth.HashPassword(updateUser.Password)
	if err != nil {
		fmt.Printf("error hashing password: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	updatedUser, err := u.ApiCfg.Db.UpdateUserByID(r.Context(), database.UpdateUserByIDParams{
		Email:          updateUser.Email,
		HashedPassword: hashedPassword,
		ID:             user.ID,
	})
	if err != nil {
		fmt.Printf("error updating user: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userDto, err := json.Marshal(dto.UserResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	})
	if err != nil {
		fmt.Printf("error updating user: %s\n", err)
		errDto := dto.Response{
			Error: err.Error(),
		}
		errorResponse, _ := json.Marshal(errDto)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(userDto)
}
