package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/errs"
	"github.com/ruslanDantsov/gophermart/internal/handler/middleware"
	"github.com/ruslanDantsov/gophermart/internal/model/business"
	"github.com/ruslanDantsov/gophermart/internal/model/entity"
	"time"
)

type IWithdrawRepository interface {
	Save(ctx context.Context, withdraw entity.Withdraw) (*entity.Withdraw, error)
	GetAllWithdrawDetailsByUser(ctx context.Context, userID uuid.UUID) ([]business.WithdrawDetail, error)
}

type IOrderCreatorService interface {
	AddOrder(ctx context.Context, orderCreateCommand command.OrderCreateCommand) (*entity.Order, error)
}

type WithdrawService struct {
	OrderCreatorService         IOrderCreatorService
	WithdrawRepository          IWithdrawRepository
	AccrualAggregatorRepository AccrualAggregatorRepository
}

func NewWithdrawService(orderCreatorService IOrderCreatorService, withdrawRepository IWithdrawRepository, accrualAggregatorRepository AccrualAggregatorRepository) *WithdrawService {
	return &WithdrawService{
		OrderCreatorService:         orderCreatorService,
		WithdrawRepository:          withdrawRepository,
		AccrualAggregatorRepository: accrualAggregatorRepository,
	}
}

func (s *WithdrawService) AddWithdraw(ctx context.Context, withdrawCreateCommand command.WithdrawCreateCommand, authUserID uuid.UUID) (*entity.Withdraw, error) {

	totalAccrual, err := s.AccrualAggregatorRepository.GetTotalAccrualByUser(ctx, authUserID)
	if withdrawCreateCommand.Sum > totalAccrual {
		return nil, errs.New(errs.NotEnoughAccrual, errs.NotEnoughAccrual, err)
	}
	if err != nil {
		return nil, err
	}

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

func (s *WithdrawService) GetWithdrawDetails(ctx context.Context) ([]business.WithdrawDetail, error) {
	userID := ctx.Value(middleware.CtxUserIDKey{}).(uuid.UUID)
	withdraws, err := s.WithdrawRepository.GetAllWithdrawDetailsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return withdraws, nil
}
