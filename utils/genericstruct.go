package utils

import (
	"errors"
	"reflect"
)

var (
	ErrNoSuchField  error = errors.New("No such field in struct")
	ErrCantSetValue error = errors.New("Cannot set field value")
	ErrInvalidType  error = errors.New("Provided value type didn't match obj field type")
)

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return ErrNoSuchField
	}

	if !structFieldValue.CanSet() {
		return ErrCantSetValue
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return ErrInvalidType
	}

	structFieldValue.Set(val)
	return nil
}

func FillStruct(s interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
