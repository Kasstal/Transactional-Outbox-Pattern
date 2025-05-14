package mockdb

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MockPool struct{}

func (m *MockPool) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	return &pgxpool.Conn{}, nil
}

func (m *MockPool) Close() {}
