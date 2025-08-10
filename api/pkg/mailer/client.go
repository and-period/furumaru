//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package mailer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/sendgrid/sendgrid-go"
)

var (
	ErrInvalidArgument  = errors.New("mailer: invalid argument")
	ErrUnauthenticated  = errors.New("mailer: unauthenticated")
	ErrPermissionDenied = errors.New("mailer: permission denied")
	ErrPayloadTooLong   = errors.New("mailer: payload too long")
	ErrNotFound         = errors.New("mailer: not found")
	ErrInternal         = errors.New("mailer: internal")
	ErrCanceled         = errors.New("mailer: canceled")
	ErrUnavailable      = errors.New("mailer: unavailable")
	ErrTimeout          = errors.New("mailer: timeout")
	ErrUnknown          = errors.New("mailer: unknown")
)

type Client interface {
	// システム通知 (単一宛先)
	SendFromInfo(ctx context.Context, emailID, toName, toAddress string, substitutions map[string]interface{}) error
	// 任意の送信元からの通知 (複数宛先)
	MultiSend(ctx context.Context, emailID, fromName, fromAddress string, ps []*Personalization) error
	// システム通知 (複数宛先)
	MultiSendFromInfo(ctx context.Context, emailID string, ps []*Personalization) error
}

type Params struct {
	APIKey      string
	FromName    string
	FromAddress string
	TemplateMap map[string]string
}

type client struct {
	now         func() time.Time
	client      *sendgrid.Client
	fromName    string
	fromAddress string
	templateMap map[string]string
}

type options struct{}

type Option func(*options)

// NewClient - メール送信用クライアントの生成
func NewClient(params *Params, opts ...Option) Client {
	dopts := &options{}
	for i := range opts {
		opts[i](dopts)
	}
	return &client{
		now:         jst.Now,
		client:      sendgrid.NewSendClient(params.APIKey),
		fromName:    params.FromName,
		fromAddress: params.FromAddress,
		templateMap: params.TemplateMap,
	}
}

/**
 * private method
 */
func (c *client) mailError(e error) error {
	if e == nil {
		return nil
	}
	slog.Error("Failed to send mail", log.Error(e))

	switch {
	case errors.Is(e, context.Canceled):
		return fmt.Errorf("%w: %s", ErrCanceled, e.Error())
	case errors.Is(e, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", ErrTimeout, e.Error())
	}

	var serr *SendGridError
	if !errors.As(e, &serr) {
		return fmt.Errorf("%w: %s", ErrUnknown, e.Error())
	}

	switch serr.Code {
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", ErrInvalidArgument, serr.Error())
	case http.StatusUnauthorized:
		return fmt.Errorf("%w: %s", ErrUnauthenticated, serr.Error())
	case http.StatusForbidden:
		return fmt.Errorf("%w: %s", ErrPermissionDenied, serr.Error())
	case http.StatusRequestEntityTooLarge:
		return fmt.Errorf("%w: %s", ErrPayloadTooLong, serr.Error())
	case http.StatusNotFound:
		return fmt.Errorf("%w: %s", ErrNotFound, serr.Error())
	case http.StatusInternalServerError:
		return fmt.Errorf("%w: %s", ErrInternal, serr.Error())
	case http.StatusBadGateway:
		return fmt.Errorf("%w: %s", ErrUnavailable, serr.Error())
	case http.StatusGatewayTimeout:
		return fmt.Errorf("%w: %s", ErrTimeout, serr.Error())
	default:
		return fmt.Errorf("%w: %s", ErrUnknown, serr.Error())
	}
}
