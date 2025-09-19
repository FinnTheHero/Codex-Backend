package db

import (
	cmn "Codex-Backend/api/common"
	"Codex-Backend/api/internal/domain"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (c *Client) CreateNovel(novel domain.Novel, ctx context.Context) error {
	return c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		const insertSQL = `INSERT INTO novels (title, author, description) VALUES ($1,$2,$3)`
		if _, err := conn.Exec(ctx, insertSQL, novel.Title, novel.Author, novel.Description); err != nil {
			return &cmn.Error{Err: fmt.Errorf("insert novel: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	})
}

func (c *Client) GetNovelById(id string, ctx context.Context) (domain.Novel, error) {
	novel := domain.Novel{}

	if err := c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		if err := conn.QueryRow(ctx, "SELECT id, title, author, description FROM novels WHERE id = $1", id).Scan(&novel.ID, &novel.Title, &novel.Author, &novel.Description); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &cmn.Error{Err: fmt.Errorf("novel not found: %w", err), Status: http.StatusNotFound}
			}
			return &cmn.Error{Err: fmt.Errorf("get novel by id: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	}); err != nil {
		return domain.Novel{}, err
	}
	return novel, nil
}

func (c *Client) GetAllNovels(ctx context.Context) ([]domain.Novel, error) {
	novels := []domain.Novel{}

	if err := c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		rows, err := conn.Query(ctx, "SELECT id, title, author, description FROM novels")
		if err != nil {
			return &cmn.Error{Err: fmt.Errorf("get all novels: %w", err), Status: http.StatusInternalServerError}
		}
		defer rows.Close()

		for rows.Next() {
			novel := domain.Novel{}
			if err := rows.Scan(&novel.ID, &novel.Title, &novel.Author, &novel.Description); err != nil {
				return &cmn.Error{Err: fmt.Errorf("scan novel row: %w", err), Status: http.StatusInternalServerError}
			}
			novels = append(novels, novel)
		}
		if err := rows.Err(); err != nil {
			return &cmn.Error{Err: fmt.Errorf("scan novel rows: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return novels, nil
}

func (c *Client) UpdateNovel(novel domain.Novel, ctx context.Context) error {
	if err := c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		query := fmt.Sprintf("UPDATE novels SET title = $1, description = $2, updated_at = $3 WHERE id = $4")
		_, err := conn.Exec(ctx, query, novel.Title, novel.Description, time.Now(), novel.ID)
		if err != nil {
			return &cmn.Error{Err: fmt.Errorf("update novel: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteNovel(novelId string, ctx context.Context) error {
	if err := c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		query := fmt.Sprintf("UPDATE novels SET deleted = $1 WHERE id = $2")
		_, err := conn.Exec(ctx, query, true, novelId)
		if err != nil {
			return &cmn.Error{Err: fmt.Errorf("delete novel: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetNovelByTitle(title string, ctx context.Context) (domain.Novel, error) {
	novel := domain.Novel{}

	if err := c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		query := fmt.Sprintf("SELECT id, title, author, description FROM novels WHERE title = $1 AND deleted = $2")
		row := conn.QueryRow(ctx, query, title, false)
		if err := row.Scan(&novel.ID, &novel.Title, &novel.Author, &novel.Description); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &cmn.Error{Err: errors.New("Novel not found"), Status: http.StatusNotFound}
			}
			return &cmn.Error{Err: fmt.Errorf("get novel by title: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	}); err != nil {
		return domain.Novel{}, err
	}
	return novel, nil
}
