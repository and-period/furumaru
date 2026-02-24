package stripe

import (
	"context"
	"errors"
	"testing"

	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/payment"
	mock "github.com/and-period/furumaru/api/mock/pkg/stripe"
	lib "github.com/stripe/stripe-go/v82"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestShowPayment(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, m *mock.MockClient)
		paymentID string
		expect    *payment.PaymentResult
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, m *mock.MockClient) {
				m.EXPECT().GetPaymentIntent(ctx, "pi_xxx").Return(&lib.PaymentIntent{
					Status: lib.PaymentIntentStatusRequiresCapture,
				}, nil)
			},
			paymentID: "pi_xxx",
			expect: &payment.PaymentResult{
				Status: entity.PaymentStatusAuthorized,
			},
		},
		{
			name: "error",
			setup: func(ctx context.Context, m *mock.MockClient) {
				m.EXPECT().GetPaymentIntent(ctx, "pi_xxx").Return(nil, errors.New("stripe error"))
			},
			paymentID: "pi_xxx",
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
			result, err := p.ShowPayment(ctx, tt.paymentID)
			if tt.expectErr != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestCapturePayment(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, m *mock.MockClient)
		paymentID string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, m *mock.MockClient) {
				m.EXPECT().Capture(ctx, "pi_xxx").Return(&lib.PaymentIntent{}, nil)
			},
			paymentID: "pi_xxx",
		},
		{
			name: "error",
			setup: func(ctx context.Context, m *mock.MockClient) {
				m.EXPECT().Capture(ctx, "pi_xxx").Return(nil, errors.New("stripe error"))
			},
			paymentID: "pi_xxx",
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
			err := p.CapturePayment(ctx, tt.paymentID)
			if tt.expectErr != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestCancelPayment(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, m *mock.MockClient)
		paymentID string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, m *mock.MockClient) {
				m.EXPECT().Cancel(ctx, "pi_xxx", lib.PaymentIntentCancellationReasonRequestedByCustomer).
					Return(&lib.PaymentIntent{}, nil)
			},
			paymentID: "pi_xxx",
		},
		{
			name: "error",
			setup: func(ctx context.Context, m *mock.MockClient) {
				m.EXPECT().Cancel(ctx, "pi_xxx", lib.PaymentIntentCancellationReasonRequestedByCustomer).
					Return(nil, errors.New("stripe error"))
			},
			paymentID: "pi_xxx",
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
			err := p.CancelPayment(ctx, tt.paymentID)
			if tt.expectErr != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestRefundPayment(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, m *mock.MockClient)
		params    *payment.RefundParams
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, m *mock.MockClient) {
				m.EXPECT().Refund(ctx, "pi_xxx", int64(500), "customer request").
					Return(&lib.Refund{}, nil)
			},
			params: &payment.RefundParams{
				PaymentID:   "pi_xxx",
				Amount:      500,
				Description: "customer request",
			},
		},
		{
			name: "error",
			setup: func(ctx context.Context, m *mock.MockClient) {
				m.EXPECT().Refund(ctx, "pi_xxx", int64(500), "customer request").
					Return(nil, errors.New("stripe error"))
			},
			params: &payment.RefundParams{
				PaymentID:   "pi_xxx",
				Amount:      500,
				Description: "customer request",
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
			err := p.RefundPayment(ctx, tt.params)
			if tt.expectErr != nil {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
