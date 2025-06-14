package order

import (
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/gophermart/internal/dto/view"
	"net/http"
)

func (h *OrderHandler) HandleGetOrders(ginContext *gin.Context) {

	orders, err := h.OrderService.GetOrders(ginContext.Request.Context())

	if err != nil {
		h.Log.Error(err.Error())
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong on request processing"})
		return
	}

	if len(orders) == 0 {
		ginContext.Status(http.StatusNoContent)
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
	ginContext.Writer.WriteHeader(http.StatusOK)
	ginContext.JSON(http.StatusOK, viewModels)
}
