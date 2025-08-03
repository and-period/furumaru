//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package ivs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/ivs"
	"github.com/aws/aws-sdk-go-v2/service/ivs/types"
)

type Client interface {
	// ######################
	// チャンネル関連
	// ######################
	// チャンネル作成
	CreateChannel(ctx context.Context, params *CreateChannelParams) (*ivs.CreateChannelOutput, error)
	// チャンネル取得
	GetChannel(ctx context.Context, params *GetChannelParams) (*types.Channel, error)
	// チャンネル削除
	DeleteChannel(ctx context.Context, params *DeleteChannelParams) error
	// ######################
	// ストリーム関連
	// ######################
	// ストリーム取得
	GetStream(ctx context.Context, params *GetStreamParams) (*types.Stream, error)
	// ストリームキー取得
	GetStreamKey(ctx context.Context, params *GetStreamKeyParams) (*types.StreamKey, error)
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
	RecordingConfigurationArn string
}

type client struct {
	ivs                       *ivs.Client
	recordingConfigurationArn *string
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

func NewClient(cfg aws.Config, params *Params, opts ...Option) Client {
	dopts := &options{
		maxRetries: retry.DefaultMaxAttempts,
		interval:   retry.DefaultMaxBackoff,
	}
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
		ivs:                       cli,
		recordingConfigurationArn: aws.String(params.RecordingConfigurationArn),
	}
}

func (c *client) streamError(err error) error {
	if err == nil {
		return nil
	}
	slog.Debug("Failed to amazon ivs api", log.Error(err))

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
