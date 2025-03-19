package magic

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type TestStruct struct {
}

func TestGetName(t *testing.T) {
	Convey("GetFuncName", t, func() {
		_, ok := GetFuncName(1)
		So(ok, ShouldBeFalse)
		name, ok := GetFuncName(func() {})
		So(ok, ShouldBeTrue)
		So(name, ShouldEqual, "github.com/yoshino-s/go-framework/magic.TestGetName.func1.1")
	})

	Convey("GetStructName", t, func() {
		_, ok := GetStructName(1)
		So(ok, ShouldBeFalse)
		name, ok := GetStructName(TestStruct{})
		fmt.Println(name, ok)
		So(ok, ShouldBeTrue)
		So(name, ShouldEqual, "TestStruct")
	})
}
