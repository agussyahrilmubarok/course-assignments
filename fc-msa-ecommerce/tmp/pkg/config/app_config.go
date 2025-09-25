package config

type AppConfig struct {
	Name    string `json:"name" mapstructure:"name"`
	Version string `json:"version" mapstructure:"version"`
	Host    string `json:"host" mapstructure:"host"`
	Port    int    `json:"port" mapstructure:"port"`
}
