package view

import (
	"github.com/google/uuid"
	"time"
)

//go:generate easyjson -all user_view_model.go
type UserViewModel struct {
	Id        uuid.UUID `json:"id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
}
