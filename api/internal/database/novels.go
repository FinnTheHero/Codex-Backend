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

func (c *Client) CreateNovelFromEpub(novel domain.Novel, chapters []domain.Chapter, ctx context.Context) error {
	chunkSize := 500
	totalChapters := len(chapters)

	return c.WithTx(ctx, func(tx pgx.Tx) error {
		var novelID string

		if err := tx.QueryRow(ctx,
			`INSERT INTO novels (title, author, description, chapter_count) VALUES ($1, $2, $3, $4) RETURNING id`,
			novel.Title, novel.Author, novel.Description, totalChapters,
		).Scan(&novelID); err != nil {
			return &cmn.Error{Err: fmt.Errorf("insert novel: %w", err), Status: http.StatusInternalServerError}
		}

		for i := 0; i < len(chapters); i += chunkSize {
			end := min(i+chunkSize, len(chapters))
			chunk := chapters[i:end]

			// Batch insert this chunk
			b := &pgx.Batch{}
			insertSQL := `INSERT INTO chapters (novel_id, title, author, description, content, chapter_index, deleted) VALUES ($1,$2,$3,$4,$5,$6,$7)`
			for _, ch := range chunk {
				b.Queue(insertSQL, novelID, ch.Title, ch.Author, ch.Description, ch.Content, ch.Index, ch.Deleted)
			}

			br := tx.SendBatch(ctx, b)
			for range chunk {
				if _, err := br.Exec(); err != nil {
					br.Close()
					return fmt.Errorf("batch exec chunk %d-%d: %w", i, end, err)
				}
			}
			br.Close()
		}

		return nil
	})
}

func (c *Client) CreateNovel(novel domain.CreateNovel, ctx context.Context) error {
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
		if err := conn.QueryRow(ctx, "SELECT id, title, author, description, deleted, created_at, updated_at FROM novels WHERE id = $1", id).Scan(&novel.ID, &novel.Title, &novel.Author, &novel.Description, &novel.Deleted, &novel.CreatedAt, &novel.UpdatedAt); err != nil {
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
		rows, err := conn.Query(ctx, "SELECT id, title, author, description, deleted, created_at, updated_at FROM novels")
		if err != nil {
			return &cmn.Error{Err: fmt.Errorf("get all novels: %w", err), Status: http.StatusInternalServerError}
		}
		defer rows.Close()

		for rows.Next() {
			novel := domain.Novel{}
			if err := rows.Scan(&novel.ID, &novel.Title, &novel.Author, &novel.Description, &novel.Deleted, &novel.CreatedAt, &novel.UpdatedAt); err != nil {
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
		query := "UPDATE novels SET title = $1, description = $2, updated_at = $3 WHERE id = $4"
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
		query := "UPDATE novels SET deleted = $1 WHERE id = $2"
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
		query := "SELECT id, title, author, description FROM novels WHERE title = $1 AND deleted = $2"
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
