package config

import (
	"strings"

	"example.com/backend/pkg/connections"
	"example.com/backend/pkg/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	App         App                    `json:"app" mapstructure:"app"`
	Logging     logger.Logging         `json:"logging" mapstructure:"logging"`
	PostgresSQL connections.PostgreSQL `json:"postgres" mapstructure:"postgres"`
	MySQL       connections.MySQL      `json:"mysql" mapstructure:"mysql"`
	Midtrans    connections.Midtrans   `json:"midtrans" mapstructure:"midtrans"`
}

func Load(path string) *Config {
	var cfg Config

	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("failed to read config")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to unmarshal config")
	}

	log.Info().
		Str("app_name", cfg.App.Name).
		Int("app_port", cfg.App.Port).
		Msg("configuration loaded successfully")

	return &cfg
}
