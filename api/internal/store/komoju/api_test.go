package komoju

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testAPIRequest struct {
	Message string `json:"message"`
}

type testAPIResponse struct {
	Message string `json:"message"`
}

func TestAPIClient_Do(t *testing.T) {
	t.Parallel()
	opts := []Option{
		WithLogger(zap.NewNop()),
		WithMaxRetries(1),
		WithDebugMode(true),
	}
	client := NewAPIClient(&http.Client{}, "basic-id", "secret", opts...)
	params := &APIParams{
		Host:   "example.com",
		Path:   "/hoge",
		Method: http.MethodPost,
		Body:   &testAPIRequest{},
	}
	err := client.Do(t.Context(), params, nil)
	require.Error(t, err)
}

func TestAPIClient_do(t *testing.T) {
	t.Parallel()
	client := NewAPIClient(&http.Client{}, "basic-id", "secret")
	params := &APIParams{
		Host:   "example.com",
		Path:   "/hoge",
		Method: http.MethodPost,
		Body:   &testAPIRequest{},
	}
	err := client.do(t.Context(), params, nil)
	require.Error(t, err)
}

func TestAPIClient_isRetryable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		err    error
		expect bool
	}{
		{
			name:   "non error",
			err:    nil,
			expect: false,
		},
		{
			name:   "too many request",
			err:    ErrTooManyRequest,
			expect: true,
		},
		{
			name:   "bad gateway",
			err:    ErrBadGateway,
			expect: true,
		},
		{
			name:   "gateway timeout",
			err:    ErrGatewayTimeout,
			expect: true,
		},
		{
			name:   "other error",
			err:    assert.AnError,
			expect: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := NewAPIClient(&http.Client{}, "basic-id", "secret")
			actual := client.isRetryable(tt.err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAPIClient_statusCheck(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status int
		body   interface{}
		expect error
	}{
		{
			name:   "ok",
			status: http.StatusOK,
			expect: nil,
		},
		{
			name:   "accepted",
			status: http.StatusAccepted,
			expect: nil,
		},
		{
			name:   "accepted",
			status: http.StatusAccepted,
			expect: nil,
		},
		{
			name:   "no content",
			status: http.StatusNoContent,
			expect: nil,
		},
		{
			name:   "too many request",
			status: http.StatusTooManyRequests,
			expect: ErrTooManyRequest,
		},
		{
			name:   "bad gateway",
			status: http.StatusBadGateway,
			expect: ErrBadGateway,
		},
		{
			name:   "gateway timeout",
			status: http.StatusGatewayTimeout,
			expect: ErrGatewayTimeout,
		},
		{
			name:   "komoju error",
			status: http.StatusUnauthorized,
			body: map[string]map[string]string{
				"error": {
					"message": "A required parameter (amount) is missing",
					"code":    "missing_parameter",
					"param":   "amount",
				},
			},
			expect: &Error{
				Method:  http.MethodPost,
				Route:   "/hoge",
				Status:  401,
				Code:    "missing_parameter",
				Message: "A required parameter (amount) is missing",
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := NewAPIClient(&http.Client{}, "basic-id", "secret")
			var buf []byte
			if tt.body != nil {
				buf, _ = json.Marshal(tt.body)
			}
			params := &APIParams{
				Method: http.MethodPost,
				Path:   "/hoge",
			}
			res := &http.Response{
				Request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/hoge"},
				},
				StatusCode: tt.status,
				Body:       io.NopCloser(bytes.NewBuffer(buf)),
			}
			actual := client.statusCheck(res, params)
			assert.Equal(t, actual, tt.expect)
		})
	}
}

func TestAPIClient_request(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	client := NewAPIClient(&http.Client{}, "basic-id", "secret")
	params := &APIParams{
		Host:   "example.com",
		Path:   "/hoge",
		Method: http.MethodPost,
		Body:   &testAPIRequest{},
	}
	res, err := client.request(ctx, params)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestAPIClient_bind(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		body   interface{}
		out    interface{}
		hasErr bool
	}{
		{
			name: "success",
			body: map[string]map[string]string{
				"error": {
					"message": "A required parameter (amount) is missing",
					"code":    "missing_parameter",
					"param":   "amount",
				},
			},
			out:    &ErrorResponse{},
			hasErr: false,
		},
		{
			name:   "success without output",
			body:   "",
			out:    nil,
			hasErr: false,
		},
		{
			name:   "failed to decode",
			body:   "",
			out:    &testAPIResponse{},
			hasErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := NewAPIClient(&http.Client{}, "basic-id", "secret")
			var buf []byte
			if tt.body != nil {
				buf, _ = json.Marshal(tt.body)
			}
			res := &http.Response{
				Request: &http.Request{
					Method: http.MethodPost,
					URL:    &url.URL{Path: "/hoge"},
				},
				Body: io.NopCloser(bytes.NewBuffer(buf)),
			}
			err := client.bind(tt.out, res)
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}
