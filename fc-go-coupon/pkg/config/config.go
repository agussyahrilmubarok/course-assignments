package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `json:"server" mapstructure:"server"`
	Postgres PostgresConfig `json:"postgres" mapstructure:"postgres"`
}

func LoadConfig(location string) (*Config, error) {
	v := viper.New()

	dir := filepath.Dir(location)
	filename := filepath.Base(location)
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	v.AddConfigPath(dir)
	v.SetConfigName(name)
	if ext != "" {
		v.SetConfigType(strings.TrimPrefix(ext, "."))
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
