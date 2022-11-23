package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFilterAccessOrder(t *testing.T) {
	t.Parallel()

	in := &store.GetOrderInput{
		OrderID: "order-id",
	}
	order := &sentity.Order{
		ID:                "order-id",
		UserID:            "user-id",
		CoordinatorID:     "coordinator-id",
		ScheduleID:        "schedule-id",
		PaymentStatus:     sentity.PaymentStatusCaptured,
		FulfillmentStatus: sentity.FulfillmentStatusFulfilled,
		CancelType:        sentity.CancelTypeUnknown,
		CancelReason:      "",
		OrderItems: sentity.OrderItems{
			{
				ID:         "item-id",
				OrderID:    "order-id",
				ProductID:  "product-id",
				Price:      100,
				Quantity:   1,
				Weight:     1000,
				WeightUnit: sentity.WeightUnitGram,
				CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
		},
		OrderPayment: sentity.OrderPayment{
			ID:             "payment-id",
			TransactionID:  "transaction-id",
			OrderID:        "order-id",
			PromotionID:    "promotion-id",
			PaymentID:      "payment-id",
			PaymentType:    sentity.PaymentTypeCard,
			Subtotal:       100,
			Discount:       0,
			ShippingCharge: 500,
			Tax:            60,
			Total:          660,
			Lastname:       "&.",
			Firstname:      "スタッフ",
			PostalCode:     "1000014",
			Prefecture:     "東京都",
			City:           "千代田区",
			AddressLine1:   "永田町1-7-1",
			AddressLine2:   "",
			PhoneNumber:    "+819012345678",
			CreatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		OrderFulfillment: sentity.OrderFulfillment{
			ID:              "fulfillment-id",
			OrderID:         "order-id",
			ShippingID:      "shipping-id",
			TrackingNumber:  "",
			ShippingCarrier: sentity.ShippingCarrierUnknown,
			ShippingMethod:  sentity.DeliveryTypeNormal,
			BoxSize:         sentity.ShippingSize60,
			BoxCount:        1,
			WeightTotal:     1000,
			Lastname:        "&.",
			Firstname:       "スタッフ",
			PostalCode:      "1000014",
			Prefecture:      "東京都",
			City:            "千代田区",
			AddressLine1:    "永田町1-7-1",
			AddressLine2:    "",
			PhoneNumber:     "+819012345678",
			CreatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		OrderActivities: sentity.OrderActivities{
			{
				ID:        "event-id",
				OrderID:   "order-id",
				UserID:    "user-id",
				EventType: sentity.OrderEventTypeUnknown,
				Detail:    "支払いが完了しました。",
			},
		},
		CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		options []testOption
		expect  int
	}{
		{
			name:    "administrator success",
			setup:   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			expect:  http.StatusOK,
		},
		{
			name: "coordinator success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetOrder(gomock.Any(), in).Return(order, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			expect:  http.StatusOK,
		},
		{
			name: "coordinator forbidden",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetOrder(gomock.Any(), in).Return(order, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator)},
			expect:  http.StatusForbidden,
		},
		{
			name: "coordinator failed to get order",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetOrder(gomock.Any(), in).Return(nil, assert.AnError)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			expect:  http.StatusInternalServerError,
		},
		{
			name:    "forbidden order",
			setup:   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			options: []testOption{withRole(uentity.AdminRoleProducer)},
			expect:  http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const route, path = "/orders/:orderId", "/orders/order-id"
			testMiddleware(t, tt.setup, route, path, tt.expect, func(h *handler) gin.HandlerFunc {
				return h.filterAccessOrder
			}, tt.options...)
		})
	}
}

func TestListOrders(t *testing.T) {
	t.Parallel()

	ordersIn := &store.ListOrdersInput{
		Limit:  20,
		Offset: 0,
		Orders: []*store.ListOrdersOrder{},
	}
	orders := sentity.Orders{
		{
			ID:                "order-id",
			UserID:            "user-id",
			CoordinatorID:     "coordinator-id",
			ScheduleID:        "schedule-id",
			PaymentStatus:     sentity.PaymentStatusCaptured,
			FulfillmentStatus: sentity.FulfillmentStatusFulfilled,
			CancelType:        sentity.CancelTypeUnknown,
			CancelReason:      "",
			OrderItems: sentity.OrderItems{
				{
					ID:         "item-id",
					OrderID:    "order-id",
					ProductID:  "product-id",
					Price:      100,
					Quantity:   1,
					Weight:     1000,
					WeightUnit: sentity.WeightUnitGram,
					CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
					UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				},
			},
			OrderPayment: sentity.OrderPayment{
				ID:             "payment-id",
				TransactionID:  "transaction-id",
				OrderID:        "order-id",
				PromotionID:    "promotion-id",
				PaymentID:      "payment-id",
				PaymentType:    sentity.PaymentTypeCard,
				Subtotal:       100,
				Discount:       0,
				ShippingCharge: 500,
				Tax:            60,
				Total:          660,
				Lastname:       "&.",
				Firstname:      "スタッフ",
				PostalCode:     "1000014",
				Prefecture:     "東京都",
				City:           "千代田区",
				AddressLine1:   "永田町1-7-1",
				AddressLine2:   "",
				PhoneNumber:    "+819012345678",
				CreatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			OrderFulfillment: sentity.OrderFulfillment{
				ID:              "fulfillment-id",
				OrderID:         "order-id",
				ShippingID:      "shipping-id",
				TrackingNumber:  "",
				ShippingCarrier: sentity.ShippingCarrierUnknown,
				ShippingMethod:  sentity.DeliveryTypeNormal,
				BoxSize:         sentity.ShippingSize60,
				BoxCount:        1,
				WeightTotal:     1000,
				Lastname:        "&.",
				Firstname:       "スタッフ",
				PostalCode:      "1000014",
				Prefecture:      "東京都",
				City:            "千代田区",
				AddressLine1:    "永田町1-7-1",
				AddressLine2:    "",
				PhoneNumber:     "+819012345678",
				CreatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			OrderActivities: sentity.OrderActivities{
				{
					ID:        "event-id",
					OrderID:   "order-id",
					UserID:    "user-id",
					EventType: sentity.OrderEventTypeUnknown,
					Detail:    "支払いが完了しました。",
				},
			},
			CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	usersIn := &user.MultiGetUsersInput{
		UserIDs: []string{"user-id"},
	}
	users := uentity.Users{
		{
			ID:         "user-id",
			Registered: true,
			Customer: uentity.Customer{
				UserID:        "user-id",
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
				CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			Member: uentity.Member{
				UserID:       "user-id",
				AccountID:    "account-id",
				CognitoID:    "cognito-id",
				Username:     "username",
				ProviderType: uentity.ProviderTypeEmail,
				Email:        "test-user@and-period.jp",
				PhoneNumber:  "+819012345678",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
				CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
				VerifiedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
			CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	productsIn := &store.MultiGetProductsInput{
		ProductIDs: []string{"product-id"},
	}
	products := sentity.Products{
		{
			ID:              "product-id",
			TypeID:          "product-type-id",
			ProducerID:      "producer-id",
			Name:            "新鮮なじゃがいも",
			Description:     "新鮮なじゃがいもをお届けします。",
			Public:          true,
			Inventory:       100,
			Weight:          1300,
			WeightUnit:      sentity.WeightUnitGram,
			Item:            1,
			ItemUnit:        "袋",
			ItemDescription: "1袋あたり100gのじゃがいも",
			Media: sentity.MultiProductMedia{
				{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
				{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
			},
			Price:            400,
			DeliveryType:     sentity.DeliveryTypeNormal,
			Box60Rate:        50,
			Box80Rate:        40,
			Box100Rate:       30,
			OriginPrefecture: "滋賀県",
			OriginCity:       "彦根市",
			CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}

	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		options []testOption
		query   string
		expect  *testResponse
	}{
		{
			name: "success administrator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListOrders(gomock.Any(), ordersIn).Return(orders, int64(1), nil)
				mocks.user.EXPECT().MultiGetUsers(gomock.Any(), usersIn).Return(users, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			options: []testOption{withRole(uentity.AdminRoleAdministrator)},
			query:   "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.OrdersResponse{
					Orders: []*response.Order{
						{
							ID:         "order-id",
							UserID:     "user-id",
							UserName:   "&. スタッフ",
							ScheduleID: "schedule-id",
							Payment: &response.OrderPayment{
								TransactionID:  "transaction-id",
								PromotionID:    "promotion-id",
								PaymentID:      "payment-id",
								PaymentType:    2,
								Status:         4,
								Subtotal:       100,
								Discount:       0,
								ShippingCharge: 500,
								Tax:            60,
								Total:          660,
								Lastname:       "&.",
								Firstname:      "スタッフ",
								PostalCode:     "1000014",
								Prefecture:     "東京都",
								City:           "千代田区",
								AddressLine1:   "永田町1-7-1",
								AddressLine2:   "",
								PhoneNumber:    "+819012345678",
							},
							Fulfillment: &response.OrderFulfillment{
								TrackingNumber:  "",
								Status:          2,
								ShippingCarrier: 0,
								ShippingMethod:  1,
								BoxSize:         1,
								BoxCount:        1,
								WeightTotal:     1.0,
								Lastname:        "&.",
								Firstname:       "スタッフ",
								PostalCode:      "1000014",
								Prefecture:      "東京都",
								City:            "千代田区",
								AddressLine1:    "永田町1-7-1",
								AddressLine2:    "",
								PhoneNumber:     "+819012345678",
							},
							Refund: &response.OrderRefund{
								Canceled: false,
								Type:     0,
								Reason:   "",
							},
							Items: []*response.OrderItem{
								{
									ProductID: "product-id",
									Name:      "新鮮なじゃがいも",
									Price:     100,
									Quantity:  1,
									Weight:    1.0,
									Media: []*response.ProductMedia{
										{
											URL:         "https://and-period.jp/thumbnail01.png",
											IsThumbnail: true,
											Images:      []*response.Image{},
										},
										{
											URL:         "https://and-period.jp/thumbnail02.png",
											IsThumbnail: false,
											Images:      []*response.Image{},
										},
									},
								},
							},
							OrderedAt:   0,
							PaidAt:      0,
							DeliveredAt: 0,
							CanceledAt:  0,
							CreatedAt:   1640962800,
							UpdatedAt:   1640962800,
						},
					},
					Total: 1,
				},
			},
		},
		{
			name: "success coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				ordersIn := &store.ListOrdersInput{
					CoordinatorID: "coordinator-id",
					Limit:         20,
					Offset:        0,
					Orders:        []*store.ListOrdersOrder{},
				}
				mocks.store.EXPECT().ListOrders(gomock.Any(), ordersIn).Return(sentity.Orders{}, int64(0), nil)
			},
			options: []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			query:   "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.OrdersResponse{
					Orders: []*response.Order{},
					Total:  0,
				},
			},
		},
		{
			name:  "invalid limit",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?limit=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:  "invalid offset",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?offset=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:  "invalid orders",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?offset=paymentStatus,-fulfilmentStatus,other",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to list orders",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListOrders(gomock.Any(), ordersIn).Return(nil, int64(0), assert.AnError)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to multi get users",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListOrders(gomock.Any(), ordersIn).Return(orders, int64(1), nil)
				mocks.user.EXPECT().MultiGetUsers(gomock.Any(), usersIn).Return(nil, assert.AnError)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to multi get products",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListOrders(gomock.Any(), ordersIn).Return(orders, int64(1), nil)
				mocks.user.EXPECT().MultiGetUsers(gomock.Any(), usersIn).Return(users, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(nil, assert.AnError)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/orders%s"
			path := fmt.Sprintf(format, tt.query)
			testGet(t, tt.setup, tt.expect, path, tt.options...)
		})
	}
}

func TestGetOrder(t *testing.T) {
	t.Parallel()

	orderIn := &store.GetOrderInput{
		OrderID: "order-id",
	}
	order := &sentity.Order{
		ID:                "order-id",
		UserID:            "user-id",
		CoordinatorID:     "coordinator-id",
		ScheduleID:        "schedule-id",
		PaymentStatus:     sentity.PaymentStatusCaptured,
		FulfillmentStatus: sentity.FulfillmentStatusFulfilled,
		CancelType:        sentity.CancelTypeUnknown,
		CancelReason:      "",
		OrderItems: sentity.OrderItems{
			{
				ID:         "item-id",
				OrderID:    "order-id",
				ProductID:  "product-id",
				Price:      100,
				Quantity:   1,
				Weight:     1000,
				WeightUnit: sentity.WeightUnitGram,
				CreatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
				UpdatedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
			},
		},
		OrderPayment: sentity.OrderPayment{
			ID:             "payment-id",
			TransactionID:  "transaction-id",
			OrderID:        "order-id",
			PromotionID:    "promotion-id",
			PaymentID:      "payment-id",
			PaymentType:    sentity.PaymentTypeCard,
			Subtotal:       100,
			Discount:       0,
			ShippingCharge: 500,
			Tax:            60,
			Total:          660,
			Lastname:       "&.",
			Firstname:      "スタッフ",
			PostalCode:     "1000014",
			Prefecture:     "東京都",
			City:           "千代田区",
			AddressLine1:   "永田町1-7-1",
			AddressLine2:   "",
			PhoneNumber:    "+819012345678",
			CreatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		OrderFulfillment: sentity.OrderFulfillment{
			ID:              "fulfillment-id",
			OrderID:         "order-id",
			ShippingID:      "shipping-id",
			TrackingNumber:  "",
			ShippingCarrier: sentity.ShippingCarrierUnknown,
			ShippingMethod:  sentity.DeliveryTypeNormal,
			BoxSize:         sentity.ShippingSize60,
			BoxCount:        1,
			WeightTotal:     1000,
			Lastname:        "&.",
			Firstname:       "スタッフ",
			PostalCode:      "1000014",
			Prefecture:      "東京都",
			City:            "千代田区",
			AddressLine1:    "永田町1-7-1",
			AddressLine2:    "",
			PhoneNumber:     "+819012345678",
			CreatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:       jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		OrderActivities: sentity.OrderActivities{
			{
				ID:        "event-id",
				OrderID:   "order-id",
				UserID:    "user-id",
				EventType: sentity.OrderEventTypeUnknown,
				Detail:    "支払いが完了しました。",
			},
		},
		CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	userIn := &user.GetUserInput{
		UserID: "user-id",
	}
	u := &uentity.User{
		ID:         "user-id",
		Registered: true,
		Customer: uentity.Customer{
			UserID:        "user-id",
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "すたっふ",
			PostalCode:    "1000014",
			Prefecture:    "東京都",
			City:          "千代田区",
			AddressLine1:  "永田町1-7-1",
			AddressLine2:  "",
			CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		Member: uentity.Member{
			UserID:       "user-id",
			AccountID:    "account-id",
			CognitoID:    "cognito-id",
			Username:     "username",
			ProviderType: uentity.ProviderTypeEmail,
			Email:        "test-user@and-period.jp",
			PhoneNumber:  "+819012345678",
			ThumbnailURL: "https://and-period.jp/thumbnail.png",
			CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
			VerifiedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	productsIn := &store.MultiGetProductsInput{
		ProductIDs: []string{"product-id"},
	}
	products := sentity.Products{
		{
			ID:              "product-id",
			TypeID:          "product-type-id",
			ProducerID:      "producer-id",
			Name:            "新鮮なじゃがいも",
			Description:     "新鮮なじゃがいもをお届けします。",
			Public:          true,
			Inventory:       100,
			Weight:          1300,
			WeightUnit:      sentity.WeightUnitGram,
			Item:            1,
			ItemUnit:        "袋",
			ItemDescription: "1袋あたり100gのじゃがいも",
			Media: sentity.MultiProductMedia{
				{URL: "https://and-period.jp/thumbnail01.png", IsThumbnail: true},
				{URL: "https://and-period.jp/thumbnail02.png", IsThumbnail: false},
			},
			Price:            400,
			DeliveryType:     sentity.DeliveryTypeNormal,
			Box60Rate:        50,
			Box80Rate:        40,
			Box100Rate:       30,
			OriginPrefecture: "滋賀県",
			OriginCity:       "彦根市",
			CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}

	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		orderID string
		expect  *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetOrder(gomock.Any(), orderIn).Return(order, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), userIn).Return(u, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			orderID: "order-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.OrderResponse{
					Order: &response.Order{
						ID:         "order-id",
						UserID:     "user-id",
						UserName:   "&. スタッフ",
						ScheduleID: "schedule-id",
						Payment: &response.OrderPayment{
							TransactionID:  "transaction-id",
							PromotionID:    "promotion-id",
							PaymentID:      "payment-id",
							PaymentType:    2,
							Status:         4,
							Subtotal:       100,
							Discount:       0,
							ShippingCharge: 500,
							Tax:            60,
							Total:          660,
							Lastname:       "&.",
							Firstname:      "スタッフ",
							PostalCode:     "1000014",
							Prefecture:     "東京都",
							City:           "千代田区",
							AddressLine1:   "永田町1-7-1",
							AddressLine2:   "",
							PhoneNumber:    "+819012345678",
						},
						Fulfillment: &response.OrderFulfillment{
							TrackingNumber:  "",
							Status:          2,
							ShippingCarrier: 0,
							ShippingMethod:  1,
							BoxSize:         1,
							BoxCount:        1,
							WeightTotal:     1.0,
							Lastname:        "&.",
							Firstname:       "スタッフ",
							PostalCode:      "1000014",
							Prefecture:      "東京都",
							City:            "千代田区",
							AddressLine1:    "永田町1-7-1",
							AddressLine2:    "",
							PhoneNumber:     "+819012345678",
						},
						Refund: &response.OrderRefund{
							Canceled: false,
							Type:     0,
							Reason:   "",
						},
						Items: []*response.OrderItem{
							{
								ProductID: "product-id",
								Name:      "新鮮なじゃがいも",
								Price:     100,
								Quantity:  1,
								Weight:    1.0,
								Media: []*response.ProductMedia{
									{
										URL:         "https://and-period.jp/thumbnail01.png",
										IsThumbnail: true,
										Images:      []*response.Image{},
									},
									{
										URL:         "https://and-period.jp/thumbnail02.png",
										IsThumbnail: false,
										Images:      []*response.Image{},
									},
								},
							},
						},
						OrderedAt:   0,
						PaidAt:      0,
						DeliveredAt: 0,
						CanceledAt:  0,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
			},
		},
		{
			name: "failed to get order",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetOrder(gomock.Any(), orderIn).Return(nil, assert.AnError)
			},
			orderID: "order-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get user",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetOrder(gomock.Any(), orderIn).Return(order, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), userIn).Return(nil, assert.AnError)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(products, nil)
			},
			orderID: "order-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to multi get products",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetOrder(gomock.Any(), orderIn).Return(order, nil)
				mocks.user.EXPECT().GetUser(gomock.Any(), userIn).Return(u, nil)
				mocks.store.EXPECT().MultiGetProducts(gomock.Any(), productsIn).Return(nil, assert.AnError)
			},
			orderID: "order-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/orders/%s"
			path := fmt.Sprintf(format, tt.orderID)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}
