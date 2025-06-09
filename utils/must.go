package utils

import (
	"fmt"
	"reflect"
)

func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

func Must2[T1 any, T2 any](obj1 T1, obj2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return obj1, obj2
}

func MustNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func interfaceHasNilValue(actual any) bool {
	value := reflect.ValueOf(actual)
	kind := value.Kind()
	nilable := kind == reflect.Slice ||
		kind == reflect.Chan ||
		kind == reflect.Func ||
		kind == reflect.Ptr ||
		kind == reflect.Map

	// Careful: reflect.Value.IsNil() will panic unless it's an interface, chan, map, func, slice, or ptr
	// Reference: http://golang.org/pkg/reflect/#Value.IsNil
	return nilable && value.IsNil()
}

func NoNil(obj ...interface{}) bool {
	for _, o := range obj {
		if o == nil || interfaceHasNilValue(o) {
			return false
		}
	}
	return true
}

func MustNoNil(obj ...interface{}) {
	for _, o := range obj {
		if o == nil {
			panic(fmt.Errorf("should be %T not nil", o))
		}
	}
}
