package config

import "github.com/spf13/viper"

type Config struct {
	App struct {
		Name   string `mapstructure:"name"`
		Host   string `mapstructure:"host"`
		Port   int    `mapstructure:"port"`
		Env    string `mapstructure:"env"`
		Metric struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
		} `json:"metric"`
	} `json:"app"`

	Postgres struct {
		Host            string `mapstructure:"host"`
		Port            int    `mapstructure:"port"`
		User            string `mapstructure:"user"`
		Password        string `mapstructure:"password"`
		DbName          string `mapstructure:"dbname"`
		SslMode         string `mapstructure:"sslmode"`
		MaxOpenConns    int    `mapstructure:"max_open_conns"`
		MaxIdleConns    int    `mapstructure:"max_idle_conns"`
		ConnMaxLifetime string `mapstructure:"conn_max_lifetime"` // Example: "1h", "30m"
	} `mapstructure:"postgres"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	Logger struct {
		Level    string `mapstructure:"level"`    // Example: "info", "debug"
		Filepath string `mapstructure:"filepath"` // Example: "logs/account.log"
	} `mapstructure:"logger"`

	OTEL struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `json:"otel"`
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
