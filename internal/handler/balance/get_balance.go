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

type BalanceGetter interface {
	GetBalance(ctx context.Context, userID uuid.UUID) (*business.Balance, error)
}

type BalanceHandler struct {
	log            zap.Logger
	balanceService BalanceGetter
}

func NewBalanceHandler(log *zap.Logger, balanceService BalanceGetter) *BalanceHandler {
	return &BalanceHandler{
		log:            *log,
		balanceService: balanceService,
	}
}

func (h *BalanceHandler) HandleGetBalance(ginContext *gin.Context) {
	currentUserID := ginContext.Request.Context().Value(middleware.CtxUserIDKey{}).(uuid.UUID)
	balance, err := h.balanceService.GetBalance(ginContext.Request.Context(), currentUserID)

	if err != nil {
		h.log.Error(err.Error())
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong on request processing"})
		return
	}

	viewModel := view.BalanceViewModel{
		Current:   balance.Total,
		Withdrawn: balance.Withdrawn,
	}

	ginContext.Header("Content-Type", "application/json")
	ginContext.JSON(http.StatusOK, viewModel)
}
