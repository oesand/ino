package shorten

import (
	"context"
	"fmt"
	"reflect"
)

// ScanRows scans all rows from rows into a slice of values of type T.
// For pointer-to-struct types, column names are mapped to struct fields using
// field names or `ino` tags.
func ScanRows[T any](rows Rows) ([]T, error) {
	typ := reflect.TypeFor[T]()

	var values []T
	if typ.Kind() == reflect.Pointer {
		elem := typ.Elem()
		if elem.Kind() != reflect.Struct {
			return nil, fmt.Errorf("shorten: expects type %s to be pointer to struct", typ)
		}

		columns := rows.Columns()
		if len(columns) == 0 {
			return nil, fmt.Errorf("shorten: no columns found")
		}

		var i int
		for rows.Next() {
			value, fields, err := scanStruct(elem, columns)
			if err != nil {
				return nil, err
			}

			err = rows.Scan(fields...)
			if err != nil {
				return nil, fmt.Errorf("shorten: scan row [%d]: %w", i, err)
			}
			i++

			values = append(values, value.(T))
		}
	} else {
		var i int
		for rows.Next() {
			var item T
			err := rows.Scan(&item)
			if err != nil {
				return nil, fmt.Errorf("shorten: scan row [%d]: %w", i, err)
			}
			i++

			values = append(values, item)
		}
	}

	err := rows.Close()
	if err != nil {
		return nil, err
	}
	return values, nil
}

// ScanRow scans a single row from rows into a value of type T. For pointer-to-struct
// types, column names are mapped to struct fields using field names or `ino` tags.
func ScanRow[T any](rows Rows) (T, error) {
	typ := reflect.TypeFor[T]()

	var result T
	if typ.Kind() == reflect.Pointer {
		elem := typ.Elem()
		if elem.Kind() != reflect.Struct {
			return result, fmt.Errorf("shorten: expects type %s to be pointer to struct", typ)
		}

		if !rows.Next() {
			return result, nil
		}

		columns := rows.Columns()
		if len(columns) == 0 {
			return result, fmt.Errorf("shorten: no columns found")
		}

		value, fields, err := scanStruct(elem, columns)
		if err != nil {
			return result, err
		}

		err = rows.Scan(fields...)
		if err != nil {
			return result, fmt.Errorf("shorten: scan: %w", err)
		}

		result = value.(T)
	} else {
		if !rows.Next() {
			return result, nil
		}

		err := rows.Scan(&result)
		if err != nil {
			return result, fmt.Errorf("shorten: scan: %w", err)
		}
	}

	err := rows.Close()
	return result, err
}

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
