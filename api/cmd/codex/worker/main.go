package main

import (
	cmn "Codex-Backend/api/internal/common"
	firestore_services "Codex-Backend/api/internal/usecases/collections"
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
	riverClient := cmn.InitializeRiverClient(context.Background())
	if err := riverClient.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	workers := river.NewWorkers()
	river.AddWorker(workers, &EPUBWorker{})
}

type EPUBWorker struct {
	river.WorkerDefaults[worker.ProcessEPUBArgs]
}

func (w *EPUBWorker) Work(ctx context.Context, job *river.Job[worker.ProcessEPUBArgs]) error {
	return firestore_services.CreateNovelFromEPUB(job.Args.File, ctx)
}
