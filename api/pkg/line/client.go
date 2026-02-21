//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package line

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
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
	ErrCanceled          = errors.New("line: canceled")
	ErrUnavailable       = errors.New("line: unavailable")
	ErrTimeout           = errors.New("line: timeout")
	ErrUnknown           = errors.New("line: unknown")
)

type Client interface {
	PushMessage(ctx context.Context, messages ...messaging_api.MessageInterface) error
}

type Params struct {
	RoomID string
	Token  string
	Secret string
}

type client struct {
	now    func() time.Time
	api    *messaging_api.MessagingApiAPI
	roomID string
}

type options struct{}

type Option func(*options)

// NewClient - LINE API接続用クライアントの生成
func NewClient(params *Params, opts ...Option) (Client, error) {
	dopts := &options{}
	for i := range opts {
		opts[i](dopts)
	}
	api, err := messaging_api.NewMessagingApiAPI(params.Token)
	if err != nil {
		return nil, err
	}
	return &client{
		now:    jst.Now,
		api:    api,
		roomID: params.RoomID,
	}, nil
}

func (c *client) PushMessage(ctx context.Context, messages ...messaging_api.MessageInterface) error {
	req := &messaging_api.PushMessageRequest{
		To:       c.roomID,
		Messages: messages,
	}
	resp, _, err := c.api.WithContext(ctx).PushMessageWithHttpInfo(req, "")
	return c.apiError(err, resp)
}

func (c *client) apiError(e error, resp *http.Response) error {
	if e == nil {
		return nil
	}
	slog.Error("Failed to send line api", log.Error(e))

	switch {
	case errors.Is(e, context.Canceled):
		return fmt.Errorf("%w: %s", ErrCanceled, e.Error())
	case errors.Is(e, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", ErrTimeout, e.Error())
	}

	if resp != nil {
		return c.mapStatusCode(resp.StatusCode, e)
	}

	return fmt.Errorf("%w: %s", ErrUnknown, e.Error())
}

func (c *client) mapStatusCode(code int, e error) error {
	switch code {
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", ErrInvalidArgument, e.Error())
	case http.StatusUnauthorized:
		return fmt.Errorf("%w: %s", ErrUnauthenticated, e.Error())
	case http.StatusForbidden:
		return fmt.Errorf("%w: %s", ErrPermissionDenied, e.Error())
	case http.StatusNotFound:
		return fmt.Errorf("%w: %s", ErrNotFound, e.Error())
	case http.StatusConflict:
		return fmt.Errorf("%w: %s", ErrAlreadyExists, e.Error())
	case http.StatusRequestEntityTooLarge:
		return fmt.Errorf("%w: %s", ErrPayloadTooLong, e.Error())
	case http.StatusTooManyRequests:
		return fmt.Errorf("%w: %s", ErrResourceExhausted, e.Error())
	case http.StatusInternalServerError:
		return fmt.Errorf("%w: %s", ErrInternal, e.Error())
	case http.StatusBadGateway:
		return fmt.Errorf("%w: %s", ErrUnavailable, e.Error())
	case http.StatusGatewayTimeout:
		return fmt.Errorf("%w: %s", ErrTimeout, e.Error())
	default:
		return fmt.Errorf("%w: %s", ErrUnknown, e.Error())
	}
}
