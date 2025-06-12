package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ruslanDantsov/gophermart/internal/infrastructure/storage/postgre"
	"github.com/ruslanDantsov/gophermart/internal/model"
	"github.com/ruslanDantsov/gophermart/internal/repository/query"
	"time"
)

type UserRepository struct {
	storage *postgre.PostgreStorage
}

func NewUserRepository(storage *postgre.PostgreStorage) *UserRepository {
	return &UserRepository{storage: storage}
}

func (r *UserRepository) Save(ctx context.Context, userData model.UserData) error {

	_, err := r.storage.Conn.Exec(ctx,
		query.InsertOrUpdateUserData,
		userData.ID,
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

func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*model.UserData, error) {
	var (
		existingID        uuid.UUID
		existingLogin     string
		existingPassword  string
		existingCreatedAt time.Time
	)

	err := r.storage.Conn.QueryRow(ctx,
		query.FindUserByLogin,
		login).
		Scan(&existingID, &existingLogin, &existingPassword, &existingCreatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.NoDataFound:
				return nil, fmt.Errorf("user not found: %w", err)
			default:
				return nil, fmt.Errorf("postgresql error on searching user data (code %s): %w", pgErr.Code, err)
			}
		} else {
			return nil, fmt.Errorf("error on on searching user data: %w", err)
		}
	}

	userData := &model.UserData{
		ID:        existingID,
		Login:     existingLogin,
		Password:  existingPassword,
		CreatedAt: existingCreatedAt,
	}

	return userData, nil
}
