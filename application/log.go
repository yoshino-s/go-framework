package application

import (
	"runtime/debug"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
	"github.com/yoshino-s/go-framework/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	Level string `mapstructure:"level"`
	File  string `mapstructure:"file"`
	Debug bool   `mapstructure:"debug"`
}

func (l *logConfiguration) Register(flagSet *pflag.FlagSet) {
	flagSet.String("log.level", "info", "log level")
	flagSet.String("log.file", "", "log file path")
	flagSet.Bool("log.debug", false, "log debug")
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

	var config zap.Config

	if common.IsDev() || c.Debug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig = log.NewColoredDevelopmentEncoderConfig()
	} else {
		config = zap.NewProductionConfig()
	}
	if c.File != "" {
		config.OutputPaths = []string{c.File}
	}

	if c.Level != "" {
		level, err := zapcore.ParseLevel(c.Level)
		if err != nil {
			panic(err)
		}
		config.Level = zap.NewAtomicLevelAt(level)
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	(*l.logger) = logger
}
