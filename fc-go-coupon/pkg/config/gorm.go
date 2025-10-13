package config

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresConfig struct {
	Host         string        `json:"host" mapstructure:"host"`
	Port         int           `json:"port" mapstructure:"port"`
	User         string        `json:"user" mapstructure:"user"`
	Password     string        `json:"password" mapstructure:"password"`
	Database     string        `json:"database" mapstructure:"database"`
	SslMode      string        `json:"sslmode" mapstructure:"sslmode"`
	TimeZone     string        `json:"time_zone" mapstructure:"time_zone"`
	MaxConns     int           `json:"max_conns" mapstructure:"max_conns"`
	MaxIdleConns int           `json:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxLifeTime  time.Duration `json:"max_lifetime" mapstructure:"max_lifetime"`
	LogLevel     string        `json:"log_level" mapstructure:"log_level"`
}

func (p *PostgresConfig) getDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		p.Host, p.User, p.Password, p.Database, p.Port, p.SslMode, p.TimeZone,
	)
}

func (p *PostgresConfig) getLogLevel() logger.LogLevel {
	switch strings.ToLower(p.LogLevel) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn", "warning":
		return logger.Warn
	default:
		return logger.Info // default
	}
}

func NewPostgres(cfg *PostgresConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.getDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(cfg.getLogLevel()),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if cfg.MaxConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.MaxLifeTime)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
