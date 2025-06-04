package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type CommonHandler struct {
	Log zap.Logger
}

func NewCommonHandler(log zap.Logger) *CommonHandler {
	return &CommonHandler{
		Log: log,
	}
}

func (h *CommonHandler) ServeHTTP(ginContext *gin.Context) {
	h.Log.Error(fmt.Sprintf("Request is unsupported: url: %v; method: %v",
		ginContext.Request.RequestURI,
		ginContext.Request.Method))
	ginContext.String(http.StatusNotFound, "Request is unsupported")
}
