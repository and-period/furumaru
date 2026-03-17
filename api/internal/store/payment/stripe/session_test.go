package stripe

import (
	"context"
	"errors"
	"testing"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/payment"
	mock "github.com/and-period/furumaru/api/mock/pkg/stripe"
	pkgstripe "github.com/and-period/furumaru/api/pkg/stripe"
	lib "github.com/stripe/stripe-go/v82"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateSession(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, m *mock.MockClient)
		params    *payment.CreateSessionParams
		expect    *payment.CreateSessionResult
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, m *mock.MockClient) {
				m.EXPECT().CreateCheckoutSession(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, in *pkgstripe.CreateCheckoutSessionParams) (*lib.CheckoutSession, error) {
						assert.Equal(t, int64(1000), in.Amount)
						assert.Equal(t, "jpy", in.Currency)
						assert.Equal(t, "order-id", in.Description)
						assert.Equal(t, "test@example.com", in.Email)

						return &lib.CheckoutSession{
							ID:  "cs_xxx",
							URL: "https://checkout.stripe.com/pay/cs_xxx",
							PaymentIntent: &lib.PaymentIntent{
								ID: "pi_xxx",
							},
						}, nil
					})
			},
			params: &payment.CreateSessionParams{
				OrderID:     "order-id",
				Amount:      1000,
				CallbackURL: "http://example.com/callback",
				Customer: &payment.CreateSessionCustomer{
					Email: "test@example.com",
				},
			},
			expect: &payment.CreateSessionResult{
				SessionID:  "pi_xxx",
				SessionURL: "https://checkout.stripe.com/pay/cs_xxx",
			},
		},
		{
			name: "error",
			setup: func(ctx context.Context, m *mock.MockClient) {
				m.EXPECT().CreateCheckoutSession(ctx, gomock.Any()).Return(nil, errors.New("stripe error"))
			},
			params: &payment.CreateSessionParams{
				OrderID:     "order-id",
				Amount:      1000,
				CallbackURL: "http://example.com/callback",
				Customer: &payment.CreateSessionCustomer{
					Email: "test@example.com",
				},
			},
			expectErr: errors.New("stripe error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			client := mock.NewMockClient(ctrl)
			ctx := context.Background()
			tt.setup(ctx, client)
			p := &provider{client: client}
			result, err := p.CreateSession(ctx, tt.params)
			if tt.expectErr != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestGetSession(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, m *mock.MockClient)
		sessionID string
		expect    *payment.GetSessionResult
		expectErr error
	}{
		{
			name: "success - requires_capture",
			setup: func(ctx context.Context, m *mock.MockClient) {
				m.EXPECT().GetPaymentIntent(ctx, "pi_xxx").Return(&lib.PaymentIntent{
					Status: lib.PaymentIntentStatusRequiresCapture,
				}, nil)
			},
			sessionID: "pi_xxx",
			expect: &payment.GetSessionResult{
				PaymentStatus: entity.PaymentStatusAuthorized,
			},
		},
		{
			name: "success - succeeded",
			setup: func(ctx context.Context, m *mock.MockClient) {
				m.EXPECT().GetPaymentIntent(ctx, "pi_xxx").Return(&lib.PaymentIntent{
					Status: lib.PaymentIntentStatusSucceeded,
				}, nil)
			},
			sessionID: "pi_xxx",
			expect: &payment.GetSessionResult{
				PaymentStatus: entity.PaymentStatusCaptured,
			},
		},
		{
			name: "error",
			setup: func(ctx context.Context, m *mock.MockClient) {
				m.EXPECT().GetPaymentIntent(ctx, "pi_xxx").Return(nil, errors.New("stripe error"))
			},
			sessionID: "pi_xxx",
			expectErr: errors.New("stripe error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			client := mock.NewMockClient(ctrl)
			ctx := context.Background()
			tt.setup(ctx, client)
			p := &provider{client: client}
			result, err := p.GetSession(ctx, tt.sessionID)
			if tt.expectErr != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestIsSessionFailed_Provider(t *testing.T) {
	t.Parallel()
	p := &provider{}
	assert.False(t, p.IsSessionFailed(nil))
	assert.True(t, p.IsSessionFailed(errors.New("some error")))
	assert.False(t, p.IsSessionFailed(ErrNotSupported))
}
