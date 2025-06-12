package order

import (
	"fmt"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
	"net/http"
)

type OrderHandler struct {
	Log zap.Logger
}

func NewOrderHandler(log *zap.Logger) *OrderHandler {
	return &OrderHandler{
		Log: *log,
	}
}

func (h *OrderHandler) HandleRegisterOrder(ginContext *gin.Context) {
	contentType := ginContext.GetHeader("Content-Type")
	if contentType != "text/plain" {
		h.Log.Info(fmt.Sprintf("Unsupported content type: %s ", contentType))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported content type"})
		return
	}

	orderNumber, err := ginContext.GetRawData()
	h.Log.Info(fmt.Sprintf("Order %s has been registered", orderNumber))

	if err != nil {
		h.Log.Info(fmt.Sprintf("Invalid request body: %s ", err.Error()))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := goluhn.Validate(string(orderNumber)); err != nil {
		h.Log.Info(fmt.Sprintf("Invalid order number: %s ", err.Error()))
		ginContext.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid order number"})
		return
	}

	ginContext.Header("Content-Type", "application/json")
	ginContext.Writer.WriteHeader(http.StatusAccepted)

}
