//go:generate go tool mockgen -source=$GOFILE -package mock_$GOPACKAGE -destination=./../../mock/pkg/$GOPACKAGE/$GOFILE
package postalcode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Client interface {
	// Search - 住所検索
	Search(ctx context.Context, code string) (*PostalCode, error)
}

var (
	ErrInvalidArgument  = errors.New("postalcode: invalid argument")
	ErrNotFound         = errors.New("postalcode: not found")
	ErrInternal         = errors.New("postalcode: internal")
	ErrCanceled         = errors.New("postalcode: canceled")
	ErrUnavailable      = errors.New("postalcode: unavailable")
	ErrDeadlineExceeded = errors.New("postalcode: deadline exceeded")
	ErrUnknown          = errors.New("postalcode: unknown")
	ErrTimeout          = errors.New("postalcode: timeout")
)

type client struct {
	client     *http.Client
	logger     *zap.Logger
	maxRetries int
}

type options struct {
	maxRetries int
	logger     *zap.Logger
}

type Option func(*options)

func WithMaxRetries(maxRetries int) Option {
	return func(opts *options) {
		opts.maxRetries = maxRetries
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func NewClient(c *http.Client, opts ...Option) Client {
	dopts := &options{
		maxRetries: 3,
		logger:     zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &client{
		client:     c,
		logger:     dopts.logger,
		maxRetries: dopts.maxRetries,
	}
}

func (c *client) do(req *http.Request, out interface{}) error {
	var (
		res *http.Response
		err error
	)
	req.Header.Set("Content-Type", "application/json")
	for retries := 0; retries < c.maxRetries; retries++ {
		res, err = c.client.Do(req)
		if err != nil || !c.retryable(res.StatusCode) {
			break
		}
	}
	if err != nil {
		return err
	}
	defer res.Body.Close() //nolint:errcheck
	if err := c.checkStatus(res); err != nil {
		return err
	}
	return c.bind(out, res)
}

func (c *client) retryable(status int) bool {
	switch status {
	case http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func (c *client) checkStatus(res *http.Response) error {
	if res == nil {
		return ErrUnknown
	}
	status := res.StatusCode
	if status < 400 {
		return nil
	}
	aerr := &apiError{}
	if err := c.bind(aerr, res); err != nil {
		return err
	}
	return aerr
}

func (c *client) bind(out interface{}, res *http.Response) error {
	return json.NewDecoder(res.Body).Decode(out)
}

type apiError struct {
	Status  int    `json:"status"`  // HTTPステータスコード
	Message string `json:"message"` // エラーメッセージ
}

func (e *apiError) Error() string {
	return fmt.Sprintf("postalcode: api requets error. status=%d,message=%s", e.Status, e.Message)
}

func (c *client) newError(err error) error {
	if err == nil {
		return nil
	}
	c.logger.Debug("Failed to postalcode api", zap.Error(err))

	switch {
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("%w: %s", ErrCanceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", ErrTimeout, err.Error())
	}

	e := &apiError{}
	if errors.As(err, &e) {
		switch e.Status {
		case http.StatusBadRequest:
			return fmt.Errorf("%w: %s", ErrInvalidArgument, e.Error())
		case http.StatusNotFound:
			return fmt.Errorf("%w: %s", ErrNotFound, e.Error())
		case http.StatusBadGateway, http.StatusServiceUnavailable:
			return fmt.Errorf("%w: %s", ErrUnavailable, e.Error())
		case http.StatusGatewayTimeout:
			return fmt.Errorf("%w: %s", ErrDeadlineExceeded, e.Error())
		default:
			return fmt.Errorf("%w: %s", ErrInternal, e.Error())
		}
	}

	return fmt.Errorf("%w: %s", ErrUnknown, err.Error())
}
