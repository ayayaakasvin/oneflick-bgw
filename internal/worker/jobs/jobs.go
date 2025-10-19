package jobs

import (
	// "github.com/ayayaakasvin/trends-updater/internal/models/inner"
	"context"
	"os/exec"

	"github.com/ayayaakasvin/trends-updater/internal/worker"
	"github.com/sirupsen/logrus"
)

type CustomJobs struct {
// 	eventRepo 	inner.EventRepository
// 	cache 		inner.Cache


	logger 		*logrus.Logger
}

func NewCustomJobs(
	// er 		inner.EventRepository,
	// cc 		inner.Cache,
	lg 	*logrus.Logger,
) *CustomJobs {
	return &CustomJobs{
		// eventRepo: er,
		// cache: cc,
		logger: lg,
	}
}

func (h *CustomJobs) PrintEvery5s() worker.JobFunc {
	return func(ctx context.Context) error {
		h.logger.Info("Printing 5s...")
		return nil
	}
}

func (h *CustomJobs) PrintEvery10m() worker.JobFunc {
	return func(ctx context.Context) error {
		h.logger.Info("Printing 10m...")
		return nil
	}
}

func (h *CustomJobs) CancelContext() worker.JobFunc {
	return func(ctx context.Context) error {
		h.logger.Info("Cancelling context")
		return nil
	}
}

func (h *CustomJobs) PrintRandomWord3lenEvery5s() worker.JobFunc {
	return func(ctx context.Context) error {
		cmd := exec.Command("/home/ayayasvin/go/bin/random", "string", "-len=3", "-u")

		out, err := cmd.CombinedOutput()
		if err != nil {
			h.logger.Infof("Error: %v", err)
		}
		h.logger.Info(string(out))
		return nil
	}
}