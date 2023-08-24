package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/pkg/medialive"
	"github.com/and-period/furumaru/api/pkg/sfn"
	"go.uber.org/zap"
)

type Scheduler interface {
	Run(ctx context.Context, target time.Time) error
	Lambda(ctx context.Context) error
}

type Params struct {
	StepFunction      sfn.StepFunction
	MediaLive         medialive.MediaLive
	WaitGroup         *sync.WaitGroup
	Database          *database.Database
	Environment       string
	ArchiveBucketName string
}

type options struct {
	logger      *zap.Logger
	concurrency int64
}

type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithConcurrency(concurrency int64) Option {
	return func(opts *options) {
		opts.concurrency = concurrency
	}
}
