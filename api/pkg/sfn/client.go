//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package sfn

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

type StepFunction interface {
	// ステートマシンの実行
	StartExecution(ctx context.Context, input interface{}) error
}

type Params struct {
	StateMachineARN string
}

type client struct {
	sfn *sfn.Client
	arn *string
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

func NewStepFunction(cfg aws.Config, params *Params, opts ...Option) StepFunction {
	dopts := &options{
		maxRetries: retry.DefaultMaxAttempts,
		interval:   retry.DefaultMaxBackoff,
	}
	for i := range opts {
		opts[i](dopts)
	}
	cli := sfn.NewFromConfig(cfg, func(o *sfn.Options) {
		o.Retryer = retry.NewStandard(func(o *retry.StandardOptions) {
			o.MaxAttempts = dopts.maxRetries
			o.MaxBackoff = dopts.interval
		})
	})
	return &client{
		sfn: cli,
		arn: aws.String(params.StateMachineARN),
	}
}

func (c *client) StartExecution(ctx context.Context, input interface{}) error {
	buf, err := json.Marshal(input)
	if err != nil {
		return err
	}
	in := &sfn.StartExecutionInput{
		StateMachineArn: c.arn,
		Input:           aws.String(string(buf)),
	}
	_, err = c.sfn.StartExecution(ctx, in)
	return err
}
