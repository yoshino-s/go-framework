package configuration

import (
	"runtime/debug"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ Configuration = (*logConfiguration)(nil)
var LogConfiguration = &logConfiguration{}

type logConfiguration struct{}

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

func (*logConfiguration) Register(flagSet *pflag.FlagSet) {
	flagSet.String("log.level", "info", "log level")
	flagSet.String("log.file", "", "log file path")
	if err := viper.BindPFlags(flagSet); err != nil {
		panic(err)
	}
	Register(LogConfiguration)
}

func (*logConfiguration) Read() {
	logFile := viper.GetString("log.file")
	logLevel := viper.GetString("log.level")

	if isInTest() {
		return
	}

	var config zap.Config

	if common.Version == "dev" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig = log.NewColoredDevelopmentEncoderConfig()
	} else {
		config = zap.NewProductionConfig()
	}
	if logFile != "" {
		config.OutputPaths = []string{logFile}
	}

	if logLevel != "" {
		level, err := zapcore.ParseLevel(logLevel)
		if err != nil {
			panic(err)
		}
		config.Level = zap.NewAtomicLevelAt(level)
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
