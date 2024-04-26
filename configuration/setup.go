package configuration

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var registeredConfigurations []Configuration

func Register(c Configuration) {
	registeredConfigurations = append(registeredConfigurations, c)
}

func Setup(
	name string,
) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetEnvPrefix(name)

	viper.SetConfigName(name)
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", name))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", name))
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	for _, c := range registeredConfigurations {
		c.Read()
	}
}
