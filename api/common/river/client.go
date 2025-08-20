package queue

import (
	"Codex-Backend/api/internal/services/worker"
	"context"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

var (
	riverClient *river.Client[pgx.Tx]
	riverOnce   sync.Once
)

func InitializeRiverClient(ctx context.Context, workers *river.Workers) *river.Client[pgx.Tx] {
	dbPool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
		Logger: logger,
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 10},
		},
		MaxAttempts: 3,
		Workers:     workers,
	})
	if err != nil {
		panic(err)
	}

	return riverClient
}

func GetRiverClient(ctx context.Context) *river.Client[pgx.Tx] {
	riverOnce.Do(func() {
		log.Println("Initializing River client...")

		workers := river.NewWorkers()
		river.AddWorker(workers, &worker.EPUBWorker{})

		riverClient = InitializeRiverClient(ctx, workers)

		log.Println("River client initialized successfully")
	})

	return riverClient
}
