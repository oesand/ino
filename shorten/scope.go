package shorten

import (
	"context"
	"database/sql"

	"github.com/oesand/ino/internal"
)

var contextTxKey = internal.CtxKey{Key: "shorten/tx"}

func Scope(ctx context.Context) (context.Context, *TxScope) {
	return ScopeOptions(ctx, false, sql.LevelReadCommitted)
}

func ScopeOptions(ctx context.Context, requireNew bool, level sql.IsolationLevel) (context.Context, *TxScope) {
	var scope *TxScope
	if !requireNew {
		scope, _ = ctx.Value(contextTxKey).(*TxScope)
	}

	if scope == nil {
		scope = &TxScope{
			level: level,
		}

		ctx = context.WithValue(ctx, contextTxKey, scope)
	}

	return ctx, scope
}

func SuppressScope(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextTxKey, nil)
}

type TxScope struct {
	level sql.IsolationLevel

	tx     Tx
	commit bool
}

func (scope *TxScope) Commit() {
	scope.commit = true
}

func (scope *TxScope) End() error {
	if scope.tx == nil {
		return nil
	}

	if r := recover(); r != nil {
		_ = scope.tx.Rollback()
		panic(r)
	}

	if scope.commit {
		return scope.tx.Commit()
	}
	return scope.tx.Rollback()
}
