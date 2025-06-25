package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/errs"
	"github.com/ruslanDantsov/gophermart/internal/handler/middleware"
	"github.com/ruslanDantsov/gophermart/internal/infrastructure/storage/postgre"
	"github.com/ruslanDantsov/gophermart/internal/model/entity"
	"github.com/ruslanDantsov/gophermart/internal/repository/query"
)

type OrderRepository struct {
	storage *postgre.PostgreStorage
}

func NewOrderRepository(storage *postgre.PostgreStorage) *OrderRepository {
	return &OrderRepository{storage: storage}
}

func (r *OrderRepository) Save(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	currentUserID := ctx.Value(middleware.CtxUserIDKey{})
	tx, err := r.storage.Conn.Begin(ctx)
	if err != nil {
		return nil, errs.New(errs.Generic, "failed to begin transaction ", err)
	}
	defer tx.Rollback(ctx)

	var existingUserID uuid.UUID
	err = tx.QueryRow(ctx,
		query.FindUserByOrderNumber,
		order.Number,
	).Scan(&existingUserID)

	switch {
	case errors.Is(err, sql.ErrNoRows):
	case err != nil:
		return nil, errs.New(errs.Generic, "failed to execute query", err)
	case existingUserID == currentUserID:
		return nil, errs.New(errs.OrderAddedByCurrentUser, "order already added by current user", err)
	default:
		return nil, errs.New(errs.OrderAddedByAnotherUser, "order already added by another user", err)
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
		return nil, errs.New(errs.Generic, "failed to execute query ", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errs.New(errs.Generic, "failed to commit transaction ", err)
	}

	return order, nil
}

func (r *OrderRepository) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	var orders []entity.Order

	rows, err := r.storage.Conn.Query(ctx, query.GetAllOrdersByUser, userID)

	if err != nil {
		return nil, errs.New(errs.Generic, "failed to execute query ", err)
	}

	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		err := rows.Scan(
			&order.ID,
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.CreatedAt,
			&order.UserID,
		)
		if err != nil {
			return nil, errs.New(errs.Generic, "failed to scan order ", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.New(errs.Generic, "rows iteration error ", err)
	}

	return orders, nil
}

func (r *OrderRepository) GetUnprocessedOrders(ctx context.Context) ([]string, error) {
	var numbers []string

	rows, err := r.storage.Conn.Query(ctx, query.GetUnprocessedOrderNumbers)

	if err != nil {
		return nil, errs.New(errs.Generic, "failed to execute query ", err)
	}

	defer rows.Close()

	for rows.Next() {
		var number string
		err := rows.Scan(
			&number,
		)
		if err != nil {
			return nil, errs.New(errs.Generic, "failed to scan orders ", err)
		}
		numbers = append(numbers, number)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.New(errs.Generic, "rows iteration error ", err)
	}

	return numbers, nil

}

func (r *OrderRepository) GetTotalAccrualByUser(ctx context.Context, userID uuid.UUID) (float64, error) {
	var totalAccrual float64
	err := r.storage.Conn.QueryRow(ctx,
		query.GetTotalAccrualByUser,
		userID).
		Scan(&totalAccrual)

	if err != nil {
		return 0, errs.New(errs.Generic, "failed to execute query ", err)
	}

	return totalAccrual, nil
}

func (r *OrderRepository) UpdateAccrualData(ctx context.Context, number string, accrual float64, status string) error {
	_, err := r.storage.Conn.Exec(ctx,
		query.UpdateAccrualData,
		status,
		accrual,
		number)
	if err != nil {
		return errs.New(errs.Generic, "failed to execute query ", err)
	}

	return nil
}
