package shorten

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

type mockRows struct {
	columns []string
	rows    [][]any
	pos     int
	closed  bool
}

func (m *mockRows) Columns() []string {
	return m.columns
}

func (m *mockRows) Next() bool {
	if m.pos < len(m.rows) {
		m.pos++
		return true
	}
	return false
}

func (m *mockRows) Scan(dest ...any) error {
	if m.pos == 0 || m.pos > len(m.rows) {
		return fmt.Errorf("no current row")
	}

	row := m.rows[m.pos-1]
	if len(dest) != len(row) {
		return fmt.Errorf("expected %d destinations, got %d", len(row), len(dest))
	}

	for i, value := range row {
		d := reflect.ValueOf(dest[i])
		if d.Kind() != reflect.Pointer || d.IsNil() {
			return fmt.Errorf("destination %d must be non-nil pointer", i)
		}
		elem := d.Elem()
		val := reflect.ValueOf(value)
		if !val.Type().AssignableTo(elem.Type()) {
			return fmt.Errorf("cannot assign %T to %s", value, elem.Type())
		}
		elem.Set(val)
	}
	return nil
}

func (m *mockRows) Close() error {
	m.closed = true
	return nil
}

type mockExec struct {
	queryFunc func(ctx context.Context, query string, args ...any) (Rows, error)
}

func (m *mockExec) Prepare(ctx context.Context, query string) (Stmt, error) {
	return nil, nil
}

func (m *mockExec) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	return 0, nil
}

func (m *mockExec) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	return m.queryFunc(ctx, query, args...)
}

func (m *mockExec) Release() error {
	return nil
}

type mockTx struct {
	commitCount   int
	rollbackCount int
	commitErr     error
	rollbackErr   error
}

func (m *mockTx) Prepare(ctx context.Context, query string) (Stmt, error) {
	return nil, nil
}

func (m *mockTx) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	return 0, nil
}

func (m *mockTx) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	return nil, nil
}

func (m *mockTx) Release() error {
	return nil
}

func (m *mockTx) Commit() error {
	m.commitCount++
	return m.commitErr
}

func (m *mockTx) Rollback() error {
	m.rollbackCount++
	return m.rollbackErr
}

type mockFactory struct {
	conn      Exec
	tx        Tx
	connCount int
	txCount   int
	connErr   error
	txErr     error
}

func (m *mockFactory) getConn(ctx context.Context) (Exec, error) {
	m.connCount++
	return m.conn, m.connErr
}

func (m *mockFactory) getTx(ctx context.Context, level sql.IsolationLevel) (Tx, error) {
	m.txCount++
	return m.tx, m.txErr
}

func TestGet_UsesTransactionScope(t *testing.T) {
	ctx, scope := Scope(context.Background())
	tx := &mockTx{}
	factory := &mockFactory{tx: tx}

	exec, err := Get(ctx, factory)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exec != tx {
		t.Fatalf("expected transaction from Get, got %T", exec)
	}
	if scope.tx != tx {
		t.Fatalf("expected transaction stored in scope")
	}
	if factory.txCount != 1 {
		t.Fatalf("expected getTx called once, got %d", factory.txCount)
	}
}

func TestGet_UsesConnectionWithoutScope(t *testing.T) {
	conn := &mockExec{}
	factory := &mockFactory{conn: conn}

	exec, err := Get(context.Background(), factory)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exec != conn {
		t.Fatalf("expected direct connection from Get, got %T", exec)
	}
	if factory.connCount != 1 {
		t.Fatalf("expected getConn called once, got %d", factory.connCount)
	}
}

func TestScopeOptions_RequireNewAndSuppress(t *testing.T) {
	ctx := context.Background()
	ctxWithScope, scope := ScopeOptions(ctx, false, sql.LevelSerializable)
	if scope == nil {
		t.Fatal("expected scope")
	}
	if scope.level != sql.LevelSerializable {
		t.Fatalf("expected isolation level %v, got %v", sql.LevelSerializable, scope.level)
	}

	_, nextScope := ScopeOptions(ctxWithScope, true, sql.LevelRepeatableRead)
	if nextScope == nil {
		t.Fatal("expected new scope")
	}
	if nextScope == scope {
		t.Fatal("expected requireNew to create a new scope")
	}

	suppressed := SuppressScope(ctxWithScope)
	if suppressed == ctxWithScope {
		t.Fatal("expected suppress scope to return a new context value")
	}
	if got, _ := suppressed.Value(contextTxKey).(*TxScope); got != nil {
		t.Fatal("expected suppressed context to have nil transaction scope")
	}
}

func TestTxScope_EndBehavior(t *testing.T) {
	tx := &mockTx{}
	scope := &TxScope{tx: tx}

	if err := scope.End(); err != nil {
		t.Fatalf("unexpected rollback error: %v", err)
	}
	if tx.rollbackCount != 1 {
		t.Fatalf("expected rollback once, got %d", tx.rollbackCount)
	}

	tx.rollbackCount = 0
	scope.Commit()
	if err := scope.End(); err != nil {
		t.Fatalf("unexpected commit error: %v", err)
	}
	if tx.commitCount != 1 {
		t.Fatalf("expected commit once, got %d", tx.commitCount)
	}
}

func TestTxScope_EndDuringPanicRollsBack(t *testing.T) {
	tx := &mockTx{}
	scope := &TxScope{tx: tx}

	panicked := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()

		defer func() {
			_ = scope.End()
		}()

		panic("boom")
	}()

	if !panicked {
		t.Fatal("expected panic to be propagated")
	}
	if tx.rollbackCount != 1 {
		t.Fatalf("expected rollback once during panic, got %d", tx.rollbackCount)
	}
}

func TestScanRows_PrimitiveType(t *testing.T) {
	rows := &mockRows{
		columns: []string{"value"},
		rows:    [][]any{{1}, {2}, {3}},
	}

	values, err := ScanRows[int](rows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(values) != 3 {
		t.Fatalf("expected 3 values, got %d", len(values))
	}
	if values[0] != 1 || values[1] != 2 || values[2] != 3 {
		t.Fatalf("unexpected values: %#v", values)
	}
	if !rows.closed {
		t.Fatal("expected rows to be closed")
	}
}

func TestScanRows_PointerToStruct(t *testing.T) {
	type User struct {
		ID   int    `ino:"id"`
		Name string `ino:"name"`
	}

	rows := &mockRows{
		columns: []string{"id", "name"},
		rows:    [][]any{{1, "alice"}},
	}

	values, err := ScanRows[*User](rows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(values) != 1 {
		t.Fatalf("expected 1 value, got %d", len(values))
	}
	if values[0].ID != 1 || values[0].Name != "alice" {
		t.Fatalf("unexpected struct value: %#v", values[0])
	}
}

func TestScanRow_NoRowsReturnsZero(t *testing.T) {
	rows := &mockRows{columns: []string{"value"}, rows: [][]any{}}

	value, err := ScanRow[int](rows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != 0 {
		t.Fatalf("expected zero value, got %v", value)
	}
}

func TestQueryAndQuerySingle(t *testing.T) {
	exec := &mockExec{
		queryFunc: func(ctx context.Context, query string, args ...any) (Rows, error) {
			if query == "select values" {
				return &mockRows{columns: []string{"value"}, rows: [][]any{{10}, {20}}}, nil
			}
			return &mockRows{columns: []string{"value"}, rows: [][]any{}}, nil
		},
	}

	values, err := Query[int](exec, context.Background(), "select values")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(values) != 2 || values[0] != 10 || values[1] != 20 {
		t.Fatalf("unexpected query values: %#v", values)
	}

	value, err := QuerySingle[int](exec, context.Background(), "select none")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != 0 {
		t.Fatalf("expected zero value from QuerySingle with no rows, got %v", value)
	}
}
