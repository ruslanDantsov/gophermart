package postgre

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type ctxTxKey struct{}

func ContextWithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, ctxTxKey{}, tx)
}

func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(ctxTxKey{}).(pgx.Tx)
	return tx, ok
}

type DBExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

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

func (s *PostgreStorage) GetExecutor(ctx context.Context) DBExecutor {
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}
	return s.Conn
}

func (s *PostgreStorage) WithTx(ctx context.Context, fn func(ctx context.Context, db DBExecutor) error) error {
	if existingTx, ok := TxFromContext(ctx); ok {
		return fn(ctx, existingTx)
	}

	tx, err := s.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	ctx = ContextWithTx(ctx, tx)

	err = fn(ctx, tx)

	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			s.Log.Error("rollback failed", zap.Error(rollbackErr))
		}
		return err
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		return fmt.Errorf("commit tx: %w", commitErr)
	}

	return nil
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
