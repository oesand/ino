package shorten

import (
	"fmt"
	"reflect"
	"sync"
)

var structFields sync.Map

func scanStruct(typ reflect.Type, columns []string) (any, []any, error) {
	if typ.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("mapper: expects type %s to be struct", typ)
	}

	var (
		index  map[string][]int
		values = make([]any, 0, len(columns))
	)

	if idx, found := structFields.Load(typ); found {
		index = idx.(map[string][]int)
	} else {
		index = structIdx(typ)
		structFields.Store(typ, index)
	}

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

func structIdx(t reflect.Type) map[string][]int {
	fields := make(map[string][]int)
	for i := 0; i < t.NumField(); i++ {
		var (
			f    = t.Field(i)
			name = f.Name
		)
		if tn := f.Tag.Get("ino"); len(tn) != 0 {
			name = tn
		}
		switch {
		case name == "-", len(f.PkgPath) != 0 && !f.Anonymous:
			continue
		}
		switch {
		case f.Anonymous:
			if f.Type.Kind() != reflect.Ptr {
				for k, idx := range structIdx(f.Type) {
					fields[k] = append(f.Index, idx...)
				}
			}
		default:
			fields[name] = f.Index
		}
	}
	return fields
}
