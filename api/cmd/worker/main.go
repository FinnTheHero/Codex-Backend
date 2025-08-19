package main

import (
	cmn "Codex-Backend/api/internal/common"
	"Codex-Backend/api/internal/usecases/worker"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/riverqueue/river"
)

func init() {
	if mode := os.Getenv("GIN_MODE"); mode == "debug" {
		cmn.LoadEnvVariables()
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workers := river.NewWorkers()
	river.AddWorker(workers, &worker.EPUBWorker{})

	riverClient := cmn.InitializeRiverClient(ctx, workers)
	if err := riverClient.Start(ctx); err != nil {
		log.Fatal("Failed to start River client:", err)
	}

	log.Println("River worker started successfully, waiting for jobs...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, starting graceful shutdown...")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := riverClient.Stop(shutdownCtx); err != nil {
		log.Printf("Error during River client shutdown: %v", err)
	} else {
		log.Println("River worker stopped gracefully")
	}

	select {}
}
