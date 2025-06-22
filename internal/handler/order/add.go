package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/errs"
	"github.com/ruslanDantsov/gophermart/internal/model/entity"
	"go.uber.org/zap"
	"net/http"
)

type IOrderCreatorService interface {
	AddOrder(ctx context.Context, orderCreateCommand command.OrderCreateCommand) (*entity.Order, error)
}

type IOrderGetterService interface {
	GetOrders(ctx context.Context) ([]entity.Order, error)
}

type OrderHandler struct {
	Log                 zap.Logger
	OrderCreatorService IOrderCreatorService
	OrderGetterService  IOrderGetterService
}

func NewOrderHandler(log *zap.Logger, orderCreatorService IOrderCreatorService, orderGetterService IOrderGetterService) *OrderHandler {
	return &OrderHandler{
		Log:                 *log,
		OrderCreatorService: orderCreatorService,
		OrderGetterService:  orderGetterService,
	}
}

func (h *OrderHandler) HandleRegisterOrder(ginContext *gin.Context) {
	contentType := ginContext.GetHeader("Content-Type")
	if contentType != "text/plain" {
		h.Log.Error(fmt.Sprintf("Unsupported content type: %s ", contentType))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported content type"})
		return
	}

	orderNumber, err := ginContext.GetRawData()

	if err != nil {
		h.Log.Error(fmt.Sprintf("Invalid request body: %s ", err.Error()))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	_, err = h.OrderCreatorService.AddOrder(ginContext.Request.Context(), command.OrderCreateCommand{Number: string(orderNumber)})

	if err != nil {
		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case errs.OrderAddedByCurrentUser:
				ginContext.Writer.WriteHeader(http.StatusOK)
			case errs.OrderAddedByAnotherUser:
				ginContext.Writer.WriteHeader(http.StatusConflict)
			case errs.InvalidOrderNumber:
				ginContext.JSON(http.StatusUnprocessableEntity, gin.H{"error": appErr.Message})
			default:
				ginContext.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Message})
			}
			h.Log.Error(fmt.Sprintf(appErr.Message+", description: %s ", err.Error()))
			return
		}
	}

	ginContext.Header("Content-Type", "application/json")
	ginContext.Writer.WriteHeader(http.StatusAccepted)
}
