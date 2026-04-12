package shorten

import (
	"reflect"
	"strings"
	"sync"
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

func TestStructIdx_TagHandling(t *testing.T) {
	type TestStruct struct {
		ID   int    `ino:"user_id"`
		Name string `ino:"-"`
		Age  int
	}

	index := structIdx(reflect.TypeOf(TestStruct{}))

	if idx, ok := index["user_id"]; !ok {
		t.Error("expected 'user_id' from tag")
	} else if len(idx) != 1 {
		t.Errorf("expected flat index for top-level field, got %v", idx)
	}

	if _, ok := index["Name"]; ok {
		t.Error("Name should be ignored due to tag '-'")
	}

	if _, ok := index["Age"]; !ok {
		t.Error("expected 'Age'")
	}
}

func TestStructIdx_EmbeddedStruct(t *testing.T) {
	type Address struct {
		City string
	}

	type User struct {
		Address
		Name string
	}

	index := structIdx(reflect.TypeOf(User{}))

	if _, ok := index["City"]; !ok {
		t.Error("expected 'City' from embedded struct")
	}

	if idx, ok := index["City"]; !ok {
		t.Fatal("City not found")
	} else if len(idx) != 2 {
		t.Errorf("expected index [1,0] for City, got %v", idx)
	}
}

func TestStructIdx_UnexportedField(t *testing.T) {
	type TestStruct struct {
		Public  string
		private string // unexported
	}

	index := structIdx(reflect.TypeOf(TestStruct{}))

	if _, ok := index["private"]; ok {
		t.Error("unexported field should be ignored")
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

func TestStructIdx_EmptyStruct(t *testing.T) {
	type Empty struct{}

	index := structIdx(reflect.TypeOf(Empty{}))
	if len(index) != 0 {
		t.Errorf("empty struct should have empty index map, got %d fields", len(index))
	}
}

func BenchmarkScanStruct_CacheHit(b *testing.B) {
	type TestStruct struct {
		ID   int
		Name string
	}

	columns := []string{"ID", "Name"}

	// Заполняем cache
	structFields.Store(reflect.TypeOf(TestStruct{}), map[string][]int{
		"ID":   {},
		"Name": {1},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = scanStruct(reflect.TypeOf(&TestStruct{}), columns)
	}
}

func BenchmarkScanStruct_CacheMiss(b *testing.B) {
	type TestStruct struct {
		ID   int
		Name string
	}

	columns := []string{"ID", "Name"}
	structFields = sync.Map{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = scanStruct(reflect.TypeOf(&TestStruct{}), columns)
	}
}
