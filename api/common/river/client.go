package queue

import (
	"Codex-Backend/api/internal/service/worker"
	"context"
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

var (
	riverClient *river.Client[pgx.Tx]
	riverOnce   sync.Once
)

// func InitializeRiverClient(ctx context.Context, workers *river.Workers) *river.Client[pgx.Tx] {
// 	dbPool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
// 	if err != nil {
// 		panic(err)
// 	}

// 	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
// 		Level: slog.LevelInfo,
// 	}))

// 	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
// 		Logger: logger,
// 		Queues: map[string]river.QueueConfig{
// 			river.QueueDefault: {MaxWorkers: 10},
// 		},
// 		MaxAttempts: 3,
// 		Workers:     workers,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	return riverClient
// }

func InitializeRiverClient(ctx context.Context, workers *river.Workers) *river.Client[pgx.Tx] {
	var dbPool *pgxpool.Pool
	var err error

	// Retry connection with exponential backoff
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		dbPool, err = pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
		if err == nil {
			// Test the connection
			if pingErr := dbPool.Ping(ctx); pingErr == nil {
				log.Printf("Successfully connected to database on attempt %d", i+1)
				break
			} else {
				log.Printf("Database ping failed on attempt %d: %v", i+1, pingErr)
				err = pingErr
			}
		} else {
			log.Printf("Failed to create connection pool on attempt %d: %v", i+1, err)
		}

		if i == maxRetries-1 {
			panic(fmt.Sprintf("Failed to connect to database after %d attempts: %v", maxRetries, err))
		}

		// Wait before retrying (exponential backoff)
		waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
		log.Printf("Retrying database connection in %v...", waitTime)
		time.Sleep(waitTime)
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
