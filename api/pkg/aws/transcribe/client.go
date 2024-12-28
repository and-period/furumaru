package translate

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/transcribe"
	"go.uber.org/zap"
)

//nolint:lll
type Client interface {
	GetTranscriptionJob(ctx context.Context, in *transcribe.GetTranscriptionJobInput, opts ...func(*transcribe.Options)) (*transcribe.GetTranscriptionJobOutput, error)
	StartTranscriptionJob(ctx context.Context, in *transcribe.StartTranscriptionJobInput, opts ...func(*transcribe.Options)) (*transcribe.StartTranscriptionJobOutput, error)
}

type client struct {
	*transcribe.Client
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
	cli := transcribe.NewFromConfig(cfg, func(o *transcribe.Options) {
		o.Retryer = retry.NewStandard(func(so *retry.StandardOptions) {
			so.MaxAttempts = dopts.maxRetries
			so.MaxBackoff = dopts.interval
		})
	})
	return &client{
		Client: cli,
		logger: dopts.logger,
	}
}
