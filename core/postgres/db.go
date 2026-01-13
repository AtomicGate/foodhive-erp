package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var ErrConnectionFailed = errors.New("connection failed")
var ErrPingFailed = errors.New("ping failed")

type (
	Rows       = pgx.Rows
	CommandTag = pgconn.CommandTag
	Executor   interface {
		Query(ctx context.Context, query string, args ...interface{}) Rows
		QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
		Exec(ctx context.Context, query string, args ...interface{}) (CommandTag, error)
	}
	Transaction interface {
		Executor
		Commit(ctx context.Context) error
		Rollback(ctx context.Context) error
	}
	Connection interface {
		Executor
		BeginTx(ctx context.Context) (Transaction, error)
	}
	impl struct {
		pool *pgxpool.Pool
	}
	txImpl struct {
		tx pgx.Tx
	}
)

func (i *impl) Query(ctx context.Context, query string, args ...interface{}) Rows {
	// In case of error, the rows returned will be errRows which can be checked with rows.Err()
	rows, _ := i.pool.Query(ctx, query, args...)
	return rows
}

func (i *impl) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	row := i.pool.QueryRow(ctx, query, args...)
	return row
}

func (i *impl) Exec(ctx context.Context, query string, args ...interface{}) (CommandTag, error) {
	return i.pool.Exec(ctx, query, args...)
}

func (i *impl) BeginTx(ctx context.Context) (Transaction, error) {
	tx, err := i.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &txImpl{tx: tx}, nil
}

func (t *txImpl) Query(ctx context.Context, query string, args ...interface{}) Rows {
	rows, _ := t.tx.Query(ctx, query, args...)
	return rows
}

func (t *txImpl) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return t.tx.QueryRow(ctx, query, args...)
}

func (t *txImpl) Exec(ctx context.Context, query string, args ...interface{}) (CommandTag, error) {
	return t.tx.Exec(ctx, query, args...)
}

func (t *txImpl) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *txImpl) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func New(config string) (Connection, error) {
	// Use connection pool instead of single connection
	pool, err := pgxpool.New(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConnectionFailed, err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPingFailed, err)
	}

	return &impl{
		pool: pool,
	}, nil
}
