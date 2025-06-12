package model

import (
	"github.com/google/uuid"
	"time"
)

type UserData struct {
	ID        uuid.UUID
	Login     string
	Password  string
	CreatedAt time.Time
}
