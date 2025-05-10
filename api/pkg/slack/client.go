//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package slack

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

var (
	ErrInvalidArgument = errors.New("slack: invalid argument")
	ErrUnauthenticated = errors.New("slack: unauthenticated")
	ErrInternal        = errors.New("slack: internal")
	ErrCanceled        = errors.New("slack: canceled")
	ErrTimeout         = errors.New("slack: timeout")
	ErrUnknown         = errors.New("slack: unknown")
)

type Client interface {
	SendMessage(ctx context.Context, options ...slack.MsgOption) error
}

type Params struct {
	Token     string
	ChannelID string
}

type client struct {
	now       func() time.Time
	client    *slack.Client
	logger    *zap.Logger
	channelID string
}

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func NewClient(params *Params, opts ...Option) Client {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &client{
		now:       jst.Now,
		client:    slack.New(params.Token),
		logger:    dopts.logger,
		channelID: params.ChannelID,
	}
}

func (c *client) SendMessage(ctx context.Context, options ...slack.MsgOption) error {
	//nolint:dogsled
	_, _, _, err := c.client.SendMessageContext(ctx, c.channelID, options...)
	return c.slackError(err)
}

func (c *client) slackError(err error) error {
	if err == nil {
		return nil
	}
	c.logger.Error("Failed to send slack api", zap.Error(err))

	switch {
	case errors.Is(err, slack.ErrParametersMissing),
		errors.Is(err, slack.ErrBlockIDNotUnique),
		errors.Is(err, slack.ErrInvalidConfiguration),
		errors.Is(err, slack.ErrMissingHeaders):
		return fmt.Errorf("%w: %s", ErrInvalidArgument, err.Error())
	case errors.Is(err, slack.ErrExpiredTimestamp):
		return fmt.Errorf("%w: %s", ErrUnauthenticated, err.Error())
	case errors.Is(err, context.Canceled),
		errors.Is(err, slack.ErrAlreadyDisconnected),
		errors.Is(err, slack.ErrRTMDisconnected),
		errors.Is(err, slack.ErrRTMGoodbye):
		return fmt.Errorf("%w: %s", ErrCanceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", ErrTimeout, err.Error())
	case errors.Is(err, slack.ErrRTMDeadman):
		return fmt.Errorf("%w: %s", ErrInternal, err.Error())
	default:
		return fmt.Errorf("%w: %s", ErrUnknown, err.Error())
	}
}
