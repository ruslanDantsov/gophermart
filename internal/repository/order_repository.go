package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/errs"
	"github.com/ruslanDantsov/gophermart/internal/infrastructure/storage/postgre"
	"github.com/ruslanDantsov/gophermart/internal/model"
	"github.com/ruslanDantsov/gophermart/internal/repository/query"
)

type OrderRepository struct {
	storage *postgre.PostgreStorage
}

func NewOrderRepository(storage *postgre.PostgreStorage) *OrderRepository {
	return &OrderRepository{storage: storage}
}

func (r *OrderRepository) Save(ctx context.Context, order *model.Order) (*model.Order, error) {
	currentUserId := ctx.Value("userId")
	tx, err := r.storage.Conn.Begin(ctx)
	if err != nil {
		return nil, errs.New(errs.GENERIC, "failed to begin transaction ", err)
	}
	defer tx.Rollback(ctx)

	var existingUserID uuid.UUID
	err = tx.QueryRow(ctx,
		query.FindUserByOrderNumber,
		order.Number,
	).Scan(&existingUserID)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		// Order number is not yet used — proceed to insert
	case err != nil:
		return nil, errs.New(errs.GENERIC, "failed to execute query", err)
	case existingUserID == currentUserId:
		return nil, errs.New(errs.ORDER_ADDED_BY_CURRENT_USER, "order already added by current user", err)
	default:
		return nil, errs.New(errs.ORDER_ADDED_BY_ANOTHER_USER, "order already added by another user", err)
	}

	_, err = tx.Exec(ctx,
		query.InsertOrder,
		order.ID,
		order.Number,
		order.Status,
		order.Accrual,
		order.CreatedAt,
		order.UserID)

	if err != nil {
		return nil, errs.New(errs.GENERIC, "failed to execute query ", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errs.New(errs.GENERIC, "failed to commit transaction ", err)
	}

	return order, nil
}

func (r *OrderRepository) GetAllByUserId(ctx context.Context, userId uuid.UUID) ([]model.Order, error) {
	var orders []model.Order

	rows, err := r.storage.Conn.Query(ctx, query.GetAllOrdersByUser, userId)

	if err != nil {
		return nil, errs.New(errs.GENERIC, "failed to execute query ", err)
	}

	defer rows.Close()

	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.CreatedAt,
			&order.UserID,
		)
		if err != nil {
			return nil, errs.New(errs.GENERIC, "failed to scan order ", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.New(errs.GENERIC, "rows iteration error ", err)
	}

	return orders, nil
}
