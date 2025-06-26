package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/model/business"
)

type AccrualAggregatorRepository interface {
	GetTotalAccrualByUser(ctx context.Context, userID uuid.UUID) (float64, error)
}

type WithdrawnAggregatorRepository interface {
	GetTotalWithdrawByUser(ctx context.Context, userID uuid.UUID) (float64, error)
}
type BalanceService struct {
	accrualAggregatorRepository   AccrualAggregatorRepository
	withdrawnAggregatorRepository WithdrawnAggregatorRepository
}

func NewBalanceService(accrualAggregatorRepository AccrualAggregatorRepository, withdrawnAggregatorRepository WithdrawnAggregatorRepository) *BalanceService {
	return &BalanceService{
		accrualAggregatorRepository:   accrualAggregatorRepository,
		withdrawnAggregatorRepository: withdrawnAggregatorRepository,
	}
}

func (s *BalanceService) GetBalance(ctx context.Context, userID uuid.UUID) (*business.Balance, error) {
	totalAccrual, err := s.accrualAggregatorRepository.GetTotalAccrualByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	totalWithdrawn, err := s.withdrawnAggregatorRepository.GetTotalWithdrawByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &business.Balance{
		Total:     totalAccrual - totalWithdrawn,
		Withdrawn: totalWithdrawn,
	}, nil
}
