package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/infrastructure/storage/postgre"
)

type WithdrawnRepository struct {
	storage *postgre.PostgreStorage
}

func NewWithdrawnRepository(storage *postgre.PostgreStorage) *WithdrawnRepository {
	return &WithdrawnRepository{storage: storage}
}

func (r *WithdrawnRepository) GetTotalWithdrawnByUser(ctx context.Context, userID uuid.UUID) (float64, error) {
	//TODO: Implement logic
	return 50, nil
}
