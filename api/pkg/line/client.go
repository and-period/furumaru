//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
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
	ErrCanceled          = errors.New("cognito: canceled")
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
		roomID: params.RoomID,
	}, nil
}

func (c *client) PushMessage(_ context.Context, messages ...linebot.SendingMessage) error {
	_, err := c.client.PushMessage(c.roomID, messages...).Do()
	return c.lineError(err)
}

func (c *client) lineError(e error) error {
	if e == nil {
		return nil
	}
	c.logger.Error("Failed to send line api", zap.Error(e))

	switch {
	case errors.Is(e, context.Canceled):
		return fmt.Errorf("%w: %s", ErrCanceled, e.Error())
	case errors.Is(e, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", ErrTimeout, e.Error())
	}

	var aerr *linebot.APIError
	if !errors.As(e, &aerr) {
		return fmt.Errorf("%w: %s", ErrUnknown, e.Error())
	}

	switch aerr.Code {
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", ErrInvalidArgument, aerr.Error())
	case http.StatusUnauthorized:
		return fmt.Errorf("%w: %s", ErrUnauthenticated, aerr.Error())
	case http.StatusForbidden:
		return fmt.Errorf("%w: %s", ErrPermissionDenied, aerr.Error())
	case http.StatusNotFound:
		return fmt.Errorf("%w: %s", ErrNotFound, aerr.Error())
	case http.StatusConflict:
		return fmt.Errorf("%w: %s", ErrAlreadyExists, aerr.Error())
	case http.StatusRequestEntityTooLarge:
		return fmt.Errorf("%w: %s", ErrPayloadTooLong, aerr.Error())
	case http.StatusTooManyRequests:
		return fmt.Errorf("%w: %s", ErrResourceExhausted, aerr.Error())
	case http.StatusInternalServerError:
		return fmt.Errorf("%w: %s", ErrInternal, aerr.Error())
	case http.StatusBadGateway:
		return fmt.Errorf("%w: %s", ErrUnavailable, aerr.Error())
	case http.StatusGatewayTimeout:
		return fmt.Errorf("%w: %s", ErrTimeout, aerr.Error())
	default:
		return fmt.Errorf("%w: %s", ErrUnknown, aerr.Error())
	}
}
