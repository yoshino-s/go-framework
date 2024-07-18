package configuration

import (
	"encoding/json"
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
	Format string
	Path   string
}

type generateConfiguration struct {
	Config generateConfig
}

func (*generateConfiguration) Register(flagSet *pflag.FlagSet) {
	flagSet.Bool("generate-config.enable", false, "generate config enable")
	flagSet.String("generate-config.format", "yaml", "generate config format, one of json, yaml, env")
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

func walk(v any, prefix string, env map[string]string) {
	switch v := v.(type) {
	case map[string]interface{}:
		for key, value := range v {
			walk(value, fmt.Sprintf("%s.%s", prefix, key), env)
		}
	case []interface{}:
		for i, value := range v {
			walk(value, fmt.Sprintf("%s.%d", prefix, i), env)
		}
	default:
		env[prefix] = fmt.Sprintf("%v", v)
	}
}

func marshalEnv() ([]byte, error) {
	v := viper.GetViper()
	prefix := v.GetEnvPrefix()

	if prefix != "" {
		prefix = fmt.Sprintf("%s_", strings.ToUpper(prefix))
	}

	value, ok := magic.GetUnexported(v, "envKeyReplacer")
	if !ok {
		return nil, fmt.Errorf("envKeyReplacer not found")
	}

	envKeyReplacer, ok := value.(*viper.StringReplacer)
	if !ok {
		return nil, fmt.Errorf("envKeyReplacer should be *viper.StringReplacer, but got %T", value)
	}

	env := make(map[string]string)

	walk(v.AllSettings(), "", env)

	var content []byte

	for key, value := range env {
		key = strings.ToUpper(prefix + (*envKeyReplacer).Replace(key[1:]))
		content = append(content, []byte(fmt.Sprintf("%s=%s\n", key, value))...)
	}

	return content, nil
}

func (c *generateConfiguration) Read() {
	err := common.DecodeFromMapstructure(viper.AllSettings()["generate-config"], &c.Config)
	if err != nil {
		panic(err)
	}

	if c.Config.Enable {
		viper.Set("generate-config", nil)
		var content []byte

		if c.Config.Format == "json" {
			content = common.Must(json.MarshalIndent(viper.AllSettings(), "", "  "))
		} else if c.Config.Format == "yaml" {
			content = []byte(common.Must(marshalYaml()))
		} else {
			content = common.Must(marshalEnv())
		}

		if c.Config.Path == "" {
			fmt.Print(string(content))
		} else {
			common.MustNoError(os.WriteFile(c.Config.Path, content, 0644))
		}

		os.Exit(0)
	}
}
