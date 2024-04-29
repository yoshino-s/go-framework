package magic

import (
	"reflect"
	"unsafe"
)

func GetUnexported(
	v any,
	fieldName string,
) (any, bool) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	value = value.FieldByName(fieldName)
	if !value.IsValid() {
		return nil, false
	}
	if !value.CanAddr() {
		value = value.Elem()
	}
	if !value.CanAddr() {
		return nil, false
	}
	value = reflect.NewAt(value.Type(), unsafe.Pointer(value.UnsafeAddr()))
	return value.Interface(), true
}
