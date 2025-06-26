package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type CommonHandler struct {
	log zap.Logger
}

func NewCommonHandler(log *zap.Logger) *CommonHandler {
	return &CommonHandler{
		log: *log,
	}
}

func (h *CommonHandler) HandleUnsupportedRequest(ginContext *gin.Context) {
	h.log.Warn(fmt.Sprintf("Request is unsupported: url: %v; method: %v",
		ginContext.Request.RequestURI,
		ginContext.Request.Method))
	ginContext.String(http.StatusNotFound, "Request is unsupported")
}
