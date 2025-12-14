package komoju

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/and-period/furumaru/api/pkg/log"
)

const defaultMaxRetries = 3

type options struct {
	maxRetries int64
	debugMode  bool
}

type Option func(*options)

func WithMaxRetries(maxRetries int64) Option {
	return func(opts *options) {
		opts.maxRetries = maxRetries
	}
}

func WithDebugMode(enable bool) Option {
	return func(opts *options) {
		opts.debugMode = enable
	}
}

type transport struct {
	base http.RoundTripper
	opts *options
}

var _ http.RoundTripper = (*transport)(nil)

func (t *transport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	if t.opts.debugMode {
		in, _ := httputil.DumpRequest(req, true)
		defer func() {
			var out []byte
			if res != nil {
				out, _ = httputil.DumpResponse(res, true)
			}
			slog.Debug("Send komoju request", slog.String("input", string(in)), slog.String("output", string(out)))
		}()
	}
	res, err = t.base.RoundTrip(req)
	return res, err
}

type APIClient struct {
	client *http.Client
	opts   *options
	apiKey string
}

func NewAPIClient(client *http.Client, basicID, secret string, opts ...Option) *APIClient {
	dopts := &options{
		maxRetries: defaultMaxRetries,
		debugMode:  false,
	}
	for i := range opts {
		opts[i](dopts)
	}
	base := client.Transport
	if base == nil {
		base = http.DefaultTransport
	}
	client.Transport = &transport{
		base: base,
		opts: dopts,
	}
	auth := fmt.Sprintf("%s:%s", basicID, secret)
	return &APIClient{
		client: client,
		opts:   dopts,
		apiKey: base64.StdEncoding.EncodeToString([]byte(auth)),
	}
}

func (c *APIClient) Do(ctx context.Context, params *APIParams, res interface{}) (err error) {
	fn := func() error {
		return c.do(ctx, params, res)
	}
	retry := backoff.NewExponentialBackoff(c.opts.maxRetries)
	return backoff.Retry(ctx, retry, fn, backoff.WithRetryablel(c.isRetryable))
}

func (c *APIClient) do(ctx context.Context, params *APIParams, out interface{}) error {
	res, err := c.request(ctx, params)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer c.closeResponseBody(res)
	if err := c.statusCheck(res, params); err != nil {
		return err
	}
	return c.bind(out, res)
}

func (c *APIClient) isRetryable(err error) bool {
	return errors.Is(err, ErrTooManyRequest) || errors.Is(err, ErrBadGateway) || errors.Is(err, ErrGatewayTimeout)
}

type ErrorResponse struct {
	Data *ErrorData `json:"error"`
}

type ErrorData struct {
	Param   string `json:"param"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (c *APIClient) statusCheck(res *http.Response, params *APIParams) error {
	switch res.StatusCode {
	case http.StatusOK, http.StatusAccepted, http.StatusNoContent:
		return nil // 正常系
	case http.StatusTooManyRequests:
		return ErrTooManyRequest
	case http.StatusBadGateway:
		return ErrBadGateway
	case http.StatusGatewayTimeout:
		return ErrGatewayTimeout
	}
	out := &ErrorResponse{}
	if err := c.bind(out, res); err != nil {
		return err
	}
	slog.Warn("Received komoju error",
		slog.Int("status", res.StatusCode),
		slog.String("method", res.Request.Method),
		slog.String("path", res.Request.URL.Path),
		slog.Any("input", params.Body),
		slog.String("output.param", out.Data.Param),
		slog.String("output.code", out.Data.Code),
		slog.String("output.detail", out.Data.Message),
	)
	return &Error{
		Method:  params.Method,
		Route:   params.Path,
		Status:  res.StatusCode,
		Code:    ErrCode(out.Data.Code),
		Message: out.Data.Message,
	}
}

func (c *APIClient) request(ctx context.Context, params *APIParams) (*http.Response, error) {
	req, err := params.newHTTPRequest(ctx)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", c.apiKey))
	res, err := c.client.Do(req)
	if err == nil {
		return res, nil
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return nil, fmt.Errorf("komoju: failed api request: %w", err)
}

func (c *APIClient) bind(out interface{}, res *http.Response) error {
	if out == nil {
		return nil
	}
	err := json.NewDecoder(res.Body).Decode(out)
	if err == nil {
		return nil
	}
	body, _ := io.ReadAll(res.Body)
	slog.Error("Failed to decode komoju response body",
		slog.Int("status", res.StatusCode),
		slog.String("method", res.Request.Method),
		slog.String("path", res.Request.URL.Path),
		slog.String("body", string(body)),
		log.Error(err),
	)
	return fmt.Errorf("komoju: failed to decode body: %w", err)
}

func (c *APIClient) closeResponseBody(res *http.Response) error {
	//nolint:errcheck
	io.Copy(io.Discard, res.Body)
	return res.Body.Close()
}

type APIParams struct {
	Host   string
	Method string
	Path   string
	Params []interface{}
	Body   interface{}
}

func (p *APIParams) newHTTPRequest(ctx context.Context) (*http.Request, error) {
	u, err := url.ParseRequestURI(p.Host + fmt.Sprintf(p.Path, p.Params...))
	if err != nil {
		return nil, fmt.Errorf("komoju: failed to parse request uri: %w", err)
	}
	if p.Body == nil {
		return http.NewRequestWithContext(ctx, p.Method, u.String(), nil)
	}
	body, err := json.Marshal(p.Body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, p.Method, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("komoju: failed to new request: %w", err)
	}
	return req, nil
}
