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
	var totalWithdrawn float64
	err := r.storage.Conn.QueryRow(ctx,
		query.GetTotalWithdrawByUser,
		userID).
		Scan(&totalWithdrawn)

	if err != nil {
		return 0, errs.New(errs.Generic, "failed to execute query ", err)
	}

	return totalWithdrawn, nil
}

func (r *WithdrawnRepository) Save(ctx context.Context, withdraw entity.Withdraw) (*entity.Withdraw, error) {
	_, err := r.storage.Conn.Exec(ctx,
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
	var withdraws []business.WithdrawDetail

	rows, err := r.storage.Conn.Query(ctx, query.GetAllWithdrawDetailsByUser, userID)

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
