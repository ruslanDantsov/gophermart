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
	GetTotalWithdrawnByUser(ctx context.Context, userID uuid.UUID) (float64, error)
}
type BalanceService struct {
	AccrualAggregatorRepository   AccrualAggregatorRepository
	WithdrawnAggregatorRepository WithdrawnAggregatorRepository
}

func NewBalanceService(accrualAggregatorRepository AccrualAggregatorRepository, withdrawnAggregatorRepository WithdrawnAggregatorRepository) *BalanceService {
	return &BalanceService{
		AccrualAggregatorRepository:   accrualAggregatorRepository,
		WithdrawnAggregatorRepository: withdrawnAggregatorRepository,
	}
}

func (s *BalanceService) GetBalance(ctx context.Context, userID uuid.UUID) (*business.Balance, error) {
	totalAccrual, err := s.AccrualAggregatorRepository.GetTotalAccrualByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	totalWithdrawn, err := s.WithdrawnAggregatorRepository.GetTotalWithdrawnByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &business.Balance{
		Accrual:   totalAccrual,
		Withdrawn: totalWithdrawn,
	}, nil
}
