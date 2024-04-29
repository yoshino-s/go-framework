package configuration

import "github.com/spf13/pflag"

var _ Configuration = &CombinationConfiguration{}

type CombinationConfiguration []Configuration

func (c *CombinationConfiguration) Register(set *pflag.FlagSet) {
	for _, v := range *c {
		v.Register(set)
	}
}

func (c *CombinationConfiguration) Read() {
	for _, v := range *c {
		v.Read()
	}
}
