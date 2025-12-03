package config

import "time"

type App struct {
	Name  string `json:"name" mapstructure:"name"`
	Port  int    `json:"port" mapstructure:"port"`
	Level string `json:"level" mapstructure:"level"`
	Token struct {
		SecretKey string        `json:"secret_key" mapstructure:"secret_key"`
		ExpiresAt time.Duration `json:"expires_at" mapstructure:"expires_at"`
	} `json:"token" mapstructure:"token"`
	Clients []string `json:"clients" mapstructure:"clients"`
}
