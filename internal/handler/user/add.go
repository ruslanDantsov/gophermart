package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/dto/view"
	"github.com/ruslanDantsov/gophermart/internal/model/entity"
	"github.com/ruslanDantsov/gophermart/internal/service"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type IUserService interface {
	AddUser(ctx context.Context, userCreateCommand command.UserCreateCommand) (*entity.UserData, error)
	FindByLoginAndPassword(ctx context.Context, login string, password string) (*entity.UserData, error)
}

type IAuthUtils interface {
	GenerateJWT(id uuid.UUID, username string) (*service.TokenResult, error)
}

type UserHandler struct {
	Log         zap.Logger
	UserService IUserService
	AuthService IAuthUtils
}

func NewUserHandler(log *zap.Logger, userService IUserService, authService IAuthUtils) *UserHandler {
	return &UserHandler{
		Log:         *log,
		UserService: userService,
		AuthService: authService,
	}
}

func (h *UserHandler) HandleRegisterUser(ginContext *gin.Context) {
	//TODO: check for content type

	var userCreateCommand command.UserCreateCommand
	if err := ginContext.ShouldBindJSON(&userCreateCommand); err != nil {
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
			h.Log.Error("Failed to parse register user body request: " + err.Error())
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
	}

	userData, err := h.UserService.AddUser(ginContext.Request.Context(), userCreateCommand)
	if err != nil {
		h.Log.Error("Failed to save user: " + err.Error())
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ginContext.Header("Content-Type", "application/json")
	ginContext.Writer.WriteHeader(http.StatusOK)

	userViewModel := view.UserViewModel{
		ID:        userData.ID,
		Login:     userData.Login,
		CreatedAt: userData.CreatedAt,
	}
	_, err = easyjson.MarshalToWriter(userViewModel, ginContext.Writer)

	if err != nil {
		h.Log.Error(fmt.Sprintf("error on marshal user data response %v", err))
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong on marshal user data response"})
	}

}
