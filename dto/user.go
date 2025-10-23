package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email,"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}
