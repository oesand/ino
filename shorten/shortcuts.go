package shorten

import "fmt"

func ScanRowsStruct[T any](rows Rows) ([]*T, error) {
	defer rows.Close()

	columns := rows.Columns()
	if len(columns) == 0 {
		return nil, fmt.Errorf("mapper: no columns found")
	}

	var i int
	var values []*T
	for rows.Next() {
		value, fields, err := scanStruct[T](columns)
		if err != nil {
			return nil, err
		}

		err = rows.Scan(fields...)
		if err != nil {
			return nil, fmt.Errorf("scan rows [%d]: %w", i, err)
		}
		i++

		values = append(values, value)
	}
	return values, nil
}

func ScanRowStruct[T any](rows Rows) (*T, error) {
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	columns := rows.Columns()
	if len(columns) == 0 {
		return nil, fmt.Errorf("mapper: no columns found")
	}

	value, fields, err := scanStruct[T](columns)
	if err != nil {
		return nil, err
	}

	err = rows.Scan(fields...)
	if err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	return value, nil
}
