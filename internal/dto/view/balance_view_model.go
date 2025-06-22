package view

//go:generate easyjson -all balance_view_model.go
type BalanceViewModel struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
