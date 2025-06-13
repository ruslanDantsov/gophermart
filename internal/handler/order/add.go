package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/errs"
	"github.com/ruslanDantsov/gophermart/internal/model"
	"go.uber.org/zap"
	"net/http"
)

type IOrderService interface {
	AddOrder(ctx context.Context, orderCreateCommand command.OrderCreateCommand) (*model.Order, error)
}

type OrderHandler struct {
	Log          zap.Logger
	OrderService IOrderService
}

func NewOrderHandler(log *zap.Logger, orderService IOrderService) *OrderHandler {
	return &OrderHandler{
		Log:          *log,
		OrderService: orderService,
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

	if err := goluhn.Validate(string(orderNumber)); err != nil {
		h.Log.Error(fmt.Sprintf("Invalid order number: %s ", err.Error()))
		ginContext.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid order number"})
		return
	}

	_, err = h.OrderService.AddOrder(ginContext.Request.Context(), command.OrderCreateCommand{Number: string(orderNumber)})

	if err != nil {
		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case errs.ORDER_ADDED_BY_CURRENT_USER:
				ginContext.Writer.WriteHeader(http.StatusOK)
			case errs.ORDER_ADDED_BY_ANOTHER_USER:
				ginContext.Writer.WriteHeader(http.StatusConflict)
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
