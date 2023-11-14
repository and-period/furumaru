package service

import (
	"context"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/store/komoju"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		Lastname:          "&.",
		Firstname:         "購入者",
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
				checkoutmocks(mocks, t, now, entity.PaymentMethodTypeCreditCard)
				mocks.komojuSession.EXPECT().OrderCreditCard(gomock.Any(), params).Return(session, nil)
			},
			input: &store.CheckoutCreditCardInput{
				CheckoutDetail: store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1540,
				},
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
			name: "failed to order credit card",
			setup: func(ctx context.Context, mocks *mocks) {
				checkoutmocks(mocks, t, now, entity.PaymentMethodTypeCreditCard)
				mocks.komojuSession.EXPECT().OrderCreditCard(gomock.Any(), params).Return(nil, assert.AnError)
			},
			input: &store.CheckoutCreditCardInput{
				CheckoutDetail: store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1540,
				},
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
				checkoutmocks(mocks, t, now, entity.PaymentMethodTypePayPay)
				mocks.komojuSession.EXPECT().OrderPayPay(gomock.Any(), params).Return(session, nil)
			},
			input: &store.CheckoutPayPayInput{
				CheckoutDetail: store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1540,
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
				checkoutmocks(mocks, t, now, entity.PaymentMethodTypePayPay)
				mocks.komojuSession.EXPECT().OrderPayPay(gomock.Any(), params).Return(nil, assert.AnError)
			},
			input: &store.CheckoutPayPayInput{
				CheckoutDetail: store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1540,
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

func checkoutmocks(
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
			Items:         entity.CartItems{{ProductID: "product-id", Quantity: 2}},
			CoordinatorID: "coordinator-id",
		}},
		ExpiredAt: now.Add(defaultCartTTL),
		CreatedAt: now,
		UpdatedAt: now,
	}
	params := &komoju.CreateSessionParams{
		Amount:       1540,
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
			AddressRevisionID: 1,
			Status:            entity.PaymentStatusPending,
			TransactionID:     "transaction-id",
			MethodType:        methodType,
			Subtotal:          1000,
			Discount:          100,
			ShippingFee:       500,
			Tax:               140,
			Total:             1540,
			OrderedAt:         now,
		},
		OrderFulfillments: entity.OrderFulfillments{
			{
				AddressRevisionID: 1,
				Status:            entity.FulfillmentStatusUnfulfilled,
				TrackingNumber:    "",
				ShippingCarrier:   entity.ShippingCarrierUnknown,
				ShippingType:      entity.ShippingTypeNormal,
				BoxNumber:         1,
				BoxSize:           entity.ShippingSize80,
			},
		},
		OrderItems: entity.OrderItems{
			{
				ProductRevisionID: 1,
				Quantity:          2,
			},
		},
		UserID:        "user-id",
		CoordinatorID: "coordinator-id",
		PromotionID:   "promotion-id",
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
				PhoneNumber:    "+819012345678",
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
		Get(gomock.Any(), "promotion-id").
		Return(&entity.Promotion{
			ID:           "promotion-id",
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
		}, nil)
	m.db.Product.EXPECT().
		MultiGet(gomock.Any(), []string{"product-id"}).
		Return(entity.Products{
			{
				ID:        "product-id",
				Name:      "じゃがいも",
				Inventory: 30,
				ProductRevision: entity.ProductRevision{
					ID:        1,
					ProductID: "product-id",
					Price:     500,
				},
			},
		}, nil)
	m.komojuSession.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, in *komoju.CreateSessionParams) (*komoju.SessionResponse, error) {
			params.OrderID = in.OrderID // ignore
			require.Equal(t, in, params)
			return &komoju.SessionResponse{ID: "transaction-id"}, nil
		})
	m.db.Order.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, in *entity.Order) error {
			order.ID = in.ID // ignore
			order.OrderPayment.OrderID = in.ID
			require.Len(t, in.OrderFulfillments, len(order.OrderFulfillments))
			for i := range order.OrderFulfillments {
				order.OrderFulfillments[i].OrderID = in.ID
				order.OrderFulfillments[i].ID = in.OrderFulfillments[i].ID
			}
			require.Len(t, in.OrderItems, len(order.OrderItems))
			for i := range order.OrderItems {
				order.OrderItems[i].OrderID = in.ID
				order.OrderItems[i].FulfillmentID = in.OrderItems[i].FulfillmentID
			}
			require.Equal(t, order, in)
			return nil
		})
	m.db.Product.EXPECT().MultiGet(gomock.Any(), []string{}).Return(entity.Products{}, nil).AnyTimes()
	m.cache.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(assert.AnError).AnyTimes()
}

func TestNotifyPaymentCompleted(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &database.UpdateOrderPaymentParams{
		Status:    entity.PaymentStatusAuthorized,
		PaymentID: "payment-id",
		IssuedAt:  now,
	}
	tests := []struct {
		name   string
		setup  func(ctx context.Context, mocks *mocks)
		input  *store.NotifyPaymentCompletedInput
		expect error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().UpdatePaymentStatus(ctx, "order-id", params).Return(nil)
			},
			input: &store.NotifyPaymentCompletedInput{
				OrderID:   "order-id",
				PaymentID: "payment-id",
				Status:    entity.PaymentStatusAuthorized,
				IssuedAt:  now,
			},
			expect: nil,
		},
		{
			name:   "invalid argument",
			setup:  func(ctx context.Context, mocks *mocks) {},
			input:  &store.NotifyPaymentCompletedInput{},
			expect: exception.ErrInvalidArgument,
		},
		{
			name: "not updatable",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().UpdatePaymentStatus(ctx, "order-id", params).Return(database.ErrFailedPrecondition)
			},
			input: &store.NotifyPaymentCompletedInput{
				OrderID:   "order-id",
				PaymentID: "payment-id",
				Status:    entity.PaymentStatusAuthorized,
				IssuedAt:  now,
			},
			expect: nil,
		},
		{
			name: "failed to update payment status",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Order.EXPECT().UpdatePaymentStatus(ctx, "order-id", params).Return(assert.AnError)
			},
			input: &store.NotifyPaymentCompletedInput{
				OrderID:   "order-id",
				PaymentID: "payment-id",
				Status:    entity.PaymentStatusAuthorized,
				IssuedAt:  now,
			},
			expect: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.NotifyPaymentCompleted(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expect)
		}))
	}
}

func TestCheckout(t *testing.T) {
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
			PhoneNumber:    "+819012345678",
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
			return nil
		}
		mocks.cache.EXPECT().Get(gomock.Any(), &entity.Cart{SessionID: sessionID}).DoAndReturn(fn)
	}
	sparams := func() *komoju.CreateSessionParams {
		return &komoju.CreateSessionParams{
			Amount:       1540,
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
	}
	session := &komoju.SessionResponse{ID: "transaction-id"}
	sessionmocks := func(mocks *mocks, params *komoju.CreateSessionParams, session *komoju.SessionResponse, err error) {
		fn := func(_ context.Context, in *komoju.CreateSessionParams) (*komoju.SessionResponse, error) {
			params.OrderID = in.OrderID // ignore
			require.Equal(t, in, params)
			return session, err
		}
		mocks.komojuSession.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(fn)
	}
	order := func() *entity.Order {
		return &entity.Order{
			OrderPayment: entity.OrderPayment{
				AddressRevisionID: 1,
				Status:            entity.PaymentStatusPending,
				TransactionID:     "transaction-id",
				MethodType:        entity.PaymentMethodTypeCreditCard,
				Subtotal:          1000,
				Discount:          100,
				ShippingFee:       500,
				Tax:               140,
				Total:             1540,
				OrderedAt:         now,
			},
			OrderFulfillments: entity.OrderFulfillments{
				{
					AddressRevisionID: 1,
					Status:            entity.FulfillmentStatusUnfulfilled,
					TrackingNumber:    "",
					ShippingCarrier:   entity.ShippingCarrierUnknown,
					ShippingType:      entity.ShippingTypeNormal,
					BoxNumber:         1,
					BoxSize:           entity.ShippingSize80,
				},
			},
			OrderItems: entity.OrderItems{
				{
					ProductRevisionID: 1,
					Quantity:          2,
				},
			},
			UserID:        "user-id",
			CoordinatorID: "coordinator-id",
			PromotionID:   "promotion-id",
		}
	}
	ordermocks := func(mocks *mocks, order *entity.Order, err error) {
		fn := func(_ context.Context, in *entity.Order) error {
			order.ID = in.ID // ignore
			order.OrderPayment.OrderID = in.ID
			require.Len(t, in.OrderFulfillments, len(order.OrderFulfillments))
			for i := range order.OrderFulfillments {
				order.OrderFulfillments[i].OrderID = in.ID
				order.OrderFulfillments[i].ID = in.OrderFulfillments[i].ID
			}
			require.Len(t, in.OrderItems, len(order.OrderItems))
			for i := range order.OrderItems {
				order.OrderItems[i].OrderID = in.ID
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
				sessionmocks(mocks, sparams(), session, nil)
				ordermocks(mocks, order(), nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().Get(gomock.Any(), "promotion-id").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(30), nil)
				mocks.db.Product.EXPECT().MultiGet(gomock.Any(), []string{}).Return(entity.Products{}, nil)
				mocks.cache.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1540,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
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
			name: "failed to get user",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(nil, assert.AnError)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().Get(gomock.Any(), "promotion-id").Return(promotion, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
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
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(nil, assert.AnError).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().Get(gomock.Any(), "promotion-id").Return(promotion, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
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
				mocks.db.Promotion.EXPECT().Get(gomock.Any(), "promotion-id").Return(promotion, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
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
				mocks.db.Promotion.EXPECT().Get(gomock.Any(), "promotion-id").Return(nil, assert.AnError)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
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
				promotion := &entity.Promotion{Public: false}
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().Get(gomock.Any(), "promotion-id").Return(promotion, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
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
				mocks.db.Promotion.EXPECT().Get(gomock.Any(), "promotion-id").Return(promotion, nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
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
			name: "failed to insufficient stock",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().Get(gomock.Any(), "promotion-id").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(nil, assert.AnError)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
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
			name: "failed to insufficient stock",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().Get(gomock.Any(), "promotion-id").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(0), nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
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
				mocks.db.Promotion.EXPECT().Get(gomock.Any(), "promotion-id").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(30), nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1000,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
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
				sessionmocks(mocks, sparams(), nil, assert.AnError)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().Get(gomock.Any(), "promotion-id").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(30), nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1540,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
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
				sessionmocks(mocks, sparams(), session, nil)
				ordermocks(mocks, order(), assert.AnError)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().Get(gomock.Any(), "promotion-id").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(30), nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1540,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
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
			name: "failed to callback function",
			setup: func(ctx context.Context, mocks *mocks) {
				cartmocks(mocks, cart.SessionID, cart, nil)
				sessionmocks(mocks, sparams(), session, nil)
				ordermocks(mocks, order(), nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), customerIn).Return(customer, nil)
				mocks.user.EXPECT().GetAddress(gomock.Any(), addressIn).Return(address, nil).Times(2)
				mocks.db.Shipping.EXPECT().GetByCoordinatorID(gomock.Any(), "coordinator-id").Return(shipping, nil)
				mocks.db.Promotion.EXPECT().Get(gomock.Any(), "promotion-id").Return(promotion, nil)
				mocks.db.Product.EXPECT().MultiGet(ctx, []string{"product-id"}).Return(products(30), nil)
			},
			params: &checkoutParams{
				payload: &store.CheckoutDetail{
					UserID:            "user-id",
					SessionID:         "session-id",
					CoordinatorID:     "coordinator-id",
					BoxNumber:         0,
					PromotionID:       "promotion-id",
					BillingAddressID:  "address-id",
					ShippingAddressID: "address-id",
					CallbackURL:       "http://example.com/callback",
					Total:             1540,
				},
				paymentMethodType: entity.PaymentMethodTypeCreditCard,
				payFn: func(ctx context.Context, sessionID string, params *entity.NewOrderParams) (*komoju.OrderSessionResponse, error) {
					return nil, assert.AnError
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
