package inner

import (
	"context"

	"github.com/ayayaakasvin/trends-updater/internal/models"
)

type EventRepository interface {
	eventOperations

	Close() error
	Ping() 	error
}

type eventOperations interface {
	FetchUpdateTrending(ctx context.Context) ([]models.EventStats, error)
	ArchiveOldEvents(ctx context.Context) (int, error)
	UpdateEventStatuses(ctx context.Context) (int64, error)
}