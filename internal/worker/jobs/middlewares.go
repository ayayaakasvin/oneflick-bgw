package jobs

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/ayayaakasvin/trends-updater/internal/worker"
	"github.com/sirupsen/logrus"
)

func (cj *CustomJobs) WithTimeLogging(id string) worker.Middleware {
	return func(next worker.JobHandler) worker.JobHandler {
		return worker.JobFunc(func(ctx context.Context) error {
			start := time.Now()

			err := next.Run(ctx)

			duration := time.Since(start)

			// Structured log entry
			cj.logger.WithFields(logrus.Fields{
				"job":      id,
				"duration": duration,
				"error":    err,
			}).Info("Job execution completed")

			return err
		})
	}
}

const ColorRed = "\033[31m"

func (cj *CustomJobs) WithRecover() worker.Middleware {
	return func(next worker.JobHandler) worker.JobHandler {
		return worker.JobFunc(func(ctx context.Context) error {
			defer func ()  {
				if rec := recover(); rec != nil {
					msg := fmt.Sprintf("panic recovered: %v", rec)
					cj.logger.Errorf("%s%s: %s", ColorRed, msg, debug.Stack())
				}
			}()

			err := next.Run(ctx)
			return err
		})
	}
}