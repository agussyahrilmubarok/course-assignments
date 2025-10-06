package config

type ServerConfig struct {
	Name        string `json:"name" mapstructure:"name"`
	Port        int    `json:"port" mapstructure:"port"`
	Environment string `json:"environment" mapstructure:"environment"`
	LogLevel    string `json:"log_level" mapstructure:"log_level"`
	LogFilepath string `json:"log_filepath" mapstructure:"log_filepath"`
}
