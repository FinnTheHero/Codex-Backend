package db

import (
	cmn "Codex-Backend/api/common"
	"Codex-Backend/api/internal/domain"
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (c *Client) CreateUser(user domain.User, ctx context.Context) error {
	return c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		const insertSQL = `INSERT INTO users (id, email, username, type, password) VALUES ($1,$2,$3,$4,$5)`
		if _, err := conn.Exec(ctx, insertSQL, user.ID, user.Email, user.Username, user.Type, user.Password); err != nil {
			return &cmn.Error{Err: fmt.Errorf("insert user: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	})
}

func (c *Client) GetUserByEmail(email string, ctx context.Context) (domain.User, error) {
	var user domain.User
	if err := c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		const selectSQL = `SELECT id, email, username, type, password FROM users WHERE email = $1`
		if err := conn.QueryRow(ctx, selectSQL, email).Scan(&user.ID, &user.Email, &user.Username, &user.Type, &user.Password); err != nil {
			return &cmn.Error{Err: fmt.Errorf("select user by email: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	}); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (c *Client) GetUserById(userId string, ctx context.Context) (domain.User, error) {
	var user domain.User
	if err := c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		const selectSQL = `SELECT id, email, username, type, password FROM users WHERE id = $1`
		if err := conn.QueryRow(ctx, selectSQL, userId).Scan(&user.ID, &user.Email, &user.Username, &user.Type, &user.Password); err != nil {
			return &cmn.Error{Err: fmt.Errorf("select user by id: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	}); err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (c *Client) GetAllUsers(ctx context.Context) (*[]domain.User, error) {
	var users []domain.User
	if err := c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		const selectSQL = `SELECT id, email, username, type, password FROM users`
		rows, err := conn.Query(ctx, selectSQL)
		if err != nil {
			return &cmn.Error{Err: fmt.Errorf("select all users: %w", err), Status: http.StatusInternalServerError}
		}
		defer rows.Close()

		users, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (domain.User, error) {
			var user domain.User

			err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Type, &user.Password)
			if err != nil {
				return domain.User{}, err
			}
			return user, nil
		})
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &users, nil
}

func (c *Client) UpdateUser(user domain.User, ctx context.Context) error {
	return c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		const updateSQL = `UPDATE users SET email = $1, username = $2, type = $3, password = $4 WHERE id = $5`
		_, err := conn.Exec(ctx, updateSQL, user.Email, user.Username, user.Type, user.Password, user.ID)
		if err != nil {
			return &cmn.Error{Err: fmt.Errorf("update user: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	})
}

func (c *Client) DeleteUser(id string, ctx context.Context) error {
	return c.WithConn(ctx, func(conn *pgxpool.Conn) error {
		const deleteSQL = `UPDATE users SET deleted = $1 WHERE id = $2`
		_, err := conn.Exec(ctx, deleteSQL, true, id)
		if err != nil {
			return &cmn.Error{Err: fmt.Errorf("delete user: %w", err), Status: http.StatusInternalServerError}
		}
		return nil
	})
}
