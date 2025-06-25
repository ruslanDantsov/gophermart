package balance

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/dto/view"
	"github.com/ruslanDantsov/gophermart/internal/handler/middleware"
	"github.com/ruslanDantsov/gophermart/internal/model/business"
	"go.uber.org/zap"
	"net/http"
)

type IBalanceService interface {
	GetBalance(ctx context.Context, userID uuid.UUID) (*business.Balance, error)
}

type BalanceHandler struct {
	Log            zap.Logger
	BalanceService IBalanceService
}

func NewBalanceHandler(log *zap.Logger, balanceService IBalanceService) *BalanceHandler {
	return &BalanceHandler{
		Log:            *log,
		BalanceService: balanceService,
	}
}

func (h *BalanceHandler) HandleGetBalance(ginContext *gin.Context) {
	currentUserID := ginContext.Request.Context().Value(middleware.CtxUserIDKey{}).(uuid.UUID)
	balance, err := h.BalanceService.GetBalance(ginContext.Request.Context(), currentUserID)

	if err != nil {
		h.Log.Error(err.Error())
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong on request processing"})
		return
	}

	viewModel := view.BalanceViewModel{
		Current:   balance.Total,
		Withdrawn: balance.Withdrawn,
	}

	ginContext.Header("Content-Type", "application/json")
	ginContext.Writer.WriteHeader(http.StatusOK)
	ginContext.JSON(http.StatusOK, viewModel)
}
