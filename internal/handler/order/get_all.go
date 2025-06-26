package order

import (
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/gophermart/internal/dto/view"
	"net/http"
)

func (h *OrderHandler) HandleGetOrders(ginContext *gin.Context) {
	orders, err := h.orderGetterService.GetOrders(ginContext.Request.Context())

	if err != nil {
		h.log.Error(err.Error())
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong on request processing"})
		return
	}

	if len(orders) == 0 {
		ginContext.JSON(http.StatusNoContent, gin.H{})
		return
	}

	viewModels := make([]view.OrderViewModel, len(orders))
	for i, order := range orders {
		viewModels[i] = view.OrderViewModel{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.CreatedAt,
		}
	}

	ginContext.Header("Content-Type", "application/json")
	ginContext.JSON(http.StatusOK, viewModels)
}
