package pricing

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Host string `mapstructure:"host"`
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

	Consul struct {
		Address  string `mapstructure:"address"`
		WaitTime string `mapstructure:"wait_time"` // Example: "15m", "1h"
	}
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
