package worker

import (
	firestore_services "Codex-Backend/api/internal/services/collections"
	"context"

	"github.com/riverqueue/river"
)

type EPUBWorker struct {
	river.WorkerDefaults[ProcessEPUBArgs]
}

func (w *EPUBWorker) Work(ctx context.Context, job *river.Job[ProcessEPUBArgs]) error {
	return firestore_services.CreateNovelFromEPUB(job.Args.File, ctx)
}

type ProcessEPUBArgs struct {
	File []byte `json:"file"`
}

func (ProcessEPUBArgs) Kind() string { return "process_epub" }
