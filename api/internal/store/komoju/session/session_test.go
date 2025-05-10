package session

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
	"go.uber.org/zap"
)

const captureMode = "manual"

type testResponse struct {
	body interface{}
	err  error
}

type handler func(w http.ResponseWriter, r *http.Request)

type clientCaller func(ctx context.Context, client komoju.Session) (interface{}, error)

func testClient(handler handler, expect *testResponse, fn clientCaller) func(t *testing.T) {
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
		params := &Params{
			Host:         host.String(),
			Logger:       logger,
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			CaptureMode:  captureMode,
		}
		client := NewClient(ts.Client(), params, komoju.WithLogger(logger))
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
	logger := zap.NewNop()
	params := &Params{
		Logger:       logger,
		Host:         "http://example.com",
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	}
	expect := &client{
		client: komoju.NewAPIClient(cli, "client-id", "client-secret", komoju.WithLogger(logger), komoju.WithMaxRetries(1)),
		logger: logger,
		host:   "http://example.com",
	}
	actual := NewClient(cli, params, komoju.WithLogger(logger), komoju.WithMaxRetries(1)).(*client)
	assert.Equal(t, expect, actual)
}

func TestSession_Get(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/sessions/session-id", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
	}
	tests := []struct {
		name      string
		handler   func(w http.ResponseWriter, r *http.Request)
		sessionID string
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
				body := &komoju.SessionResponse{
					ID:            "session-id",
					Resource:      "session",
					Mode:          sessionMode,
					Amount:        1000,
					Currency:      defaultCurrency,
					SessionURL:    "http://example.com/sessions/session-id",
					ReturnURL:     "http://example.com/done",
					DefaultLocale: defaultLocale,
					PaymentMethods: []*komoju.PaymentMethod{
						{
							Type:             string(komoju.PaymentTypePayPay),
							Offsite:          false,
							AdditionalFields: []string{},
							Amount:           1000,
							Currency:         defaultCurrency,
							ExchangeRate:     1,
							HashedGateway:    "hash",
						},
					},
					CreatedAt:   now,
					CompletedAt: now,
					Status:      komoju.SessionStatusCompleted,
					Expired:     true,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: sessionMode,
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
						Currency:            defaultCurrency,
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      0,
						Locale:              defaultLocale,
						Session:             "session-id",
						Refunds:             []*komoju.Refund{},
						RefundRequests:      []*komoju.RefundRequest{},
					},
					SecureToken: &komoju.SecureToken{
						ID:                 "secure-token",
						VerificationStatus: "ERRORED",
						Reason:             "3D Secure 2.0 authentication failed. Please verify your information and try again.",
						CreatedAt:          now,
					},
					PaymentData: &komoju.PaymentData{
						Capture: captureMode,
					},
					LineItems: []*komoju.SessionLineItem{
						{
							Description: "芽が出たじゃがいも",
							Amount:      500,
							Quantity:    2,
						},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			sessionID: "session-id",
			expect: &testResponse{
				body: &komoju.SessionResponse{
					ID:            "session-id",
					Resource:      "session",
					Mode:          sessionMode,
					Amount:        1000,
					Currency:      defaultCurrency,
					SessionURL:    "http://example.com/sessions/session-id",
					ReturnURL:     "http://example.com/done",
					DefaultLocale: defaultLocale,
					PaymentMethods: []*komoju.PaymentMethod{
						{
							Type:          string(komoju.PaymentTypePayPay),
							Offsite:       false,
							Amount:        1000,
							Currency:      defaultCurrency,
							ExchangeRate:  1,
							HashedGateway: "hash",
						},
					},
					CreatedAt:   now,
					CompletedAt: now,
					Status:      komoju.SessionStatusCompleted,
					Expired:     true,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: sessionMode,
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
						Currency:            defaultCurrency,
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      0,
						Locale:              defaultLocale,
						Session:             "session-id",
						Refunds:             []*komoju.Refund{},
						RefundRequests:      []*komoju.RefundRequest{},
					},
					SecureToken: &komoju.SecureToken{
						ID:                 "secure-token",
						VerificationStatus: "ERRORED",
						Reason:             "3D Secure 2.0 authentication failed. Please verify your information and try again.",
						CreatedAt:          now,
					},
					PaymentData: &komoju.PaymentData{
						Capture: captureMode,
					},
					LineItems: []*komoju.SessionLineItem{
						{
							Description: "芽が出たじゃがいも",
							Amount:      500,
							Quantity:    2,
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
			sessionID: "session-id",
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodGet,
					Route:   "/api/v1/sessions/%s",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.Get(ctx, tt.sessionID)
			}))
		})
	}
}

func TestSession_Create(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/sessions", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &createSessionRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &createSessionRequest{
			ReturnURL:          "http://example.com/callback",
			Mode:               "payment",
			Amount:             1000,
			Currency:           "JPY",
			Email:              "test@example.com",
			ExpiresInSeconds:   7200,
			ExternalCustomerID: "customer-id",
			PaymentTypes:       []string{"paypay"},
			DefaultLocale:      "ja",
			PaymentData: &createSessionPaymentData{
				Amount:              1000,
				Currency:            "JPY",
				ExternalOrderNumber: "order-id",
				Name:                "&. 利用者",
				NameKana:            "あんどどっと りようしゃ",
				ShippingAddress: &createSessionAddress{
					ZipCode:        "1000014",
					Country:        "日本",
					State:          "東京都",
					City:           "千代田区",
					StreetAddress1: "永田町1-7-1",
					StreetAddress2: "",
				},
				BillingAddress: &createSessionAddress{
					ZipCode:        "1000014",
					Country:        "日本",
					State:          "東京都",
					City:           "千代田区",
					StreetAddress1: "永田町1-7-1",
					StreetAddress2: "",
				},
				Capture: "manual",
			},
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.CreateSessionParams
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
				body := &komoju.SessionResponse{
					ID:            "session-id",
					Resource:      "session",
					Mode:          sessionMode,
					Amount:        1000,
					Currency:      defaultCurrency,
					SessionURL:    "http://example.com/sessions/session-id",
					ReturnURL:     "http://example.com/done",
					DefaultLocale: defaultLocale,
					PaymentMethods: []*komoju.PaymentMethod{
						{
							Type:             string(komoju.PaymentTypePayPay),
							Offsite:          false,
							AdditionalFields: []string{},
							Amount:           1000,
							Currency:         defaultCurrency,
							ExchangeRate:     1,
							HashedGateway:    "hash",
						},
					},
					CreatedAt:   now,
					CompletedAt: now,
					Status:      komoju.SessionStatusCompleted,
					Expired:     true,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: sessionMode,
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
						Currency:            defaultCurrency,
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      0,
						Locale:              defaultLocale,
						Session:             "session-id",
						Refunds:             []*komoju.Refund{},
						RefundRequests:      []*komoju.RefundRequest{},
					},
					PaymentData: &komoju.PaymentData{
						Capture: captureMode,
					},
					LineItems: []*komoju.SessionLineItem{
						{
							Description: "芽が出たじゃがいも",
							Amount:      500,
							Quantity:    2,
						},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.CreateSessionParams{
				OrderID:     "order-id",
				Amount:      1000,
				CallbackURL: "http://example.com/callback",
				PaymentTypes: []komoju.PaymentType{
					komoju.PaymentTypePayPay,
				},
				Customer: &komoju.CreateSessionCustomer{
					ID:       "customer-id",
					Name:     "&. 利用者",
					NameKana: "あんどどっと りようしゃ",
					Email:    "test@example.com",
				},
				BillingAddress: &komoju.CreateSessionAddress{
					ZipCode:      "1000014",
					Prefecture:   "東京都",
					City:         "千代田区",
					AddressLine1: "永田町1-7-1",
					AddressLine2: "",
				},
				ShippingAddress: &komoju.CreateSessionAddress{
					ZipCode:      "1000014",
					Prefecture:   "東京都",
					City:         "千代田区",
					AddressLine1: "永田町1-7-1",
					AddressLine2: "",
				},
			},
			expect: &testResponse{
				body: &komoju.SessionResponse{
					ID:            "session-id",
					Resource:      "session",
					Mode:          sessionMode,
					Amount:        1000,
					Currency:      defaultCurrency,
					SessionURL:    "http://example.com/sessions/session-id",
					ReturnURL:     "http://example.com/done",
					DefaultLocale: defaultLocale,
					PaymentMethods: []*komoju.PaymentMethod{
						{
							Type:          string(komoju.PaymentTypePayPay),
							Offsite:       false,
							Amount:        1000,
							Currency:      defaultCurrency,
							ExchangeRate:  1,
							HashedGateway: "hash",
						},
					},
					CreatedAt:   now,
					CompletedAt: now,
					Status:      komoju.SessionStatusCompleted,
					Expired:     true,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: sessionMode,
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
						Currency:            defaultCurrency,
						CapturedAt:          now,
						ExternalOrderNumber: "order-id",
						CreatedAt:           now,
						AmountRefunded:      0,
						Locale:              defaultLocale,
						Session:             "session-id",
						Refunds:             []*komoju.Refund{},
						RefundRequests:      []*komoju.RefundRequest{},
					},
					PaymentData: &komoju.PaymentData{
						Capture: captureMode,
					},
					LineItems: []*komoju.SessionLineItem{
						{
							Description: "芽が出たじゃがいも",
							Amount:      500,
							Quantity:    2,
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
			params: &komoju.CreateSessionParams{
				OrderID:     "order-id",
				Amount:      1000,
				CallbackURL: "http://example.com/callback",
				PaymentTypes: []komoju.PaymentType{
					komoju.PaymentTypePayPay,
				},
				Customer: &komoju.CreateSessionCustomer{
					ID:       "customer-id",
					Name:     "&. 利用者",
					NameKana: "あんどどっと りようしゃ",
					Email:    "test@example.com",
				},
				BillingAddress: &komoju.CreateSessionAddress{
					ZipCode:      "1000014",
					Prefecture:   "東京都",
					City:         "千代田区",
					AddressLine1: "永田町1-7-1",
					AddressLine2: "",
				},
				ShippingAddress: &komoju.CreateSessionAddress{
					ZipCode:      "1000014",
					Prefecture:   "東京都",
					City:         "千代田区",
					AddressLine1: "永田町1-7-1",
					AddressLine2: "",
				},
			},
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/sessions",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.Create(ctx, tt.params)
			}))
		})
	}
}

func TestSession_Cancel(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/sessions/session-id/cancel", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
	}
	tests := []struct {
		name      string
		handler   func(w http.ResponseWriter, r *http.Request)
		sessionID string
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
				body := &komoju.SessionResponse{
					ID:            "session-id",
					Resource:      "session",
					Mode:          sessionMode,
					Amount:        1000,
					Currency:      defaultCurrency,
					SessionURL:    "http://example.com/sessions/session-id",
					ReturnURL:     "http://example.com/done",
					DefaultLocale: defaultLocale,
					PaymentMethods: []*komoju.PaymentMethod{
						{
							Type:             string(komoju.PaymentTypePayPay),
							Offsite:          false,
							AdditionalFields: []string{},
							Amount:           1000,
							Currency:         defaultCurrency,
							ExchangeRate:     1,
							HashedGateway:    "hash",
						},
					},
					CreatedAt:   now,
					CancelledAt: now,
					Status:      komoju.SessionStatusCancelled,
					Expired:     false,
					PaymentData: &komoju.PaymentData{
						Capture: captureMode,
					},
					LineItems: []*komoju.SessionLineItem{
						{
							Description: "芽が出たじゃがいも",
							Amount:      500,
							Quantity:    2,
						},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			sessionID: "session-id",
			expect: &testResponse{
				body: &komoju.SessionResponse{
					ID:            "session-id",
					Resource:      "session",
					Mode:          sessionMode,
					Amount:        1000,
					Currency:      defaultCurrency,
					SessionURL:    "http://example.com/sessions/session-id",
					ReturnURL:     "http://example.com/done",
					DefaultLocale: defaultLocale,
					PaymentMethods: []*komoju.PaymentMethod{
						{
							Type:          string(komoju.PaymentTypePayPay),
							Offsite:       false,
							Amount:        1000,
							Currency:      defaultCurrency,
							ExchangeRate:  1,
							HashedGateway: "hash",
						},
					},
					CreatedAt:   now,
					CancelledAt: now,
					Status:      komoju.SessionStatusCancelled,
					Expired:     false,
					PaymentData: &komoju.PaymentData{
						Capture: captureMode,
					},
					LineItems: []*komoju.SessionLineItem{
						{
							Description: "芽が出たじゃがいも",
							Amount:      500,
							Quantity:    2,
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
			sessionID: "session-id",
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/sessions/%s/cancel",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.Cancel(ctx, tt.sessionID)
			}))
		})
	}
}

func TestSession_ExecuteCreditCard(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/sessions/session-id/pay", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &orderCreditCardRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &orderCreditCardRequest{
			Capture: "manual",
			PaymentDetails: &creditCardDetails{
				Type:              "credit_card",
				Email:             "test@example.com",
				Number:            "4100000000005000",
				Month:             "12",
				Year:              "2023",
				VerificationValue: "123",
				Name:              "AND TARO",
				// ThreeDSecure:      true,
			},
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.OrderCreditCardParams
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
				body := &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "session",
						Status:   komoju.PaymentStatusPending,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:           "credit_card",
							Email:          "test@example.com",
							Brand:          "visa",
							LastFourDigits: "5000",
							Month:          12,
							Year:           2023,
						},
						PaymentMethodFee:   0,
						Total:              1000,
						Currency:           defaultCurrency,
						CreatedAt:          now,
						AmountRefunded:     0,
						Locale:             defaultLocale,
						Session:            "session-id",
						CustomerFamilyName: "&.",
						CustomerGivenName:  "利用者",
						Refunds:            []*komoju.Refund{},
						RefundRequests:     []*komoju.RefundRequest{},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.OrderCreditCardParams{
				SessionID:         "session-id",
				Number:            "4100000000005000",
				Month:             12,
				Year:              2023,
				VerificationValue: "123",
				Email:             "test@example.com",
				Name:              "AND TARO",
			},
			expect: &testResponse{
				body: &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "session",
						Status:   komoju.PaymentStatusPending,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:           "credit_card",
							Email:          "test@example.com",
							Brand:          "visa",
							LastFourDigits: "5000",
							Month:          12,
							Year:           2023,
						},
						PaymentMethodFee:   0,
						Total:              1000,
						Currency:           defaultCurrency,
						CreatedAt:          now,
						AmountRefunded:     0,
						Locale:             defaultLocale,
						Session:            "session-id",
						CustomerFamilyName: "&.",
						CustomerGivenName:  "利用者",
						Refunds:            []*komoju.Refund{},
						RefundRequests:     []*komoju.RefundRequest{},
					},
				},
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
			params: &komoju.OrderCreditCardParams{
				SessionID:         "session-id",
				Number:            "4100000000005000",
				Month:             12,
				Year:              2023,
				VerificationValue: "123",
				Email:             "test@example.com",
				Name:              "AND TARO",
			},
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/sessions/%s/pay",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.OrderCreditCard(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecuteBankTransfer(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/sessions/session-id/pay", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &orderBankTransferRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &orderBankTransferRequest{
			Capture: "manual",
			PaymentDetails: &bankTransferDetails{
				Type:           "bank_transfer",
				Email:          "test@example.com",
				PhoneNumber:    "09012341234",
				FamilyName:     "&.",
				GivenName:      "利用者",
				FamilyNameKana: "あんどどっと",
				GivenNameKana:  "りようしゃ",
			},
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.OrderBankTransferParams
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
				body := &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:              "payment-id",
						Resource:        "session",
						Status:          komoju.PaymentStatusAuthorized,
						Amount:          1000,
						Tax:             0,
						PaymentDeadline: now.AddDate(0, 7, 0),
						PaymentDetails: &komoju.PaymentDetails{
							Type:              "bank_transfer",
							Email:             "test@example.com",
							OrderID:           "order-id",
							BankName:          "三井住友銀行",
							AccountBranchName: "東京中央支店",
							AccountNumber:     "1234567",
							AccountType:       "普通預金",
							AccountName:       "振込先名",
							InstructionsURL:   "http://example.com/instructions",
							PaymentDeadline:   now.AddDate(0, 7, 0),
						},
						PaymentMethodFee:   0,
						Total:              1000,
						Currency:           defaultCurrency,
						CreatedAt:          now,
						AmountRefunded:     0,
						Locale:             defaultLocale,
						Session:            "session-id",
						CustomerFamilyName: "&.",
						CustomerGivenName:  "利用者",
						Refunds:            []*komoju.Refund{},
						RefundRequests:     []*komoju.RefundRequest{},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.OrderBankTransferParams{
				SessionID:     "session-id",
				Email:         "test@example.com",
				PhoneNumber:   "09012341234",
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
			},
			expect: &testResponse{
				body: &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:              "payment-id",
						Resource:        "session",
						Status:          komoju.PaymentStatusAuthorized,
						Amount:          1000,
						Tax:             0,
						PaymentDeadline: now.AddDate(0, 7, 0),
						PaymentDetails: &komoju.PaymentDetails{
							Type:              "bank_transfer",
							Email:             "test@example.com",
							OrderID:           "order-id",
							BankName:          "三井住友銀行",
							AccountBranchName: "東京中央支店",
							AccountNumber:     "1234567",
							AccountType:       "普通預金",
							AccountName:       "振込先名",
							InstructionsURL:   "http://example.com/instructions",
							PaymentDeadline:   now.AddDate(0, 7, 0),
						},
						PaymentMethodFee:   0,
						Total:              1000,
						Currency:           defaultCurrency,
						CreatedAt:          now,
						AmountRefunded:     0,
						Locale:             defaultLocale,
						Session:            "session-id",
						CustomerFamilyName: "&.",
						CustomerGivenName:  "利用者",
						Refunds:            []*komoju.Refund{},
						RefundRequests:     []*komoju.RefundRequest{},
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
			params: &komoju.OrderBankTransferParams{
				SessionID:     "session-id",
				Email:         "test@example.com",
				PhoneNumber:   "09012341234",
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
			},
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/sessions/%s/pay",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.OrderBankTransfer(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecuteKonbini(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/sessions/session-id/pay", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &orderKonbiniRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &orderKonbiniRequest{
			Capture: "manual",
			PaymentDetails: &konbiniDetails{
				Type:  "konbini",
				Store: "seven-eleven",
				Email: "test@example.com",
			},
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.OrderKonbiniParams
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
				body := &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusCompleted,
					Payment: &komoju.PaymentInfo{
						ID:              "payment-id",
						Resource:        "session",
						Status:          komoju.PaymentStatusAuthorized,
						Amount:          1000,
						Tax:             0,
						PaymentDeadline: now.AddDate(0, 7, 0),
						PaymentDetails: &komoju.PaymentDetails{
							Type:            "konbini",
							Store:           "seven-eleven",
							Email:           "test@example.com",
							Receipt:         "receipt",
							InstructionsURL: "http://example.com/instructions",
						},
						PaymentMethodFee: 190,
						Total:            1190,
						Currency:         defaultCurrency,
						CreatedAt:        now,
						AmountRefunded:   0,
						Locale:           defaultLocale,
						Session:          "session-id",
						Refunds:          []*komoju.Refund{},
						RefundRequests:   []*komoju.RefundRequest{},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.OrderKonbiniParams{
				SessionID: "session-id",
				Store:     komoju.KonbiniTypeSevenEleven,
				Email:     "test@example.com",
			},
			expect: &testResponse{
				body: &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusCompleted,
					Payment: &komoju.PaymentInfo{
						ID:              "payment-id",
						Resource:        "session",
						Status:          komoju.PaymentStatusAuthorized,
						Amount:          1000,
						Tax:             0,
						PaymentDeadline: now.AddDate(0, 7, 0),
						PaymentDetails: &komoju.PaymentDetails{
							Type:            "konbini",
							Store:           "seven-eleven",
							Email:           "test@example.com",
							Receipt:         "receipt",
							InstructionsURL: "http://example.com/instructions",
						},
						PaymentMethodFee: 190,
						Total:            1190,
						Currency:         defaultCurrency,
						CreatedAt:        now,
						AmountRefunded:   0,
						Locale:           defaultLocale,
						Session:          "session-id",
						Refunds:          []*komoju.Refund{},
						RefundRequests:   []*komoju.RefundRequest{},
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
			params: &komoju.OrderKonbiniParams{
				SessionID: "session-id",
				Store:     komoju.KonbiniTypeSevenEleven,
				Email:     "test@example.com",
			},
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/sessions/%s/pay",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.OrderKonbini(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecutePayPay(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/sessions/session-id/pay", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &orderPayPayRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &orderPayPayRequest{
			Capture: "manual",
			PaymentDetails: &paypayDetails{
				Type: "paypay",
			},
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.OrderPayPayParams
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
				body := &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "payment",
						Status:   komoju.PaymentStatusPending,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        "paypay",
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee: 0,
						Total:            1000,
						Currency:         defaultCurrency,
						CreatedAt:        now,
						AmountRefunded:   0,
						Locale:           defaultLocale,
						Session:          "session-id",
						Refunds:          []*komoju.Refund{},
						RefundRequests:   []*komoju.RefundRequest{},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.OrderPayPayParams{
				SessionID: "session-id",
			},
			expect: &testResponse{
				body: &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "payment",
						Status:   komoju.PaymentStatusPending,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        "paypay",
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee: 0,
						Total:            1000,
						Currency:         defaultCurrency,
						CreatedAt:        now,
						AmountRefunded:   0,
						Locale:           defaultLocale,
						Session:          "session-id",
						Refunds:          []*komoju.Refund{},
						RefundRequests:   []*komoju.RefundRequest{},
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
			params: &komoju.OrderPayPayParams{
				SessionID: "session-id",
			},
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/sessions/%s/pay",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.OrderPayPay(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecuteLinePay(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/sessions/session-id/pay", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &orderPayPayRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &orderPayPayRequest{
			Capture: "manual",
			PaymentDetails: &paypayDetails{
				Type: "linepay",
			},
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.OrderLinePayParams
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
				body := &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "payment",
						Status:   komoju.PaymentStatusPending,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:               "linepay",
							Email:              "test@example.com",
							TransactionKey:     "transaction-key",
							RedirectURL:        "http://example.com/redirect",
							PaymentURLApp:      "http://example.com/payment",
							PaymentAccessToken: "access-token",
						},
						PaymentMethodFee: 0,
						Total:            1000,
						Currency:         defaultCurrency,
						CreatedAt:        now,
						AmountRefunded:   0,
						Locale:           defaultLocale,
						Session:          "session-id",
						Refunds:          []*komoju.Refund{},
						RefundRequests:   []*komoju.RefundRequest{},
					},
					AppURL: "http://example.com/app",
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.OrderLinePayParams{
				SessionID: "session-id",
			},
			expect: &testResponse{
				body: &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "payment",
						Status:   komoju.PaymentStatusPending,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:               "linepay",
							Email:              "test@example.com",
							TransactionKey:     "transaction-key",
							RedirectURL:        "http://example.com/redirect",
							PaymentURLApp:      "http://example.com/payment",
							PaymentAccessToken: "access-token",
						},
						PaymentMethodFee: 0,
						Total:            1000,
						Currency:         defaultCurrency,
						CreatedAt:        now,
						AmountRefunded:   0,
						Locale:           defaultLocale,
						Session:          "session-id",
						Refunds:          []*komoju.Refund{},
						RefundRequests:   []*komoju.RefundRequest{},
					},
					AppURL: "http://example.com/app",
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
			params: &komoju.OrderLinePayParams{
				SessionID: "session-id",
			},
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/sessions/%s/pay",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.OrderLinePay(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecuteMerpay(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/sessions/session-id/pay", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &orderMerpayRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &orderMerpayRequest{
			Capture: "manual",
			PaymentDetails: &merpayDetails{
				Type: "merpay",
			},
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.OrderMerpayParams
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
				body := &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "payment",
						Status:   komoju.PaymentStatusPending,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:           "merpay",
							Email:          "test@example.com",
							TransactionKey: "transaction-key",
							RedirectURL:    "http://example.com/redirect",
						},
						PaymentMethodFee: 0,
						Total:            1000,
						Currency:         defaultCurrency,
						CreatedAt:        now,
						AmountRefunded:   0,
						Locale:           defaultLocale,
						Session:          "session-id",
						Refunds:          []*komoju.Refund{},
						RefundRequests:   []*komoju.RefundRequest{},
					},
					AppURL: "http://example.com/app",
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.OrderMerpayParams{
				SessionID: "session-id",
			},
			expect: &testResponse{
				body: &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "payment",
						Status:   komoju.PaymentStatusPending,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:           "merpay",
							Email:          "test@example.com",
							TransactionKey: "transaction-key",
							RedirectURL:    "http://example.com/redirect",
						},
						PaymentMethodFee: 0,
						Total:            1000,
						Currency:         defaultCurrency,
						CreatedAt:        now,
						AmountRefunded:   0,
						Locale:           defaultLocale,
						Session:          "session-id",
						Refunds:          []*komoju.Refund{},
						RefundRequests:   []*komoju.RefundRequest{},
					},
					AppURL: "http://example.com/app",
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
			params: &komoju.OrderMerpayParams{
				SessionID: "session-id",
			},
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/sessions/%s/pay",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.OrderMerpay(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecuteRakutenPay(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/sessions/session-id/pay", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &orderRakutenPayRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &orderRakutenPayRequest{
			Capture: "manual",
			PaymentDetails: &rakutenPayDetails{
				Type: "rakutenpay",
			},
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.OrderRakutenPayParams
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
				body := &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "payment",
						Status:   komoju.PaymentStatusPending,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        "rakutenpay",
							Email:       "test@example.com",
							ChargeKey:   "charge-key",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee: 0,
						Total:            1000,
						Currency:         defaultCurrency,
						CreatedAt:        now,
						AmountRefunded:   0,
						Locale:           defaultLocale,
						Session:          "session-id",
						Refunds:          []*komoju.Refund{},
						RefundRequests:   []*komoju.RefundRequest{},
					},
					AppURL: "http://example.com/app",
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.OrderRakutenPayParams{
				SessionID: "session-id",
			},
			expect: &testResponse{
				body: &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "payment",
						Status:   komoju.PaymentStatusPending,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        "rakutenpay",
							Email:       "test@example.com",
							ChargeKey:   "charge-key",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee: 0,
						Total:            1000,
						Currency:         defaultCurrency,
						CreatedAt:        now,
						AmountRefunded:   0,
						Locale:           defaultLocale,
						Session:          "session-id",
						Refunds:          []*komoju.Refund{},
						RefundRequests:   []*komoju.RefundRequest{},
					},
					AppURL: "http://example.com/app",
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
			params: &komoju.OrderRakutenPayParams{
				SessionID: "session-id",
			},
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/sessions/%s/pay",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.OrderRakutenPay(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecuteAUPay(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/sessions/session-id/pay", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &orderRakutenPayRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &orderRakutenPayRequest{
			Capture: "manual",
			PaymentDetails: &rakutenPayDetails{
				Type: "aupay",
			},
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.OrderAUPayParams
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
				body := &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "payment",
						Status:   komoju.PaymentStatusPending,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        "aupay",
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee: 0,
						Total:            1000,
						Currency:         defaultCurrency,
						CreatedAt:        now,
						AmountRefunded:   0,
						Locale:           defaultLocale,
						Session:          "session-id",
						Refunds:          []*komoju.Refund{},
						RefundRequests:   []*komoju.RefundRequest{},
					},
					AppURL: "http://example.com/app",
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.OrderAUPayParams{
				SessionID: "session-id",
			},
			expect: &testResponse{
				body: &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:       "payment-id",
						Resource: "payment",
						Status:   komoju.PaymentStatusPending,
						Amount:   1000,
						Tax:      0,
						PaymentDetails: &komoju.PaymentDetails{
							Type:        "aupay",
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee: 0,
						Total:            1000,
						Currency:         defaultCurrency,
						CreatedAt:        now,
						AmountRefunded:   0,
						Locale:           defaultLocale,
						Session:          "session-id",
						Refunds:          []*komoju.Refund{},
						RefundRequests:   []*komoju.RefundRequest{},
					},
					AppURL: "http://example.com/app",
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
			params: &komoju.OrderAUPayParams{
				SessionID: "session-id",
			},
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/sessions/%s/pay",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.OrderAUPay(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecutePaidy(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/sessions/session-id/pay", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &orderPaidyRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &orderPaidyRequest{
			Capture: "manual",
			PaymentDetails: &paidyDetails{
				Type:         "paidy",
				CustomerName: "&. 利用者",
				Email:        "test@example.com",
			},
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.OrderPaidyParams
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
				body := &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:              "payment-id",
						Resource:        "session",
						Status:          komoju.PaymentStatusAuthorized,
						Amount:          1000,
						Tax:             0,
						PaymentDeadline: now.AddDate(0, 7, 0),
						PaymentDetails: &komoju.PaymentDetails{
							Type:        "paidy",
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee:   0,
						Total:              1000,
						Currency:           defaultCurrency,
						CreatedAt:          now,
						AmountRefunded:     0,
						Locale:             defaultLocale,
						Session:            "session-id",
						CustomerFamilyName: "&.",
						CustomerGivenName:  "利用者",
						Refunds:            []*komoju.Refund{},
						RefundRequests:     []*komoju.RefundRequest{},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.OrderPaidyParams{
				SessionID: "session-id",
				Name:      "&. 利用者",
				Email:     "test@example.com",
			},
			expect: &testResponse{
				body: &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:              "payment-id",
						Resource:        "session",
						Status:          komoju.PaymentStatusAuthorized,
						Amount:          1000,
						Tax:             0,
						PaymentDeadline: now.AddDate(0, 7, 0),
						PaymentDetails: &komoju.PaymentDetails{
							Type:        "paidy",
							Email:       "test@example.com",
							RedirectURL: "http://example.com/redirect",
						},
						PaymentMethodFee:   0,
						Total:              1000,
						Currency:           defaultCurrency,
						CreatedAt:          now,
						AmountRefunded:     0,
						Locale:             defaultLocale,
						Session:            "session-id",
						CustomerFamilyName: "&.",
						CustomerGivenName:  "利用者",
						Refunds:            []*komoju.Refund{},
						RefundRequests:     []*komoju.RefundRequest{},
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
			params: &komoju.OrderPaidyParams{
				SessionID: "session-id",
				Name:      "&. 利用者",
				Email:     "test@example.com",
			},
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/sessions/%s/pay",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.OrderPaidy(ctx, tt.params)
			}))
		})
	}
}

func TestSession_ExecutePayEasy(t *testing.T) {
	t.Parallel()
	now := time.Date(2023, 10, 5, 18, 30, 0, 0, time.UTC)
	requestFn := func(t *testing.T, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/sessions/session-id/pay", r.URL.Path)
		// request header
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic Y2xpZW50LWlkOmNsaWVudC1zZWNyZXQ=", r.Header.Get("Authorization"))
		// request body
		body := &orderPayEasyRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)
		require.NoError(t, err)
		expect := &orderPayEasyRequest{
			Capture: "manual",
			PaymentDetails: &payEasyDetails{
				Type:           "pay_easy",
				Email:          "test@example.com",
				PhoneNumber:    "09012341234",
				FamilyName:     "&.",
				GivenName:      "利用者",
				FamilyNameKana: "あんどどっと",
				GivenNameKana:  "りようしゃ",
			},
		}
		assert.Equal(t, expect, body)
	}
	tests := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		params  *komoju.OrderPayEasyParams
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
				body := &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:              "payment-id",
						Resource:        "session",
						Status:          komoju.PaymentStatusAuthorized,
						Amount:          1000,
						Tax:             0,
						PaymentDeadline: now.AddDate(0, 7, 0),
						PaymentDetails: &komoju.PaymentDetails{
							Type:            "pay_easy",
							Email:           "test@example.com",
							BankID:          "bank-id",
							CustomerID:      "customer-id",
							ConfirmationID:  "confirmation-id",
							InstructionsURL: "http://example.com/instructions",
							PaymentURL:      "http://example.com",
						},
						PaymentMethodFee:   0,
						Total:              1000,
						Currency:           defaultCurrency,
						CreatedAt:          now,
						AmountRefunded:     0,
						Locale:             defaultLocale,
						Session:            "session-id",
						CustomerFamilyName: "&.",
						CustomerGivenName:  "利用者",
						Refunds:            []*komoju.Refund{},
						RefundRequests:     []*komoju.RefundRequest{},
					},
				}
				err := json.NewEncoder(w).Encode(body)
				assert.NoError(t, err)
			},
			params: &komoju.OrderPayEasyParams{
				SessionID:     "session-id",
				Email:         "test@example.com",
				PhoneNumber:   "09012341234",
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
			},
			expect: &testResponse{
				body: &komoju.OrderSessionResponse{
					RedirectURL: "http://example.com/redirect",
					Status:      komoju.SessionStatusPending,
					Payment: &komoju.PaymentInfo{
						ID:              "payment-id",
						Resource:        "session",
						Status:          komoju.PaymentStatusAuthorized,
						Amount:          1000,
						Tax:             0,
						PaymentDeadline: now.AddDate(0, 7, 0),
						PaymentDetails: &komoju.PaymentDetails{
							Type:            "pay_easy",
							Email:           "test@example.com",
							BankID:          "bank-id",
							CustomerID:      "customer-id",
							ConfirmationID:  "confirmation-id",
							InstructionsURL: "http://example.com/instructions",
							PaymentURL:      "http://example.com",
						},
						PaymentMethodFee:   0,
						Total:              1000,
						Currency:           defaultCurrency,
						CreatedAt:          now,
						AmountRefunded:     0,
						Locale:             defaultLocale,
						Session:            "session-id",
						CustomerFamilyName: "&.",
						CustomerGivenName:  "利用者",
						Refunds:            []*komoju.Refund{},
						RefundRequests:     []*komoju.RefundRequest{},
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
			params: &komoju.OrderPayEasyParams{
				SessionID:     "session-id",
				Email:         "test@example.com",
				PhoneNumber:   "09012341234",
				Lastname:      "&.",
				Firstname:     "利用者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "りようしゃ",
			},
			expect: &testResponse{
				err: &komoju.Error{
					Method:  http.MethodPost,
					Route:   "/api/v1/sessions/%s/pay",
					Status:  http.StatusNotFound,
					Code:    komoju.ErrCodeNotFound,
					Message: "The requested resource could not be found.",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Run(tt.name, testClient(tt.handler, tt.expect, func(ctx context.Context, client komoju.Session) (interface{}, error) {
				return client.OrderPayEasy(ctx, tt.params)
			}))
		})
	}
}
