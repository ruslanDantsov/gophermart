package entity

import (
	"github.com/google/uuid"
	"time"
)

const (
	OrderNewStatus        = "NEW"
	OrderProcessingStatus = "PROCESSING"
	OrderInvalidStatus    = " INVALID"
	OrderProcessedStatus  = "PROCESSED"
)

type Order struct {
	ID        uuid.UUID
	Number    string
	Status    string
	Accrual   float64
	CreatedAt time.Time
	UserID    uuid.UUID
}
