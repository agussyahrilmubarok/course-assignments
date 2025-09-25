package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name        string        `mapstructure:"name"`
		Version     string        `mapstructure:"version"`
		Host        string        `mapstructure:"host"`
		Port        int           `mapstructure:"port"`
		IdleTimeout time.Duration `mapstructure:"idleTimeout"`
	} `mapstructure:"app"`

	Mongo struct {
		URI    string `mapstructure:"uri"`
		DBName string `mapstructure:"dbName"`
	} `mapstructure:"mongo"`

	JWT struct {
		SecretKey string        `mapstructure:"secretKey"`
		Expiry    time.Duration `mapstructure:"expiry"` // e.g: "10h", "30m"
	} `mapstructure:"jwt"`

	Log struct {
		Level         string `mapstructure:"level"`         // debug, info, warn, error
		Output        string `mapstructure:"output"`        // stdout, file, both
		FilePath      string `mapstructure:"filePath"`      // e.g: logs/app.log
		PrettyConsole bool   `mapstructure:"prettyConsole"` // true: human-friendly console
	} `mapstructure:"log"`
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
