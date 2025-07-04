package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ruslanDantsov/gophermart/internal/dto/view"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

type MockUnprocessedOrderService struct {
	mock.Mock
}

func (m *MockUnprocessedOrderService) GetUnprocessedOrders(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockUnprocessedOrderService) UpdateAccrualData(ctx context.Context, number string, accrual float64, status string) error {
	args := m.Called(ctx, number, accrual, status)
	return args.Error(0)
}

type MockAccrualClient struct {
	mock.Mock
}

func (m *MockAccrualClient) GetAccrualData(ctx context.Context, orderID string) (*view.AccrualResponse, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).(*view.AccrualResponse), args.Error(1)
}

func TestAccrualOrderService_ProcessOrders(t *testing.T) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	t.Run("successfully processes orders", func(t *testing.T) {
		mockUnprocessedService := new(MockUnprocessedOrderService)
		mockAccrualClient := new(MockAccrualClient)

		orderNumbers := []string{"123", "456"}
		mockUnprocessedService.On("GetUnprocessedOrders", ctx).Return(orderNumbers, nil)

		mockAccrualClient.On("GetAccrualData", ctx, "123").Return(&view.AccrualResponse{
			Order:   "123",
			Accrual: 100.0,
			Status:  "PROCESSED",
		}, nil)

		mockAccrualClient.On("GetAccrualData", ctx, "456").Return(&view.AccrualResponse{
			Order:   "456",
			Accrual: 50.0,
			Status:  "INVALID",
		}, nil)

		mockUnprocessedService.On("UpdateAccrualData", ctx, "123", 100.0, "PROCESSED").Return(nil)
		mockUnprocessedService.On("UpdateAccrualData", ctx, "456", 50.0, "INVALID").Return(nil)

		svc := NewAccrualOrderService(mockUnprocessedService, mockAccrualClient, logger)
		svc.ProcessOrders(ctx)

		mockUnprocessedService.AssertExpectations(t)
		mockAccrualClient.AssertExpectations(t)
	})

	t.Run("handles GetAccrualData error", func(t *testing.T) {
		mockUnprocessedService := new(MockUnprocessedOrderService)
		mockAccrualClient := new(MockAccrualClient)

		orderNumbers := []string{"123"}
		mockUnprocessedService.On("GetUnprocessedOrders", ctx).Return(orderNumbers, nil)
		mockAccrualClient.On("GetAccrualData", ctx, "123").Return((*view.AccrualResponse)(nil), errors.New("accrual service error"))

		svc := NewAccrualOrderService(mockUnprocessedService, mockAccrualClient, logger)
		svc.ProcessOrders(ctx)

		mockUnprocessedService.AssertExpectations(t)
		mockAccrualClient.AssertExpectations(t)
	})

	t.Run("skips REGISTERED orders", func(t *testing.T) {
		mockUnprocessedService := new(MockUnprocessedOrderService)
		mockAccrualClient := new(MockAccrualClient)

		orderNumbers := []string{"123"}
		mockUnprocessedService.On("GetUnprocessedOrders", ctx).Return(orderNumbers, nil)
		mockAccrualClient.On("GetAccrualData", ctx, "123").Return(&view.AccrualResponse{
			Order:   "123",
			Accrual: 0,
			Status:  "REGISTERED",
		}, nil)

		svc := NewAccrualOrderService(mockUnprocessedService, mockAccrualClient, logger)
		svc.ProcessOrders(ctx)

		mockUnprocessedService.AssertExpectations(t)
		mockAccrualClient.AssertExpectations(t)
	})

	t.Run("handles UpdateAccrualData error", func(t *testing.T) {
		mockUnprocessedService := new(MockUnprocessedOrderService)
		mockAccrualClient := new(MockAccrualClient)

		orderNumbers := []string{"123"}
		mockUnprocessedService.On("GetUnprocessedOrders", ctx).Return(orderNumbers, nil)
		mockAccrualClient.On("GetAccrualData", ctx, "123").Return(&view.AccrualResponse{
			Order:   "123",
			Accrual: 100.0,
			Status:  "PROCESSED",
		}, nil)
		mockUnprocessedService.On("UpdateAccrualData", ctx, "123", 100.0, "PROCESSED").Return(errors.New("update error"))

		svc := NewAccrualOrderService(mockUnprocessedService, mockAccrualClient, logger)
		svc.ProcessOrders(ctx)

		mockUnprocessedService.AssertExpectations(t)
		mockAccrualClient.AssertExpectations(t)
	})
}
