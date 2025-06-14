package view

import (
	"time"
)

//go:generate easyjson -all order_view_model.go
type OrderViewModel struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}
