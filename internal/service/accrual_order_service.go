package service

import (
	"context"
	"github.com/ruslanDantsov/gophermart/internal/dto/view"
	"go.uber.org/zap"
)

type UnprocessedOrderService interface {
	GetUnprocessedOrders(ctx context.Context) ([]string, error)
	UpdateAccrualData(ctx context.Context, number string, accrual float64, status string) error
}

type AccrualClient interface {
	GetAccrualData(ctx context.Context, orderID string) (*view.AccrualResponse, error)
}

type AccrualOrderService struct {
	unprocessedOrderService UnprocessedOrderService
	accrualClient           AccrualClient
	log                     *zap.Logger
}

func NewAccrualOrderService(unprocessedOrderService UnprocessedOrderService, accrualClient AccrualClient, log *zap.Logger) *AccrualOrderService {
	return &AccrualOrderService{
		unprocessedOrderService: unprocessedOrderService,
		accrualClient:           accrualClient,
		log:                     log,
	}
}

func (s *AccrualOrderService) ProcessOrders(ctx context.Context) {
	s.log.Info("Starting process for updating accrual data ...")

	unprocessedOrderNumbers, err := s.unprocessedOrderService.GetUnprocessedOrders(ctx)
	if err != nil {
		s.log.Error("Something went wrong on getting orders",
			zap.String("error", err.Error()),
		)
	}

	processedOrderCount := 0
	for _, orderNumber := range unprocessedOrderNumbers {
		accrualResponse, err := s.accrualClient.GetAccrualData(ctx, orderNumber)
		if err != nil {
			s.log.Error("Something went wrong on handling request to Accrual service",
				zap.String("error", err.Error()),
			)
			continue
		}

		if accrualResponse == nil {
			s.log.Error("Something went wrong on handling request to Accrual service: Blank response")
			continue
		}

		if accrualResponse.Status == view.AccrualOrderRegisteredStatus {
			continue
		}

		err = s.unprocessedOrderService.UpdateAccrualData(ctx, accrualResponse.Order, accrualResponse.Accrual, accrualResponse.Status)
		if err != nil {
			s.log.Error("Something went wrong on updating accrual data for order",
				zap.String("error", err.Error()),
			)
			continue
		}
		processedOrderCount++
	}
	s.log.Info("Accrual data has been updated",
		zap.Int("processed_orders", processedOrderCount),
	)

	s.log.Info("Process for updating accrual data has been finished")
}
