package db

import (
	cmn "Codex-Backend/api/common"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type seekCursor struct {
	Index int64  `json:"idx"`
	ID    string `json:"id"`
}

// WithConn acquires a connection from the pool, runs fn(conn) and releases it.
// fn receives *pgxpool.Conn (you can call .Exec/.QueryRow on it).
func (c *Client) WithConn(ctx context.Context, fn func(conn *pgxpool.Conn) error) error {
	if c == nil || c.Pool == nil {
		return &cmn.Error{Err: errors.New("postgres client not initialized"), Status: http.StatusInternalServerError}
	}
	acq, err := c.Pool.Acquire(ctx)
	if err != nil {
		return &cmn.Error{Err: fmt.Errorf("failed to acquire conn: %w", err), Status: http.StatusInternalServerError}
	}
	defer acq.Release()
	return fn(acq)
}

// WithTx runs fn inside a transaction. It ensures proper rollback on error/panic and commits on success.
func (c *Client) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	if c == nil || c.Pool == nil {
		return &cmn.Error{Err: errors.New("postgres client not initialized"), Status: http.StatusInternalServerError}
	}

	acq, err := c.Pool.Acquire(ctx)
	if err != nil {
		return &cmn.Error{Err: fmt.Errorf("acquire conn for tx: %w", err), Status: http.StatusInternalServerError}
	}
	defer acq.Release()

	tx, err := acq.Begin(ctx)
	if err != nil {
		return &cmn.Error{Err: fmt.Errorf("begin tx: %w", err), Status: http.StatusInternalServerError}
	}

	// ensure rollback if fn fails or panic happens
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return &cmn.Error{Err: fmt.Errorf("commit tx: %w", err), Status: http.StatusInternalServerError}
	}
	return nil
}

func encodeCursor(c seekCursor) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func decodeCursor(encoded string) (seekCursor, error) {
	if encoded == "" {
		return seekCursor{Index: -1, ID: ""}, nil // special empty cursor
	}
	b, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return seekCursor{}, err
	}
	var sc seekCursor
	if err := json.Unmarshal(b, &sc); err != nil {
		return seekCursor{}, err
	}
	return sc, nil
}
