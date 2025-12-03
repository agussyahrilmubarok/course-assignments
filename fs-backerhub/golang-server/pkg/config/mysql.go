package config

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQL struct {
	Host            string        `json:"host" mapstructure:"host"`
	Port            int           `json:"port" mapstructure:"port"`
	Name            string        `json:"name" mapstructure:"name"`
	Username        string        `json:"username" mapstructure:"username"`
	Password        string        `json:"password" mapstructure:"password"`
	MaxOpenConns    int           `json:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns" mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
}

func NewGormMySQL(cfg MySQL) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
