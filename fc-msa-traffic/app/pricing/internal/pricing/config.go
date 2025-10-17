package pricing

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Port int    `mapstructure:"port"`
		Env  string `mapstructure:"env"`
	} `json:"app"`

	Redis struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"redis"`

	Logger struct {
		Level    string `mapstructure:"level"`    // Example: "info", "debug"
		Filepath string `mapstructure:"filepath"` // Example: "logs/pricing.log"
	} `mapstructure:"logger"`
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

func NewRedis(cfg *Config) (*redis.Client, error) {
	rdbAddr := fmt.Sprintf("%v:%v", cfg.Redis.Host, cfg.Redis.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     rdbAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}

func NewZerolog(cfg *Config) (zerolog.Logger, error) {
	logDir := filepath.Dir(cfg.Logger.Filepath)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return zerolog.Logger{}, err
	}

	logFile, err := os.OpenFile(cfg.Logger.Filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return zerolog.Logger{}, err
	}

	multi := zerolog.MultiLevelWriter(os.Stdout, logFile)

	level, err := zerolog.ParseLevel(cfg.Logger.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	logger := zerolog.New(multi).
		With().
		Timestamp().
		Logger()

	return logger, nil
}
