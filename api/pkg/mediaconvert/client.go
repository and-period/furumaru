//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package mediaconvert

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert/types"
)

type MediaConvert interface {
	CreateJob(ctx context.Context, template string, settings *types.JobSettings) error
}

type Params struct {
	RoleARN  string
	Endpoint string
}

type client struct {
	convert *mediaconvert.Client
	role    *string
}

type options struct {
	maxRetries int
	interval   time.Duration
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

func NewMediaConvert(cfg aws.Config, params *Params, opts ...Option) MediaConvert {
	dopts := &options{
		maxRetries: retry.DefaultMaxAttempts,
		interval:   retry.DefaultMaxBackoff,
	}
	for i := range opts {
		opts[i](dopts)
	}
	cli := mediaconvert.NewFromConfig(cfg, func(o *mediaconvert.Options) {
		o.BaseEndpoint = aws.String(params.Endpoint)
		o.Retryer = retry.NewStandard(func(o *retry.StandardOptions) {
			o.MaxAttempts = dopts.maxRetries
			o.MaxBackoff = dopts.interval
		})
	})
	return &client{
		convert: cli,
		role:    aws.String(params.RoleARN),
	}
}
