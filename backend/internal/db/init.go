package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type ctxKey string

const connPoolKey ctxKey = "pgx_conn_pool"
const txKey ctxKey = "pgx_tx"

func InitDBPool() (*pgxpool.Pool, func(), error) {
	pool, err := pgxpool.New(context.Background(), os.Getenv("COCKROACH_CONNECTION_STRING"))
	if err != nil {
		return nil, nil, err
	}

	return pool, func() {
		pool.Close()
	}, nil
}

func WithConnPool(ctx context.Context, pool *pgxpool.Pool) context.Context {
	return context.WithValue(ctx, connPoolKey, pool)
}

func WithTx(ctx context.Context, tx *pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetConn(ctx context.Context) (DBExecutor, bool) {
	if tx, ok := ctx.Value(txKey).(*pgx.Tx); ok {
		return *tx, true
	}
	if pool, ok := ctx.Value(connPoolKey).(*pgxpool.Pool); ok {
		return pool, true
	}
	return nil, false
}
