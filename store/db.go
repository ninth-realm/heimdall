package store

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	TxBeginner
	UserRepository
	PasswordRepository
}

type TxBeginner interface {
	BeginTx(xtx context.Context) (*sqlx.Tx, error)
}

type WorkFn[T any] func(*sqlx.Tx) (T, error)

func RunUnitOfWork[T any](
	ctx context.Context,
	repo TxBeginner,
	fn WorkFn[T],
) (T, error) {
	txn, err := repo.BeginTx(ctx)
	if err != nil {
		return *new(T), err
	}

	val, err := fn(txn)
	if err != nil {
		txn.Rollback()
		return *new(T), err
	}

	txn.Commit()

	return val, nil
}

type QueryOptions struct {
	Ctx context.Context
	Txn *sqlx.Tx
}

func (q *QueryOptions) Context() context.Context {
	if q.Ctx == nil {
		return context.Background()
	}

	return q.Ctx
}
