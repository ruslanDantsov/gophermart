package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/model"
	"time"
)

type IOrderRepository interface {
	Save(ctx context.Context, order *model.Order) (*model.Order, error)
	GetAllByUserId(ctx context.Context, userId uuid.UUID) ([]model.Order, error)
}
type OrderService struct {
	OrderRepository IOrderRepository
}

func NewOrderService(orderRepository IOrderRepository) *OrderService {
	return &OrderService{
		OrderRepository: orderRepository,
	}
}

func (s *OrderService) AddOrder(ctx context.Context, orderCreateCommand command.OrderCreateCommand) (*model.Order, error) {
	userId := ctx.Value("userId").(uuid.UUID)

	rawOrder := &model.Order{
		ID:        uuid.New(),
		Number:    orderCreateCommand.Number,
		Status:    model.ORDER_NEW_STATUS,
		Accrual:   0,
		CreatedAt: time.Now(),
		UserID:    userId,
	}

	if _, err := s.OrderRepository.Save(ctx, rawOrder); err != nil {
		return nil, err
	}

	return rawOrder, nil
}

func (s *OrderService) GetOrders(ctx context.Context) ([]model.Order, error) {
	userId := ctx.Value("userId").(uuid.UUID)
	orders, err := s.OrderRepository.GetAllByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
