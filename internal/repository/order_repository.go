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
	var savedOrder *entity.Order

	err := r.storage.WithTx(ctx, func(ctx context.Context, db postgre.DBExecutor) error {
		currentUserID := ctx.Value(middleware.CtxUserIDKey{})

		var existingUserID uuid.UUID
		err := db.QueryRow(ctx,
			query.FindUserByOrderNumber,
			order.Number,
		).Scan(&existingUserID)

		switch {
		case errors.Is(err, sql.ErrNoRows):
		case err != nil:
			return errs.New(errs.Generic, "failed to execute query", err)
		case existingUserID == currentUserID:
			return errs.New(errs.OrderAddedByCurrentUser, "order already added by current user", err)
		default:
			return errs.New(errs.OrderAddedByAnotherUser, "order already added by another user", err)
		}

		_, err = db.Exec(ctx,
			query.InsertOrder,
			order.ID,
			order.Number,
			order.Status,
			order.Accrual,
			order.CreatedAt,
			order.UserID)

		if err != nil {
			return errs.New(errs.Generic, "failed to execute query ", err)
		}

		savedOrder = order
		return nil
	})

	return savedOrder, err
}

func (r *OrderRepository) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	db := r.storage.GetExecutor(ctx)

	var orders []entity.Order

	rows, err := db.Query(ctx, query.GetAllOrdersByUser, userID)

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
	db := r.storage.GetExecutor(ctx)

	var numbers []string

	rows, err := db.Query(ctx, query.GetUnprocessedOrderNumbers)

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
	db := r.storage.GetExecutor(ctx)

	var totalAccrual float64
	err := db.QueryRow(ctx,
		query.GetTotalAccrualByUser,
		userID).
		Scan(&totalAccrual)

	if err != nil {
		return 0, errs.New(errs.Generic, "failed to execute query ", err)
	}

	return totalAccrual, nil
}

func (r *OrderRepository) UpdateAccrualData(ctx context.Context, number string, accrual float64, status string) error {
	db := r.storage.GetExecutor(ctx)

	_, err := db.Exec(ctx,
		query.UpdateAccrualData,
		status,
		accrual,
		number)
	if err != nil {
		return errs.New(errs.Generic, "failed to execute query ", err)
	}

	return nil
}
