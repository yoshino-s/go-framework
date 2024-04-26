package http

type Feature uint16

const (
	FeatureNone    Feature = 0
	FeatureVersion Feature = 1 << iota
	FeatureHealth
	FeatureReady
	FeatureMetrics

	FeatureAll = FeatureVersion | FeatureHealth | FeatureReady | FeatureMetrics
)

func (f Feature) Has(flag Feature) bool {
	return f&flag != 0
}

func (f Feature) Add(flag Feature) Feature {
	return f | flag
}

func (f Feature) Remove(flag Feature) Feature {
	return f &^ flag
}

type Config struct {
	Log        bool    `mapstructure:"log"`
	Debug      bool    `mapstructure:"debug"`
	Feature    Feature `mapstructure:"feature"`
	ListenAddr string  `mapstructure:"addr"`
}

var DefaultConfig = Config{
	Log:        true,
	Debug:      false,
	Feature:    FeatureAll,
	ListenAddr: ":8080",
}
