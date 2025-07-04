package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/errs"
	"github.com/ruslanDantsov/gophermart/internal/infrastructure/storage/postgre"
	"github.com/ruslanDantsov/gophermart/internal/model/business"
	"github.com/ruslanDantsov/gophermart/internal/model/entity"
	"github.com/ruslanDantsov/gophermart/internal/repository/query"
)

type WithdrawnRepository struct {
	storage *postgre.PostgreStorage
}

func NewWithdrawnRepository(storage *postgre.PostgreStorage) *WithdrawnRepository {
	return &WithdrawnRepository{storage: storage}
}

func (r *WithdrawnRepository) GetTotalWithdrawByUser(ctx context.Context, userID uuid.UUID) (float64, error) {
	db := r.storage.GetExecutor(ctx)

	var totalWithdrawn float64
	err := db.QueryRow(ctx,
		query.GetTotalWithdrawByUser,
		userID).
		Scan(&totalWithdrawn)

	if err != nil {
		return 0, errs.New(errs.Generic, "failed to execute query ", err)
	}

	return totalWithdrawn, nil
}

func (r *WithdrawnRepository) Save(ctx context.Context, withdraw entity.Withdraw) (*entity.Withdraw, error) {
	db := r.storage.GetExecutor(ctx)

	_, err := db.Exec(ctx,
		query.InsertWithdraw,
		withdraw.ID,
		withdraw.Sum,
		withdraw.CreatedAt,
		withdraw.OrderID)

	if err != nil {
		return nil, errs.New(errs.Generic, "failed to execute query ", err)
	}

	return &withdraw, nil
}

func (r *WithdrawnRepository) GetAllWithdrawDetailsByUser(ctx context.Context, userID uuid.UUID) ([]business.WithdrawDetail, error) {
	db := r.storage.GetExecutor(ctx)

	var withdraws []business.WithdrawDetail

	rows, err := db.Query(ctx, query.GetAllWithdrawDetailsByUser, userID)
	if err != nil {
		return nil, errs.New(errs.Generic, "failed to execute query ", err)
	}

	defer rows.Close()

	for rows.Next() {
		var withdraw business.WithdrawDetail
		err := rows.Scan(
			&withdraw.OrderNumber,
			&withdraw.Sum,
			&withdraw.CreatedAt,
		)
		if err != nil {
			return nil, errs.New(errs.Generic, "failed to scan withdraws ", err)
		}
		withdraws = append(withdraws, withdraw)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.New(errs.Generic, "rows iteration error ", err)
	}

	return withdraws, nil
}
