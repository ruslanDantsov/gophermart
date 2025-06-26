package service

import (
	"context"
	"github.com/ruslanDantsov/gophermart/internal/dto/view"
	"go.uber.org/zap"
)

type IUnprocessedOrderService interface {
	GetUnprocessedOrders(ctx context.Context) ([]string, error)
	UpdateAccrualData(ctx context.Context, number string, accrual float64, status string) error
}

type IAccrualClient interface {
	GetAccrualData(ctx context.Context, orderID string) (*view.AccrualResponse, error)
}

type AccrualOrderService struct {
	UnprocessedOrderService IUnprocessedOrderService
	OrderStatusClient       IAccrualClient
	Log                     *zap.Logger
}

func NewAccrualOrderService(unprocessedOrderService IUnprocessedOrderService, orderStatusClient IAccrualClient, log *zap.Logger) *AccrualOrderService {
	return &AccrualOrderService{
		UnprocessedOrderService: unprocessedOrderService,
		OrderStatusClient:       orderStatusClient,
		Log:                     log,
	}
}

func (s *AccrualOrderService) ProcessOrders(ctx context.Context) {
	s.Log.Info("Starting process for updating accrual data ...")
	unprocessedOrderNumbers, err := s.UnprocessedOrderService.GetUnprocessedOrders(ctx)
	if err != nil {
		s.Log.Error("Something went wrong on getting orders",
			zap.String("error", err.Error()),
		)
	}

	processedOrderCount := 0
	for _, orderNumber := range unprocessedOrderNumbers {
		accrualResponse, err := s.OrderStatusClient.GetAccrualData(ctx, orderNumber)
		if err != nil {
			s.Log.Error("Something went wrong on handling request to Accrual service",
				zap.String("error", err.Error()),
			)
			continue
		}

		if accrualResponse == nil {
			s.Log.Error("Something went wrong on handling request to Accrual service: Blank response")
			continue
		}

		if accrualResponse.Status == "REGISTERED" {
			continue
		}

		err = s.UnprocessedOrderService.UpdateAccrualData(ctx, accrualResponse.Order, accrualResponse.Accrual, accrualResponse.Status)
		if err != nil {
			s.Log.Error("Something went wrong on updating accrual data for order",
				zap.String("error", err.Error()),
			)
			continue
		}
		processedOrderCount++
	}
	s.Log.Info("Accrual data has been updated",
		zap.Int("processed_orders", processedOrderCount),
	)
	s.Log.Info("Process for updating accrual data has been finished")
}
