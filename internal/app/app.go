package app

import (
	"context"
	"sync"
	"time"

	"github.com/ayayaakasvin/trends-updater/internal/worker"
	"github.com/ayayaakasvin/trends-updater/internal/worker/jobs"
	"github.com/sirupsen/logrus"
)

type BackgroundUpdater struct {
	logger 	*logrus.Logger

	worker	*worker.Worker
}

func NewBU(lg *logrus.Logger, wg *sync.WaitGroup, ctx context.Context) *BackgroundUpdater {
	return &BackgroundUpdater{
		logger: lg,
		worker: worker.NewWorker(lg, wg, ctx),
	}
}

func (b *BackgroundUpdater) RunApplication() {
	jobs := jobs.NewCustomJobs(b.logger)

	b.worker.Submit(jobs.PrintEvery5s(), time.Second * 5, false)

	b.worker.Submit(jobs.PrintRandomWord3lenEvery5s(), time.Second * 5, true)

	b.worker.Run()
}

func (b *BackgroundUpdater) Shutdown() {
	b.worker.Shutdown()
}