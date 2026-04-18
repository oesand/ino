package shorten

import (
	"reflect"
	"strings"
	"testing"
)

func TestScanStruct_SimpleStruct(t *testing.T) {
	type TestStruct struct {
		ID   int    `ino:"id"`
		Name string `ino:"name"`
	}

	columns := []string{"id", "name"}

	result, values, err := scanStruct(reflect.TypeFor[TestStruct](), columns)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	testStruct := result.(*TestStruct)

	if ptr, ok := values[0].(*int); !ok {
		t.Errorf("first value should be *int, got %T", values[0])
	} else {
		*ptr = 123
	}

	if ptr, ok := values[1].(*string); !ok {
		t.Errorf("second value should be *string, got %T", values[1])
	} else {
		*ptr = "test"
	}

	if testStruct.ID != 123 {
		t.Errorf("expected ID=123, got %d", testStruct.ID)
	}
	if testStruct.Name != "test" {
		t.Errorf("expected Name=test, got %s", testStruct.Name)
	}
}

func TestScanStruct_MissingColumn(t *testing.T) {
	type TestStruct struct {
		ID   int    `ino:"id"`
		Name string `ino:"name"`
	}

	columns := []string{"id", "missing"}

	_, _, err := scanStruct(reflect.TypeFor[TestStruct](), columns)
	if err == nil {
		t.Fatal("expected error for missing column")
	}

	expected := `mapper: missing destination name "missing" in shorten.TestStruct`
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("expected error containing %q, got %q", expected, err.Error())
	}
}

func TestScanStruct_InvalidType(t *testing.T) {
	columns := []string{"field"}

	_, _, err := scanStruct(reflect.TypeOf(42), columns)
	if err == nil {
		t.Fatal("expected error for invalid type")
	}

	_, _, err = scanStruct(reflect.TypeOf(struct{}{}), columns)
	if err == nil {
		t.Fatal("expected error for invalid type")
	}
}

func TestScanStruct_EmptyColumns(t *testing.T) {
	type TestStruct struct {
		ID int
	}

	_, _, err := scanStruct(reflect.TypeFor[TestStruct](), []string{})
	if err != nil {
		t.Errorf("unexpected error for empty columns: %v", err)
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
