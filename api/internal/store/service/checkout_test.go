package service

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCheckoutState(t *testing.T) {
	t.Parallel()
	now := jst.Date(2022, 10, 10, 18, 30, 0, 0)
	order := func(status entity.PaymentStatus) *entity.Order {
		return &entity.Order{
			ID:              "order-id",
			SessionID:       "session-id",
			UserID:          "user-id",
			PromotionID:     "",
			CoordinatorID:   "coordinator-id",
			Type:            entity.OrderTypeProduct,
			Status:          entity.OrderStatusUnpaid,
			ShippingMessage: "ご注文ありがとうございます！商品到着まで今しばらくお待ち下さい。",
			CreatedAt:       now,
			UpdatedAt:       now,
			OrderPayment: entity.OrderPayment{
				OrderID:           "order-id",
				AddressRevisionID: 1,
				TransactionID:     "transaction-id",
				Status:            status,
				MethodType:        entity.PaymentMethodTypeCreditCard,
				Subtotal:          1100,
				Discount:          0,
				ShippingFee:       500,
				Tax:               145,
				Total:             1400,
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			OrderFulfillments: entity.OrderFulfillments{
				{
					ID:                "fulfillment-id",
					OrderID:           "order-id",
					AddressRevisionID: 1,
					Status:            entity.FulfillmentStatusUnfulfilled,
					TrackingNumber:    "",
					ShippingCarrier:   entity.ShippingCarrierUnknown,
					ShippingType:      entity.ShippingTypeNormal,
					BoxNumber:         1,
					BoxSize:           entity.ShippingSize60,
					CreatedAt:         now,
					UpdatedAt:         now,
				},
			},
			OrderItems: []*entity.OrderItem{
				{
					FulfillmentID:     "fufillment-id",
					ProductRevisionID: 1,
					OrderID:           "order-id",
					Quantity:          1,
					CreatedAt:         now,
					UpdatedAt:         now,
				},
				{
					FulfillmentID:     "fufillment-id",
					ProductRevisionID: 2,
					OrderID:           "order-id",
					Quantity:          2,
					CreatedAt:         now,
					UpdatedAt:         now,
				},
			},
		}
	}
	session := &komoju.SessionResponse{
		ID: "transaction-id",
		Payment: &komoju.PaymentInfo{
			Status: komoju.PaymentStatusAuthorized,
		},
	}
	tests := []struct {
		name          string
		setup         func(ctx context.Context, mocks *mocks)
		input         *store.GetCheckoutStateInput
		expectOrderID string
		expectStatus  entity.PaymentStatus
		expectErr     error
	}{
		{
			name: "success when authorized",
			setup: func(ctx context.Context, mocks *mocks) {
				order := order(entity.PaymentStatusAuthorized)
				mocks.db.Order.EXPECT().GetByTransactionID(ctx, "user-id", "transaction-id").Return(order, nil)
			},
			input: &store.GetCheckoutStateInput{
				UserID:        "user-id",
				TransactionID: "transaction-id",
			},
			expectOrderID: "order-id",
			expectStatus:  entity.PaymentStatusAuthorized,
			expectErr:     nil,
		},
		{
			name: "success when pending",
			setup: func(ctx context.Context, mocks *mocks) {
				order := order(entity.PaymentStatusPending)
				mocks.db.Order.EXPECT().GetByTransactionID(ctx, "user-id", "transaction-id").Return(order, nil)
				mocks.komojuSession.EXPECT().Get(ctx, "transaction-id").Return(session, nil)
			},
			input: &store.GetCheckoutStateInput{
				UserID:        "user-id",
				TransactionID: "transaction-id",
			},
			expectOrderID: "order-id",
			expectStatus:  entity.PaymentStatusAuthorized,
			expectErr:     nil,
		},
		{
			name: "success with session id",
			setup: func(ctx context.Context, mocks *mocks) {
				order := order(entity.PaymentStatusAuthorized)
				mocks.db.Order.EXPECT().GetByTransactionIDWithSessionID(ctx, "session-id", "transaction-id").Return(order, nil)
			},
			input: &store.GetCheckoutStateInput{
				SessionID:     "session-id",
				TransactionID: "transaction-id",
			},
			expectOrderID: "order-id",
			expectStatus:  entity.PaymentStatusAuthorized,
			expectErr:     nil,
		},
		{
			name:          "invalid argument",
			setup:         func(ctx context.Context, mocks *mocks) {},
			input:         &store.GetCheckoutStateInput{},
			expectOrderID: "",
			expectStatus:  entity.PaymentStatusUnknown,
			expectErr:     exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().GetByTransactionID(ctx, "user-id", "transaction-id").Return(nil, assert.AnError)
			},
			input: &store.GetCheckoutStateInput{
				UserID:        "user-id",
				TransactionID: "transaction-id",
			},
			expectOrderID: "",
			expectStatus:  entity.PaymentStatusUnknown,
			expectErr:     exception.ErrInternal,
		},
		{
			name: "failed to get session",
			setup: func(ctx context.Context, mocks *mocks) {
				order := order(entity.PaymentStatusPending)
				mocks.db.Order.EXPECT().GetByTransactionID(ctx, "user-id", "transaction-id").Return(order, nil)
				mocks.komojuSession.EXPECT().Get(ctx, "transaction-id").Return(nil, assert.AnError)
			},
			input: &store.GetCheckoutStateInput{
				UserID:        "user-id",
				TransactionID: "transaction-id",
			},
			expectOrderID: "order-id",
			expectStatus:  entity.PaymentStatusUnknown,
			expectErr:     exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			orderID, status, err := service.GetCheckoutState(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expectOrderID, orderID)
			assert.Equal(t, tt.expectStatus, status)
		}))
	}
}

func TestCheckoutCreditCard(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &komoju.OrderCreditCardParams{
		SessionID:         "transaction-id",
		Number:            "4111111111111111",
		Month:             12,
		Year:              2024,
		VerificationValue: "123",
		Email:             "test@example.com",
		Name:              "AND USER",
	}
	session := &komoju.OrderSessionResponse{
		RedirectURL: "http://example.com/redirect",
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CheckoutCreditCardInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				checkoutProductMocks(mocks, t, now, entity.PaymentMethodTypeCreditCard)
				mocks.komojuSession.EXPECT().OrderCreditCard(gomock.Any(), params).Return(session, nil)
			},
			input: &store.CheckoutCreditCardInput{
				CheckoutDetail: store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					UserID:           "user-id",
					SessionID:        "session-id",
					RequestID:        "order-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
				Name:              "AND USER",
				Number:            "4111111111111111",
				Month:             12,
				Year:              2024,
				VerificationValue: "123",
			},
			expect:    "http://example.com/redirect",
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CheckoutCreditCardInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name:  "invalid credit card detail",
			setup: func(ctx context.Context, mocks *mocks) {},
			input: &store.CheckoutCreditCardInput{
				CheckoutDetail: store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					UserID:           "user-id",
					SessionID:        "session-id",
					RequestID:        "order-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
				Name:              "AND USER",
				Number:            "4111111111111111",
				Month:             12,
				Year:              2020,
				VerificationValue: "123",
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to order credit card",
			setup: func(ctx context.Context, mocks *mocks) {
				checkoutProductMocks(mocks, t, now, entity.PaymentMethodTypeCreditCard)
				mocks.komojuSession.EXPECT().OrderCreditCard(gomock.Any(), params).Return(nil, assert.AnError)
			},
			input: &store.CheckoutCreditCardInput{
				CheckoutDetail: store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					UserID:           "user-id",
					SessionID:        "session-id",
					RequestID:        "order-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
				Name:              "AND USER",
				Number:            "4111111111111111",
				Month:             12,
				Year:              2024,
				VerificationValue: "123",
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.CheckoutCreditCard(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}, withNow(now)))
	}
}

func TestCheckoutPayPay(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &komoju.OrderPayPayParams{
		SessionID: "transaction-id",
	}
	session := &komoju.OrderSessionResponse{
		RedirectURL: "http://example.com/redirect",
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CheckoutPayPayInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				checkoutProductMocks(mocks, t, now, entity.PaymentMethodTypePayPay)
				mocks.komojuSession.EXPECT().OrderPayPay(gomock.Any(), params).Return(session, nil)
			},
			input: &store.CheckoutPayPayInput{
				CheckoutDetail: store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					UserID:           "user-id",
					SessionID:        "session-id",
					RequestID:        "order-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
			},
			expect:    "http://example.com/redirect",
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CheckoutPayPayInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to order credit card",
			setup: func(ctx context.Context, mocks *mocks) {
				checkoutProductMocks(mocks, t, now, entity.PaymentMethodTypePayPay)
				mocks.komojuSession.EXPECT().OrderPayPay(gomock.Any(), params).Return(nil, assert.AnError)
			},
			input: &store.CheckoutPayPayInput{
				CheckoutDetail: store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					UserID:           "user-id",
					SessionID:        "session-id",
					RequestID:        "order-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.CheckoutPayPay(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}, withNow(now)))
	}
}

func TestCheckoutLinePay(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &komoju.OrderLinePayParams{
		SessionID: "transaction-id",
	}
	session := &komoju.OrderSessionResponse{
		RedirectURL: "http://example.com/redirect",
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CheckoutLinePayInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				checkoutProductMocks(mocks, t, now, entity.PaymentMethodTypeLinePay)
				mocks.komojuSession.EXPECT().OrderLinePay(gomock.Any(), params).Return(session, nil)
			},
			input: &store.CheckoutLinePayInput{
				CheckoutDetail: store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					UserID:           "user-id",
					SessionID:        "session-id",
					RequestID:        "order-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
			},
			expect:    "http://example.com/redirect",
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CheckoutLinePayInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to order credit card",
			setup: func(ctx context.Context, mocks *mocks) {
				checkoutProductMocks(mocks, t, now, entity.PaymentMethodTypeLinePay)
				mocks.komojuSession.EXPECT().OrderLinePay(gomock.Any(), params).Return(nil, assert.AnError)
			},
			input: &store.CheckoutLinePayInput{
				CheckoutDetail: store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					UserID:           "user-id",
					SessionID:        "session-id",
					RequestID:        "order-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.CheckoutLinePay(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}, withNow(now)))
	}
}

func TestCheckoutMerpay(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &komoju.OrderMerpayParams{
		SessionID: "transaction-id",
	}
	session := &komoju.OrderSessionResponse{
		RedirectURL: "http://example.com/redirect",
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CheckoutMerpayInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				checkoutProductMocks(mocks, t, now, entity.PaymentMethodTypeMerpay)
				mocks.komojuSession.EXPECT().OrderMerpay(gomock.Any(), params).Return(session, nil)
			},
			input: &store.CheckoutMerpayInput{
				CheckoutDetail: store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					UserID:           "user-id",
					SessionID:        "session-id",
					RequestID:        "order-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
			},
			expect:    "http://example.com/redirect",
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CheckoutMerpayInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to order credit card",
			setup: func(ctx context.Context, mocks *mocks) {
				checkoutProductMocks(mocks, t, now, entity.PaymentMethodTypeMerpay)
				mocks.komojuSession.EXPECT().OrderMerpay(gomock.Any(), params).Return(nil, assert.AnError)
			},
			input: &store.CheckoutMerpayInput{
				CheckoutDetail: store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					UserID:           "user-id",
					SessionID:        "session-id",
					RequestID:        "order-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.CheckoutMerpay(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}, withNow(now)))
	}
}

func TestCheckoutRakutenPay(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &komoju.OrderRakutenPayParams{
		SessionID: "transaction-id",
	}
	session := &komoju.OrderSessionResponse{
		RedirectURL: "http://example.com/redirect",
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CheckoutRakutenPayInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				checkoutProductMocks(mocks, t, now, entity.PaymentMethodTypeRakutenPay)
				mocks.komojuSession.EXPECT().OrderRakutenPay(gomock.Any(), params).Return(session, nil)
			},
			input: &store.CheckoutRakutenPayInput{
				CheckoutDetail: store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					UserID:           "user-id",
					SessionID:        "session-id",
					RequestID:        "order-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
			},
			expect:    "http://example.com/redirect",
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CheckoutRakutenPayInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to order credit card",
			setup: func(ctx context.Context, mocks *mocks) {
				checkoutProductMocks(mocks, t, now, entity.PaymentMethodTypeRakutenPay)
				mocks.komojuSession.EXPECT().OrderRakutenPay(gomock.Any(), params).Return(nil, assert.AnError)
			},
			input: &store.CheckoutRakutenPayInput{
				CheckoutDetail: store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					UserID:           "user-id",
					SessionID:        "session-id",
					RequestID:        "order-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.CheckoutRakutenPay(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}, withNow(now)))
	}
}

func TestCheckoutAUPay(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &komoju.OrderAUPayParams{
		SessionID: "transaction-id",
	}
	session := &komoju.OrderSessionResponse{
		RedirectURL: "http://example.com/redirect",
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *store.CheckoutAUPayInput
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				checkoutProductMocks(mocks, t, now, entity.PaymentMethodTypeAUPay)
				mocks.komojuSession.EXPECT().OrderAUPay(gomock.Any(), params).Return(session, nil)
			},
			input: &store.CheckoutAUPayInput{
				CheckoutDetail: store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					UserID:           "user-id",
					SessionID:        "session-id",
					RequestID:        "order-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
			},
			expect:    "http://example.com/redirect",
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &store.CheckoutAUPayInput{},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to order credit card",
			setup: func(ctx context.Context, mocks *mocks) {
				checkoutProductMocks(mocks, t, now, entity.PaymentMethodTypeAUPay)
				mocks.komojuSession.EXPECT().OrderAUPay(gomock.Any(), params).Return(nil, assert.AnError)
			},
			input: &store.CheckoutAUPayInput{
				CheckoutDetail: store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					UserID:           "user-id",
					SessionID:        "session-id",
					RequestID:        "order-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.CheckoutAUPay(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}, withNow(now)))
	}
}

func checkoutProductMocks(
	m *mocks,
	t *testing.T,
	now time.Time,
	methodType entity.PaymentMethodType,
) {
	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for _, val := range codes.PrefectureValues {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := entity.ShippingRates{
		{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
	}
	cart := &entity.Cart{
		SessionID: "session-id",
		Baskets: entity.CartBaskets{{
			BoxNumber:     1,
			BoxType:       entity.ShippingTypeNormal,
			BoxSize:       entity.ShippingSize80,
			BoxRate:       80,
			Items:         entity.CartItems{{ProductID: "product-id", Quantity: 2}},
			CoordinatorID: "coordinator-id",
		}},
		ExpiredAt: now.Add(defaultCartTTL),
		CreatedAt: now,
		UpdatedAt: now,
	}
	params := &komoju.CreateSessionParams{
		OrderID:      "order-id",
		Amount:       1400,
		CallbackURL:  "http://example.com/callback",
		PaymentTypes: entity.NewKomojuPaymentTypes(methodType),
		Customer: &komoju.CreateSessionCustomer{
			ID:    "user-id",
			Name:  "&. 購入者",
			Email: "test@example.com",
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
	}
	order := &entity.Order{
		OrderPayment: entity.OrderPayment{
			OrderID:           "order-id",
			AddressRevisionID: 1,
			Status:            entity.PaymentStatusPending,
			TransactionID:     "transaction-id",
			MethodType:        methodType,
			Subtotal:          1000,
			Discount:          100,
			ShippingFee:       500,
			Tax:               127,
			Total:             1400,
			OrderedAt:         now,
		},
		OrderFulfillments: entity.OrderFulfillments{
			{
				OrderID:           "order-id",
				AddressRevisionID: 1,
				Status:            entity.FulfillmentStatusUnfulfilled,
				TrackingNumber:    "",
				ShippingCarrier:   entity.ShippingCarrierUnknown,
				ShippingType:      entity.ShippingTypeNormal,
				BoxNumber:         1,
				BoxSize:           entity.ShippingSize80,
				BoxRate:           80,
			},
		},
		OrderItems: entity.OrderItems{
			{
				OrderID:           "order-id",
				ProductRevisionID: 1,
				Quantity:          2,
			},
		},
		ID:              "order-id",
		SessionID:       "session-id",
		UserID:          "user-id",
		CoordinatorID:   "coordinator-id",
		PromotionID:     "promotion-id",
		Type:            entity.OrderTypeProduct,
		Status:          entity.OrderStatusUnpaid,
		ShippingMessage: "ご注文ありがとうございます！商品到着まで今しばらくお待ち下さい。",
	}

	m.user.EXPECT().
		GetUser(gomock.Any(), &user.GetUserInput{
			UserID: "user-id",
		}).
		Return(&uentity.User{
			Member: uentity.Member{
				UserID:       "user-id",
				CognitoID:    "cognito-id",
				AccountID:    "account-id",
				Username:     "username",
				ProviderType: uentity.ProviderTypeEmail,
				Email:        "test@example.com",
				PhoneNumber:  "+819012345678",
				ThumbnailURL: "",
			},
			ID:         "user-id",
			Registered: true,
			CreatedAt:  now,
			UpdatedAt:  now,
		}, nil)
	m.user.EXPECT().
		GetAddress(gomock.Any(), &user.GetAddressInput{
			UserID:    "user-id",
			AddressID: "address-id",
		}).
		Return(&uentity.Address{
			AddressRevision: uentity.AddressRevision{
				ID:             1,
				AddressID:      "address-id",
				Lastname:       "&.",
				Firstname:      "購入者",
				PostalCode:     "1000014",
				Prefecture:     "東京都",
				PrefectureCode: 13,
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "090-1234-1234",
			},
			ID:        "address-id",
			UserID:    "user-id",
			CreatedAt: now,
			UpdatedAt: now,
		}, nil).Times(2)
	m.cache.EXPECT().
		Get(gomock.Any(), &entity.Cart{SessionID: "session-id"}).
		DoAndReturn(func(_ context.Context, in *entity.Cart) error {
			in.Baskets = cart.Baskets
			in.ExpiredAt = now.Add(defaultCartTTL)
			in.CreatedAt = now
			in.UpdatedAt = now
			return nil
		})
	m.db.Shipping.EXPECT().
		GetByCoordinatorID(gomock.Any(), "coordinator-id").
		Return(&entity.Shipping{
			ShippingRevision: entity.ShippingRevision{
				ShippingID:        "coordinator-id",
				Box60Rates:        rates,
				Box60Frozen:       800,
				Box80Rates:        rates,
				Box80Frozen:       800,
				Box100Rates:       rates,
				Box100Frozen:      800,
				HasFreeShipping:   true,
				FreeShippingRates: 3000,
			},
			ID:            "coordinator-id",
			CoordinatorID: "coordinator-id",
			CreatedAt:     now,
			UpdatedAt:     now,
		}, nil)
	m.db.Promotion.EXPECT().
		GetByCode(gomock.Any(), "code1234").
		Return(&entity.Promotion{
			ID:           "promotion-id",
			Status:       entity.PromotionStatusEnabled,
			Title:        "プロモーションタイトル",
			Description:  "プロモーションの詳細です。",
			Public:       true,
			PublishedAt:  now.AddDate(0, -1, 0),
			DiscountType: entity.DiscountTypeRate,
			DiscountRate: 10,
			Code:         "code1234",
			CodeType:     entity.PromotionCodeTypeAlways,
			StartAt:      now.AddDate(0, -1, 0),
			EndAt:        now.AddDate(0, 1, 0),
		}, nil)
	m.db.Product.EXPECT().
		MultiGet(gomock.Any(), []string{"product-id"}).
		Return(entity.Products{
			{
				ID:        "product-id",
				Name:      "じゃがいも",
				Inventory: 30,
				Public:    true,
				Status:    entity.ProductStatusForSale,
				ProductRevision: entity.ProductRevision{
					ID:        1,
					ProductID: "product-id",
					Price:     500,
				},
			},
		}, nil)
	m.komojuSession.EXPECT().Create(gomock.Any(), params).Return(&komoju.SessionResponse{ID: "transaction-id"}, nil)
	m.db.Order.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, in *entity.Order) error {
			require.Len(t, in.OrderFulfillments, len(order.OrderFulfillments))
			for i := range order.OrderFulfillments {
				order.OrderFulfillments[i].ID = in.OrderFulfillments[i].ID
			}
			require.Len(t, in.OrderItems, len(order.OrderItems))
			for i := range order.OrderItems {
				order.OrderItems[i].FulfillmentID = in.OrderItems[i].FulfillmentID
			}
			require.Equal(t, order, in)
			return nil
		})
	m.db.Product.EXPECT().MultiGet(gomock.Any(), []string{}).Return(entity.Products{}, nil).AnyTimes()
	m.cache.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(assert.AnError).AnyTimes()
	m.db.Product.EXPECT().DecreaseInventory(gomock.Any(), int64(1), int64(2)).Return(nil).AnyTimes()
}

func TestCheckoutProduct(t *testing.T) {
	t.Parallel()
	now := time.Now()
	shikoku := []int32{
		codes.PrefectureValues["tokushima"],
		codes.PrefectureValues["kagawa"],
		codes.PrefectureValues["ehime"],
		codes.PrefectureValues["kochi"],
	}
	set := set.New(shikoku...)
	others := make([]int32, 0, 47-len(shikoku))
	for _, val := range codes.PrefectureValues {
		if set.Contains(val) {
			continue
		}
		others = append(others, val)
	}
	rates := entity.ShippingRates{
		{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
	}
	customerIn := &user.GetUserInput{
		UserID: "user-id",
	}
	customer := &uentity.User{
		Member: uentity.Member{
			UserID:       "user-id",
			CognitoID:    "cognito-id",
			AccountID:    "account-id",
			Username:     "username",
			ProviderType: uentity.ProviderTypeEmail,
			Email:        "test@example.com",
			PhoneNumber:  "+819012345678",
			ThumbnailURL: "",
		},
		ID:         "user-id",
		Registered: true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	addressIn := &user.GetAddressInput{
		UserID:    "user-id",
		AddressID: "address-id",
	}
	address := &uentity.Address{
		AddressRevision: uentity.AddressRevision{
			ID:             1,
			AddressID:      "address-id",
			Lastname:       "&.",
			Firstname:      "購入者",
			PostalCode:     "1000014",
			Prefecture:     "東京都",
			PrefectureCode: 13,
			City:           "千代田区",
			AddressLine1:   "永田町1-7-1",
			AddressLine2:   "",
			PhoneNumber:    "090-1234-1234",
		},
		ID:        "address-id",
		UserID:    "user-id",
		CreatedAt: now,
		UpdatedAt: now,
	}
	shipping := &entity.Shipping{
		ShippingRevision: entity.ShippingRevision{
			ShippingID:        "coordinator-id",
			Box60Rates:        rates,
			Box60Frozen:       800,
			Box80Rates:        rates,
			Box80Frozen:       800,
			Box100Rates:       rates,
			Box100Frozen:      800,
			HasFreeShipping:   true,
			FreeShippingRates: 3000,
		},
		ID:            "coordinator-id",
		CoordinatorID: "coordinator-id",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	promotion := &entity.Promotion{
		ID:           "promotion-id",
		Status:       entity.PromotionStatusEnabled,
		Title:        "プロモーションタイトル",
		Description:  "プロモーションの詳細です。",
		Public:       true,
		PublishedAt:  now.AddDate(0, -1, 0),
		DiscountType: entity.DiscountTypeRate,
		DiscountRate: 10,
		Code:         "testcode",
		CodeType:     entity.PromotionCodeTypeAlways,
		StartAt:      now.AddDate(0, -1, 0),
		EndAt:        now.AddDate(0, 1, 0),
	}
	products := func(inventory int64) entity.Products {
		return entity.Products{
			{
				ID:        "product-id",
				Name:      "じゃがいも",
				Inventory: inventory,
				Public:    true,
				Status:    entity.ProductStatusForSale,
				ProductRevision: entity.ProductRevision{
					ID:        1,
					ProductID: "product-id",
					Price:     500,
				},
			},
		}
	}
	cart := &entity.Cart{
		SessionID: "session-id",
		Baskets: entity.CartBaskets{{
			BoxNumber:     1,
			BoxType:       entity.ShippingTypeNormal,
			BoxSize:       entity.ShippingSize80,
			BoxRate:       80,
			Items:         entity.CartItems{{ProductID: "product-id", Quantity: 2}},
			CoordinatorID: "coordinator-id",
		}},
		ExpiredAt: now.Add(defaultCartTTL),
		CreatedAt: now,
		UpdatedAt: now,
	}
	cartmocks := func(mocks *mocks, sessionID string, cart *entity.Cart, err error) {
		fn := func(_ context.Context, in *entity.Cart) error {
			in.Baskets = cart.Baskets
			in.ExpiredAt = now.Add(defaultCartTTL)
			in.CreatedAt = now
			in.UpdatedAt = now
			return err
		}
		mocks.cache.EXPECT().Get(gomock.Any(), &entity.Cart{SessionID: sessionID}).DoAndReturn(fn)
	}
	sparams := &komoju.CreateSessionParams{
		OrderID:      "order-id",
		Amount:       1400,
		CallbackURL:  "http://example.com/callback",
		PaymentTypes: []komoju.PaymentType{komoju.PaymentTypeCreditCard},
		Customer: &komoju.CreateSessionCustomer{
			ID:    "user-id",
			Name:  "&. 購入者",
			Email: "test@example.com",
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
	}
	session := &komoju.SessionResponse{
		ID:        "transaction-id",
		ReturnURL: "https://example.com",
	}
	order := func() *entity.Order {
		return &entity.Order{
			OrderPayment: entity.OrderPayment{
				OrderID:           "order-id",
				AddressRevisionID: 1,
				Status:            entity.PaymentStatusPending,
				TransactionID:     "transaction-id",
				MethodType:        entity.PaymentMethodTypeCreditCard,
				Subtotal:          1000,
				Discount:          100,
				ShippingFee:       500,
				Tax:               127,
				Total:             1400,
				OrderedAt:         now,
			},
			OrderFulfillments: entity.OrderFulfillments{
				{
					OrderID:           "order-id",
					AddressRevisionID: 1,
					Status:            entity.FulfillmentStatusUnfulfilled,
					TrackingNumber:    "",
					ShippingCarrier:   entity.ShippingCarrierUnknown,
					ShippingType:      entity.ShippingTypeNormal,
					BoxNumber:         1,
					BoxSize:           entity.ShippingSize80,
					BoxRate:           80,
				},
			},
			OrderItems: entity.OrderItems{
				{
					OrderID:           "order-id",
					ProductRevisionID: 1,
					Quantity:          2,
				},
			},
			ID:              "order-id",
			SessionID:       "session-id",
			UserID:          "user-id",
			CoordinatorID:   "coordinator-id",
			PromotionID:     "promotion-id",
			Type:            entity.OrderTypeProduct,
			Status:          entity.OrderStatusUnpaid,
			ShippingMessage: "ご注文ありがとうございます！商品到着まで今しばらくお待ち下さい。",
		}
	}
	ordermocks := func(mocks *mocks, order *entity.Order, err error) {
		fn := func(_ context.Context, in *entity.Order) error {
			require.Len(t, in.OrderFulfillments, len(order.OrderFulfillments))
			for i := range order.OrderFulfillments {
				order.OrderFulfillments[i].ID = in.OrderFulfillments[i].ID
			}
			require.Len(t, in.OrderItems, len(order.OrderItems))
			for i := range order.OrderItems {
				order.OrderItems[i].FulfillmentID = in.OrderItems[i].FulfillmentID
			}
			require.Equal(t, order, in)
			return err
		}
		mocks.db.Order.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(fn)
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		params    *checkoutParams
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				ordermocks(mocks, order(), nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(30), nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{}).Return(entity.Products{}, nil)
				mocks.komojuSession.EXPECT().Create(gomock.Any(), sparams).Return(session, nil)
				mocks.cache.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(assert.AnError)
				mocks.db.Product.EXPECT().DecreaseInventory(gomock.Any(), int64(1), int64(2)).Return(nil).AnyTimes()
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "http://example.com/redirect",
			expectErr: nil,
		},
		{
			name: "success without payment",
			setup: func(ctx context.Context, mocks *mocks) {
				products := entity.Products{
					{
						ID:        "product-id",
						Name:      "じゃがいも",
						Inventory: 30,
						Public:    true,
						Status:    entity.ProductStatusForSale,
						ProductRevision: entity.ProductRevision{
							ID:        1,
							ProductID: "product-id",
							Price:     0,
						},
					},
				}
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products, nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{}).Return(entity.Products{}, nil)
				mocks.db.Order.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, in *entity.Order) error {
					assert.Equal(t, int64(0), in.OrderPayment.Total)
					return nil
				})
				mocks.cache.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(assert.AnError)
				mocks.db.Product.EXPECT().DecreaseInventory(gomock.Any(), int64(1), int64(2)).Return(nil).AnyTimes()
				mocks.messenger.EXPECT().NotifyOrderAuthorized(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            0,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "http://example.com/callback",
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{},
					Type:                  entity.OrderTypeProduct,
					RequestID:             "order-id",
					UserID:                "user-id",
					SessionID:             "session-id",
					PromotionCode:         "code1234",
					BillingAddressID:      "address-id",
					CallbackURL:           "http://example.com/callback",
					Total:                 1400,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(nil, assert.AnError)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get address",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(nil, assert.AnError).Times(2)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get shipping",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(nil, assert.AnError)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(nil, assert.AnError)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to disable promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				promotion := &entity.Promotion{Status: entity.PromotionStatusPrivate}
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to no traget baskets",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, &entity.Cart{}, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get products",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(nil, assert.AnError)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to unmatch products count",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(entity.Products{}, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to insufficient stock",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(0), nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to checksum",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(30), nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create session",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(30), nil)
				mocks.komojuSession.EXPECT().Create(gomock.Any(), sparams).Return(nil, assert.AnError)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create order",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				ordermocks(mocks, order(), assert.AnError)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(30), nil)
				mocks.komojuSession.EXPECT().Create(gomock.Any(), sparams).Return(session, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create order when unprocessable entity error",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				ordermocks(mocks, order(), nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(30), nil)
				mocks.komojuSession.EXPECT().Create(gomock.Any(), sparams).Return(session, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					return nil, &komoju.Error{Status: http.StatusUnprocessableEntity, Code: komoju.ErrCodeUnprocessableEntity}
				},
			},
			expect:    "https://example.com?session_id=transaction-id",
			expectErr: nil,
		},
		{
			name: "failed to callback function",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				ordermocks(mocks, order(), nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(30), nil)
				mocks.komojuSession.EXPECT().Create(gomock.Any(), sparams).Return(session, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            1400,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					return nil, assert.AnError
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create order without payment",
			setup: func(ctx context.Context, mocks *mocks) {
				products := entity.Products{
					{
						ID:        "product-id",
						Name:      "じゃがいも",
						Inventory: 30,
						Public:    true,
						Status:    entity.ProductStatusForSale,
						ProductRevision: entity.ProductRevision{
							ID:        1,
							ProductID: "product-id",
							Price:     0,
						},
					},
				}
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products, nil)
				mocks.db.Order.EXPECT().Create(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutProductDetail: store.CheckoutProductDetail{
						CoordinatorID:     "coordinator-id",
						BoxNumber:         0,
						ShippingAddressID: "address-id",
					},
					Type:             entity.OrderTypeProduct,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            0,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.checkout(ctx, tt.params)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}, withNow(now)))
	}
}

func TestCheckoutExperience(t *testing.T) {
	t.Parallel()
	now := time.Now()
	customerIn := &user.GetUserInput{
		UserID: "user-id",
	}
	customer := &uentity.User{
		Member: uentity.Member{
			UserID:       "user-id",
			CognitoID:    "cognito-id",
			AccountID:    "account-id",
			Username:     "username",
			ProviderType: uentity.ProviderTypeEmail,
			Email:        "test@example.com",
			PhoneNumber:  "+819012345678",
			ThumbnailURL: "",
		},
		ID:         "user-id",
		Registered: true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	addressIn := &user.GetAddressInput{
		UserID:    "user-id",
		AddressID: "address-id",
	}
	address := &uentity.Address{
		AddressRevision: uentity.AddressRevision{
			ID:             1,
			AddressID:      "address-id",
			Lastname:       "&.",
			Firstname:      "購入者",
			PostalCode:     "1000014",
			Prefecture:     "東京都",
			PrefectureCode: 13,
			City:           "千代田区",
			AddressLine1:   "永田町1-7-1",
			AddressLine2:   "",
			PhoneNumber:    "090-1234-1234",
		},
		ID:        "address-id",
		UserID:    "user-id",
		CreatedAt: now,
		UpdatedAt: now,
	}
	promotion := &entity.Promotion{
		ID:           "promotion-id",
		Status:       entity.PromotionStatusEnabled,
		Title:        "プロモーションタイトル",
		Description:  "プロモーションの詳細です。",
		Public:       true,
		PublishedAt:  now.AddDate(0, -1, 0),
		DiscountType: entity.DiscountTypeRate,
		DiscountRate: 10,
		Code:         "testcode",
		CodeType:     entity.PromotionCodeTypeAlways,
		StartAt:      now.AddDate(0, -1, 0),
		EndAt:        now.AddDate(0, 1, 0),
	}
	experience := &entity.Experience{
		ID:            "experience-id",
		CoordinatorID: "coordinator-id",
		ProducerID:    "producer-id",
		TypeID:        "experience-type-id",
		Title:         "じゃがいも収穫",
		Description:   "じゃがいもを収穫する体験です。",
		Public:        true,
		SoldOut:       false,
		Status:        entity.ExperienceStatusAccepting,
		ThumbnailURL:  "http://example.com/thumbnail.png",
		Media: []*entity.ExperienceMedia{
			{URL: "http://example.com/thumbnail01.png", IsThumbnail: true},
			{URL: "http://example.com/thumbnail02.png", IsThumbnail: false},
		},
		RecommendedPoints: []string{
			"じゃがいもを収穫する楽しさを体験できます。",
			"新鮮なじゃがいもを持ち帰ることができます。",
		},
		PromotionVideoURL:  "http://example.com/promotion.mp4",
		Duration:           60,
		Direction:          "彦根駅から徒歩10分",
		BusinessOpenTime:   "1000",
		BusinessCloseTime:  "1800",
		HostPostalCode:     "5220061",
		HostPrefecture:     "滋賀県",
		HostPrefectureCode: 25,
		HostCity:           "彦根市",
		HostAddressLine1:   "金亀町１−１",
		HostAddressLine2:   "",
		HostLongitude:      136.251739,
		HostLatitude:       35.276833,
		StartAt:            now.AddDate(0, 0, -1),
		EndAt:              now.AddDate(0, 0, 1),
		ExperienceRevision: entity.ExperienceRevision{
			ID:                    1,
			ExperienceID:          "experience-id",
			PriceAdult:            1000,
			PriceJuniorHighSchool: 800,
			PriceElementarySchool: 600,
			PricePreschool:        400,
			PriceSenior:           700,
			CreatedAt:             now,
			UpdatedAt:             now,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	sparams := &komoju.CreateSessionParams{
		OrderID:      "order-id",
		Amount:       3240,
		CallbackURL:  "http://example.com/callback",
		PaymentTypes: []komoju.PaymentType{komoju.PaymentTypeCreditCard},
		Customer: &komoju.CreateSessionCustomer{
			ID:    "user-id",
			Name:  "&. 購入者",
			Email: "test@example.com",
		},
		BillingAddress: &komoju.CreateSessionAddress{
			ZipCode:      "1000014",
			Prefecture:   "東京都",
			City:         "千代田区",
			AddressLine1: "永田町1-7-1",
			AddressLine2: "",
		},
	}
	session := &komoju.SessionResponse{
		ID:        "transaction-id",
		ReturnURL: "https://example.com",
	}
	order := func() *entity.Order {
		return &entity.Order{
			OrderPayment: entity.OrderPayment{
				OrderID:           "order-id",
				AddressRevisionID: 1,
				Status:            entity.PaymentStatusPending,
				TransactionID:     "transaction-id",
				MethodType:        entity.PaymentMethodTypeCreditCard,
				Subtotal:          3600,
				Discount:          360,
				ShippingFee:       0,
				Tax:               294,
				Total:             3240,
				OrderedAt:         now,
			},
			OrderExperience: entity.OrderExperience{
				OrderID:               "order-id",
				ExperienceRevisionID:  1,
				AdultCount:            2,
				JuniorHighSchoolCount: 2,
				ElementarySchoolCount: 0,
				PreschoolCount:        0,
				SeniorCount:           0,
				Remarks: entity.OrderExperienceRemarks{
					Transportation: "車で伺います。",
					RequestedDate:  jst.Date(2024, 1, 2, 0, 0, 0, 0),
					RequestedTime:  jst.Date(0, 1, 1, 18, 30, 0, 0),
				},
			},
			ID:            "order-id",
			SessionID:     "session-id",
			UserID:        "user-id",
			CoordinatorID: "coordinator-id",
			PromotionID:   "promotion-id",
			Type:          entity.OrderTypeExperience,
			Status:        entity.OrderStatusUnpaid,
		}
	}
	ordermocks := func(mocks *mocks, order *entity.Order, err error) {
		fn := func(_ context.Context, in *entity.Order) error {
			require.Equal(t, order, in)
			return err
		}
		mocks.db.Order.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(fn)
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		params    *checkoutParams
		expect    string
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				ordermocks(mocks, order(), nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Experience.EXPECT().Get(gomock.Any(), "experience-id").Return(experience, nil)
				mocks.komojuSession.EXPECT().Create(gomock.Any(), sparams).Return(session, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            3240,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "http://example.com/redirect",
			expectErr: nil,
		},
		{
			name: "success without payment",
			setup: func(ctx context.Context, mocks *mocks) {
				experience := &entity.Experience{
					ID:     "experience-id",
					Status: entity.ExperienceStatusAccepting,
					ExperienceRevision: entity.ExperienceRevision{
						ID:                    1,
						PriceAdult:            0,
						PriceJuniorHighSchool: 0,
						PriceElementarySchool: 0,
						PricePreschool:        0,
						PriceSenior:           0,
					},
				}
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Experience.EXPECT().Get(gomock.Any(), "experience-id").Return(experience, nil)
				mocks.db.Order.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, in *entity.Order) error {
					assert.Equal(t, int64(0), in.OrderPayment.Total)
					return nil
				})
				mocks.messenger.EXPECT().NotifyOrderAuthorized(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            0,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "http://example.com/callback",
			expectErr: nil,
		},
		{
			name:  "invalid argument",
			setup: func(ctx context.Context, mocks *mocks) {},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{},
					Type:                     entity.OrderTypeExperience,
					RequestID:                "order-id",
					UserID:                   "user-id",
					SessionID:                "session-id",
					PromotionCode:            "code1234",
					BillingAddressID:         "address-id",
					CallbackURL:              "http://example.com/callback",
					Total:                    3240,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get user",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(nil, assert.AnError)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            3240,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get address",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(nil, assert.AnError)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            3240,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get experience",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Experience.EXPECT().Get(gomock.Any(), "experience-id").Return(nil, assert.AnError)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            3240,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(nil, assert.AnError)
				mocks.db.Experience.EXPECT().Get(gomock.Any(), "experience-id").Return(experience, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            3240,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to disable promotion",
			setup: func(ctx context.Context, mocks *mocks) {
				promotion := &entity.Promotion{Status: entity.PromotionStatusPrivate}
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Experience.EXPECT().Get(gomock.Any(), "experience-id").Return(experience, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            3240,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "experience is not accepting",
			setup: func(ctx context.Context, mocks *mocks) {
				experience := &entity.Experience{Status: entity.ExperienceStatusArchived}
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Experience.EXPECT().Get(gomock.Any(), "experience-id").Return(experience, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            3240,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrFailedPrecondition,
		},
		{
			name: "failed to checksum",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Experience.EXPECT().Get(gomock.Any(), "experience-id").Return(experience, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            10000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to create session",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Experience.EXPECT().Get(gomock.Any(), "experience-id").Return(experience, nil)
				mocks.komojuSession.EXPECT().Create(gomock.Any(), sparams).Return(nil, assert.AnError)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            3240,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create order",
			setup: func(ctx context.Context, mocks *mocks) {
				ordermocks(mocks, order(), assert.AnError)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Experience.EXPECT().Get(gomock.Any(), "experience-id").Return(experience, nil)
				mocks.komojuSession.EXPECT().Create(gomock.Any(), sparams).Return(session, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            3240,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create order when unprocessable entity error",
			setup: func(ctx context.Context, mocks *mocks) {
				ordermocks(mocks, order(), nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Experience.EXPECT().Get(gomock.Any(), "experience-id").Return(experience, nil)
				mocks.komojuSession.EXPECT().Create(gomock.Any(), sparams).Return(session, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            3240,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					return nil, &komoju.Error{Status: http.StatusUnprocessableEntity, Code: komoju.ErrCodeUnprocessableEntity}
				},
			},
			expect:    "https://example.com?session_id=transaction-id",
			expectErr: nil,
		},
		{
			name: "failed to callback function",
			setup: func(ctx context.Context, mocks *mocks) {
				ordermocks(mocks, order(), nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Experience.EXPECT().Get(gomock.Any(), "experience-id").Return(experience, nil)
				mocks.komojuSession.EXPECT().Create(gomock.Any(), sparams).Return(session, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            3240,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					return nil, assert.AnError
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create order without payment",
			setup: func(ctx context.Context, mocks *mocks) {
				experience := &entity.Experience{
					ID:     "experience-id",
					Status: entity.ExperienceStatusAccepting,
					ExperienceRevision: entity.ExperienceRevision{
						ID:                    1,
						PriceAdult:            0,
						PriceJuniorHighSchool: 0,
						PriceElementarySchool: 0,
						PricePreschool:        0,
						PriceSenior:           0,
					},
				}
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil)
				mocks.db.Promotion.EXPECT().GetByCode(gomock.Any(), "code1234").Return(promotion, nil)
				mocks.db.Experience.EXPECT().Get(gomock.Any(), "experience-id").Return(experience, nil)
				mocks.db.Order.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					CheckoutExperienceDetail: store.CheckoutExperienceDetail{
						ExperienceID:          "experience-id",
						AdultCount:            2,
						JuniorHighSchoolCount: 2,
						ElementarySchoolCount: 0,
						PreschoolCount:        0,
						SeniorCount:           0,
						Transportation:        "車で伺います。",
						RequestedDate:         "20240102",
						RequestedTime:         "1830",
					},
					Type:             entity.OrderTypeExperience,
					RequestID:        "order-id",
					UserID:           "user-id",
					SessionID:        "session-id",
					PromotionCode:    "code1234",
					BillingAddressID: "address-id",
					CallbackURL:      "http://example.com/callback",
					Total:            0,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *checkoutDetailParams) (*komoju.OrderSessionResponse, error) {
					res := &komoju.OrderSessionResponse{
						RedirectURL: "http://example.com/redirect",
					}
					return res, nil
				},
			},
			expect:    "",
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			actual, err := service.checkout(ctx, tt.params)
			assert.ErrorIs(t, err, tt.expectErr)
			assert.Equal(t, tt.expect, actual)
		}, withNow(now)))
	}
}
