package configuration

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var _ Configuration = &generateConfiguration{}
var GenerateConfiguration = &generateConfiguration{}

type generateConfiguration struct{}

func (*generateConfiguration) Register(flagSet *pflag.FlagSet) {
	flagSet.Bool("generate-config", false, "generate config")
	flagSet.String("generate-config-path", "", "generate config path")
	if err := viper.BindPFlags(flagSet); err != nil {
		panic(err)
	}
	Register(GenerateConfiguration)
}

func (*generateConfiguration) Read() {
	if viper.GetBool("generate-config") {
		p := viper.GetString("generate-config-path")
		viper.Set("generate-config-path", "")
		viper.Set("generate-config", false)
		if p != "" {
			err := viper.WriteConfigAs(viper.GetString("generate-config-path"))
			if err != nil {
				panic(err)
			}
		} else {
			tmpfile, err := os.CreateTemp("", "config*.yaml")
			if err != nil {
				panic(err)
			}
			if err = viper.WriteConfigAs(tmpfile.Name()); err != nil {
				panic(err)
			}
			_, err = tmpfile.Seek(0, 0)
			if err != nil {
				panic(err)
			}
			b, err := io.ReadAll(tmpfile)
			if err != nil {
				panic(err)
			}
			fmt.Print(string(b))
			if err := tmpfile.Close(); err != nil {
				panic(err)
			}
			if err := os.Remove(tmpfile.Name()); err != nil {
				panic(err)
			}
		}
		os.Exit(0)
	}
}
