package dto

import "github.com/google/uuid"

type UserWebhook struct {
	Event string          `json:"event"`
	Data  UserWebhookData `json:"data"`
}

type UserWebhookData struct {
	UserID uuid.UUID `json:"user_id"`
}
