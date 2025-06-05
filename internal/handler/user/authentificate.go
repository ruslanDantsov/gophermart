package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *UserHandler) HandleAuthentication(ginContext *gin.Context) {
	ginContext.String(http.StatusOK, "User has been authenticated")
}
