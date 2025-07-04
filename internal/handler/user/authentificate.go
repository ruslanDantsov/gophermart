package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"net/http"
	"strings"
)

func (h *UserHandler) HandleAuthentication(ginContext *gin.Context) {
	contentType := ginContext.GetHeader("Content-Type")
	if contentType != "application/json" {
		h.log.Error(fmt.Sprintf("Unsupported content type: %s ", contentType))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported content type"})
		return
	}

	var authCommand command.UserAuthCommand
	if err := ginContext.ShouldBindJSON(&authCommand); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := make([]string, 0, len(ve))
			for _, fe := range ve {
				field := fe.Field()
				tag := fe.Tag()
				errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' is %s", field, tag))
			}
			h.log.Error("Validation failed: " + strings.Join(errorMessages, ", "))
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
			return
		} else {
			h.log.Error("Failed to parse Auth user body request: " + err.Error())
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
	}

	userData, err := h.userManager.FindByLoginAndPassword(ginContext.Request.Context(), authCommand.Login, authCommand.Password)
	if err != nil {
		h.log.Error(fmt.Sprintf("User %s not found: %s", authCommand.Login, err.Error()))
		ginContext.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("User %s not found", authCommand.Login)})
		return
	}

	tokenResult, err := h.authManager.GenerateJWT(userData.ID, userData.Login)
	if err != nil {
		h.log.Error("Failed to generate token: " + err.Error())
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ginContext.Header("Authorization", "Bearer "+tokenResult.AccessToken)
	ginContext.JSON(http.StatusOK, gin.H{
		"access_token": tokenResult.AccessToken,
		"expires_in":   tokenResult.ExpiresIn,
		"token_type":   "Bearer",
	})
}
