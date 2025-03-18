package log

import (
	"testing"

	"go.uber.org/zap"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHack(t *testing.T) {
	Convey("TestHack", t, func() {
		logger := zap.NewNop()
		SetLoggerName(logger, "test")
		So(logger.Name(), ShouldEqual, "test")
	})
}
