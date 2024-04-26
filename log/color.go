package log

import (
	"time"

	"github.com/fatih/color"
	"go.uber.org/zap/zapcore"
)

func coloredISO8601TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	layout := "2006-01-02T15:04:05.000Z0700"

	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(color.BlueString((t.Format(layout))))
}

func ColoredShortCallerEncoder(
	caller zapcore.EntryCaller,
	enc zapcore.PrimitiveArrayEncoder,
) {
	enc.AppendString(color.CyanString(caller.TrimmedPath()))
}

func NewColoredDevelopmentEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     coloredISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   ColoredShortCallerEncoder,
	}
}
