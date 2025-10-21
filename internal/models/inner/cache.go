package inner

import (
	"context"
	"time"
)

type Cache interface {
	cacheOperations

	Close() error
}

type cacheOperations interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) 		error
}