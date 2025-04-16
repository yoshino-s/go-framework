package application

import (
	"os"
	"runtime/debug"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ configuration.Configuration = (*logConfiguration)(nil)

type logConfiguration struct {
	logger **zap.Logger
}

func isInTest() bool {
	stacks := strings.Split(string(debug.Stack()), "\n")
	for _, line := range stacks {
		if strings.HasPrefix(line, "\t") {
			path := strings.Split(strings.TrimSpace(line), ":")[0]
			if strings.HasSuffix(path, "_test.go") {
				return true
			}
		}
	}
	return false
}

type logConfig struct {
	Level  string `mapstructure:"level"`
	File   string `mapstructure:"file"`
	Format string `mapstructure:"format"`
	Rotate struct {
		Enable     bool `mapstructure:"enable"`
		MaxSize    int  `mapstructure:"max_size"`
		MaxAge     int  `mapstructure:"max_age"`
		MaxBackups int  `mapstructure:"max_backups"`
	} `mapstructure:"rotate"`
}

func (l *logConfiguration) Register(flagSet *pflag.FlagSet) {
	flagSet.String("log.level", "info", "log level")
	flagSet.String("log.file", "", "log file path")
	flagSet.String("log.format", "", "log format, one of json, console, empty for default (console for dev, json for prod)")
	flagSet.Bool("log.rotate.enable", false, "enable log rotation")
	flagSet.Int("log.rotate.max_size", 500, "max size of log file in MB")
	flagSet.Int("log.rotate.max_age", 28, "max age of log file in days")
	flagSet.Int("log.rotate.max_backups", 3, "max number of log file backups")
	if err := viper.BindPFlags(flagSet); err != nil {
		panic(err)
	}
	configuration.Register(l)
}

func (l *logConfiguration) Read() {
	var c logConfig
	err := common.DecodeFromMapstructure(viper.AllSettings()["log"], &c)
	if err != nil {
		panic(err)
	}

	if isInTest() {
		return
	}

	// var config zap.Config

	level := zapcore.WarnLevel
	if c.Level != "" {
		level, err = zapcore.ParseLevel(c.Level)
		if err != nil {
			panic(err)
		}
	}

	cores := make([]zapcore.Core, 0)

	// build console core
	consoleEncoder := zapcore.NewConsoleEncoder(log.NewColoredDevelopmentEncoderConfig())
	if c.Format == "json" || (c.Format == "" && !common.IsDev()) {
		consoleEncoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}
	cores = append(cores, zapcore.NewCore(
		consoleEncoder,
		zapcore.Lock(os.Stdout),
		level,
	))

	// build file core
	if c.File != "" {
		var sink zapcore.WriteSyncer
		if c.Rotate.Enable {
			sink = zapcore.AddSync(&lumberjack.Logger{
				Filename:   c.File,
				MaxSize:    c.Rotate.MaxSize,
				MaxBackups: c.Rotate.MaxBackups,
				MaxAge:     c.Rotate.MaxAge,
			})
		} else {
			sink, _, err = zap.Open(c.File)
			if err != nil {
				panic(err)
			}
		}
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			sink,
			level,
		)
		cores = append(cores, fileCore)
	}

	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(zapcore.WarnLevel))

	(*l.logger) = logger
}
