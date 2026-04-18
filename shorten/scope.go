package shorten

import (
	"context"
	"database/sql"

	"github.com/oesand/ino/internal"
)

var contextTxKey = internal.CtxKey{Key: "shorten/tx"}

var DefaultLevel = sql.LevelDefault

// Scope returns a context containing a transaction scope. It reuses an existing
// scope when one is already present in the context.
func Scope(ctx context.Context) (context.Context, *TxScope) {
	return ScopeOptions(ctx, false, DefaultLevel)
}

// ScopeOptions returns a context containing a transaction scope with the given
// isolation level. If requireNew is true, a fresh scope is created regardless of
// any existing scope in the context.
func ScopeOptions(ctx context.Context, requireNew bool, level sql.IsolationLevel) (context.Context, *TxScope) {
	scope := &TxScope{
		level: level,
	}

	parent, _ := ctx.Value(contextTxKey).(*TxScope)
	if requireNew || parent == nil {
		ctx = context.WithValue(ctx, contextTxKey, scope)
	}

	return ctx, scope
}

// SuppressScope returns a context that explicitly removes any active transaction scope.
func SuppressScope(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextTxKey, nil)
}

// TxScope controls commit/rollback behavior for a transaction lifecycle stored in context.
type TxScope struct {
	level sql.IsolationLevel

	tx     Tx
	commit bool
}

// Commit marks the current transaction scope so that End will commit instead
// of rolling back.
func (scope *TxScope) Commit() {
	scope.commit = true
}

// End finalizes the scoped transaction. If Commit was called, the transaction
// is committed; otherwise it is rolled back. If End is called during a panic,
// the transaction is rolled back and the panic is rethrown.
func (scope *TxScope) End(err *error) {
	if scope.tx == nil {
		return
	}

	if r := recover(); r != nil {
		_ = scope.tx.Rollback()
		panic(r)
	}

	var ierr error
	if scope.commit {
		ierr = scope.tx.Commit()
	} else {
		ierr = scope.tx.Rollback()
	}

	if err != nil && *err == nil {
		*err = ierr
	}
}
