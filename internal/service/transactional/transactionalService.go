package service

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transactional interface {
	ExecTx(context.Context, func(ctx context.Context) error) error
}

type TransactionService struct {
	pool *pgxpool.Pool
}

func NewTransactionService(pool *pgxpool.Pool) Transactional {
	return &TransactionService{pool: pool}
}
func (t *TransactionService) ExecTx(ctx context.Context, fn func(ctx context.Context) error) error {

	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {

		return err
	}

	/*defer func() {
		if err != nil {

		}
	}()*/
	txCtx := NewTxContext(ctx, tx)
	err = fn(txCtx)
	if err != nil {
		rbErr := tx.Rollback(ctx)
		if rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit(ctx)

}

type TxContextKey struct{}

func NewTxContext(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxContextKey{}, tx)
}

func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(TxContextKey{}).(pgx.Tx)
	return tx, ok
}
