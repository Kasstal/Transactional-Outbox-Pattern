package service

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	db "orders-center/db/sqlc"
)

type TransactionService struct {
	pool *pgxpool.Pool
}

func NewTransactionService(pool *pgxpool.Pool) *TransactionService {
	return &TransactionService{pool: pool}
}
func (t *TransactionService) ExecTx(ctx context.Context, fn func(tx *db.Queries) error) error {
	// Получаем соединение из пула
	conn, err := t.pool.Acquire(ctx)

	if err != nil {
		return err
	}
	// Обязательно освобождаем соединение после выполнения операции

	// Создаем транзакцию на этом соединении
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {

		return err
	}

	// Создаем запросы с транзакцией
	q := db.New(tx)

	// Выполняем функцию, переданную в ExecTx
	err = fn(q)
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			conn.Release()
		} else {
			tx.Commit(ctx)
			conn.Release()
		}
	}()
	return err
}

func (t *TransactionService) NewConnection() (*pgxpool.Conn, error) {
	conn, err := t.pool.Acquire(context.Background())
	return conn, err
}
