package shorten

import (
	"context"
	"database/sql"
)

type Rows interface {
	Columns() []string
	Next() bool
	Scan(dest ...any) error
	Close() error
}

type Stmt interface {
	Exec(ctx context.Context, args ...any) (int64, error)
	Query(ctx context.Context, args ...any) (Rows, error)
	Close() error
}

type Exec interface {
	Prepare(ctx context.Context, query string) (Stmt, error)
	Exec(ctx context.Context, query string, args ...any) (int64, error)
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	Release() error
}

type Tx interface {
	Exec
	Commit() error
	Rollback() error
}

type Factory interface {
	getConn(ctx context.Context) (Exec, error)
	getTx(ctx context.Context, level sql.IsolationLevel) (Tx, error)
}

func Get(ctx context.Context, factory Factory) (Exec, error) {
	scope, _ := ctx.Value(contextTxKey).(*TxScope)
	if scope != nil {
		if exec := scope.tx; exec != nil {
			return exec, nil
		}

		tx, err := factory.getTx(ctx, scope.level)
		if err != nil {
			return nil, err
		}
		scope.tx = tx
		return tx, nil
	}

	conn, err := factory.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
