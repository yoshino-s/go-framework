package configuration

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/magic"
)

var _ Configuration = &generateConfiguration{}
var GenerateConfiguration = &generateConfiguration{}

type generateConfig struct {
	Enable bool
	Schema bool
	Path   string
}

type generateConfiguration struct {
	Config generateConfig
}

func (*generateConfiguration) Register(flagSet *pflag.FlagSet) {
	flagSet.Bool("generate-config.enable", false, "generate config enable")
	flagSet.String("generate-config.path", "", "generate config path")
	common.MustNoError(viper.BindPFlags(flagSet))
	Register(GenerateConfiguration)
}

func marshalYaml() (string, error) {
	v := viper.GetViper()

	value, ok := magic.GetUnexported(v, "pflags")
	if !ok {
		return "", fmt.Errorf("pflags not found")
	}
	pflags, ok := value.(*map[string]viper.FlagValue)
	if !ok {
		return "", fmt.Errorf("pflags not found")
	}

	comments := make(map[string]string)

	for _, value := range *pflags {
		value, ok := magic.GetUnexported(value, "flag")
		if !ok {
			continue
		}
		flag, ok := value.(*pflag.Flag)
		if !ok {
			continue
		}
		if strings.HasPrefix(flag.Name, "generate-config.") {
			continue
		}
		comments[fmt.Sprintf("$.%s", flag.Name)] = flag.Usage
	}

	return magic.MarshalYamlWithComments(v.AllSettings(), comments)
}

func (c *generateConfiguration) Read() {
	err := common.DecodeFromMapstructure(viper.AllSettings()["generate-config"], &c.Config)
	if err != nil {
		panic(err)
	}

	if c.Config.Enable {
		viper.Set("generate-config", nil)

		content := []byte(common.Must(marshalYaml()))

		if c.Config.Path == "" {
			fmt.Print(string(content))
		} else {
			common.MustNoError(os.WriteFile(c.Config.Path, content, 0644))
		}

		os.Exit(0)
	}
}
