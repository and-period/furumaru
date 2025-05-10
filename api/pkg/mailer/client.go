//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package mailer

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/sendgrid/sendgrid-go"
	"go.uber.org/zap"
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
	logger      *zap.Logger
	fromName    string
	fromAddress string
	templateMap map[string]string
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

// NewClient - メール送信用クライアントの生成
func NewClient(params *Params, opts ...Option) Client {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &client{
		now:         jst.Now,
		client:      sendgrid.NewSendClient(params.APIKey),
		logger:      dopts.logger,
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
	c.logger.Error("Failed to send mail", zap.Error(e))

	switch {
	case errors.Is(e, context.Canceled):
		return fmt.Errorf("%w: %s", ErrCanceled, e.Error())
	case errors.Is(e, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", ErrTimeout, e.Error())
	}

	err, ok := e.(*SendGridError)
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
	case http.StatusRequestEntityTooLarge:
		return fmt.Errorf("%w: %s", ErrPayloadTooLong, err.Error())
	case http.StatusNotFound:
		return fmt.Errorf("%w: %s", ErrNotFound, err.Error())
	case http.StatusInternalServerError:
		return fmt.Errorf("%w: %s", ErrInternal, err.Error())
	case http.StatusBadGateway:
		return fmt.Errorf("%w: %s", ErrUnavailable, err.Error())
	case http.StatusGatewayTimeout:
		return fmt.Errorf("%w: %s", ErrTimeout, err.Error())
	default:
		return fmt.Errorf("%w: %s", ErrUnknown, err.Error())
	}
}
