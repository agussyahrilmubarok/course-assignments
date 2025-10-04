package config

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

func NewPostgres(cfg *Config, log zerolog.Logger) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Name,
		cfg.Postgres.SSLMode,
		cfg.Postgres.TimeZone,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse PostgreSQL DSN")
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// Connection pool settings
	poolCfg.MaxConns = int32(cfg.Postgres.MaxOpenConns)
	poolCfg.MinConns = int32(cfg.Postgres.MaxIdleConns)
	poolCfg.MaxConnLifetime = time.Duration(cfg.Postgres.ConnMaxLifetime) * time.Second

	// Connect
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create PostgreSQL connection pool")
		return nil, fmt.Errorf("connect: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("Unable to ping PostgreSQL")
		return nil, fmt.Errorf("ping: %w", err)
	}

	log.Info().
		Str("host", cfg.Postgres.Host).
		Int("port", cfg.Postgres.Port).
		Str("database", cfg.Postgres.Name).
		Msg("Connected to PostgreSQL successfully")

	return pool, nil
}
