package magic

import (
	"reflect"
	"runtime"
)

func unpackPtr(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		return unpackPtr(v.Elem())
	}
	return v
}

func GetFuncName(f any) (string, bool) {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return "", false
	}
	if v.IsNil() {
		return "", false
	}
	t := v.Pointer()
	return runtime.FuncForPC(t).Name(), true
}

func GetStructName(s any) (string, bool) {
	v := reflect.ValueOf(s)
	v = unpackPtr(v)
	if v.Kind() != reflect.Struct {
		return "", false
	}
	t := v.Type()
	return t.Name(), true
}
