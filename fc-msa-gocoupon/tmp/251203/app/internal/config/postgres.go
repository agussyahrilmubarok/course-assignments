package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, cfg *Config) (*Postgres, error) {
	dbConfig, err := pgxpool.ParseConfig(cfg.Postgres.URL)
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConns = cfg.Postgres.MaxConns
	dbConfig.MinConns = cfg.Postgres.MinConns
	dbConfig.MaxConnLifetime = cfg.Postgres.MaxConnLifetime
	dbConfig.HealthCheckPeriod = cfg.Postgres.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Postgres{
		Pool: pool,
	}, nil
}

func (d *Postgres) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}
