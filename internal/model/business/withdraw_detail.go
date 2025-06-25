package business

import "time"

type WithdrawDetail struct {
	OrderNumber string
	Sum         float64
	CreatedAt   time.Time
}
