package command

type OrderCreateCommand struct {
	Number string `json:"number" binding:"required"`
}
