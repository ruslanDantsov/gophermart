package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ruslanDantsov/gophermart/internal/config"
	"github.com/ruslanDantsov/gophermart/internal/handler"
	"go.uber.org/zap"
	"net/http"
)

type GophermartApp struct {
	cfg           *config.Config
	logger        *zap.Logger
	commonHandler *handler.CommonHandler
}

func NewGophermartApp(cfg *config.Config, log *zap.Logger) (*GophermartApp, error) {
	commonHandler := handler.NewCommonHandler(*log)

	return &GophermartApp{
		cfg:           cfg,
		logger:        log,
		commonHandler: commonHandler,
	}, nil
}

func (app *GophermartApp) Run() error {
	router := gin.Default()

	router.POST("/value/", app.getMetricHandler.GetJSON)

	router.NoRoute(app.commonHandler.ServeHTTP)

	return http.ListenAndServe(app.cfg.Address, router)

}
