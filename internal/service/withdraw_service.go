package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/dto/command"
	"github.com/ruslanDantsov/gophermart/internal/errs"
	"github.com/ruslanDantsov/gophermart/internal/handler/middleware"
	"github.com/ruslanDantsov/gophermart/internal/infrastructure/storage/postgre"
	"github.com/ruslanDantsov/gophermart/internal/model/business"
	"github.com/ruslanDantsov/gophermart/internal/model/entity"
	"time"
)

type WithdrawRepository interface {
	Save(ctx context.Context, withdraw entity.Withdraw) (*entity.Withdraw, error)
	GetAllWithdrawDetailsByUser(ctx context.Context, userID uuid.UUID) ([]business.WithdrawDetail, error)
}

type OrderCreator interface {
	AddOrder(ctx context.Context, orderCreateCommand command.OrderCreateCommand) (*entity.Order, error)
}

type WithdrawService struct {
	orderCreatorService         OrderCreator
	withdrawRepository          WithdrawRepository
	accrualAggregatorRepository AccrualAggregatorRepository
	storage                     *postgre.PostgreStorage
}

func NewWithdrawService(orderCreatorService OrderCreator, withdrawRepository WithdrawRepository, accrualAggregatorRepository AccrualAggregatorRepository, storage *postgre.PostgreStorage) *WithdrawService {
	return &WithdrawService{
		orderCreatorService:         orderCreatorService,
		withdrawRepository:          withdrawRepository,
		accrualAggregatorRepository: accrualAggregatorRepository,
		storage:                     storage,
	}
}

func (s *WithdrawService) AddWithdraw(ctx context.Context, withdrawCreateCommand command.WithdrawCreateCommand, authUserID uuid.UUID) (*entity.Withdraw, error) {
	var savedWithdraw *entity.Withdraw

	err := s.storage.WithTx(ctx, func(ctx context.Context) error {
		totalAccrual, err := s.accrualAggregatorRepository.GetTotalAccrualByUser(ctx, authUserID)
		if err != nil {
			return errs.New(errs.Generic, "failed to get total accrual", err)
		}

		if withdrawCreateCommand.Sum > totalAccrual {
			return errs.New(errs.NotEnoughAccrual, "not enough accrual", nil)
		}

		orderCreateCommand := command.OrderCreateCommand{Number: withdrawCreateCommand.Order}
		order, err := s.orderCreatorService.AddOrder(ctx, orderCreateCommand)
		if err != nil {
			return err
		}

		rawWithdraw := entity.Withdraw{
			ID:        uuid.New(),
			OrderID:   order.ID,
			CreatedAt: time.Now(),
			Sum:       withdrawCreateCommand.Sum,
		}

		savedWithdraw, err = s.withdrawRepository.Save(ctx, rawWithdraw)
		if err != nil {
			return err
		}

		return nil
	})

	return savedWithdraw, err
}

func (s *WithdrawService) GetWithdrawDetails(ctx context.Context) ([]business.WithdrawDetail, error) {
	userID := ctx.Value(middleware.CtxUserIDKey{}).(uuid.UUID)
	withdraws, err := s.withdrawRepository.GetAllWithdrawDetailsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return withdraws, nil
}
