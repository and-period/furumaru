package entity

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/set"
	"github.com/stretchr/testify/assert"
)

func TestOrderProduct(t *testing.T) {
	t.Parallel()
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
	rates := ShippingRates{
		{Number: 1, Name: "四国", Price: 250, PrefectureCodes: shikoku},
		{Number: 2, Name: "その他", Price: 500, PrefectureCodes: others},
	}
	tests := []struct {
		name   string
		params *NewProductOrderParams
		expect *Order
		hasErr bool
	}{
		{
			name: "success",
			params: &NewProductOrderParams{
				OrderID:       "order-id",
				SessionID:     "session-id",
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				Customer: &entity.User{
					Member: entity.Member{
						UserID:       "user-id",
						CognitoID:    "cognito-id",
						AccountID:    "account-id",
						Username:     "username",
						ProviderType: entity.ProviderTypeEmail,
						Email:        "test@example.com",
						PhoneNumber:  "+819012345678",
						ThumbnailURL: "",
					},
					ID:         "user-id",
					Registered: true,
				},
				BillingAddress: &entity.Address{
					AddressRevision: entity.AddressRevision{
						ID:             1,
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ID:     "address-id",
					UserID: "user-id",
				},
				ShippingAddress: &entity.Address{
					AddressRevision: entity.AddressRevision{
						ID:             1,
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ID:     "address-id",
					UserID: "user-id",
				},
				Shipping: &Shipping{
					ID:            "coordinator-id",
					CoordinatorID: "coordinator-id",
					ShippingRevision: ShippingRevision{
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
				},
				Baskets: CartBaskets{
					{
						BoxNumber: 1,
						BoxType:   ShippingTypeNormal,
						BoxSize:   ShippingSize60,
						BoxRate:   80,
						Items: []*CartItem{
							{
								ProductID: "product-id01",
								Quantity:  1,
							},
							{
								ProductID: "product-id02",
								Quantity:  2,
							},
						},
						CoordinatorID: "coordinator-id",
					},
				},
				Products: Products{
					{
						ID:   "product-id01",
						Name: "じゃがいも",
						ProductRevision: ProductRevision{
							ID:        1,
							ProductID: "product-id01",
							Price:     500,
						},
					},
					{
						ID:   "product-id02",
						Name: "人参",
						ProductRevision: ProductRevision{
							ID:        2,
							ProductID: "product-id02",
							Price:     1980,
						},
					},
				},
				PaymentMethodType: PaymentMethodTypeCreditCard,
				Promotion: &Promotion{
					Title:        "プロモーションタイトル",
					Description:  "プロモーションの詳細です。",
					Public:       true,
					DiscountType: DiscountTypeRate,
					DiscountRate: 10,
					Code:         "excode01",
					CodeType:     PromotionCodeTypeAlways,
					StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
				},
			},
			expect: &Order{
				OrderPayment: OrderPayment{
					OrderID:           "order-id",
					AddressRevisionID: 1,
					Status:            PaymentStatusPending,
					TransactionID:     "",
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          4460,
					Discount:          446,
					ShippingFee:       0,
					Tax:               364,
					Total:             4014,
				},
				OrderFulfillments: OrderFulfillments{
					{
						OrderID:           "order-id",
						AddressRevisionID: 1,
						Status:            FulfillmentStatusUnfulfilled,
						TrackingNumber:    "",
						ShippingCarrier:   ShippingCarrierUnknown,
						ShippingType:      ShippingTypeNormal,
						BoxNumber:         1,
						BoxSize:           ShippingSize60,
						BoxRate:           80,
					},
				},
				OrderItems: OrderItems{
					{
						ProductRevisionID: 1,
						OrderID:           "order-id",
						Quantity:          1,
					},
					{
						ProductRevisionID: 2,
						OrderID:           "order-id",
						Quantity:          2,
					},
				},
				ID:              "order-id",
				SessionID:       "session-id",
				UserID:          "user-id",
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				PromotionID:     "",
				Type:            OrderTypeProduct,
				Status:          OrderStatusUnpaid,
				ShippingMessage: "ご注文ありがとうございます！商品到着まで今しばらくお待ち下さい。",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewProductOrder(tt.params)
			if tt.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			for _, f := range actual.OrderFulfillments {
				f.ID = "" // ignore
			}
			for _, i := range actual.OrderItems {
				i.FulfillmentID = "" // ignore
			}
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestNewExperienceOrder(t *testing.T) {
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name   string
		params *NewExperienceOrderParams
		expect *Order
		hasErr bool
	}{
		{
			name: "success",
			params: &NewExperienceOrderParams{
				OrderID:       "order-id",
				SessionID:     "session-id",
				ShopID:        "shop-id",
				CoordinatorID: "coordinator-id",
				Customer: &entity.User{
					Member: entity.Member{
						UserID:       "user-id",
						CognitoID:    "cognito-id",
						AccountID:    "account-id",
						Username:     "username",
						ProviderType: entity.ProviderTypeEmail,
						Email:        "test@example.com",
						PhoneNumber:  "+819012345678",
						ThumbnailURL: "",
					},
					ID:         "user-id",
					Registered: true,
				},
				BillingAddress: &entity.Address{
					AddressRevision: entity.AddressRevision{
						ID:             1,
						AddressID:      "address-id",
						Lastname:       "&.",
						Firstname:      "購入者",
						PostalCode:     "1000014",
						PrefectureCode: 13,
						City:           "千代田区",
						AddressLine1:   "永田町1-7-1",
						AddressLine2:   "",
						PhoneNumber:    "090-1234-1234",
					},
					ID:     "address-id",
					UserID: "user-id",
				},
				Experience: &Experience{
					ID:                 "experience-id",
					CoordinatorID:      "coordinator-id",
					ProducerID:         "producer-id",
					TypeID:             "experience-type-id",
					Title:              "じゃがいも収穫",
					Description:        "じゃがいもを収穫する体験",
					Public:             true,
					SoldOut:            false,
					Status:             ExperienceStatusAccepting,
					ThumbnailURL:       "http://example.com/thumbnail.png",
					Media:              MultiExperienceMedia{{URL: "http://example.com/thumbnail.png", IsThumbnail: true}},
					RecommendedPoints:  []string{"ポイント1", "ポイント2"},
					PromotionVideoURL:  "http://example.com/promotion.mp4",
					HostPrefecture:     "東京都",
					HostPrefectureCode: 13,
					HostCity:           "千代田区",
					ExperienceRevision: ExperienceRevision{ExperienceID: "experience-id"},
					StartAt:            now.AddDate(0, 0, -1),
					EndAt:              now.AddDate(0, 0, 1),
					CreatedAt:          now,
					UpdatedAt:          now,
				},
				PaymentMethodType: PaymentMethodTypeCreditCard,
				Promotion: &Promotion{
					ID:           "promotion-id",
					Description:  "プロモーションの詳細です。",
					Public:       true,
					DiscountType: DiscountTypeRate,
					DiscountRate: 10,
					Code:         "excode01",
					CodeType:     PromotionCodeTypeAlways,
					StartAt:      jst.Date(2022, 8, 1, 0, 0, 0, 0),
					EndAt:        jst.Date(2022, 9, 1, 0, 0, 0, 0),
				},
				AdultCount:            2,
				JuniorHighSchoolCount: 1,
				ElementarySchoolCount: 3,
				PreschoolCount:        0,
				SeniorCount:           0,
				Transportation:        "電車",
				RequetsedDate:         "20221231",
				RequetsedTime:         "1000",
			},
			expect: &Order{
				OrderPayment: OrderPayment{
					OrderID:           "order-id",
					AddressRevisionID: 1,
					Status:            PaymentStatusPending,
					MethodType:        PaymentMethodTypeCreditCard,
					Subtotal:          0,
					Discount:          0,
					ShippingFee:       0,
					Tax:               0,
					Total:             0,
				},
				OrderExperience: OrderExperience{
					OrderID:               "order-id",
					ExperienceRevisionID:  0,
					AdultCount:            2,
					JuniorHighSchoolCount: 1,
					ElementarySchoolCount: 3,
					PreschoolCount:        0,
					SeniorCount:           0,
					Remarks: OrderExperienceRemarks{
						Transportation: "電車",
						RequestedDate:  jst.Date(2022, 12, 31, 0, 0, 0, 0),
						RequestedTime:  jst.Date(0, 1, 1, 10, 0, 0, 0),
					},
				},
				ID:              "order-id",
				SessionID:       "session-id",
				UserID:          "user-id",
				ShopID:          "shop-id",
				CoordinatorID:   "coordinator-id",
				PromotionID:     "promotion-id",
				Type:            OrderTypeExperience,
				Status:          OrderStatusUnpaid,
				ShippingMessage: "",
			},
			hasErr: false,
		},
		{
			name: "invalid payment method",
			params: &NewExperienceOrderParams{
				OrderID:               "order-id",
				SessionID:             "session-id",
				CoordinatorID:         "coordinator-id",
				Customer:              &entity.User{ID: "user-id"},
				BillingAddress:        &entity.Address{ID: "address-id"},
				Experience:            &Experience{ID: "experience-id"},
				PaymentMethodType:     PaymentMethodType(-1),
				Promotion:             &Promotion{ID: "promotion-id"},
				AdultCount:            2,
				JuniorHighSchoolCount: 1,
				ElementarySchoolCount: 3,
				PreschoolCount:        0,
				SeniorCount:           0,
				Transportation:        "train",
				RequetsedDate:         "2022-12-31",
				RequetsedTime:         "10:00",
			},
			expect: nil,
			hasErr: true,
		},
		{
			name: "missing customer",
			params: &NewExperienceOrderParams{
				OrderID:               "order-id",
				SessionID:             "session-id",
				CoordinatorID:         "coordinator-id",
				Customer:              nil,
				BillingAddress:        &entity.Address{ID: "address-id"},
				Experience:            &Experience{ID: "experience-id"},
				PaymentMethodType:     PaymentMethodTypeCreditCard,
				Promotion:             &Promotion{ID: "promotion-id"},
				AdultCount:            2,
				JuniorHighSchoolCount: 1,
				ElementarySchoolCount: 3,
				PreschoolCount:        0,
				SeniorCount:           0,
				Transportation:        "train",
				RequetsedDate:         "2022-12-31",
				RequetsedTime:         "10:00",
			},
			expect: nil,
			hasErr: true,
		},
		{
			name: "missing billing address",
			params: &NewExperienceOrderParams{
				OrderID:               "order-id",
				SessionID:             "session-id",
				CoordinatorID:         "coordinator-id",
				Customer:              &entity.User{ID: "user-id"},
				BillingAddress:        nil,
				Experience:            &Experience{ID: "experience-id"},
				PaymentMethodType:     PaymentMethodTypeCreditCard,
				Promotion:             &Promotion{ID: "promotion-id"},
				AdultCount:            2,
				JuniorHighSchoolCount: 1,
				ElementarySchoolCount: 3,
				PreschoolCount:        0,
				SeniorCount:           0,
				Transportation:        "train",
				RequetsedDate:         "2022-12-31",
				RequetsedTime:         "10:00",
			},
			expect: nil,
			hasErr: true,
		},
		{
			name: "missing experience",
			params: &NewExperienceOrderParams{
				OrderID:               "order-id",
				SessionID:             "session-id",
				CoordinatorID:         "coordinator-id",
				Customer:              &entity.User{ID: "user-id"},
				BillingAddress:        &entity.Address{ID: "address-id"},
				Experience:            nil,
				PaymentMethodType:     PaymentMethodTypeCreditCard,
				Promotion:             &Promotion{ID: "promotion-id"},
				AdultCount:            2,
				JuniorHighSchoolCount: 1,
				ElementarySchoolCount: 3,
				PreschoolCount:        0,
				SeniorCount:           0,
				Transportation:        "train",
				RequetsedDate:         "2022-12-31",
				RequetsedTime:         "10:00",
			},
			expect: nil,
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := NewExperienceOrder(tt.params)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestOrder_SetPaymentStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		status PaymentStatus
		expect OrderStatus
	}{
		{
			name:   "unpaid",
			status: PaymentStatusPending,
			expect: OrderStatusUnpaid,
		},
		{
			name:   "waiting",
			status: PaymentStatusAuthorized,
			expect: OrderStatusWaiting,
		},
		{
			name:   "preparing",
			status: PaymentStatusCaptured,
			expect: OrderStatusPreparing,
		},
		{
			name:   "canceled",
			status: PaymentStatusCanceled,
			expect: OrderStatusCanceled,
		},
		{
			name:   "refunded",
			status: PaymentStatusRefunded,
			expect: OrderStatusRefunded,
		},
		{
			name:   "failed",
			status: PaymentStatusFailed,
			expect: OrderStatusFailed,
		},
		{
			name:   "unknown",
			status: PaymentStatusUnknown,
			expect: OrderStatusUnknown,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			order := &Order{}
			order.SetPaymentStatus(tt.status)
			assert.Equal(t, tt.expect, order.Status)
		})
	}
}

func TestOrder_SetFulfillmentStatus(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		fulfillments  OrderFulfillments
		fulfillmentID string
		status        FulfillmentStatus
		expect        OrderStatus
	}{
		{
			name: "shipped",
			fulfillments: OrderFulfillments{{
				ID:     "fulfillment-id",
				Status: FulfillmentStatusUnfulfilled,
			}},
			fulfillmentID: "fulfillment-id",
			status:        FulfillmentStatusFulfilled,
			expect:        OrderStatusShipped,
		},
		{
			name: "preparing",
			fulfillments: OrderFulfillments{{
				ID:     "fulfillment-id",
				Status: FulfillmentStatusUnfulfilled,
			}},
			fulfillmentID: "other-id",
			status:        FulfillmentStatusFulfilled,
			expect:        OrderStatusPreparing,
		},
		{
			name: "unfulfilled",
			fulfillments: OrderFulfillments{{
				ID:     "fulfillment-id",
				Status: FulfillmentStatusUnfulfilled,
			}},
			fulfillmentID: "fulfillment-id",
			status:        FulfillmentStatusUnfulfilled,
			expect:        OrderStatusPreparing,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			order := &Order{OrderFulfillments: tt.fulfillments}
			order.SetFulfillmentStatus(tt.fulfillmentID, tt.status)
			assert.Equal(t, tt.expect, order.Status)
		})
	}
}

func TestOrder_SetTransaction(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name          string
		order         *Order
		transactionID string
		now           time.Time
		expect        *Order
	}{
		{
			name: "transaction with positive total",
			order: &Order{
				ID: "order-id",
				OrderPayment: OrderPayment{
					Total: 1000,
				},
			},
			transactionID: "transaction-id",
			now:           now,
			expect: &Order{
				ID: "order-id",
				OrderPayment: OrderPayment{
					Total:         1000,
					TransactionID: "transaction-id",
					OrderedAt:     now,
				},
			},
		},
		{
			name: "transaction with zero total",
			order: &Order{
				ID: "order-id",
				OrderPayment: OrderPayment{
					Total: 0,
				},
			},
			transactionID: "transaction-id",
			now:           now,
			expect: &Order{
				ID:     "order-id",
				Status: OrderStatusPreparing,
				OrderPayment: OrderPayment{
					Total:         0,
					TransactionID: "order-id",
					MethodType:    PaymentMethodTypeNone,
					Status:        PaymentStatusCaptured,
					OrderedAt:     now,
					PaidAt:        now,
					CapturedAt:    now,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.order.SetTransaction(tt.transactionID, tt.now)
			assert.Equal(t, tt.expect, tt.order)
		})
	}
}

func TestOrder_Completed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		order  *Order
		expect bool
	}{
		{
			name:   "unpaid",
			order:  &Order{Status: OrderStatusUnpaid},
			expect: false,
		},
		{
			name:   "waiting",
			order:  &Order{Status: OrderStatusWaiting},
			expect: false,
		},
		{
			name:   "preparing",
			order:  &Order{Status: OrderStatusPreparing},
			expect: false,
		},
		{
			name:   "shipped",
			order:  &Order{Status: OrderStatusShipped},
			expect: false,
		},
		{
			name:   "completed",
			order:  &Order{Status: OrderStatusCompleted},
			expect: true,
		},
		{
			name:   "canceled",
			order:  &Order{Status: OrderStatusCanceled},
			expect: true,
		},
		{
			name:   "refunded",
			order:  &Order{Status: OrderStatusRefunded},
			expect: true,
		},
		{
			name:   "failed",
			order:  &Order{Status: OrderStatusFailed},
			expect: true,
		},
		{
			name:   "nil",
			order:  nil,
			expect: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.order.Completed())
		})
	}
}

func TestOrder_EnableAction(t *testing.T) {
	t.Parallel()
	type want struct {
		capturable  bool
		preservable bool
		completable bool
		cancelable  bool
		refundable  bool
	}
	tests := []struct {
		name  string
		order *Order
		want  want
	}{
		{
			name: "payment pending",
			order: &Order{
				ID:                "order-id",
				OrderPayment:      OrderPayment{Status: PaymentStatusPending},
				OrderFulfillments: OrderFulfillments{{Status: FulfillmentStatusUnfulfilled}},
				CompletedAt:       time.Time{},
			},
			want: want{
				capturable:  false,
				preservable: false,
				completable: false,
				cancelable:  true,
				refundable:  false,
			},
		},
		{
			name: "payment authorized",
			order: &Order{
				ID:                "order-id",
				OrderPayment:      OrderPayment{Status: PaymentStatusAuthorized},
				OrderFulfillments: OrderFulfillments{{Status: FulfillmentStatusUnfulfilled}},
				CompletedAt:       time.Time{},
			},
			want: want{
				capturable:  true,
				preservable: false,
				completable: false,
				cancelable:  true,
				refundable:  false,
			},
		},
		{
			name: "payment captured and unfulfilled",
			order: &Order{
				ID:                "order-id",
				OrderPayment:      OrderPayment{Status: PaymentStatusCaptured},
				OrderFulfillments: OrderFulfillments{{Status: FulfillmentStatusUnfulfilled}},
				CompletedAt:       time.Time{},
			},
			want: want{
				capturable:  false,
				preservable: true,
				completable: false,
				cancelable:  false,
				refundable:  true,
			},
		},
		{
			name: "payment captured and fulfilled",
			order: &Order{
				ID:                "order-id",
				OrderPayment:      OrderPayment{Status: PaymentStatusCaptured},
				OrderFulfillments: OrderFulfillments{{Status: FulfillmentStatusFulfilled}},
				CompletedAt:       time.Time{},
			},
			want: want{
				capturable:  false,
				preservable: true,
				completable: true,
				cancelable:  false,
				refundable:  true,
			},
		},
		{
			name: "payment captured and already completed",
			order: &Order{
				ID:                "order-id",
				OrderPayment:      OrderPayment{Status: PaymentStatusCaptured},
				OrderFulfillments: OrderFulfillments{{Status: FulfillmentStatusFulfilled}},
				CompletedAt:       time.Now(),
			},
			want: want{
				capturable:  false,
				preservable: false,
				completable: false,
				cancelable:  false,
				refundable:  true,
			},
		},
		{
			name: "payment canceled",
			order: &Order{
				ID:                "order-id",
				OrderPayment:      OrderPayment{Status: PaymentStatusCanceled},
				OrderFulfillments: OrderFulfillments{{Status: FulfillmentStatusUnfulfilled}},
				CompletedAt:       time.Time{},
			},
			want: want{
				capturable:  false,
				preservable: false,
				completable: false,
				cancelable:  false,
				refundable:  false,
			},
		},
		{
			name: "payment refunded",
			order: &Order{
				ID:                "order-id",
				OrderPayment:      OrderPayment{Status: PaymentStatusRefunded},
				OrderFulfillments: OrderFulfillments{{Status: FulfillmentStatusFulfilled}},
				CompletedAt:       time.Time{},
			},
			want: want{
				capturable:  false,
				preservable: false,
				completable: false,
				cancelable:  false,
				refundable:  false,
			},
		},
		{
			name: "payment failed",
			order: &Order{
				ID:                "order-id",
				OrderPayment:      OrderPayment{Status: PaymentStatusFailed},
				OrderFulfillments: OrderFulfillments{{Status: FulfillmentStatusFulfilled}},
				CompletedAt:       time.Time{},
			},
			want: want{
				capturable:  false,
				preservable: false,
				completable: false,
				cancelable:  false,
				refundable:  false,
			},
		},
		{
			name:  "empty",
			order: nil,
			want: want{
				capturable:  false,
				preservable: false,
				completable: false,
				cancelable:  false,
				refundable:  false,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want.capturable, tt.order.Capturable())
			assert.Equal(t, tt.want.preservable, tt.order.Preservable())
			assert.Equal(t, tt.want.completable, tt.order.Completable())
			assert.Equal(t, tt.want.cancelable, tt.order.Cancelable())
			assert.Equal(t, tt.want.refundable, tt.order.Refundable())
		})
	}
}

func TestOrders_IDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []string
	}{
		{
			name: "success",
			orders: Orders{
				{
					ID:                "order-id",
					OrderPayment:      OrderPayment{OrderID: "order-id"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
			},
			expect: []string{"order-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.IDs())
		})
	}
}

func TestOrders_UserIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []string
	}{
		{
			name: "success",
			orders: Orders{
				{
					ID:                "order-id",
					UserID:            "user-id",
					CoordinatorID:     "coordinator-id",
					OrderPayment:      OrderPayment{OrderID: "order-id"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
			},
			expect: []string{"user-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.UserIDs())
		})
	}
}

func TestOrders_CoordinatorID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []string
	}{
		{
			name: "success",
			orders: Orders{
				{
					ID:                "order-id",
					UserID:            "user-id",
					CoordinatorID:     "coordinator-id",
					OrderPayment:      OrderPayment{OrderID: "order-id"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
			},
			expect: []string{"coordinator-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.CoordinatorIDs())
		})
	}
}

func TestOrders_PromotionIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []string
	}{
		{
			name: "success",
			orders: Orders{
				{
					ID:                "order-id01",
					UserID:            "user-id",
					CoordinatorID:     "coordinator-id",
					PromotionID:       "promotion-id",
					OrderPayment:      OrderPayment{OrderID: "order-id"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
				{
					ID:                "order-id02",
					UserID:            "user-id",
					CoordinatorID:     "coordinator-id",
					OrderPayment:      OrderPayment{OrderID: "order-id"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
			},
			expect: []string{"promotion-id"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.PromotionIDs())
		})
	}
}

func TestOrders_AddressRevisionIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []int64
	}{
		{
			name: "success",
			orders: Orders{
				{
					ID:                "order-id",
					UserID:            "user-id",
					OrderPayment:      OrderPayment{OrderID: "order-id", AddressRevisionID: 1},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id", AddressRevisionID: 2}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
			},
			expect: []int64{1, 2},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.orders.AddressRevisionIDs())
		})
	}
}

func TestOrders_ProductRevisionIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []int64
	}{
		{
			name: "success",
			orders: Orders{
				{
					ID:                "order-id",
					UserID:            "user-id",
					OrderPayment:      OrderPayment{OrderID: "order-id", AddressRevisionID: 1},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id", AddressRevisionID: 2}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
			},
			expect: []int64{1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.orders.ProductRevisionIDs())
		})
	}
}

func TestOrders_ExperienceRevisionIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders Orders
		expect []int64
	}{
		{
			name: "success",
			orders: Orders{
				{
					ID:                "order-id01",
					UserID:            "user-id",
					CoordinatorID:     "coordinator-id",
					PromotionID:       "promotion-id",
					OrderPayment:      OrderPayment{OrderID: "order-id"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
					OrderExperience:   OrderExperience{OrderID: "order-id", ExperienceRevisionID: 1},
				},
				{
					ID:                "order-id02",
					UserID:            "user-id",
					CoordinatorID:     "coordinator-id",
					OrderPayment:      OrderPayment{OrderID: "order-id"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id"}},
					OrderItems:        OrderItems{{OrderID: "order-id", ProductRevisionID: 1}},
				},
			},
			expect: []int64{1},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.ExperienceRevisionIDs())
		})
	}
}

func TestOrders_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		orders       Orders
		payments     map[string]*OrderPayment
		fulfillments map[string]OrderFulfillments
		items        map[string]OrderItems
		experiences  map[string]*OrderExperience
		expect       Orders
	}{
		{
			name: "success",
			orders: Orders{
				{ID: "order-id01"},
				{ID: "order-id02"},
				{ID: "order-id03"},
				{ID: "order-id04"},
			},
			payments: map[string]*OrderPayment{
				"order-id01": {OrderID: "order-id01"},
				"order-id02": {OrderID: "order-id02"},
				"order-id04": {OrderID: "order-id04"},
			},
			fulfillments: map[string]OrderFulfillments{
				"order-id01": {{OrderID: "order-id01"}},
				"order-id02": {{OrderID: "order-id02"}},
			},
			items: map[string]OrderItems{
				"order-id01": {{OrderID: "order-id01", ProductRevisionID: 1}},
				"order-id02": {{OrderID: "order-id02", ProductRevisionID: 1}},
			},
			experiences: map[string]*OrderExperience{
				"order-id04": {OrderID: "order-id04", ExperienceRevisionID: 1},
			},
			expect: Orders{
				{
					ID:                "order-id01",
					OrderPayment:      OrderPayment{OrderID: "order-id01"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id01"}},
					OrderItems:        OrderItems{{OrderID: "order-id01", ProductRevisionID: 1}},
				},
				{
					ID:                "order-id02",
					OrderPayment:      OrderPayment{OrderID: "order-id02"},
					OrderFulfillments: OrderFulfillments{{OrderID: "order-id02"}},
					OrderItems:        OrderItems{{OrderID: "order-id02", ProductRevisionID: 1}},
				},
				{
					ID:           "order-id03",
					OrderPayment: OrderPayment{},
				},
				{
					ID:              "order-id04",
					OrderPayment:    OrderPayment{OrderID: "order-id04"},
					OrderExperience: OrderExperience{OrderID: "order-id04", ExperienceRevisionID: 1},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.orders.Fill(tt.payments, tt.fulfillments, tt.items, tt.experiences)
			assert.Equal(t, tt.expect, tt.orders)
		})
	}
}

func TestAggregatedUserOrders_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders AggregatedUserOrders
		expect map[string]*AggregatedUserOrder
	}{
		{
			name: "success",
			orders: AggregatedUserOrders{
				{
					UserID:     "user-id",
					OrderCount: 2,
					Subtotal:   3000,
					Discount:   0,
				},
			},
			expect: map[string]*AggregatedUserOrder{
				"user-id": {
					UserID:     "user-id",
					OrderCount: 2,
					Subtotal:   3000,
					Discount:   0,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.Map())
		})
	}
}

func TestAggregatedOrderPayments_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders AggregatedOrderPayments
		expect map[PaymentMethodType]*AggregatedOrderPayment
	}{
		{
			name: "success",
			orders: AggregatedOrderPayments{
				{
					PaymentMethodType: PaymentMethodTypeCreditCard,
					OrderCount:        2,
					UserCount:         1,
					SalesTotal:        6000,
				},
			},
			expect: map[PaymentMethodType]*AggregatedOrderPayment{
				PaymentMethodTypeCreditCard: {
					PaymentMethodType: PaymentMethodTypeCreditCard,
					OrderCount:        2,
					UserCount:         1,
					SalesTotal:        6000,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.Map())
		})
	}
}

func TestAggregatedOrderPayments_OrderTotal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders AggregatedOrderPayments
		expect int64
	}{
		{
			name: "success",
			orders: AggregatedOrderPayments{
				{
					PaymentMethodType: PaymentMethodTypeCreditCard,
					OrderCount:        2,
					UserCount:         1,
					SalesTotal:        6000,
				},
				{
					PaymentMethodType: PaymentMethodTypeBankTransfer,
					OrderCount:        1,
					UserCount:         1,
					SalesTotal:        3000,
				},
			},
			expect: 3,
		},
		{
			name:   "empty",
			orders: AggregatedOrderPayments{},
			expect: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.OrderTotal())
		})
	}
}

func TestAggregatedOrderPromotions_Map(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		orders AggregatedOrderPromotions
		expect map[string]*AggregatedOrderPromotion
	}{
		{
			name: "success",
			orders: AggregatedOrderPromotions{
				{
					PromotionID:   "promotion-id",
					OrderCount:    3,
					DiscountTotal: 1980,
				},
			},
			expect: map[string]*AggregatedOrderPromotion{
				"promotion-id": {
					PromotionID:   "promotion-id",
					OrderCount:    3,
					DiscountTotal: 1980,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.Map())
		})
	}
}

func TestAggregatedPeriodOrders_MapByPeriod(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name   string
		orders AggregatedPeriodOrders
		expect map[time.Time]*AggregatedPeriodOrder
	}{
		{
			name: "success",
			orders: AggregatedPeriodOrders{
				{
					Period:        now,
					OrderCount:    3,
					UserCount:     2,
					SalesTotal:    3000,
					DiscountTotal: 0,
				},
			},
			expect: map[time.Time]*AggregatedPeriodOrder{
				now: {
					Period:        now,
					OrderCount:    3,
					UserCount:     2,
					SalesTotal:    3000,
					DiscountTotal: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.orders.MapByPeriod())
		})
	}
}
