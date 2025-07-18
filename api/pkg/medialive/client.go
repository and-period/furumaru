//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package medialive

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/medialive"
	"go.uber.org/zap"
)

type MediaLive interface {
	StartChannel(ctx context.Context, channelID string) error
	StopChannel(ctx context.Context, channelID string) error
	CreateSchedule(ctx context.Context, params *CreateScheduleParams) error
	ActivateStaticImage(ctx context.Context, channelID, imageURL string) error
	DeactivateStaticImage(ctx context.Context, channelID string) error
}

type client struct {
	media  *medialive.Client
	logger *zap.Logger
	now    func() time.Time
}

type options struct {
	maxRetries int
	interval   time.Duration
	logger     *zap.Logger
}

type Option func(*options)

func WithMaxRetries(maxRetries int) Option {
	return func(opts *options) {
		opts.maxRetries = maxRetries
	}
}

func WithInterval(interval time.Duration) Option {
	return func(opts *options) {
		opts.interval = interval
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func NewMediaLive(cfg aws.Config, opts ...Option) MediaLive {
	dopts := &options{
		maxRetries: retry.DefaultMaxAttempts,
		interval:   retry.DefaultMaxBackoff,
		logger:     zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	cli := medialive.NewFromConfig(cfg, func(o *medialive.Options) {
		o.Retryer = retry.NewStandard(func(o *retry.StandardOptions) {
			o.MaxAttempts = dopts.maxRetries
			o.MaxBackoff = dopts.interval
		})
	})
	return &client{
		media:  cli,
		logger: dopts.logger,
		now:    time.Now,
	}
}
