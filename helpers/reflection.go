package helpers

import (
	"fmt"
	"reflect"
)

func GetFieldValue(data any, field string) (any, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}
	if field == "" {
		return nil, fmt.Errorf("field name is empty")
	}

	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, fmt.Errorf("data is nil pointer")
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("unsupported type: %s", val.Kind())
	}

	structField := val.FieldByName(field)
	if !structField.IsValid() {
		return nil, fmt.Errorf("field %s not found", field)
	}
	if !structField.CanInterface() {
		return nil, fmt.Errorf("field %s is unexported", field)
	}

	return structField.Interface(), nil
}
