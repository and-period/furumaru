//go:generate mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package line

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

var (
	ErrInvalidArgument   = errors.New("line: invalid argument")
	ErrUnauthenticated   = errors.New("line: unauthenticated")
	ErrPermissionDenied  = errors.New("line: permission denied")
	ErrPayloadTooLong    = errors.New("line: payload too long")
	ErrNotFound          = errors.New("line: not found")
	ErrAlreadyExists     = errors.New("line: already exists")
	ErrResourceExhausted = errors.New("line: resource exhausted")
	ErrInternal          = errors.New("line: internal")
	ErrUnavailable       = errors.New("line: unavailable")
	ErrTimeout           = errors.New("line: timeout")
	ErrUnknown           = errors.New("line: unknown")
)

type Client interface {
	PushMessage(ctx context.Context, messages ...linebot.SendingMessage) error
}

type Params struct {
	RoomID string
	Token  string
	Secret string
}

type client struct {
	now    func() time.Time
	client *linebot.Client
	logger *zap.Logger
	roomID string
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

// NewClient - LINE API接続用クライアントの生成
func NewClient(params *Params, opts ...Option) (Client, error) {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	bot, err := linebot.New(params.Secret, params.Token)
	if err != nil {
		return nil, err
	}
	return &client{
		now:    jst.Now,
		client: bot,
		logger: dopts.logger,
	}, nil
}

func (c *client) PushMessage(ctx context.Context, messages ...linebot.SendingMessage) error {
	_, err := c.client.PushMessage(c.roomID, messages...).Do()
	return c.lineError(err)
}

func (c *client) lineError(e error) error {
	if e == nil {
		return nil
	}
	c.logger.Error("Failed to send line api", zap.Error(e))

	err, ok := e.(*linebot.APIError)
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknown, e.Error())
	}

	switch err.Code {
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", ErrInvalidArgument, err.Error())
	case http.StatusUnauthorized:
		return fmt.Errorf("%w: %s", ErrUnauthenticated, err.Error())
	case http.StatusForbidden:
		return fmt.Errorf("%w: %s", ErrPermissionDenied, err.Error())
	case http.StatusNotFound:
		return fmt.Errorf("%w: %s", ErrNotFound, err.Error())
	case http.StatusConflict:
		return fmt.Errorf("%w: %s", ErrAlreadyExists, err.Error())
	case http.StatusRequestEntityTooLarge:
		return fmt.Errorf("%w: %s", ErrPayloadTooLong, err.Error())
	case http.StatusTooManyRequests:
		return fmt.Errorf("%w: %s", ErrResourceExhausted, err.Error())
	case http.StatusInternalServerError:
		return fmt.Errorf("%w: %s", ErrInvalidArgument, err.Error())
	case http.StatusBadGateway:
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	case http.StatusGatewayTimeout:
		return fmt.Errorf("%w: %s", ErrTimeout, err.Error())
	default:
		return fmt.Errorf("%w: %s", ErrUnknown, err.Error())
	}
}
