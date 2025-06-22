package entity

import (
	"github.com/google/uuid"
	"time"
)

type Withdraw struct {
	ID        uuid.UUID
	Sum       float64
	CreatedAt time.Time
	OrderID   uuid.UUID
}
