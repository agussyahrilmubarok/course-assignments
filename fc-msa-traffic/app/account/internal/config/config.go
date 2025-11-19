package config

import "time"

type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Env  string `mapstructure:"env"`
	} `mapstructure:"app"`

	Http struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"http"`

	Postgres struct {
		Host            string        `mapstructure:"host"`
		Port            int           `mapstructure:"port"`
		User            string        `mapstructure:"user"`
		Password        string        `mapstructure:"password"`
		DbName          string        `mapstructure:"dbname"`
		SslMode         string        `mapstructure:"sslmode"`
		MaxOpenConns    int           `mapstructure:"max_open_conns"`
		MaxIdleConns    int           `mapstructure:"max_idle_conns"`
		ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // Example: "1h", "30m"
	} `mapstructure:"postgres"`

	Jwt struct {
		SecretKey string `mapstructure:"secret_key"`
		ExpiresIn string `mapstructure:"expires_in"` // Example: "15m", "1h"
	} `mapstructure:"jwt"`

	Logger struct {
		Level    string `mapstructure:"level"`    // Example: "info", "debug"
		Filepath string `mapstructure:"filepath"` // Example: "logs/account.log"
	} `mapstructure:"logger"`

	Consul struct {
		Address  string `mapstructure:"address"`
		WaitTime string `mapstructure:"wait_time"` // Example: "15m", "1h"
	}
}
