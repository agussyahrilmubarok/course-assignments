package account

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Port int    `mapstructure:"port"`
		Env  string `mapstructure:"env"`
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

	Jwt struct {
		SecretKey string `mapstructure:"secret_key"`
		ExpiresIn string `mapstructure:"expires_in"` // Example: "15m", "1h"
	} `mapstructure:"jwt"`

	Logger struct {
		Level    string `mapstructure:"level"`    // Example: "info", "debug"
		Filepath string `mapstructure:"filepath"` // Example: "logs/account.log"
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

func NewPostgres(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DbName,
		cfg.Postgres.SslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Info, Warn, Error
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	sqlDB.SetMaxOpenConns(cfg.Postgres.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Postgres.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Postgres.MaxOpenConns))

	return db, nil
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
