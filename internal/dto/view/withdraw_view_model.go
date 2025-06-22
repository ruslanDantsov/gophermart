package view

import (
	"github.com/google/uuid"
	"time"
)

//go:generate easyjson -all withdraw_view_model.go
type WithdrawViewModel struct {
	OrderID     uuid.UUID `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
