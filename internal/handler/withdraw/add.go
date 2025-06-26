package withdraw

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/errs"
	"github.com/ruslanDantsov/gophermart/internal/handler/middleware"
	"github.com/ruslanDantsov/gophermart/internal/model/business"
	"github.com/ruslanDantsov/gophermart/internal/model/entity"
	"go.uber.org/zap"
	"net/http"
)

type WithdrawCreator interface {
	AddWithdraw(ctx context.Context, withdrawCreateCommand command.WithdrawCreateCommand, authUserID uuid.UUID) (*entity.Withdraw, error)
}

type WithdrawGetter interface {
	GetWithdrawDetails(ctx context.Context) ([]business.WithdrawDetail, error)
}

type WithdrawHandler struct {
	log                    zap.Logger
	withdrawCreatorService WithdrawCreator
	withdrawGetterService  WithdrawGetter
}

func NewWithdrawHandler(log *zap.Logger, withdrawCreatorService WithdrawCreator, withdrawGetterService WithdrawGetter) *WithdrawHandler {
	return &WithdrawHandler{
		log:                    *log,
		withdrawCreatorService: withdrawCreatorService,
		withdrawGetterService:  withdrawGetterService,
	}
}
func (h *WithdrawHandler) HandleAddingWithdraw(ginContext *gin.Context) {
	contentType := ginContext.GetHeader("Content-Type")
	if contentType != "application/json" {
		h.log.Error(fmt.Sprintf("Unsupported content type: %s ", contentType))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported content type"})
		return
	}

	var withdrawCreateCommand command.WithdrawCreateCommand

	if err := ginContext.ShouldBindJSON(&withdrawCreateCommand); err != nil {
		h.log.Error(fmt.Sprintf("Invalid JSON: %s", err.Error()))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	authUserID := ginContext.Request.Context().Value(middleware.CtxUserIDKey{}).(uuid.UUID)
	_, err := h.withdrawCreatorService.AddWithdraw(ginContext.Request.Context(), withdrawCreateCommand, authUserID)
	if err != nil {
		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case errs.NotEnoughAccrual:
				ginContext.JSON(http.StatusPaymentRequired, gin.H{"error": appErr.Message})
			case errs.OrderAddedByCurrentUser:
				ginContext.Writer.WriteHeader(http.StatusOK)
			case errs.OrderAddedByAnotherUser:
				ginContext.Writer.WriteHeader(http.StatusConflict)
			case errs.InvalidOrderNumber:
				ginContext.JSON(http.StatusUnprocessableEntity, gin.H{"error": appErr.Message})
			default:
				ginContext.JSON(http.StatusInternalServerError, gin.H{"error": appErr.Message})
			}
			h.log.Error(fmt.Sprintf(appErr.Message+", description: %s ", err.Error()))
			return
		}

		h.log.Error(fmt.Sprintf("Unexpected error: %s", err.Error()))
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	ginContext.Status(http.StatusOK)
}
