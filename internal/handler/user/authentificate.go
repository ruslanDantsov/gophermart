package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"net/http"
	"strings"
)

func (h *UserHandler) HandleAuthentication(ginContext *gin.Context) {
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
			h.Log.Error("Validation failed: " + strings.Join(errorMessages, ", "))
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
			return
		} else {
			h.Log.Error("Failed to parse Auth user body request: " + err.Error())
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
	}

	userData, err := h.UserService.FindByLoginAndPassword(ginContext.Request.Context(), authCommand.Login, authCommand.Password)
	if err != nil {
		h.Log.Error(fmt.Sprintf("User %s not found: %s", authCommand.Login, err.Error()))
		ginContext.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("User %s not found", authCommand.Login)})
		return
	}

	tokenResult, err := h.AuthService.GenerateJWT(uuid.New(), userData.Login)
	if err != nil {
		h.Log.Error("Failed to generate token: " + err.Error())
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"access_token": tokenResult.AccessToken,
		"expires_in":   tokenResult.ExpiresIn,
		"token_type":   "Bearer",
	})
}
