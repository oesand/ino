package shorten

import (
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

func scanStruct(typ reflect.Type, columns []string) (any, []any, error) {
	if typ.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("mapper: expects type %s to be struct", typ)
	}

	index := getStructMapping(typ)
	values := make([]any, 0, len(columns))
	val := reflect.New(typ).Elem()
	for _, name := range columns {
		idx, found := index[name]
		if !found {
			return nil, nil, fmt.Errorf("mapper: missing destination name %q in %s", name, typ)
		}
		field := val.FieldByIndex(idx)
		values = append(values, field.Addr().Interface())
	}
	return val.Addr().Interface(), values, nil
}
