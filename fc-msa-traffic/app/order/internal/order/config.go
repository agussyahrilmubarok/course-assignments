package order

import (
	"time"
)

type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		Env  string `mapstructure:"env"`
	} `json:"app"`

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

	MySQL struct {
		Host            string        `mapstructure:"host"`
		Port            int           `mapstructure:"port"`
		User            string        `mapstructure:"user"`
		Password        string        `mapstructure:"password"`
		DbName          string        `mapstructure:"dbname"`
		MaxOpenConns    int           `mapstructure:"max_open_conns"`
		MaxIdleConns    int           `mapstructure:"max_idle_conns"`
		ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // Example: "1h", "30m"
	} `mapstructure:"mysql"`

	Redis struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"redis"`

	Logger struct {
		Level    string `mapstructure:"level"`    // Example: "info", "debug"
		Filepath string `mapstructure:"filepath"` // Example: "logs/account.log"
	} `mapstructure:"logger"`

	Consul struct {
		Address  string `mapstructure:"address"`
		WaitTime string `mapstructure:"wait_time"` // Example: "15m", "1h"
	}
}
