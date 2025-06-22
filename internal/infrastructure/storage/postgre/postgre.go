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

//func (r *OrderRepository) getExecer(ctx context.Context) (postgre.Execer, error) {
//	tx := r.Storage.GetTx(ctx)
//
//	var execer postgre.Execer
//	if tx != nil {
//		execer = tx
//	} else {
//		conn, err := r.Storage.Conn.Acquire(ctx)
//		if err != nil {
//			return nil, err
//		}
//		defer conn.Release()
//		execer = conn
//	}
//	return execer, nil
//
//}

type Execer interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
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

type ctxTxKey struct{}

func (ps *PostgreStorage) GetTx(ctx context.Context) pgx.Tx {
	tx, _ := ctx.Value(ctxTxKey{}).(pgx.Tx)
	return tx
}

func (ps *PostgreStorage) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := ps.Conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to start tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	txCtx := context.WithValue(ctx, ctxTxKey{}, tx)

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
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
