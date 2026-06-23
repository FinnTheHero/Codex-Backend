package db

import (
	cmn "Codex-Backend/api/common"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type MigrationRunner struct {
	client *Client
}

func NewMigrationRunner(client *Client) *MigrationRunner {
	return &MigrationRunner{client: client}
}
func (mr *MigrationRunner) RunMigrations(ctx context.Context, migrationsDir string) error {
	// Ensure migration tracking table exists
	if err := mr.ensureMigrationTable(ctx); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// Get applied migrations
	applied, err := mr.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Get migration files
	files, err := mr.getMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Run pending migrations
	for _, file := range files {
		version := mr.extractVersion(file)
		if applied[version] {
			fmt.Printf("Migration %s already applied, skipping\n", version)
			continue
		}

		if err := mr.runMigration(ctx, filepath.Join(migrationsDir, file), version); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", file, err)
		}
		fmt.Printf("Applied migration: %s\n", file)
	}

	return nil
}
func (mr *MigrationRunner) ensureMigrationTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`

	_, err := mr.client.Pool.Exec(ctx, query)
	return err
}
func (mr *MigrationRunner) getAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	query := "SELECT version FROM schema_migrations"
	rows, err := mr.client.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}
func (mr *MigrationRunner) getMigrationFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".sql") {
			files = append(files, d.Name())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort files to ensure correct order
	sort.Strings(files)
	return files, nil
}
func (mr *MigrationRunner) extractVersion(filename string) string {
	// Extract version from filename like "001_initial_schema.sql"
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) > 0 {
		return strings.TrimSuffix(parts[0], ".sql")
	}
	return filename
}
func (mr *MigrationRunner) runMigration(ctx context.Context, filePath, version string) error {
	// Read SQL file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Start transaction
	tx, err := mr.client.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Execute migration
	if _, err := tx.Exec(ctx, string(content)); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record migration as applied
	if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	return nil
}
func (c *Client) EnsureSchema(ctx context.Context) error {
	if c == nil || c.Pool == nil {
		return &cmn.Error{Err: errors.New("postgres client not initialized"), Status: http.StatusInternalServerError}
	}

	runner := NewMigrationRunner(c)
	if err := runner.RunMigrations(ctx, "migrations"); err != nil {
		return &cmn.Error{Err: fmt.Errorf("migration error: %w", err), Status: http.StatusInternalServerError}
	}

	return nil
}
