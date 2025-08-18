package main

import (
	cmn "Codex-Backend/api/internal/common"
	"Codex-Backend/api/internal/usecases/worker"
	"context"
	"log"
	"os"

	"github.com/riverqueue/river"
)

func init() {
	if mode := os.Getenv("GIN_MODE"); mode == "debug" {
		cmn.LoadEnvVariables()
	}
}

func main() {
	workers := river.NewWorkers()
	river.AddWorker(workers, &worker.EPUBWorker{})
	riverClient := cmn.InitializeRiverClient(context.Background(), workers)
	if err := riverClient.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Println("River worker started successfully, waiting for jobs...")

	select {}
}
