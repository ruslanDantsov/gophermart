package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/ruslanDantsov/gophermart/internal/errs"
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

func (r *OrderRepository) FindUserIDByOrderNumber(ctx context.Context, orderNumber string) (uuid.UUID, error) {
	db := r.storage.GetExecutor(ctx)

	var userID uuid.UUID
	err := db.QueryRow(ctx,
		query.FindUserByOrderNumber,
		orderNumber,
	).Scan(&userID)

	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, nil
	}

	if err != nil {
		return uuid.Nil, errs.New(errs.Generic, "failed to execute query", err)
	}
	return userID, nil
}

func (r *OrderRepository) Save(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	db := r.storage.GetExecutor(ctx)

	_, err := db.Exec(ctx,
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

	return order, nil
}

//r.storage.WithTx(ctx, func(ctx context.Context, db postgre.DBExecutor) error {

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
