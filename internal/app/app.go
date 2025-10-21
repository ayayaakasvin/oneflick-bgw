package app

import (
	"context"
	"sync"
	"time"

	"github.com/ayayaakasvin/trends-updater/internal/models/inner"
	"github.com/ayayaakasvin/trends-updater/internal/worker"
	"github.com/ayayaakasvin/trends-updater/internal/worker/jobs"
	"github.com/sirupsen/logrus"
)

type BackgroundUpdater struct {
	logger 	*logrus.Logger

	eventRepo 	inner.EventRepository
	cache 		inner.Cache

	worker	*worker.Worker
}

func NewBU(lg *logrus.Logger, wg *sync.WaitGroup, ctx context.Context, er inner.EventRepository, cc inner.Cache) *BackgroundUpdater {
	return &BackgroundUpdater{
		logger: lg,
		worker: worker.NewWorker(lg, wg, ctx),
		eventRepo: er,
		cache: cc,
	}
}

func (b *BackgroundUpdater) setup() {
	jobs := jobs.NewCustomJobs(b.eventRepo, b.cache, b.logger)

	pingID := "Ping-1-3m"
	pingHandler := worker.Chain(jobs.PingRepository(), jobs.WithRecover(), jobs.WithTimeLogging(pingID))
	pingJob := worker.NewJob(pingID, pingHandler, time.Minute * 3, true)

	updateStatusID := "Update-Status-2-10m"
	updateStatusHandler := worker.Chain(jobs.UpdateEventsStatus(), jobs.WithRecover(), jobs.WithTimeLogging(updateStatusID))
	updateStatusJob := worker.NewJob(updateStatusID, updateStatusHandler, time.Minute * 10, false)

	archiveEventsID := "Archive-Old-Events-3-1d"
	archiveEventsHandler := worker.Chain(jobs.ArchieveOldEvents(), jobs.WithRecover(), jobs.WithTimeLogging(archiveEventsID))
	archiveEventsJob := worker.NewJob(archiveEventsID, archiveEventsHandler, time.Hour * 24, false)

	b.worker.AddJob(pingJob, updateStatusJob, archiveEventsJob)
}

func (b *BackgroundUpdater) RunApplication() {
	b.setup()

	b.worker.Run()
}