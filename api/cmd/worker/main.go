package main

import (
	cmn "Codex-Backend/api/internal/common"
	queue "Codex-Backend/api/internal/common/river"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	if mode := os.Getenv("GIN_MODE"); mode == "debug" {
		cmn.LoadEnvVariables()
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	riverClient := queue.GetRiverClient(ctx)
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
}
