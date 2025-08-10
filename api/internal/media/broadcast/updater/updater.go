package updater

import (
	"context"
	"net/url"
	"sync"

	"github.com/and-period/furumaru/api/internal/media/database"
)

type Updater interface {
	Lambda(ctx context.Context, event CreatePayload) error
}

type Params struct {
	WaitGroup  *sync.WaitGroup
	Database   *database.Database
	StorageURL *url.URL
}

type options struct {
	maxRetries int64
}

type Option func(*options)

func WithMaxRetries(maxRetries int64) Option {
	return func(opts *options) {
		opts.maxRetries = maxRetries
	}
}
