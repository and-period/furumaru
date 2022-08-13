//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../../mock/pkg/firebase/$GOPACKAGE/$GOFILE
package messaging

import (
	"context"
	"errors"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/and-period/furumaru/api/pkg/backoff"
	"go.uber.org/zap"
)

var (
	ErrInvalidArgument   = errors.New("messaging: invalid argument")
	ErrUnauthenticated   = errors.New("messaging: unauthenticated")
	ErrInternal          = errors.New("messaging: internal")
	ErrCanceled          = errors.New("messaging: canceled")
	ErrNotFound          = errors.New("messaging: not found")
	ErrResourceExhausted = errors.New("messaging: resource exhausted")
	ErrUnavailable       = errors.New("messaging: unavailable")
	ErrTimeout           = errors.New("messaging: timeout")
	ErrUnknown           = errors.New("messaging: unknown")
)

type Client interface {
	// プッシュ通知 (単一宛先)
	Send(ctx context.Context, notification *Notification, token string) error
	// プッシュ通知 (複数宛先)
	MultiSend(ctx context.Context, notification *Notification, tokens ...string) (int64, int64, error)
}

type client struct {
	messageing *messaging.Client
	logger     *zap.Logger
	maxRetries int64
}

type options struct {
	logger     *zap.Logger
	maxRetries int64
}

type Option func(opts *options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithMaxRetries(maxRetries int64) Option {
	return func(opts *options) {
		opts.maxRetries = maxRetries
	}
}

func NewClient(ctx context.Context, app *firebase.App, opts ...Option) (Client, error) {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	messaging, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return &client{
		messageing: messaging,
		logger:     dopts.logger,
	}, nil
}

func (c *client) do(ctx context.Context, fn func() error) error {
	retry := backoff.NewExponentialBackoff(c.maxRetries)
	return backoff.Retry(ctx, retry, fn, backoff.WithRetryablel(c.isRetryable))
}

func (c *client) isRetryable(err error) bool {
	switch {
	case
		messaging.IsMessageRateExceeded(err),
		messaging.IsQuotaExceeded(err),
		messaging.IsMessageRateExceeded(err),
		messaging.IsQuotaExceeded(err),
		messaging.IsUnavailable(err):
		return true
	default:
		return false
	}
}

func (c *client) sendError(err error) error {
	if err == nil {
		return nil
	}
	c.logger.Debug("Failed to firebase cloud messaging api", zap.Error(err))

	switch {
	// For Context
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("%w: %s", ErrCanceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", ErrTimeout, err.Error())
	// For Firebase Cloud Messaging
	case messaging.IsInvalidArgument(err):
		return fmt.Errorf("%w: %s", ErrInvalidArgument, err.Error())
	case
		messaging.IsInvalidAPNSCredentials(err),
		messaging.IsThirdPartyAuthError(err),
		messaging.IsMismatchedCredential(err):
		return fmt.Errorf("%w: %s", ErrInternal, err.Error())
	case
		messaging.IsRegistrationTokenNotRegistered(err),
		messaging.IsUnregistered(err):
		return fmt.Errorf("%w: %s", ErrNotFound, err.Error())
	case
		messaging.IsMessageRateExceeded(err),
		messaging.IsQuotaExceeded(err):
		return fmt.Errorf("%w: %s", ErrResourceExhausted, err.Error())
	case messaging.IsUnavailable(err):
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	case messaging.IsInternal(err):
		return fmt.Errorf("%w: %s", ErrInternal, err.Error())
	default:
		return fmt.Errorf("%w: %s", ErrUnknown, err.Error())
	}
}
