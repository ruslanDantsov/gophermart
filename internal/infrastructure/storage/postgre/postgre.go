package postgre

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type PostgreStorage struct {
	Conn *pgxpool.Pool
	Log  zap.Logger
}

func NewPostgreStorage(ctx context.Context, log *zap.Logger, connectionString string) (*PostgreStorage, error) {
	if err := applyMigrations(connectionString); err != nil {
		return nil, err
	}

	conn, err := pgxpool.New(ctx, connectionString)

	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	return &PostgreStorage{
		Conn: conn,
		Log:  *log,
	}, nil
}

func applyMigrations(connectionString string) error {
	sqlDB, err := sql.Open("pgx", connectionString)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(sqlDB, "internal/infrastructure/db/migrations"); err != nil {
		return fmt.Errorf("unable to apply migrations: %w", err)
	}

	return nil
}
