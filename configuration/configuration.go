package configuration

import (
	"github.com/spf13/pflag"
)

type Configuration interface {
	Register(set *pflag.FlagSet)
	Read()
}
