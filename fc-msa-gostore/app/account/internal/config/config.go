package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`

	MongoDB struct {
		URI     string        `mapstructure:"uri"`
		DB      string        `mapstructure:"db"`
		Timeout time.Duration `mapstructure:"timeout"`
	} `mapstructure:"mongodb"`

	Logging struct {
		Filepath string `mapstructure:"filepath"`
		Level    string `mapstructure:"level"`
	} `mapstructure:"logging"`
}

func NewConfig(filepath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(filepath)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
