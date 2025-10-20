package handler

import (
	"chirpy/core/config"
	"chirpy/dto"
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
	var createUser dto.CreateUser
	err := decoder.Decode(&createUser)
	if err != nil {
		fmt.Printf("error decoding user request: %s\n", err)
		errDto := dto.Response{
			Error: err.Error(),
		}
		errorResponse, _ := json.Marshal(errDto)
		_, err = w.Write(errorResponse)
	}

	user, err := u.ApiCfg.Db.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Email:     createUser.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

	userDto, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("error creating user: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(userDto)
}
