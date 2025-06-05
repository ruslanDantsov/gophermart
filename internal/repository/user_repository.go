package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ruslanDantsov/gophermart/internal/infrastructure/storage/postgre"
	"github.com/ruslanDantsov/gophermart/internal/model"
	"github.com/ruslanDantsov/gophermart/internal/repository/query"
)

type UserRepository struct {
	storage *postgre.PostgreStorage
}

func NewUserRepository(storage *postgre.PostgreStorage) *UserRepository {
	return &UserRepository{storage: storage}
}

func (repository *UserRepository) Save(ctx context.Context, userData model.UserData) error {

	_, err := repository.storage.Conn.Exec(ctx,
		query.InsertOrUpdateUserData,
		userData.Id,
		userData.Login,
		userData.Password,
		userData.CreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return fmt.Errorf("login not unique: %w", err)
			default:
				return fmt.Errorf("postgresql error when saving user data (code %s): %w", pgErr.Code, err)
			}
		} else {
			return fmt.Errorf("error on saving user data: %w", err)
		}
	}
	return nil
}
