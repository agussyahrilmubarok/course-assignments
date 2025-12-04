package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Backend struct {
		Name  string `mapstructure:"name"`
		Port  int    `mapstructure:"port"`
		Level string `mapstructure:"string"`
		Token struct {
			SecretKey string        `mapstructure:"secret_key"`
			ExpiresAt time.Duration `mapstructure:"expires_at"`
		} `mapstructure:"token"`
	} `mapstructure:"backend"`

	Logger struct {
		Filepath string `mapstructure:"filepath"`
		Level    string `mapstructure:"level"`
	} `mapstructure:"logger"`

	Postgres Postgres `mapstructure:"postgres"`
	Midtrans Midtrans `mapstructure:"midtrans"`
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
