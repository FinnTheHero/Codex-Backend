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

// SQL query constants
const (
	listChaptersAscSQL = `
		SELECT id, title, author, description, content, chapter_index, deleted, created_at, updated_at
		FROM chapters
		WHERE novel_id = $1 AND (chapter_index, id) > ($2, $3)
		ORDER BY chapter_index ASC, id ASC
		LIMIT $4`

	listChaptersAscFirstSQL = `
		SELECT id, title, author, description, content, chapter_index, deleted, created_at, updated_at
		FROM chapters
		WHERE novel_id = $1
		ORDER BY chapter_index ASC, id ASC
		LIMIT $2`

	listChaptersDescSQL = `
		SELECT id, title, author, description, content, chapter_index, deleted, created_at, updated_at
		FROM chapters
		WHERE novel_id = $1 AND (chapter_index, id) < ($2, $3)
		ORDER BY chapter_index DESC, id DESC
		LIMIT $4`

	listChaptersDescFirstSQL = `
		SELECT id, title, author, description, content, chapter_index, deleted, created_at, updated_at
		FROM chapters
		WHERE novel_id = $1
		ORDER BY chapter_index DESC, id DESC
		LIMIT $2`
)

/*
ListChaptersSeek returns up to `limit` chapters for a novel using seek-pagination.

  - cursor: encoded cursor string from previous page (or empty for first page)

  - limit: max rows to return

  - asc: if true order by chapter_index ASC, id ASC (older -> newer); if false, DESC Returns:

  - slice of chapters

  - nextCursor: encoded cursor to use for the next page (empty if no more rows)
*/
func (c *Client) ListChaptersSeek(novelId string, limit int, cursor string, asc bool, ctx context.Context) ([]domain.Chapter, string, error) {
	if limit <= 0 {
		limit = 100
	}

	// decode cursor
	sc, err := decodeCursor(cursor)
	if err != nil {
		return nil, "", &cmn.Error{Err: fmt.Errorf("invalid cursor: %w", err), Status: http.StatusBadRequest}
	}

	var results []domain.Chapter
	if err := c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		var rows pgx.Rows

		fetchLimit := limit + 1

		if asc {
			if sc.Index == -1 { // First page
				rows, err = conn.Query(ctx, listChaptersAscFirstSQL, novelId, fetchLimit)
			} else {
				rows, err = conn.Query(ctx, listChaptersAscSQL, novelId, sc.Index, sc.ID, fetchLimit)
			}
		} else {
			if sc.Index == -1 { // First page
				rows, err = conn.Query(ctx, listChaptersDescFirstSQL, novelId, fetchLimit)
			} else {
				rows, err = conn.Query(ctx, listChaptersDescSQL, novelId, sc.Index, sc.ID, fetchLimit)
			}
		}

		if err := rows.Err(); err != nil {
			return &cmn.Error{Err: fmt.Errorf("rows error: %w", err), Status: http.StatusInternalServerError}
		}
		defer rows.Close()

		results, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.Chapter, error) {
			var chapter domain.Chapter

			err := row.Scan(&chapter.Title, &chapter.Author, &chapter.Description,
				&chapter.Content, &chapter.Index, &chapter.Deleted, &chapter.CreatedAt, &chapter.UpdatedAt)
			if err != nil {
				return domain.Chapter{}, &cmn.Error{Err: fmt.Errorf("scan ListChaptersSeek: %w", err), Status: http.StatusInternalServerError}
			}

			return chapter, nil
		})
		if err != nil {
			return &cmn.Error{Err: fmt.Errorf("collect rows: %w", err), Status: http.StatusInternalServerError}
		}

		return nil
	}); err != nil {
		return nil, "", err
	}

	var nextCursor string
	hasMore := len(results) > limit
	if hasMore {
		results = results[:limit]
	}

	if len(results) > 0 && hasMore {
		lastResult := results[len(results)-1]
		lastCursor := seekCursor{Index: int64(lastResult.Index), ID: lastResult.ID}
		nextCursor, err = encodeCursor(lastCursor)
		if err != nil {
			return nil, "", &cmn.Error{Err: fmt.Errorf("encode cursor: %w", err), Status: http.StatusInternalServerError}
		}
	}

	return results, nextCursor, nil
}

func (c *Client) CreateChapter(novelId string, chapter domain.Chapter, ctx context.Context) error {
	var newIndex int64

	if err := c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		err := c.Pool.QueryRow(ctx, `UPDATE novels SET chapter_count = chapter_count + 1, updated_at = now() WHERE id = $1 RETURNING chapter_count`, novelId).Scan(&newIndex)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &cmn.Error{Err: fmt.Errorf("novel not found: %w", err), Status: http.StatusNotFound}
			}
			return &cmn.Error{Err: fmt.Errorf("Update novels chapter_count: %w", err), Status: http.StatusInternalServerError}
		}

		// Insert chapter using newIndex
		const insertSQL = `
			INSERT INTO chapters (novel_id, title, author, description, content, chapter_index)
			VALUES ($1, $2, $3, $4, $5, $6);
			`

		if _, err = c.Pool.Exec(ctx, insertSQL,
			novelId,
			chapter.Title,
			chapter.Author,
			chapter.Description,
			chapter.Content,
			newIndex,
		); err != nil {
			return &cmn.Error{Err: fmt.Errorf("insert chapter: %w", err), Status: http.StatusInternalServerError}
		}

		return nil
	}); err != nil {
		return &cmn.Error{Err: fmt.Errorf("create chapter: %w", err), Status: http.StatusInternalServerError}
	}

	return nil
}

func (c *Client) GetChapterById(novelId string, chapterId string, ctx context.Context) (domain.Chapter, error) {
	chapter := domain.Chapter{}

	if err := c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		if err := conn.QueryRow(ctx, "SELECT id, novel_id, title, author, description, content, chapter_index, deleted, created_at, updated_at FROM chapters WHERE id = $1 AND novel_id = $2 LIMIT 1", chapterId, novelId).Scan(
			&chapter.ID, novelId, &chapter.Title, &chapter.Author, &chapter.Description, &chapter.Content, &chapter.Index, &chapter.Deleted, &chapter.CreatedAt, &chapter.UpdatedAt,
		); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return &cmn.Error{Err: errors.New("chapter not found"), Status: http.StatusNotFound}
			}
			return &cmn.Error{Err: fmt.Errorf("postgres client error - get chapter by id: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	}); err != nil {
		return domain.Chapter{}, err
	}

	return chapter, nil
}

// Use seek pagination to get chapters in batches
func (c *Client) GetAllChapters(novelId string, pageSize int, asc bool, ctx context.Context) ([]domain.Chapter, error) {
	if c == nil || c.Pool == nil {
		return nil, &cmn.Error{Err: errors.New("postgres client not initialized"), Status: http.StatusInternalServerError}
	}
	if pageSize <= 0 {
		pageSize = 500
	}

	var all []domain.Chapter
	cursor := ""
	for {
		chs, nextCursor, err := c.ListChaptersSeek(novelId, pageSize, cursor, asc, ctx)
		if err != nil {
			return nil, err
		}
		all = append(all, chs...)
		if nextCursor == "" {
			break
		}
		cursor = nextCursor
	}
	return all, nil
}

func (c *Client) UpdateChapter(novelId string, chapter domain.Chapter, ctx context.Context) error {
	if err := c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		query := "UPDATE chapters SET title = $1, description = $2, content = $3, updated_at = $4 WHERE id = $5"
		_, err := conn.Exec(ctx, query, chapter.Title, chapter.Description, chapter.Content, time.Now(), chapter.ID)
		if err != nil {
			return &cmn.Error{Err: fmt.Errorf("update chapter: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteChapter(chapterId string, ctx context.Context) error {
	if err := c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		query := "UPDATE chapters SET deleted = $1 WHERE id = $2"
		_, err := conn.Exec(ctx, query, true, chapterId)
		if err != nil {
			return &cmn.Error{Err: fmt.Errorf("delete chapter: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
