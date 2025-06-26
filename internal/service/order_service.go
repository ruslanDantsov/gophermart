package service

import (
	"context"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/errs"
	"github.com/ruslanDantsov/gophermart/internal/handler/middleware"
	"github.com/ruslanDantsov/gophermart/internal/infrastructure/storage/postgre"
	"github.com/ruslanDantsov/gophermart/internal/model/entity"
	"time"
)

type IOrderRepository interface {
	Save(ctx context.Context, order *entity.Order) (*entity.Order, error)
	GetAllByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
	GetUnprocessedOrders(ctx context.Context) ([]string, error)
	UpdateAccrualData(ctx context.Context, number string, accrual float64, status string) error
	FindUserIDByOrderNumber(ctx context.Context, orderNumber string) (uuid.UUID, error)
}
type OrderService struct {
	Storage         *postgre.PostgreStorage
	OrderRepository IOrderRepository
}

func NewOrderService(orderRepository IOrderRepository, storage *postgre.PostgreStorage) *OrderService {
	return &OrderService{
		OrderRepository: orderRepository,
		Storage:         storage,
	}
}

func (s *OrderService) AddOrder(ctx context.Context, orderCreateCommand command.OrderCreateCommand) (*entity.Order, error) {
	if err := goluhn.Validate(orderCreateCommand.Number); err != nil {
		return nil, errs.New(errs.InvalidOrderNumber, "invalid order number", err)
	}

	authUserID := ctx.Value(middleware.CtxUserIDKey{}).(uuid.UUID)

	rawOrder := &entity.Order{
		ID:        uuid.New(),
		Number:    orderCreateCommand.Number,
		Status:    entity.OrderNewStatus,
		Accrual:   0,
		CreatedAt: time.Now(),
		UserID:    authUserID,
	}

	var savedOrder *entity.Order
	err := s.Storage.WithTx(ctx, func(ctx context.Context) error {
		var existingUserID uuid.UUID
		existingUserID, err := s.OrderRepository.FindUserIDByOrderNumber(ctx, rawOrder.Number)
		if err != nil {
			return err
		}

		switch {
		case existingUserID == uuid.Nil:
		case existingUserID == authUserID:
			return errs.New(errs.OrderAddedByCurrentUser, "order already added by current user", nil)
		default:
			return errs.New(errs.OrderAddedByAnotherUser, "order already added by another user", nil)
		}

		if _, err := s.OrderRepository.Save(ctx, rawOrder); err != nil {
			return err
		}

		savedOrder = rawOrder
		return nil
	})

	return savedOrder, err
}

func (s *OrderService) GetOrders(ctx context.Context) ([]entity.Order, error) {
	userID := ctx.Value(middleware.CtxUserIDKey{}).(uuid.UUID)
	orders, err := s.OrderRepository.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderService) GetUnprocessedOrders(ctx context.Context) ([]string, error) {
	numbers, err := s.OrderRepository.GetUnprocessedOrders(ctx)
	if err != nil {
		return nil, err
	}

	return numbers, nil
}

func (s *OrderService) UpdateAccrualData(ctx context.Context, number string, accrual float64, status string) error {
	return s.OrderRepository.UpdateAccrualData(ctx, number, accrual, status)
}
