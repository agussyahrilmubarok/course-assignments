package config

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQL struct {
	Host            string        `json:"host" mapstructure:"host"`
	Port            int           `json:"port" mapstructure:"port"`
	Name            string        `json:"name" mapstructure:"name"`
	Username        string        `json:"username" mapstructure:"username"`
	Password        string        `json:"password" mapstructure:"password"`
	SSLMode         string        `json:"ssl_mode" mapstructure:"ssl_mode"`
	TimeZone        string        `json:"time_zone" mapstructure:"time_zone"`
	MaxOpenConns    int           `json:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns" mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
}

func NewGormPostgres(cfg PostgreSQL) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Username,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
		cfg.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}
