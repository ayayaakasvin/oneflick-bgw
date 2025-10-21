package worker

import (
	"context"
	"time"
)

type JobHandler interface {
	Run(ctx context.Context) error
}

type JobFunc func(ctx context.Context) error

func (f JobFunc) Run(ctx context.Context) error { return f(ctx) }

type Job struct {
	ID 				string
	operation 		JobHandler
	interval  		time.Duration

	executeOnceRun	bool
}

func NewJob(id string, operation JobHandler, tt time.Duration, executeOnceRun bool) Job {
	return Job{
		ID: 		id,
		operation: 	operation,
		interval:  	tt,
		executeOnceRun: executeOnceRun,
	}
}