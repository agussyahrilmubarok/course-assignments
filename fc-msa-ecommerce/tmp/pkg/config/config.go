package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig   `json:"app" mapstructure:"app"`
	Mongo MongoConfig `json:"mongo" mapstructure:"mongo"`
}

func Load(path string) (*Config, error) {
	var cfg Config

	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
