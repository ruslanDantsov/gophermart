package service

import (
	"context"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/errs"
	"github.com/ruslanDantsov/gophermart/internal/model/entity"
	"time"
)

type IOrderRepository interface {
	Save(ctx context.Context, order *entity.Order) (*entity.Order, error)
	GetAllByUser(ctx context.Context, userId uuid.UUID) ([]entity.Order, error)
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
		return nil, errs.New(errs.INVALID_ORDER_NUMBER, "invalid order number", err)
	}

	userId := ctx.Value("userId").(uuid.UUID)

	rawOrder := &entity.Order{
		ID:        uuid.New(),
		Number:    orderCreateCommand.Number,
		Status:    entity.ORDER_NEW_STATUS,
		Accrual:   0,
		CreatedAt: time.Now(),
		UserID:    userId,
	}

	if _, err := s.OrderRepository.Save(ctx, rawOrder); err != nil {
		return nil, err
	}

	return rawOrder, nil
}

func (s *OrderService) GetOrders(ctx context.Context) ([]entity.Order, error) {
	userId := ctx.Value("userId").(uuid.UUID)
	orders, err := s.OrderRepository.GetAllByUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
