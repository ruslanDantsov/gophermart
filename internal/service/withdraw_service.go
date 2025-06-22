package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/handler/middleware"
	"github.com/ruslanDantsov/gophermart/internal/model/entity"
	"time"
)

type IWithdrawRepository interface {
	Save(ctx context.Context, withdraw entity.Withdraw) (*entity.Withdraw, error)
	GetAllByUser(ctx context.Context, userID uuid.UUID) ([]entity.Withdraw, error)
}

type IOrderCreatorService interface {
	AddOrder(ctx context.Context, orderCreateCommand command.OrderCreateCommand) (*entity.Order, error)
}
type WithdrawService struct {
	OrderCreatorService IOrderCreatorService
	WithdrawRepository  IWithdrawRepository
}

func NewWithdrawService(orderCreatorService IOrderCreatorService, withdrawRepository IWithdrawRepository) *WithdrawService {
	return &WithdrawService{
		OrderCreatorService: orderCreatorService,
		WithdrawRepository:  withdrawRepository,
	}
}

func (s *WithdrawService) AddWithdraw(ctx context.Context, withdrawCreateCommand command.WithdrawCreateCommand) (*entity.Withdraw, error) {
	orderCreateCommand := command.OrderCreateCommand{Number: withdrawCreateCommand.Order}

	order, err := s.OrderCreatorService.AddOrder(ctx, orderCreateCommand)

	if err != nil {
		return nil, err
	}

	rawWithdraw := entity.Withdraw{
		ID:        uuid.New(),
		OrderID:   order.ID,
		CreatedAt: time.Now(),
		Sum:       withdrawCreateCommand.Sum,
	}

	withdraw, err := s.WithdrawRepository.Save(ctx, rawWithdraw)

	if err != nil {
		return nil, err
	}

	return withdraw, nil
}

func (s *WithdrawService) GetWithdraws(ctx context.Context) ([]entity.Withdraw, error) {
	userID := ctx.Value(middleware.CtxUserIDKey{}).(uuid.UUID)
	withdraws, err := s.WithdrawRepository.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return withdraws, nil
}
