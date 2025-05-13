package log

import (
	"os"
	"runtime/debug"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/utils"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	ScopeName = "github.com/yoshino-s/go-framework/log"
)

var _ configuration.Configuration = (*LogConfiguration)(nil)

type LogConfiguration struct {
	Logger **zap.Logger
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
	Levels struct {
		Console string `mapstructure:"console"`
		File    string `mapstructure:"file"`
	} `mapstructure:"levels"`
	Otel bool `mapstructure:"otel"`
}

func (l *LogConfiguration) Register(flagSet *pflag.FlagSet) {
	flagSet.String("log.level", "info", "log level")
	flagSet.String("log.file", "", "log file path")
	flagSet.String("log.format", "", "log format, one of json, console, empty for default (console for dev, json for prod)")
	flagSet.Bool("log.rotate.enable", false, "enable log rotation")
	flagSet.Int("log.rotate.max_size", 500, "max size of log file in MB")
	flagSet.Int("log.rotate.max_age", 28, "max age of log file in days")
	flagSet.Int("log.rotate.max_backups", 3, "max number of log file backups")
	flagSet.String("log.levels.console", "", "log level for console, empty for same as log.level")
	flagSet.String("log.levels.file", "", "log level for file, empty for same as log.level")
	flagSet.Bool("log.otel", false, "enable sending logs to otel collector")
	if err := viper.BindPFlags(flagSet); err != nil {
		panic(err)
	}
	configuration.Register(l)
}

func (l *LogConfiguration) Read() {
	var c logConfig
	utils.MustDecodeFromMapstructure(viper.AllSettings()["log"], &c)

	if isInTest() {
		return
	}

	// var config zap.Config

	level := zapcore.WarnLevel
	if c.Level != "" {
		level = utils.Must(zapcore.ParseLevel(c.Level))
	}

	cores := make([]zapcore.Core, 0)

	// build console core
	consoleEncoder := zapcore.NewConsoleEncoder(NewColoredDevelopmentEncoderConfig())
	if c.Format == "json" || (c.Format == "" && !common.IsDev()) {
		consoleEncoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	consoleLevel := level
	if c.Levels.Console != "" {
		consoleLevel = utils.Must(zapcore.ParseLevel(c.Levels.Console))
	}

	cores = append(cores, zapcore.NewCore(
		consoleEncoder,
		zapcore.Lock(os.Stdout),
		consoleLevel,
	))

	// build file core
	if c.File != "" {
		fileLevel := level
		if c.Levels.File != "" {
			fileLevel = utils.Must(zapcore.ParseLevel(c.Levels.File))
		}

		var sink zapcore.WriteSyncer
		if c.Rotate.Enable {
			sink = zapcore.AddSync(&lumberjack.Logger{
				Filename:   c.File,
				MaxSize:    c.Rotate.MaxSize,
				MaxBackups: c.Rotate.MaxBackups,
				MaxAge:     c.Rotate.MaxAge,
			})
		} else {
			sink, _ = utils.Must2(zap.Open(c.File))
		}
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			sink,
			fileLevel,
		)
		cores = append(cores, fileCore)
	}

	if c.Otel {
		cores = append(cores,
			otelzap.NewCore(ScopeName),
		)
	}

	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(zapcore.WarnLevel))

	(*l.Logger) = logger
}
