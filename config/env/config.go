package env

import (
	"strings"

	"github.com/spf13/viper"
)

type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetUInt64(key string) uint64
	GetFloat64(key string) float64
	GetBool(key string) bool
	Init()
}

type viperConfig struct{}

func (v *viperConfig) Init() {
	viper.SetEnvPrefix(`baseapp`)
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(`.`, `_`)
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigType(`json`)
	viper.SetConfigFile(`config.json`)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func (v *viperConfig) GetString(key string) string {
	return viper.GetString(key)
}

func (v *viperConfig) GetInt(key string) int {
	return viper.GetInt(key)
}

// GetUInt64 ...
func (v *viperConfig) GetUInt64(key string) uint64 {
	return viper.GetUint64(key)
}

// GetFloat64 ...
func (v *viperConfig) GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func (v *viperConfig) GetBool(key string) bool {
	return viper.GetBool(key)
}

func NewViperConfig() Config {
	v := &viperConfig{}
	v.Init()

	return v
}
