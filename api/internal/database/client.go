package db

import (
	"context"
	"sync"
	"time"

	cmn "Codex-Backend/api/common"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	Pool *pgxpool.Pool
}

var (
	instance *Client
	initErr  error
	once     sync.Once
)

type ClientConfig struct {
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	HealthCheckPeriod time.Duration
}

func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		MaxConns:          20,
		MinConns:          2,
		MaxConnLifetime:   time.Hour,
		HealthCheckPeriod: 30 * time.Second,
	}
}

func NewClient(ctx context.Context, connString string) (*Client, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	return &Client{Pool: pool}, nil
}

// func NewClient(ctx context.Context, connString string) (*Client, error) {
// 	return NewClientWithConfig(ctx, connString, DefaultClientConfig())
// }

func NewClientWithConfig(ctx context.Context, connString string, config ClientConfig) (*Client, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	// Apply configuration with validation
	if config.MaxConns <= 0 {
		config.MaxConns = 20
	}
	if config.MinConns <= 0 {
		config.MinConns = 2
	}
	if config.MinConns > config.MaxConns {
		config.MinConns = config.MaxConns
	}

	cfg.MaxConns = config.MaxConns
	cfg.MinConns = config.MinConns
	cfg.MaxConnLifetime = config.MaxConnLifetime
	cfg.HealthCheckPeriod = config.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Client{Pool: pool}, nil
}

func GetClient(ctx context.Context) (*Client, error) {
	once.Do(func() {
		connString := cmn.GetEnvVariable("DATABASE_URL")

		cfg, err := pgxpool.ParseConfig(connString)
		if err != nil {
			initErr = err
			return
		}

		// Production-ready defaults
		cfg.MaxConns = 20
		cfg.MinConns = 2
		cfg.MaxConnLifetime = time.Hour
		cfg.HealthCheckPeriod = 30 * time.Second

		pool, err := pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			initErr = err
			return
		}

		instance = &Client{Pool: pool}
	})

	return instance, initErr
}

func (c *Client) Close() {
	if instance != nil && instance.Pool != nil {
		instance.Pool.Close()
		instance = nil
		once = sync.Once{}
	}
}
