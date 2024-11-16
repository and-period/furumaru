package scheduler

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/mediaconvert"
	"github.com/and-period/furumaru/api/pkg/medialive"
	"github.com/and-period/furumaru/api/pkg/sfn"
	"github.com/and-period/furumaru/api/pkg/storage"
	"go.uber.org/zap"
)

type Scheduler interface {
	Run(ctx context.Context, target time.Time) error
	Lambda(ctx context.Context) error
}

type Params struct {
	StepFunction       sfn.StepFunction
	MediaLive          medialive.MediaLive
	MediaConvert       mediaconvert.MediaConvert
	WaitGroup          *sync.WaitGroup
	Database           *database.Database
	Storage            storage.Bucket
	Store              store.Service
	Environment        string
	ArchiveBucketName  string
	ConvertJobTemplate string
}

type options struct {
	logger      *zap.Logger
	concurrency int64
	storageURL  *url.URL
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

func WithStorageURL(storageURL string) Option {
	return func(opts *options) {
		url, err := url.Parse(storageURL)
		if err != nil {
			return
		}
		opts.storageURL = url
	}
}

func newArchiveHLSPath(scheduleID string) string {
	if scheduleID == "" {
		return ""
	}
	return "/" + fmt.Sprintf(entity.BroadcastArchiveHLSPath, scheduleID)
}

func newArchiveMP4Path(scheduleID string) string {
	if scheduleID == "" {
		return ""
	}
	return "/" + fmt.Sprintf(entity.BroadcastArchiveMP4Path, scheduleID)
}
