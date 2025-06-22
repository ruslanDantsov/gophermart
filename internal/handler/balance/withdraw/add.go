package withdraw

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

type IWithdrawService interface {
	AddWithdraw(ctx context.Context, withdrawCreateCommand command.WithdrawCreateCommand) (*entity.Withdraw, error)
}

type WithdrawHandler struct {
	Log             zap.Logger
	WithdrawService IWithdrawService
}

func NewWithdrawHandler(log *zap.Logger, withdrawService IWithdrawService) *WithdrawHandler {
	return &WithdrawHandler{
		Log:             *log,
		WithdrawService: withdrawService,
	}
}
func (h *WithdrawHandler) HandleWithdraw(ginContext *gin.Context) {
	contentType := ginContext.GetHeader("Content-Type")
	if contentType != "application/json" {
		h.Log.Error(fmt.Sprintf("Unsupported content type: %s ", contentType))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported content type"})
		return
	}

	var withdrawCreateCommand command.WithdrawCreateCommand

	if err := ginContext.ShouldBindJSON(&withdrawCreateCommand); err != nil {
		h.Log.Error(fmt.Sprintf("Invalid JSON: %s", err.Error()))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	_, err := h.WithdrawService.AddWithdraw(ginContext.Request.Context(), withdrawCreateCommand)
	//TODO: 402 — на счету недостаточно средств;
	if err != nil {
		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case errs.ORDER_ADDED_BY_CURRENT_USER:
				ginContext.Writer.WriteHeader(http.StatusOK)
			case errs.ORDER_ADDED_BY_ANOTHER_USER:
				ginContext.Writer.WriteHeader(http.StatusConflict)
			case errs.INVALID_ORDER_NUMBER:
				ginContext.JSON(http.StatusUnprocessableEntity, gin.H{"error": appErr.Message})
			default:
				ginContext.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Message})
			}
			h.Log.Error(fmt.Sprintf(appErr.Message+", description: %s ", err.Error()))
			return
		}

		h.Log.Error(fmt.Sprintf("Unexpected error: %s", err.Error()))
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	ginContext.Status(http.StatusOK)
}
