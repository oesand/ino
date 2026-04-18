package shorten

import (
	"context"
)

func Query[T any](exec Exec, ctx context.Context, query string, args ...any) ([]T, error) {
	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return ScanRows[T](rows)
}

func QuerySingle[T any](exec Exec, ctx context.Context, query string, args ...any) (T, error) {
	rows, err := exec.Query(ctx, query, args...)
	if err != nil {
		var zero T
		return zero, err
	}
	return ScanRow[T](rows)
}

func FireExec(factory Factory, ctx context.Context, query string, args ...any) (result int64, err error) {
	exec, err := Get(ctx, factory)
	if err != nil {
		return -1, err
	}
	defer exec.Release(&err)

	return exec.Exec(ctx, query, args...)
}

func FireQuery[T any](factory Factory, ctx context.Context, query string, args ...any) (result []T, err error) {
	exec, err := Get(ctx, factory)
	if err != nil {
		return nil, err
	}
	defer exec.Release(&err)

	return Query[T](exec, ctx, query, args...)
}

func FireQuerySingle[T any](factory Factory, ctx context.Context, query string, args ...any) (result T, err error) {
	exec, err := Get(ctx, factory)
	if err != nil {
		var zero T
		return zero, err
	}
	defer exec.Release(&err)

	return QuerySingle[T](exec, ctx, query, args...)
}
