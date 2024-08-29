package common

import "fmt"

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

func NoNil(obj ...interface{}) bool {
	for _, o := range obj {
		if o == nil {
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
