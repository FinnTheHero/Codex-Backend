package worker

import (
	"Codex-Backend/api/internal/service"
	"context"

	"github.com/riverqueue/river"
)

type EPUBWorker struct {
	river.WorkerDefaults[ProcessEPUBArgs]
}

func (w *EPUBWorker) Work(ctx context.Context, job *river.Job[ProcessEPUBArgs]) error {
	return service.CreateNovelFromEPUB(job.Args.File, ctx)
}

type ProcessEPUBArgs struct {
	File []byte `json:"file"`
}

func (ProcessEPUBArgs) Kind() string { return "process_epub" }
