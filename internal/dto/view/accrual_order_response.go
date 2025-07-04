package view

//go:generate easyjson -all accrual_order_response.go

const (
	AccrualOrderRegisteredStatus = "REGISTERED"
	AccrualOrderProcessingStatus = "PROCESSING"
	AccrualOrderInvalidStatus    = "INVALID"
	AccrualOrderProcessedStatus  = "PROCESSED"
)

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
