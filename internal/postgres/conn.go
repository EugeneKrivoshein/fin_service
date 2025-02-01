package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/EugeneKrivoshein/fin_service/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PGXProvider оборачивает pgxpool.Pool
type PGXProvider struct {
	Pool *pgxpool.Pool
}

func (p *PGXProvider) Close() {
	p.Pool.Close()
}

func NewPGXProvider() (*PGXProvider, error) {
	cfg, err := config.LoadConfig("config.env")
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных через pgx: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("база данных недоступна: %w", err)
	}

	return &PGXProvider{Pool: pool}, nil
}
