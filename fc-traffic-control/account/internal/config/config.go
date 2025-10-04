package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name  string `mapstructure:"name"`
		Port  int    `mapstructure:"port"`
		Level string `mapstructure:"level"`
	} `mapstructure:"app"`

	Postgres struct {
		Host            string `mapstructure:"host"`
		Port            int    `mapstructure:"port"`
		Name            string `mapstructure:"name"`
		User            string `mapstructure:"user"`
		Password        string `mapstructure:"password"`
		SSLMode         string `mapstructure:"ssl_mode"`
		TimeZone        string `mapstructure:"time_zone"`
		MaxOpenConns    int    `mapstructure:"max_open_conns"`
		MaxIdleConns    int    `mapstructure:"max_idle_conns"`
		ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	} `mapstructure:"postgres"`

	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`
}

// LoadEnv loads configuration from file or environment variables
func LoadEnv(configPath string) *Config {
	// Temporary minimal logger (always console)
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
	}
	log := zerolog.New(consoleWriter).With().Timestamp().Str("component", "config").Logger()

	v := viper.New()

	if configPath != "" {
		dir := filepath.Dir(configPath)
		filename := filepath.Base(configPath)
		ext := strings.TrimPrefix(filepath.Ext(filename), ".")
		name := strings.TrimSuffix(filename, filepath.Ext(filename))

		v.AddConfigPath(dir)
		v.SetConfigName(name)
		v.SetConfigType(ext)
		log.Info().Str("file", configPath).Msg("Loading configuration file")
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		log.Warn().Msg("No config path provided, using default search paths")
	}

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Warn().Err(err).Msg("Config file not found, using environment variables only")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal config")
	}

	if used := v.ConfigFileUsed(); used != "" {
		log.Info().Str("file", used).Msg("Configuration loaded successfully")
	} else {
		log.Info().Msg("Configuration loaded from environment variables")
	}

	return &cfg
}
