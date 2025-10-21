package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ayayaakasvin/trends-updater/internal/models/inner"
	"github.com/ayayaakasvin/trends-updater/internal/worker"
	"github.com/sirupsen/logrus"
)

type CustomJobs struct {
	eventRepo 	inner.EventRepository
	cache 		inner.Cache

	logger 		*logrus.Logger
}

func NewCustomJobs(
	er 			inner.EventRepository,
	cc 			inner.Cache,
	lg 			*logrus.Logger,
) *CustomJobs {
	return &CustomJobs{
		eventRepo: er,
		cache: cc,
		logger: lg,
	}
}

const (
	trendingUpdateTimeKey 	= "trending_update_time"
	trendingKey 			= "trending_events"
	trendingTTL 			= time.Minute * 10
)

func (cj *CustomJobs) PingRepository() worker.JobFunc {
	return func(ctx context.Context) error {
		start := time.Now()

		if err := cj.eventRepo.Ping(); err != nil {
			duration := time.Since(start)
			errormsg := fmt.Sprintf("Failed to ping database after %s: %s", duration, err.Error())
			cj.logger.Error(errormsg)
			return fmt.Errorf("%s", errormsg)
		}

		duration := time.Since(start)
		cj.logger.WithFields(logrus.Fields{
			"operation": "ping database",
			"duration":  duration.String(),
		}).Info("Database connection is stable")

		return nil
	}
}

func (cj *CustomJobs) UpdateTrending() worker.JobFunc {
	return func(ctx context.Context) error {
		events, err := cj.eventRepo.FetchUpdateTrending(ctx)
		if err != nil {
			return fmt.Errorf("failed to retrieve records: %s", err.Error())
		}

		if len(events) == 0 {
			cj.logger.Warn("No trending events found in database")
			return nil
		}

		jsonData, err := json.Marshal(events)
		if err != nil {
			return fmt.Errorf("failed to marshall data: %s", err.Error())
		}

		if err := cj.cache.Set(ctx, trendingKey, jsonData, trendingTTL); err != nil {
			return fmt.Errorf("failed to write events records to cache: %s", err)
		}

		now := time.Now().UTC().Format(time.RFC3339)
		if err := cj.cache.Set(ctx, trendingUpdateTimeKey, now, trendingTTL); err != nil {
			return fmt.Errorf("failed to write events trending time to cache: %s", err)
		}

		cj.logger.WithField("count", len(events)).Info("Trending events successfully updated")
		return nil
	}
}

func (cj *CustomJobs) ArchieveOldEvents() worker.JobFunc {
	return func(ctx context.Context) error {
		count, err := cj.eventRepo.ArchiveOldEvents(ctx)
		if err != nil {
			return fmt.Errorf("failed to archieve events: %s", err.Error())
		}

		cj.logger.Infof("Archived %d old events", count)
		return nil
	}
}

func (cj *CustomJobs) UpdateEventsStatus() worker.JobFunc {
	return func(ctx context.Context) error {
		count, err := cj.eventRepo.UpdateEventStatuses(ctx)
		if err != nil {
			return fmt.Errorf("failed to update event statuses: %s", err.Error())
		}

		if count > 0 {
			cj.logger.Infof("Marked %d events as finished", count)
		}
		return nil
	}
}