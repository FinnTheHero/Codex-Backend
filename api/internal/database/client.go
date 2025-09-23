package db

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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

func (c *Client) EnsureSchema(ctx context.Context) error {
	if c == nil || c.Pool == nil {
		return &cmn.Error{Err: errors.New("postgres client not initialized"), Status: http.StatusInternalServerError}
	}

	tx, err := c.Pool.Begin(ctx)
	if err != nil {
		return &cmn.Error{Err: fmt.Errorf("begin tx for schema: %w", err), Status: http.StatusInternalServerError}
	}
	defer tx.Rollback(ctx)

	stmts := []string{
		// extension for gen_random_uuid
		`CREATE EXTENSION IF NOT EXISTS pgcrypto;`,

		`CREATE TABLE IF NOT EXISTS users (
				id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
				username text NOT NULL,
				type text NOT NULL,
				email text NOT NULL,
				password text NOT NULL,
				created_at timestamptz NOT NULL DEFAULT now(),
				updated_at timestamptz NOT NULL DEFAULT now()
			);`,

		// novels with chapter_count for atomic index allocation
		`CREATE TABLE IF NOT EXISTS novels (
				id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
				title text NOT NULL,
				author text NOT NULL,
				description text NOT NULL,
				chapter_count bigint NOT NULL DEFAULT 0,
				created_at timestamptz NOT NULL DEFAULT now(),
				updated_at timestamptz NOT NULL DEFAULT now()
			);`,

		// parent partitioned chapters table (hash partition on novel_id)
		`CREATE TABLE IF NOT EXISTS chapters (
				id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
				novel_id uuid NOT NULL,
				title text NOT NULL,
				author text NOT NULL,
				description text NOT NULL,
				content text NOT NULL,
				chapter_index bigint DEFAULT 0,
				deleted boolean DEFAULT false,
				created_at timestamptz NOT NULL DEFAULT now(),
				updated_at timestamptz NOT NULL DEFAULT now()
			) PARTITION BY HASH (novel_id);`,
		// an index that supports seek pagination: novel_id, chapter_index, id
		`CREATE INDEX IF NOT EXISTS idx_chapters_novel_index_id ON chapters (novel_id, chapter_index, id);`,

		`CREATE INDEX IF NOT EXISTS idx_novels_title_id ON novels (title, id);`,
		`CREATE INDEX IF NOT EXISTS idx_users_email_id ON users (email, id);`,
	}

	for _, s := range stmts {
		if _, err := tx.Exec(ctx, s); err != nil {
			return &cmn.Error{Err: fmt.Errorf("schema creation exec: %w", err), Status: http.StatusInternalServerError}
		}
	}

	// Create partitions (idempotent)
	const partitionsCount = 16
	for i := range partitionsCount {
		stmt := fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS chapters_p%d PARTITION OF chapters FOR VALUES WITH (MODULUS %d, REMAINDER %d);`,
			i, partitionsCount, i,
		)
		if _, err := tx.Exec(ctx, stmt); err != nil {
			return &cmn.Error{Err: fmt.Errorf("creating partition %d: %w", i, err), Status: http.StatusInternalServerError}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return &cmn.Error{Err: fmt.Errorf("commit schema creation: %w", err), Status: http.StatusInternalServerError}
	}

	return nil
}
