package withdraw

import (
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/gophermart/internal/dto/view"
	"net/http"
)

func (h *WithdrawHandler) HandleGetWithdraws(ginContext *gin.Context) {
	withdraws, err := h.WithdrawGetterService.GetWithdrawDetails(ginContext.Request.Context())

	if err != nil {
		h.Log.Error(err.Error())
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong on request processing"})
		return
	}

	if len(withdraws) == 0 {
		ginContext.JSON(http.StatusNoContent, gin.H{})
		return
	}

	viewModels := make([]view.WithdrawViewModel, len(withdraws))
	for i, withdraw := range withdraws {
		viewModels[i] = view.WithdrawViewModel{
			OrderNumber: withdraw.OrderNumber,
			Sum:         withdraw.Sum,
			ProcessedAt: withdraw.CreatedAt,
		}
	}

	ginContext.Header("Content-Type", "application/json")
	ginContext.Writer.WriteHeader(http.StatusOK)
	ginContext.JSON(http.StatusOK, viewModels)
}
