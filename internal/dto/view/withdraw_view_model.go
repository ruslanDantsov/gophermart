package view

import (
	"time"
)

//go:generate easyjson -all withdraw_view_model.go
type WithdrawViewModel struct {
	OrderNumber string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
