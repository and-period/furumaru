//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package batch

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	"github.com/aws/aws-sdk-go-v2/service/batch/types"
	"go.uber.org/zap"
)

type Client interface {
	SubmitJob(ctx context.Context, params *SubmitJobParams) error
}

type SubmitJobParams struct {
	JobName       string
	JobDefinition string
	JobQueue      string
	Command       []string
}

type client struct {
	batch  *batch.Client
	logger *zap.Logger
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

func NewClient(cfg aws.Config, opts ...Option) Client {
	dopts := &options{
		maxRetries: retry.DefaultMaxAttempts,
		interval:   retry.DefaultMaxBackoff,
		logger:     zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	cli := batch.NewFromConfig(cfg, func(o *batch.Options) {
		o.Retryer = retry.NewStandard(func(so *retry.StandardOptions) {
			so.MaxAttempts = dopts.maxRetries
			so.MaxBackoff = dopts.interval
		})
	})
	return &client{
		batch:  cli,
		logger: dopts.logger,
	}
}

func (c *client) SubmitJob(ctx context.Context, params *SubmitJobParams) error {
	in := &batch.SubmitJobInput{
		JobName:       aws.String(params.JobName),
		JobDefinition: aws.String(params.JobDefinition),
		JobQueue:      aws.String(params.JobQueue),
		EcsPropertiesOverride: &types.EcsPropertiesOverride{
			TaskProperties: []types.TaskPropertiesOverride{
				{
					Containers: []types.TaskContainerOverrides{
						{
							Command: params.Command,
						},
					},
				},
			},
		},
	}
	_, err := c.batch.SubmitJob(ctx, in)
	return err
}
