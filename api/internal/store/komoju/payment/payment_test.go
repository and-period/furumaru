//nolint:paralleltest,tparallel
package payment

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/store/komoju"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testResponse struct {
	body interface{}
	err  error
}

type handler func(w http.ResponseWriter, r *http.Request)

type clientCaller func(ctx context.Context, client komoju.Payment) (interface{}, error)

func testClient(handler handler, expect *testResponse, fn clientCaller) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()
		u, err := url.Parse(ts.URL)
		require.NoError(t, err)
		host := url.URL{
			Scheme: u.Scheme,
			Host:   u.Host,
		}
		require.NoError(t, err)
		params := &Params{
			Host:         host.String(),
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		}
		client := NewClient(ts.Client(), params)
		ctx, cancel := context.WithCancel(t.Context())
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
		assert.Equal(t, expect.err, e)
	}
}

func TestClient(t *testing.T) {
	t.Parallel()
	cli := &http.Client{}
	params := &Params{
		Host:         "http://example.com",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	}
	expect := &client{
		client: komoju.NewAPIClient(cli, "client-id", "client-secret", komoju.WithMaxRetries(1)),
		host:   "http://example.com",
	}
	actual := NewClient(cli, params, komoju.WithMaxRetries(1)).(*client)
	assert.Equal(t, expect, actual)
}

func TestSession_Show(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/payments/payment-id", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
	}
	tests := []struct {
		name      string
		handler   func(w http.ResponseWriter, r *http.Request)
		paymentID string
		expect    *testResponse
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				requestFn(t, r)
				// response header
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				// response body
				body := &komoju.PaymentResponse{
					PaymentInfo: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "ja",
						Status:   komoju.PaymentStatusCaptured,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        string(komoju.PaymentTypePayPay),
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee:    0,
						Total:               1000,
						Currency:            "JPY",
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      0,
						Locale:              "ja",
						Session:             "session-id",
						Refunds:             []*komoju.Refund{},
						RefundRequests:      []*komoju.RefundRequest{},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			paymentID: "payment-id",
			expect: &testResponse{
				body: &komoju.PaymentResponse{
					PaymentInfo: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "ja",
						Status:   komoju.PaymentStatusCaptured,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        string(komoju.PaymentTypePayPay),
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee:    0,
						Total:               1000,
						Currency:            "JPY",
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      0,
						Locale:              "ja",
						Session:             "session-id",
						Refunds:             []*komoju.Refund{},
						RefundRequests:      []*komoju.RefundRequest{},
					},
				},
				err: nil,
			},
		},
		{
			name: "failed to http request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				requestFn(t, r)
				// response header
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				// response body
				body := &komoju.ErrorResponse{
					Data: &komoju.ErrorData{
						Code:    string(komoju.ErrCodeNotFound),
						Message: "The requested resource could not be found.",
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			paymentID: "payment-id",
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodGet,
					Route:   "/api/v1/payments/%s",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Payment) (interface{}, error) {
				return client.Show(ctx, tt.paymentID)
			}))
		})
	}
}

func TestSession_Capture(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/payments/payment-id/capture", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &captureRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &captureRequest{
			Amount: 0,
			Tax:    0,
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name      string
		handler   func(w http.ResponseWriter, r *http.Request)
		paymentID string
		expect    *testResponse
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				requestFn(t, r)
				// response header
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				// response body
				body := &komoju.PaymentResponse{
					PaymentInfo: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "ja",
						Status:   komoju.PaymentStatusCaptured,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        string(komoju.PaymentTypePayPay),
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee:    0,
						Total:               1000,
						Currency:            "JPY",
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      0,
						Locale:              "ja",
						Session:             "session-id",
						Refunds:             []*komoju.Refund{},
						RefundRequests:      []*komoju.RefundRequest{},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			paymentID: "payment-id",
			expect: &testResponse{
				body: &komoju.PaymentResponse{
					PaymentInfo: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "ja",
						Status:   komoju.PaymentStatusCaptured,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        string(komoju.PaymentTypePayPay),
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee:    0,
						Total:               1000,
						Currency:            "JPY",
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      0,
						Locale:              "ja",
						Session:             "session-id",
						Refunds:             []*komoju.Refund{},
						RefundRequests:      []*komoju.RefundRequest{},
					},
				},
				err: nil,
			},
		},
		{
			name: "failed to http request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				requestFn(t, r)
				// response header
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				// response body
				body := &komoju.ErrorResponse{
					Data: &komoju.ErrorData{
						Code:    string(komoju.ErrCodeNotFound),
						Message: "The requested resource could not be found.",
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			paymentID: "payment-id",
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/payments/%s/capture",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Payment) (interface{}, error) {
				return client.Capture(ctx, tt.paymentID)
			}))
		})
	}
}

func TestSession_Cancel(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/payments/payment-id/cancel", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
	}
	tests := []struct {
		name      string
		handler   func(w http.ResponseWriter, r *http.Request)
		paymentID string
		expect    *testResponse
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				requestFn(t, r)
				// response header
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				// response body
				body := &komoju.PaymentResponse{
					PaymentInfo: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "ja",
						Status:   komoju.PaymentStatusCancelled,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        string(komoju.PaymentTypePayPay),
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee:    0,
						Total:               1000,
						Currency:            "JPY",
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      0,
						Locale:              "ja",
						Session:             "session-id",
						Refunds:             []*komoju.Refund{},
						RefundRequests:      []*komoju.RefundRequest{},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			paymentID: "payment-id",
			expect: &testResponse{
				body: &komoju.PaymentResponse{
					PaymentInfo: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "ja",
						Status:   komoju.PaymentStatusCancelled,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        string(komoju.PaymentTypePayPay),
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee:    0,
						Total:               1000,
						Currency:            "JPY",
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      0,
						Locale:              "ja",
						Session:             "session-id",
						Refunds:             []*komoju.Refund{},
						RefundRequests:      []*komoju.RefundRequest{},
					},
				},
				err: nil,
			},
		},
		{
			name: "failed to http request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				requestFn(t, r)
				// response header
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				// response body
				body := &komoju.ErrorResponse{
					Data: &komoju.ErrorData{
						Code:    string(komoju.ErrCodeNotFound),
						Message: "The requested resource could not be found.",
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			paymentID: "payment-id",
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/payments/%s/cancel",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Payment) (interface{}, error) {
				return client.Cancel(ctx, tt.paymentID)
			}))
		})
	}
}

func TestSession_Refund(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/payments/payment-id/refund", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &refundRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &refundRequest{
			Amount:      1000,
			Description: "商品在庫が不足しているため",
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.RefundParams
		expect  *testResponse
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				requestFn(t, r)
				// response header
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				// response body
				body := &komoju.PaymentResponse{
					PaymentInfo: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "ja",
						Status:   komoju.PaymentStatusRefunded,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        string(komoju.PaymentTypePayPay),
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee:    0,
						Total:               1000,
						Currency:            "JPY",
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      1000,
						Locale:              "ja",
						Session:             "session-id",
						Refunds: []*komoju.Refund{
							{
								ID:         "refund-id",
								Resource:   "refund",
								Amount:     1000,
								Currency:   "JPT",
								Payment:    "payment-id",
								CreatedAt:  now,
								Chargeback: false,
							},
						},
						RefundRequests: []*komoju.RefundRequest{},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.RefundParams{
				PaymentID:   "payment-id",
				Amount:      1000,
				Description: "商品在庫が不足しているため",
			},
			expect: &testResponse{
				body: &komoju.PaymentResponse{
					PaymentInfo: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "ja",
						Status:   komoju.PaymentStatusRefunded,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        string(komoju.PaymentTypePayPay),
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee:    0,
						Total:               1000,
						Currency:            "JPY",
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      1000,
						Locale:              "ja",
						Session:             "session-id",
						Refunds: []*komoju.Refund{
							{
								ID:         "refund-id",
								Resource:   "refund",
								Amount:     1000,
								Currency:   "JPT",
								Payment:    "payment-id",
								CreatedAt:  now,
								Chargeback: false,
							},
						},
						RefundRequests: []*komoju.RefundRequest{},
					},
				},
				err: nil,
			},
		},
		{
			name: "failed to http request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				requestFn(t, r)
				// response header
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				// response body
				body := &komoju.ErrorResponse{
					Data: &komoju.ErrorData{
						Code:    string(komoju.ErrCodeNotFound),
						Message: "The requested resource could not be found.",
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.RefundParams{
				PaymentID:   "payment-id",
				Amount:      1000,
				Description: "商品在庫が不足しているため",
			},
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/payments/%s/refund",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Payment) (interface{}, error) {
				return client.Refund(ctx, tt.params)
			}))
		})
	}
}

func TestSession_RefundRequest(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/payments/payment-id/refund_request", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &refundRequestRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &refundRequestRequest{
			Amount:                  1000,
			CustomerName:            "ｱﾝﾄﾞ ﾄﾞｯﾄ",
			BankName:                "三井住友",
			BankCode:                "0009",
			BranchName:              "神保町",
			BranchNumber:            "001",
			AccountType:             "normal",
			AccountNumber:           "1234567",
			IncludePaymentMethodFee: true,
			Description:             "商品の返品があったため",
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.RefundRequestParams
		expect  *testResponse
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				requestFn(t, r)
				// response header
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				// response body
				body := &komoju.PaymentResponse{
					PaymentInfo: &komoju.PaymentInfo{
						ID:              "payment-id",
						Resource:        "ja",
						Status:          komoju.PaymentStatusCaptured,
						Amount:          1000,
						Tax:             0,
						PaymentDeadline: now,
						PaymentDetails: &komoju.PaymentDetails{
							Type:            string(komoju.PaymentTypeKonbini),
							Store:           string(komoju.KonbiniTypeSevenEleven),
							Email:           "test@example.com",
							Receipt:         "receipt",
							InstructionsURL: "http://example.com/instructions",
						},
						PaymentMethodFee:    190,
						Total:               1190,
						Currency:            "JPY",
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      1000,
						Locale:              "ja",
						Session:             "session-id",
						Refunds:             []*komoju.Refund{},
						RefundRequests: []*komoju.RefundRequest{
							{
								ID:            "refund-id",
								Payment:       "payment-id",
								CustomerName:  "ｱﾝﾄﾞ ﾄﾞｯﾄ",
								BankName:      "三井住友",
								BankCode:      "0009",
								BranchName:    "神保町",
								BranchNumber:  "001",
								AccountNumber: "1234567",
							},
						},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.RefundRequestParams{
				PaymentID:               "payment-id",
				Amount:                  1000,
				CustomerName:            "ｱﾝﾄﾞ ﾄﾞｯﾄ",
				BankName:                "三井住友",
				BankCode:                "0009",
				BranchName:              "神保町",
				BranchNumber:            "001",
				AccountType:             "normal",
				AccountNumber:           "1234567",
				IncludePaymentMethodFee: true,
				Description:             "商品の返品があったため",
			},
			expect: &testResponse{
				body: &komoju.PaymentResponse{
					PaymentInfo: &komoju.PaymentInfo{
						ID:              "payment-id",
						Resource:        "ja",
						Status:          komoju.PaymentStatusCaptured,
						Amount:          1000,
						Tax:             0,
						PaymentDeadline: now,
						PaymentDetails: &komoju.PaymentDetails{
							Type:            string(komoju.PaymentTypeKonbini),
							Store:           string(komoju.KonbiniTypeSevenEleven),
							Email:           "test@example.com",
							Receipt:         "receipt",
							InstructionsURL: "http://example.com/instructions",
						},
						PaymentMethodFee:    190,
						Total:               1190,
						Currency:            "JPY",
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      1000,
						Locale:              "ja",
						Session:             "session-id",
						Refunds:             []*komoju.Refund{},
						RefundRequests: []*komoju.RefundRequest{
							{
								ID:            "refund-id",
								Payment:       "payment-id",
								CustomerName:  "ｱﾝﾄﾞ ﾄﾞｯﾄ",
								BankName:      "三井住友",
								BankCode:      "0009",
								BranchName:    "神保町",
								BranchNumber:  "001",
								AccountNumber: "1234567",
							},
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "failed to http request",
			handler: func(w http.ResponseWriter, r *http.Request) {
				requestFn(t, r)
				// response header
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-Type", "application/json;charset=utf-8")
				// response body
				body := &komoju.ErrorResponse{
					Data: &komoju.ErrorData{
						Code:    string(komoju.ErrCodeNotFound),
						Message: "The requested resource could not be found.",
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.RefundRequestParams{
				PaymentID:               "payment-id",
				Amount:                  1000,
				CustomerName:            "ｱﾝﾄﾞ ﾄﾞｯﾄ",
				BankName:                "三井住友",
				BankCode:                "0009",
				BranchName:              "神保町",
				BranchNumber:            "001",
				AccountType:             "normal",
				AccountNumber:           "1234567",
				IncludePaymentMethodFee: true,
				Description:             "商品の返品があったため",
			},
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/payments/%s/refund_request",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Payment) (interface{}, error) {
				return client.RefundRequest(ctx, tt.params)
			}))
		})
	}
}
