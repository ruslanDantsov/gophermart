package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/gophermart/internal/config"
	"github.com/ruslanDantsov/gophermart/internal/handler"
	"github.com/ruslanDantsov/gophermart/internal/handler/middleware"
	"github.com/ruslanDantsov/gophermart/internal/handler/order"
	"github.com/ruslanDantsov/gophermart/internal/handler/user"
	"github.com/ruslanDantsov/gophermart/internal/infrastructure/storage/postgre"
	"github.com/ruslanDantsov/gophermart/internal/repository"
	"github.com/ruslanDantsov/gophermart/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type GophermartApp struct {
	cfg           *config.Config
	logger        *zap.Logger
	storage       *postgre.PostgreStorage
	commonHandler *handler.CommonHandler
	userHandler   *user.UserHandler
	orderHandler  *order.OrderHandler
}

func NewGophermartApp(ctx context.Context, cfg *config.Config, log *zap.Logger) (*GophermartApp, error) {
	storage, err := postgre.NewPostgreStorage(ctx, log, cfg.DatabaseConnection)
	if err != nil {
		return nil, err
	}

	userRepository := repository.NewUserRepository(storage)
	passwordService := &service.PasswordService{}
	userService := service.NewUserService(userRepository, passwordService)

	commonHandler := handler.NewCommonHandler(log)

	authService := service.NewAuthService(cfg.JWTSecret)
	userHandler := user.NewUserHandler(log, userService, authService)

	orderHandler := order.NewOrderHandler(log)

	return &GophermartApp{
		cfg:           cfg,
		logger:        log,
		storage:       storage,
		commonHandler: commonHandler,
		userHandler:   userHandler,
		orderHandler:  orderHandler,
	}, nil
}

func (app *GophermartApp) Run(ctx context.Context) error {
	router := gin.Default()

	router.POST("/api/user/register", app.userHandler.HandleRegisterUser)
	router.POST("/api/user/login", app.userHandler.HandleAuthentication)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(app.cfg.JWTSecret, app.logger))

	protected.POST("/api/user/orders", app.orderHandler.HandleRegisterOrder)

	router.NoRoute(app.commonHandler.HandleUnsupportedRequest)

	return http.ListenAndServe(app.cfg.Address, router)

}
