package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Worker struct {
	mu 		sync.Mutex

	wg 		*sync.WaitGroup

	logger 	*logrus.Entry
	ctx    	context.Context
	cancel 	context.CancelFunc

	jobs	[]Job
}

func NewWorker(
	logger *logrus.Logger,
	wg *sync.WaitGroup,
	parent context.Context,
) *Worker {
	ctx, cancel := context.WithCancel(parent)
	return &Worker{
		logger: logger.WithField("component", "worker"),
		wg:     wg,
		ctx:    ctx,
		cancel: cancel,
		jobs: make([]Job, 0),
	}
}

func (w *Worker) AddJob(j ...Job) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.jobs = append(w.jobs, j...)
}

func (w *Worker) Run() {
	defer w.wg.Done()

	for _, j := range w.Jobs() {
		job := j
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()

			jobLogger := w.logger.WithField("jobID", job.ID)

			// create ticker inside goroutine so we can stop it when the
			// job exits - avoids leaking tickers when goroutine returns.
			ticker := time.NewTicker(job.interval)
			defer ticker.Stop()

			if j.executeOnceRun {
				if err := job.operation.Run(w.ctx); err != nil {
						if errors.Is(err, ErrFatal) {
							w.logger.Error(err)
							w.cancel()
							return
						}

						jobLogger.Warnf("Error recieved but not fatal: %s", err.Error())
					}
			}

			for {
				select {
				case <-ticker.C:
					if err := job.operation.Run(w.ctx); err != nil {
						w.cancel()
					}
				case <-w.ctx.Done():
					jobLogger.Errorf("Critical error recieved from worker %s: %s", job.ID, w.ctx.Err())
					return
				}
			}
		}()
	}
}

// Thread-safe snapshot (for iterating)
func (w *Worker) Jobs() []Job {
	w.mu.Lock()
	defer w.mu.Unlock()
	cp := make([]Job, len(w.jobs))
	copy(cp, w.jobs)
	return cp
}