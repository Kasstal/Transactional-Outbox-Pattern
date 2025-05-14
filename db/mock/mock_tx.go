package mockdb

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// MockTx is a complete mock implementation of pgx.Tx for testing
type MockTx struct {
	CommitFunc   func(ctx context.Context) error
	RollbackFunc func(ctx context.Context) error
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	return m, nil
}

func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}

func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	return &MockBatchResults{}
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	return pgx.LargeObjects{}
}

func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}

func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error) {
	return pgconn.NewCommandTag(""), nil
}

func (m *MockTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return &MockRows{}, nil
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return &MockRow{}
}

func (m *MockTx) Conn() *pgx.Conn {
	return &pgx.Conn{}
}

func (m *MockTx) Commit(ctx context.Context) error {
	if m.CommitFunc != nil {
		return m.CommitFunc(ctx)
	}
	return nil
}

func (m *MockTx) Rollback(ctx context.Context) error {
	if m.RollbackFunc != nil {
		return m.RollbackFunc(ctx)
	}
	return nil
}

// MockBatchResults implements pgx.BatchResults for testing
type MockBatchResults struct{}

func (m *MockBatchResults) Exec() (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}

func (m *MockBatchResults) Query() (pgx.Rows, error) {
	return &MockRows{}, nil
}

func (m *MockBatchResults) QueryRow() pgx.Row {
	return &MockRow{}
}

func (m *MockBatchResults) Close() error {
	return nil
}

// MockRows implements pgx.Rows for testing
type MockRows struct {
	closeErr error
}

func (m *MockRows) Close() {
}

func (m *MockRows) Err() error {
	return nil
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("")
}

func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (m *MockRows) Next() bool {
	return false
}

func (m *MockRows) Scan(dest ...any) error {
	return nil
}

func (m *MockRows) Values() ([]any, error) {
	return nil, nil
}

func (m *MockRows) RawValues() [][]byte {
	return nil
}

func (m *MockRows) Conn() *pgx.Conn {
	return nil
}

// MockRow implements pgx.Row for testing
type MockRow struct {
	scanErr error
}

func (m *MockRow) Scan(dest ...any) error {
	if m.scanErr != nil {
		return m.scanErr
	}
	return nil
}
