package service

import (
	"context"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/errs"
	"github.com/ruslanDantsov/gophermart/internal/handler/middleware"
	"github.com/ruslanDantsov/gophermart/internal/model/entity"
	"time"
)

type IOrderRepository interface {
	Save(ctx context.Context, order *entity.Order) (*entity.Order, error)
	GetAllByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
}
type OrderService struct {
	OrderRepository IOrderRepository
}

func NewOrderService(orderRepository IOrderRepository) *OrderService {
	return &OrderService{
		OrderRepository: orderRepository,
	}
}

func (s *OrderService) AddOrder(ctx context.Context, orderCreateCommand command.OrderCreateCommand) (*entity.Order, error) {
	if err := goluhn.Validate(orderCreateCommand.Number); err != nil {
		return nil, errs.New(errs.InvalidOrderNumber, "invalid order number", err)
	}

	userID := ctx.Value(middleware.CtxUserIDKey{}).(uuid.UUID)

	rawOrder := &entity.Order{
		ID:        uuid.New(),
		Number:    orderCreateCommand.Number,
		Status:    entity.OrderNewStatus,
		Accrual:   0,
		CreatedAt: time.Now(),
		UserID:    userID,
	}

	if _, err := s.OrderRepository.Save(ctx, rawOrder); err != nil {
		return nil, err
	}

	return rawOrder, nil
}

func (s *OrderService) GetOrders(ctx context.Context) ([]entity.Order, error) {
	userID := ctx.Value(middleware.CtxUserIDKey{}).(uuid.UUID)
	orders, err := s.OrderRepository.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
