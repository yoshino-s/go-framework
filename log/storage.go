package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Storage interface {
	Write(level zapcore.Level, entry zapcore.Entry, fields ...zapcore.Field) error
}

var _ zapcore.Core = (*storageCore)(nil)

type storageCore struct {
	storage Storage
	fields  []zapcore.Field
}

func storageToCore(s Storage) *storageCore {
	return &storageCore{
		storage: s,
	}
}

func (c *storageCore) clone() *storageCore {
	return &storageCore{
		storage: c.storage,
		fields:  append([]zapcore.Field{}, c.fields...),
	}
}

func (s *storageCore) Enabled(level zapcore.Level) bool {
	return true
}

func (s *storageCore) With(fields []zapcore.Field) zapcore.Core {
	c := s.clone()
	c.fields = append(c.fields, fields...)
	return c
}

func (s *storageCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if s.Enabled(ent.Level) {
		return ce.AddCore(ent, s)
	}
	return ce
}

func (s *storageCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	// Here you can implement the logic to store logs in your storage
	// For example, you can use s.storage to save the logs

	s.storage.Write(ent.Level, ent, fields...)
	return nil
}

func (s *storageCore) Sync() error {
	return nil
}

func WrapStorage(l *zap.Logger, storage Storage) *zap.Logger {
	l = l.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(
			c,
			storageToCore(storage),
		)
	}))

	return l
}
