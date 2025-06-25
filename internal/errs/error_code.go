package errs

const (
	Generic                 = "generic"
	OrderAddedByCurrentUser = "order added by current user"
	OrderAddedByAnotherUser = "order added By another user"
	InvalidOrderNumber      = "invalid order number"
	NotEnoughAccrual        = "not enough accrual on the user account"
	OrderStatusClient       = "something went wrong on handling response from Accrual service"
)
