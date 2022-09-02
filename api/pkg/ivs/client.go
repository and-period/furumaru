//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package ivs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	ivs "github.com/aws/aws-sdk-go-v2/service/ivs"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"
	"go.uber.org/zap"
)

type Client interface {
	// ######################
	// チャンネル関連
	// ######################
	// チャンネル作成
	CreateChannel(ctx context.Context, params *CreateChannelParams) (*ivs.CreateChannelOutput, error)
}

var (
	ErrInvalidArgument   = errors.New("ivs: invalid argument")
	ErrForbidden         = errors.New("ivs: access denied")
	ErrNotFound          = errors.New("ivs: not found")
	ErrAlreadyExists     = errors.New("ivs: already exists")
	ErrInternal          = errors.New("ivs: internal")
	ErrCanceled          = errors.New("ivs: canceled")
	ErrResourceExhausted = errors.New("ivs: resource exhausted")
	ErrUnknown           = errors.New("ivs: unknown")
	ErrTimeout           = errors.New("ivs: timeout")
)

type Params struct {
	Authorized bool
}

type client struct {
	ivs    *ivs.Client
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

func NewClient(cfg aws.Config, params *Params, opts ...Option) Client {
	dopts := &options{}
	for i := range opts {
		opts[i](dopts)
	}
	cli := ivs.NewFromConfig(cfg, func(o *ivs.Options) {
		o.Retryer = retry.NewStandard(func(o *retry.StandardOptions) {
			o.MaxAttempts = dopts.maxRetries
			o.MaxBackoff = dopts.interval
		})
	})
	return &client{
		ivs:    cli,
		logger: dopts.logger,
	}
}

func (c *client) streamError(err error) error {
	if err == nil {
		return nil
	}
	c.logger.Debug("Failed to cognito api", zap.Error(err))

	switch {
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("%w: %s", ErrCanceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", ErrTimeout, err.Error())
	}

	var (
		vae *types.ValidationException
		ade *types.AccessDeniedException
		coe *types.ConflictException
		ise *types.InternalServerException
		rne *types.ResourceNotFoundException
		sqe *types.ServiceQuotaExceededException
		the *types.ThrottlingException
	)

	switch {
	case errors.As(err, &vae):
		return fmt.Errorf("%w: %s", ErrInvalidArgument, err.Error())
	case errors.As(err, &ade):
		return fmt.Errorf("%w: %s", ErrForbidden, err.Error())
	case errors.As(err, &ise):
		return fmt.Errorf("%w: %s", ErrInternal, err.Error())
	case errors.As(err, &rne):
		return fmt.Errorf("%w: %s", ErrNotFound, err.Error())
	case errors.As(err, &coe):
		return fmt.Errorf("%w: %s", ErrAlreadyExists, err.Error())
	case errors.As(err, &sqe), errors.As(err, &the):
		return fmt.Errorf("%w: %s", ErrResourceExhausted, err.Error())
	default:
		return fmt.Errorf("%w: %s", ErrUnknown, err.Error())
	}
}
