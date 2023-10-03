package session

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/and-period/furumaru/api/internal/store/komoju"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testResponse struct {
	body interface{}
	err  error
	code komoju.ErrCode
}

type handler func(w http.ResponseWriter, r *http.Request)

type sessionClientCaller func(ctx context.Context, client komoju.Session) (interface{}, error)

func testSessionClient(handler handler, expect *testResponse, fn sessionClientCaller) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()
		u, err := url.Parse(ts.URL)
		require.NoError(t, err)
		logger := zap.NewNop()
		host := url.URL{
			Scheme: u.Scheme,
			Host:   u.Host,
		}
		require.NoError(t, err)
		params := &SessionParams{
			Host:         host.String(),
			Logger:       logger,
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		}
		client := NewSessionClient(ts.Client(), params, komoju.WithLogger(logger))
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		res, err := fn(ctx, client)
		if err == nil {
			assert.NoError(t, err)
			assert.Equal(t, expect.body, res)
			return
		}
		e := &komoju.Error{}
		if !errors.As(err, &e) {
			assert.ErrorIs(t, err, expect.err)
			return
		}
		assert.Equal(t, expect.code, e.Code)
	}
}

func TestSessionClient(t *testing.T) {
	t.Parallel()
	client := &http.Client{}
	logger := zap.NewNop()
	params := &SessionParams{
		Logger:       logger,
		Host:         "http://example.com",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	}
	expect := &sessionClient{
		client: komoju.NewAPIClient(client, "client-id", "client-secret", komoju.WithLogger(logger), komoju.WithMaxRetries(1)),
		logger: logger,
		host:   "http://example.com",
	}
	actual := NewSessionClient(client, params, komoju.WithLogger(logger), komoju.WithMaxRetries(1)).(*sessionClient)
	assert.Equal(t, expect, actual)
}

func TestSession_Show(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		handler   func(w http.ResponseWriter, r *http.Request)
		sessionID string
		expect    *testResponse
	}{
		{
			name:      "not implemented",
			handler:   func(w http.ResponseWriter, r *http.Request) {},
			sessionID: "session-id",
			expect: &testResponse{
				body: nil,
				err:  komoju.ErrNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testSessionClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.Show(ctx, tt.sessionID)
			}))
		})
	}
}

func TestSession_Create(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.CreateSessionParams
		expect  *testResponse
	}{
		{
			name:    "not implemented",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			params:  &komoju.CreateSessionParams{},
			expect: &testResponse{
				body: nil,
				err:  komoju.ErrNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testSessionClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.Create(ctx, tt.params)
			}))
		})
	}
}

func TestSession_Cancel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		handler   func(w http.ResponseWriter, r *http.Request)
		sessionID string
		expect    *testResponse
	}{
		{
			name:      "not implemented",
			handler:   func(w http.ResponseWriter, r *http.Request) {},
			sessionID: "session-id",
			expect: &testResponse{
				body: nil,
				err:  komoju.ErrNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testSessionClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.Cancel(ctx, tt.sessionID)
			}))
		})
	}
}

func TestSession_ExecuteCreditCard(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.ExecuteCreditCardParams
		expect  *testResponse
	}{
		{
			name:    "not implemented",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			params:  &komoju.ExecuteCreditCardParams{},
			expect: &testResponse{
				body: nil,
				err:  komoju.ErrNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testSessionClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.ExecuteCreditCard(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecuteBankTransfer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.ExecuteBankTransferParams
		expect  *testResponse
	}{
		{
			name:    "not implemented",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			params:  &komoju.ExecuteBankTransferParams{},
			expect: &testResponse{
				body: nil,
				err:  komoju.ErrNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testSessionClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.ExecuteBankTransfer(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecuteKonbini(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.ExecuteKonbiniParams
		expect  *testResponse
	}{
		{
			name:    "not implemented",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			params:  &komoju.ExecuteKonbiniParams{},
			expect: &testResponse{
				body: nil,
				err:  komoju.ErrNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testSessionClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.ExecuteKonbini(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecutePayPay(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.ExecutePayPayParams
		expect  *testResponse
	}{
		{
			name:    "not implemented",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			params:  &komoju.ExecutePayPayParams{},
			expect: &testResponse{
				body: nil,
				err:  komoju.ErrNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testSessionClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.ExecutePayPay(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecuteLinePay(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.ExecuteLinePayParams
		expect  *testResponse
	}{
		{
			name:    "not implemented",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			params:  &komoju.ExecuteLinePayParams{},
			expect: &testResponse{
				body: nil,
				err:  komoju.ErrNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testSessionClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.ExecuteLinePay(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecuteMerpay(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.ExecuteMerpayParams
		expect  *testResponse
	}{
		{
			name:    "not implemented",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			params:  &komoju.ExecuteMerpayParams{},
			expect: &testResponse{
				body: nil,
				err:  komoju.ErrNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testSessionClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.ExecuteMerpay(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecuteRakutenPay(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.ExecuteRakutenPayParams
		expect  *testResponse
	}{
		{
			name:    "not implemented",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			params:  &komoju.ExecuteRakutenPayParams{},
			expect: &testResponse{
				body: nil,
				err:  komoju.ErrNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testSessionClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.ExecuteRakutenPay(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecuteAUPay(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.ExecuteAUPayParams
		expect  *testResponse
	}{
		{
			name:    "not implemented",
			handler: func(w http.ResponseWriter, r *http.Request) {},
			params:  &komoju.ExecuteAUPayParams{},
			expect: &testResponse{
				body: nil,
				err:  komoju.ErrNotImplemented,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testSessionClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.ExecuteAUPay(ctx, tt.params)
			}))
		})
	}
}
