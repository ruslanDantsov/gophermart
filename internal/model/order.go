package model

import (
	"github.com/google/uuid"
	"time"
)

const (
	ORDER_NEW_STATUS        = "NEW"
	ORDER_PROCESSING_STATUS = "PROCESSING"
)

type Order struct {
	ID        uuid.UUID
	Number    string
	Status    string
	Accrual   float64
	CreatedAt time.Time
	UserID    uuid.UUID
}
