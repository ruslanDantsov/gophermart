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

type OrderRepository interface {
	Save(ctx context.Context, order *entity.Order) (*entity.Order, error)
	GetAllByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
	GetUnprocessedOrders(ctx context.Context) ([]string, error)
	UpdateAccrualData(ctx context.Context, number string, accrual float64, status string) error
	FindUserIDByOrderNumber(ctx context.Context, orderNumber string) (uuid.UUID, error)
}
type OrderService struct {
	storage         *postgre.PostgreStorage
	orderRepository OrderRepository
}

func NewOrderService(orderRepository OrderRepository, storage *postgre.PostgreStorage) *OrderService {
	return &OrderService{
		orderRepository: orderRepository,
		storage:         storage,
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
	err := s.storage.WithTx(ctx, func(ctx context.Context) error {
		var existingUserID uuid.UUID
		existingUserID, err := s.orderRepository.FindUserIDByOrderNumber(ctx, rawOrder.Number)
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

		if _, err := s.orderRepository.Save(ctx, rawOrder); err != nil {
			return err
		}

		savedOrder = rawOrder
		return nil
	})

	return savedOrder, err
}

func (s *OrderService) GetOrders(ctx context.Context) ([]entity.Order, error) {
	userID := ctx.Value(middleware.CtxUserIDKey{}).(uuid.UUID)
	orders, err := s.orderRepository.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderService) GetUnprocessedOrders(ctx context.Context) ([]string, error) {
	numbers, err := s.orderRepository.GetUnprocessedOrders(ctx)
	if err != nil {
		return nil, err
	}

	return numbers, nil
}

func (s *OrderService) UpdateAccrualData(ctx context.Context, number string, accrual float64, status string) error {
	return s.orderRepository.UpdateAccrualData(ctx, number, accrual, status)
}
