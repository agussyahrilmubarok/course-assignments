package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Name string `mapstructure:"name"`
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	Logging struct {
		Filepath string `mapstructure:"filepath"`
		Level    string `mapstructure:"level"`
	} `mapstructure:"logging"`

	Postgres struct {
		URL               string        `mapstructure:"url"`
		MaxConns          int32         `mapstructure:"max_conns"`
		MinConns          int32         `mapstructure:"min_conns"`
		MaxConnLifetime   time.Duration `mapstructure:"max_conn_lifetime"`
		HealthCheckPeriod time.Duration `mapstructure:"health_check_period"`
	} `mapstructure:"postgres"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	Kafka struct {
		Brokers []string `mapstructure:"brokers"`
		GroupID string   `mapstructure:"group_id"`
	}

	Zipkin struct {
		Url string `mapstructure:"url"`
	} `mapstructure:"zipkin"`
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
