package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrUpdateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email,"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}
