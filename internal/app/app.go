package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/gophermart/internal/client"
	"github.com/ruslanDantsov/gophermart/internal/config"
	"github.com/ruslanDantsov/gophermart/internal/handler"
	"github.com/ruslanDantsov/gophermart/internal/handler/balance"
	"github.com/ruslanDantsov/gophermart/internal/handler/middleware"
	"github.com/ruslanDantsov/gophermart/internal/handler/order"
	"github.com/ruslanDantsov/gophermart/internal/handler/user"
	"github.com/ruslanDantsov/gophermart/internal/handler/withdraw"
	"github.com/ruslanDantsov/gophermart/internal/infrastructure/storage/postgre"
	"github.com/ruslanDantsov/gophermart/internal/repository"
	"github.com/ruslanDantsov/gophermart/internal/service"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type GophermartApp struct {
	cfg                 *config.Config
	logger              *zap.Logger
	accrualOrderService *service.AccrualOrderService
	storage             *postgre.PostgreStorage
	commonHandler       *handler.CommonHandler
	userHandler         *user.UserHandler
	orderHandler        *order.OrderHandler
	balanceHandler      *balance.BalanceHandler
	withdrawHandler     *withdraw.WithdrawHandler
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

	orderRepository := repository.NewOrderRepository(storage)
	orderService := service.NewOrderService(orderRepository, storage)
	orderHandler := order.NewOrderHandler(log, orderService, orderService)

	withdrawRepository := repository.NewWithdrawnRepository(storage)
	withdrawService := service.NewWithdrawService(orderService, withdrawRepository, orderRepository, storage)
	withdrawHandler := withdraw.NewWithdrawHandler(log, withdrawService, withdrawService)

	balanceService := service.NewBalanceService(orderRepository, withdrawRepository)
	balanceHandler := balance.NewBalanceHandler(log, balanceService)

	orderStatusClient := client.NewOrderStatusClient(cfg.AccrualSystemAddress)
	accrualOrderService := service.NewAccrualOrderService(orderService, orderStatusClient, log)

	return &GophermartApp{
		cfg:                 cfg,
		logger:              log,
		storage:             storage,
		commonHandler:       commonHandler,
		userHandler:         userHandler,
		orderHandler:        orderHandler,
		balanceHandler:      balanceHandler,
		withdrawHandler:     withdrawHandler,
		accrualOrderService: accrualOrderService,
	}, nil
}

func (app *GophermartApp) Run(ctx context.Context) error {
	router := gin.Default()

	router.POST("/api/user/register", app.userHandler.HandleRegisterUser)
	router.POST("/api/user/login", app.userHandler.HandleAuthentication)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(app.cfg.JWTSecret, app.logger))

	protected.POST("/api/user/orders", app.orderHandler.HandleRegisterOrder)
	protected.GET("/api/user/orders", app.orderHandler.HandleGetOrders)

	protected.GET("/api/user/balance", app.balanceHandler.HandleGetBalance)

	protected.POST("/api/user/balance/withdraw", app.withdrawHandler.HandleAddingWithdraw)
	protected.GET("/api/user/withdrawals", app.withdrawHandler.HandleGetWithdraws)

	router.NoRoute(app.commonHandler.HandleUnsupportedRequest)

	srv := &http.Server{
		Addr:    app.cfg.Address,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Fatal("Server error", zap.Error(err))
		}
	}()

	app.logger.Info("Server started")

	go func() {
		ticker := time.NewTicker(app.cfg.ReportInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				app.accrualOrderService.ProcessOrders(ctx)
			case <-ctx.Done():
				app.logger.Info("OrderStatusSyncService received shutdown signal")
				return
			}
		}
	}()

	<-ctx.Done()
	app.logger.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), app.cfg.GracefulShutdownInterval)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	app.logger.Info("Server exited properly")
	return nil
}
