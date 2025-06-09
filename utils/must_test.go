package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNoNil(t *testing.T) {
	Convey("TestNoNil", t, func() {
		var a *int
		So(a, ShouldBeNil)
		So(NoNil(a), ShouldBeFalse)
		var s *struct{}
		So(s, ShouldBeNil)
		So(NoNil(s), ShouldBeFalse)
	})
}
