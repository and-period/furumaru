package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotifyStartLive(t *testing.T) {
	t.Parallel()
	now := jst.Date(2023, 12, 25, 18, 30, 0, 0)
	scheduleIn := &store.GetScheduleInput{
		ScheduleID: "schedule-id",
	}
	schedule := &sentity.Schedule{
		ID:              "schedule-id",
		CoordinatorID:   "coordinator-id",
		Status:          sentity.ScheduleStatusLive,
		Title:           "マルシェタイトル",
		Description:     "マルシェ詳細",
		ThumbnailURL:    "",
		Thumbnails:      common.Images{},
		ImageURL:        "",
		OpeningVideoURL: "",
		Public:          true,
		Approved:        true,
		ApprovedAdminID: "",
		StartAt:         now,
		EndAt:           now.Add(time.Hour),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	coordinatorIn := &user.GetCoordinatorInput{
		CoordinatorID: "coordinator-id",
	}
	coordinator := &uentity.Coordinator{
		Admin: uentity.Admin{
			ID:            "coordinator-id",
			Lastname:      "&.",
			Firstname:     "コーディネータ",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "こーでぃねーた",
		},
		AdminID:  "coordinator-id",
		Username: "&. 担当者",
	}
	usersIn := &user.ListUsersInput{
		Limit:          200,
		Offset:         0,
		OnlyRegistered: true,
		OnlyVerified:   true,
	}
	users := uentity.Users{{ID: "user-id"}}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.NotifyStartLiveInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(schedule, nil)
				mocks.user.EXPECT().GetCoordinator(ctx, coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().ListUsers(gomock.Any(), usersIn).Return(users, int64(1), nil)
				mocks.db.ReceivedQueue.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				mocks.producer.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, b []byte) (string, error) {
						payload := &entity.WorkerPayload{}
						err := json.Unmarshal(b, payload)
						require.NoError(t, err)
						expect := &entity.WorkerPayload{
							QueueID:   payload.QueueID, // ignore
							EventType: entity.EventTypeStartLive,
							UserType:  entity.UserTypeUser,
							UserIDs:   []string{"user-id"},
							Email: &entity.MailConfig{
								TemplateID: entity.EmailTemplateIDUserStartLive,
								Substitutions: map[string]interface{}{
									"タイトル":     "マルシェタイトル",
									"コーディネータ名": "&. 担当者",
									"開催日":      "2023-12-25",
									"開始時間":     "18:30",
									"終了時間":     "19:30",
									"サイトURL":   "http://user.example.com/live/schedule-id",
								},
							},
						}
						assert.Equal(t, expect, payload)
						return "", nil
					})
			},
			input: &messenger.NotifyStartLiveInput{
				ScheduleID: "schedule-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.NotifyStartLiveInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get schedule",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(nil, assert.AnError)
			},
			input: &messenger.NotifyStartLiveInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to schedule unpublished",
			setup: func(ctx context.Context, mocks *mocks) {
				schedule := &sentity.Schedule{Approved: false}
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(schedule, nil)
			},
			input: &messenger.NotifyStartLiveInput{
				ScheduleID: "schedule-id",
			},
			expectErr: nil,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(schedule, nil)
				mocks.user.EXPECT().GetCoordinator(ctx, coordinatorIn).Return(nil, assert.AnError)
			},
			input: &messenger.NotifyStartLiveInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to list users",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(schedule, nil)
				mocks.user.EXPECT().GetCoordinator(ctx, coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().ListUsers(gomock.Any(), usersIn).Return(nil, int64(0), assert.AnError)
			},
			input: &messenger.NotifyStartLiveInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to create received queue",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(schedule, nil)
				mocks.user.EXPECT().GetCoordinator(ctx, coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().ListUsers(gomock.Any(), usersIn).Return(users, int64(1), nil)
				mocks.db.ReceivedQueue.EXPECT().Create(gomock.Any(), gomock.Any()).Return(assert.AnError)
			},
			input: &messenger.NotifyStartLiveInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetSchedule(ctx, scheduleIn).Return(schedule, nil)
				mocks.user.EXPECT().GetCoordinator(ctx, coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().ListUsers(gomock.Any(), usersIn).Return(users, int64(1), nil)
				mocks.db.ReceivedQueue.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				mocks.producer.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			input: &messenger.NotifyStartLiveInput{
				ScheduleID: "schedule-id",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.NotifyStartLive(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestNotifyOrderAuthorized(t *testing.T) {
	t.Parallel()
	now := time.Date(2024, 1, 23, 18, 30, 0, 0, time.UTC)
	orderIn := &store.GetOrderInput{
		OrderID: "order-id",
	}
	order := &sentity.Order{
		OrderPayment: sentity.OrderPayment{
			OrderID:           "order-id",
			AddressRevisionID: 1,
			Status:            sentity.PaymentStatusPending,
			TransactionID:     "transaction-id",
			PaymentID:         "payment-id",
			MethodType:        sentity.PaymentMethodTypeCreditCard,
			Subtotal:          4460,
			Discount:          446,
			ShippingFee:       0,
			Tax:               364,
			Total:             4014,
			PaidAt:            now,
		},
		OrderFulfillments: sentity.OrderFulfillments{
			{
				OrderID:           "order-id",
				AddressRevisionID: 1,
				Status:            sentity.FulfillmentStatusUnfulfilled,
				TrackingNumber:    "",
				ShippingCarrier:   sentity.ShippingCarrierUnknown,
				ShippingType:      sentity.ShippingTypeNormal,
				BoxNumber:         1,
				BoxSize:           sentity.ShippingSize60,
			},
		},
		OrderItems: sentity.OrderItems{
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
		ID:            "order-id",
		UserID:        "user-id",
		CoordinatorID: "coordinator-id",
		PromotionID:   "",
	}
	coordinatorIn := &user.GetCoordinatorInput{
		CoordinatorID: "coordinator-id",
		WithDeleted:   true,
	}
	coordinator := &uentity.Coordinator{
		Admin: uentity.Admin{
			ID:            "coordinator-id",
			Role:          uentity.AdminRoleCoordinator,
			Status:        uentity.AdminStatusActivated,
			Lastname:      "&.",
			Firstname:     "コーディネータ",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "こーでぃねーた",
			Email:         "coordinator@example.com",
		},
	}
	products := sentity.Products{
		{
			ID:           "product-id01",
			Name:         "おいしいじゃがいも",
			ThumbnailURL: "http://example.com/image01.png",
			ProductRevision: sentity.ProductRevision{
				ID:        1,
				ProductID: "product-id01",
				Price:     2000,
			},
		},
		{
			ID:           "product-id02",
			Name:         "よく茹でたカリフラワー",
			ThumbnailURL: "http://example.com/image02.png",
			ProductRevision: sentity.ProductRevision{
				ID:        2,
				ProductID: "product-id02",
				Price:     1230,
			},
		},
	}
	addresses := uentity.Addresses{
		{
			AddressRevision: uentity.AddressRevision{
				ID:             1,
				AddressID:      "address-id",
				Lastname:       "&.",
				Firstname:      "太郎",
				LastnameKana:   "あんどどっと",
				FirstnameKana:  "たろう",
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
			IsDefault: true,
		},
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.NotifyOrderAuthorizedInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetOrder(ctx, orderIn).Return(order, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProductsByRevision(gomock.Any(), gomock.Any()).Return(products, nil)
				mocks.user.EXPECT().MultiGetAddressesByRevision(gomock.Any(), gomock.Any()).Return(addresses, nil).Times(2)
				mocks.db.ReceivedQueue.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, queue *entity.ReceivedQueue) error {
						expect := &entity.ReceivedQueue{
							ID:        queue.ID, // ignore
							EventType: entity.EventTypeOrderAuthorized,
							UserType:  entity.UserTypeUser,
							UserIDs:   []string{"user-id"},
							Done:      false,
						}
						assert.Equal(t, expect, queue)
						return nil
					})
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, b []byte) (string, error) {
						payload := &entity.WorkerPayload{}
						err := json.Unmarshal(b, payload)
						require.NoError(t, err)
						expect := &entity.WorkerPayload{
							QueueID:   payload.QueueID, // ignore
							EventType: entity.EventTypeOrderAuthorized,
							UserType:  entity.UserTypeUser,
							UserIDs:   []string{"user-id"},
							Email: &entity.MailConfig{
								TemplateID: entity.EmailTemplateIDUserOrderAuthorized,
								Substitutions: map[string]interface{}{
									"注文番号":  "order-id",
									"決済方法":  "クレジットカード決済",
									"商品金額":  "4460",
									"配送手数料": "0",
									"割引金額":  "446",
									"消費税":   "364",
									"合計金額":  "4014",
									"郵便番号":  "1000014",
									"住所":    "東京都 千代田区 永田町1-7-1",
									"商品一覧": []interface{}{
										map[string]interface{}{
											"商品名":      "おいしいじゃがいも",
											"サムネイルURL": "http://example.com/image01.png",
											"購入数":      "1",
											"商品金額":     "2000",
											"合計金額":     "2000",
										},
										map[string]interface{}{
											"商品名":      "よく茹でたカリフラワー",
											"サムネイルURL": "http://example.com/image02.png",
											"購入数":      "2",
											"商品金額":     "1230",
											"合計金額":     "2460",
										},
									},
								},
							},
							Report: &entity.ReportConfig{
								TemplateID: entity.ReportTemplateIDOrderAuthorized,
								Overview:   "&. 太郎",
								Author:     "&. コーディネータ",
								Link:       "http://admin.example.com/orders/order-id",
								ReceivedAt: now.UTC(),
							},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &messenger.NotifyOrderAuthorizedInput{
				OrderID: "order-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.NotifyOrderAuthorizedInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetOrder(ctx, orderIn).Return(nil, assert.AnError)
			},
			input: &messenger.NotifyOrderAuthorizedInput{
				OrderID: "order-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to get coordinator",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetOrder(ctx, orderIn).Return(order, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, assert.AnError)
				mocks.store.EXPECT().MultiGetProductsByRevision(gomock.Any(), gomock.Any()).Return(products, nil)
				mocks.user.EXPECT().MultiGetAddressesByRevision(gomock.Any(), gomock.Any()).Return(addresses, nil).AnyTimes()
			},
			input: &messenger.NotifyOrderAuthorizedInput{
				OrderID: "order-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to multi get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetOrder(ctx, orderIn).Return(order, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProductsByRevision(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
				mocks.user.EXPECT().MultiGetAddressesByRevision(gomock.Any(), gomock.Any()).Return(addresses, nil).AnyTimes()
			},
			input: &messenger.NotifyOrderAuthorizedInput{
				OrderID: "order-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to multi get addresses",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetOrder(ctx, orderIn).Return(order, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProductsByRevision(gomock.Any(), gomock.Any()).Return(products, nil)
				mocks.user.EXPECT().MultiGetAddressesByRevision(gomock.Any(), gomock.Any()).Return(nil, assert.AnError).MinTimes(1)
			},
			input: &messenger.NotifyOrderAuthorizedInput{
				OrderID: "order-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetOrder(ctx, orderIn).Return(order, nil)
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.store.EXPECT().MultiGetProductsByRevision(gomock.Any(), gomock.Any()).Return(products, nil)
				mocks.user.EXPECT().MultiGetAddressesByRevision(gomock.Any(), gomock.Any()).Return(addresses, nil).Times(2)
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &messenger.NotifyOrderAuthorizedInput{
				OrderID: "order-id",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.NotifyOrderAuthorized(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestNotifyOrderShipped(t *testing.T) {
	t.Parallel()
	orderIn := &store.GetOrderInput{
		OrderID: "order-id",
	}
	order := &sentity.Order{
		OrderPayment: sentity.OrderPayment{
			OrderID:           "order-id",
			AddressRevisionID: 1,
			Status:            sentity.PaymentStatusPending,
			TransactionID:     "transaction-id",
			PaymentID:         "payment-id",
			MethodType:        sentity.PaymentMethodTypeCreditCard,
			Subtotal:          4460,
			Discount:          446,
			ShippingFee:       0,
			Tax:               364,
			Total:             4014,
		},
		OrderFulfillments: sentity.OrderFulfillments{
			{
				OrderID:           "order-id",
				AddressRevisionID: 1,
				Status:            sentity.FulfillmentStatusUnfulfilled,
				TrackingNumber:    "",
				ShippingCarrier:   sentity.ShippingCarrierUnknown,
				ShippingType:      sentity.ShippingTypeNormal,
				BoxNumber:         1,
				BoxSize:           sentity.ShippingSize60,
			},
		},
		OrderItems: sentity.OrderItems{
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
		UserID:          "user-id",
		CoordinatorID:   "coordinator-id",
		PromotionID:     "",
		ShippingMessage: "購入ありがとうございました",
	}
	products := sentity.Products{
		{
			ID:           "product-id01",
			Name:         "おいしいじゃがいも",
			ThumbnailURL: "http://example.com/image01.png",
			ProductRevision: sentity.ProductRevision{
				ID:        1,
				ProductID: "product-id01",
				Price:     2000,
			},
		},
		{
			ID:           "product-id02",
			Name:         "よく茹でたカリフラワー",
			ThumbnailURL: "http://example.com/image02.png",
			ProductRevision: sentity.ProductRevision{
				ID:        2,
				ProductID: "product-id02",
				Price:     1230,
			},
		},
	}
	addresses := uentity.Addresses{
		{
			AddressRevision: uentity.AddressRevision{
				ID:             1,
				AddressID:      "address-id",
				Lastname:       "&.",
				Firstname:      "太郎",
				LastnameKana:   "あんどどっと",
				FirstnameKana:  "たろう",
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
			IsDefault: true,
		},
	}
	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.NotifyOrderShippedInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetOrder(ctx, orderIn).Return(order, nil)
				mocks.store.EXPECT().MultiGetProductsByRevision(ctx, gomock.Any()).Return(products, nil)
				mocks.user.EXPECT().MultiGetAddressesByRevision(ctx, gomock.Any()).Return(addresses, nil)
				mocks.db.ReceivedQueue.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, queue *entity.ReceivedQueue) error {
						expect := &entity.ReceivedQueue{
							ID:        queue.ID, // ignore
							EventType: entity.EventTypeOrderShipped,
							UserType:  entity.UserTypeUser,
							UserIDs:   []string{"user-id"},
							Done:      false,
						}
						assert.Equal(t, expect, queue)
						return nil
					})
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, b []byte) (string, error) {
						payload := &entity.WorkerPayload{}
						err := json.Unmarshal(b, payload)
						require.NoError(t, err)
						expect := &entity.WorkerPayload{
							QueueID:   payload.QueueID, // ignore
							EventType: entity.EventTypeOrderShipped,
							UserType:  entity.UserTypeUser,
							UserIDs:   []string{"user-id"},
							Email: &entity.MailConfig{
								TemplateID: entity.EmailTemplateIDUserOrderShipped,
								Substitutions: map[string]interface{}{
									"注文番号":  "order-id",
									"決済方法":  "クレジットカード決済",
									"商品金額":  "4460",
									"配送手数料": "0",
									"割引金額":  "446",
									"消費税":   "364",
									"合計金額":  "4014",
									"郵便番号":  "1000014",
									"住所":    "東京都 千代田区 永田町1-7-1",
									"メッセージ": "購入ありがとうございました",
									"商品一覧": []interface{}{
										map[string]interface{}{
											"商品名":      "おいしいじゃがいも",
											"サムネイルURL": "http://example.com/image01.png",
											"購入数":      "1",
											"商品金額":     "2000",
											"合計金額":     "2000",
										},
										map[string]interface{}{
											"商品名":      "よく茹でたカリフラワー",
											"サムネイルURL": "http://example.com/image02.png",
											"購入数":      "2",
											"商品金額":     "1230",
											"合計金額":     "2460",
										},
									},
								},
							},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &messenger.NotifyOrderShippedInput{
				OrderID: "order-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.NotifyOrderShippedInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get order",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetOrder(ctx, orderIn).Return(nil, assert.AnError)
			},
			input: &messenger.NotifyOrderShippedInput{
				OrderID: "order-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to multi get products",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetOrder(ctx, orderIn).Return(order, nil)
				mocks.store.EXPECT().MultiGetProductsByRevision(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			input: &messenger.NotifyOrderShippedInput{
				OrderID: "order-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to send messag",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetOrder(ctx, orderIn).Return(order, nil)
				mocks.store.EXPECT().MultiGetProductsByRevision(ctx, gomock.Any()).Return(products, nil)
				mocks.user.EXPECT().MultiGetAddressesByRevision(ctx, gomock.Any()).Return(nil, assert.AnError)
			},
			input: &messenger.NotifyOrderShippedInput{
				OrderID: "order-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to send messag",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.store.EXPECT().GetOrder(ctx, orderIn).Return(order, nil)
				mocks.store.EXPECT().MultiGetProductsByRevision(ctx, gomock.Any()).Return(products, nil)
				mocks.user.EXPECT().MultiGetAddressesByRevision(ctx, gomock.Any()).Return(addresses, nil)
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			input: &messenger.NotifyOrderShippedInput{
				OrderID: "order-id",
			},
			expectErr: exception.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.NotifyOrderShipped(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestNotifyRegisterAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.NotifyRegisterAdminInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, queue *entity.ReceivedQueue) error {
						expect := &entity.ReceivedQueue{
							ID:        queue.ID, // ignore
							EventType: entity.EventTypeRegisterAdmin,
							UserType:  entity.UserTypeAdmin,
							UserIDs:   []string{"admin-id"},
							Done:      false,
						}
						assert.Equal(t, expect, queue)
						return nil
					})
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, b []byte) (string, error) {
						payload := &entity.WorkerPayload{}
						err := json.Unmarshal(b, payload)
						require.NoError(t, err)
						expect := &entity.WorkerPayload{
							QueueID:   payload.QueueID, // ignore
							EventType: entity.EventTypeRegisterAdmin,
							UserType:  entity.UserTypeAdmin,
							UserIDs:   []string{"admin-id"},
							Email: &entity.MailConfig{
								TemplateID: entity.EmailTemplateIDAdminRegister,
								Substitutions: map[string]interface{}{
									"サイトURL": "http://admin.example.com/signin",
									"パスワード":  "!Qaz2wsx",
								},
							},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &messenger.NotifyRegisterAdminInput{
				AdminID:  "admin-id",
				Password: "!Qaz2wsx",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.NotifyRegisterAdminInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			input: &messenger.NotifyRegisterAdminInput{
				AdminID:  "admin-id",
				Password: "!Qaz2wsx",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.NotifyRegisterAdmin(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestNotifyResetAdminPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.NotifyResetAdminPasswordInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, queue *entity.ReceivedQueue) error {
						expect := &entity.ReceivedQueue{
							ID:        queue.ID, // ignore
							EventType: entity.EventTypeResetAdminPassword,
							UserType:  entity.UserTypeAdmin,
							UserIDs:   []string{"admin-id"},
							Done:      false,
						}
						assert.Equal(t, expect, queue)
						return nil
					})
				mocks.producer.EXPECT().
					SendMessage(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, b []byte) (string, error) {
						payload := &entity.WorkerPayload{}
						err := json.Unmarshal(b, payload)
						require.NoError(t, err)
						expect := &entity.WorkerPayload{
							QueueID:   payload.QueueID, // ignore
							EventType: entity.EventTypeResetAdminPassword,
							UserType:  entity.UserTypeAdmin,
							UserIDs:   []string{"admin-id"},
							Email: &entity.MailConfig{
								TemplateID: entity.EmailTemplateIDAdminResetPassword,
								Substitutions: map[string]interface{}{
									"サイトURL": "http://admin.example.com/signin",
									"パスワード":  "!Qaz2wsx",
								},
							},
						}
						assert.Equal(t, expect, payload)
						return "message-id", nil
					})
			},
			input: &messenger.NotifyResetAdminPasswordInput{
				AdminID:  "admin-id",
				Password: "!Qaz2wsx",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.NotifyResetAdminPasswordInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			input: &messenger.NotifyResetAdminPasswordInput{
				AdminID:  "admin-id",
				Password: "!Qaz2wsx",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.NotifyResetAdminPassword(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}))
	}
}

func TestNotifyNotification(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 7, 21, 18, 30, 0, 0)
	notification := &entity.Notification{
		ID:    "notification-id",
		Type:  entity.NotificationTypeLive,
		Title: "お知らせ件名",
		Body:  "お知らせ内容",
		Targets: []entity.NotificationTarget{
			entity.NotificationTargetUsers,
			entity.NotificationTargetCoordinators,
			entity.NotificationTargetProducers,
		},
		CreatedBy:   "admin-id",
		PublishedAt: now,
	}
	adminIn := &user.GetAdminInput{
		AdminID: "admin-id",
	}
	admin := &uentity.Admin{
		ID:            "admin-id",
		Role:          uentity.AdminRoleAdministrator,
		Status:        uentity.AdminStatusActivated,
		Lastname:      "&.",
		Firstname:     "管理者",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "かんりしゃ",
		Email:         "test@example.com",
	}
	users := uentity.Users{{ID: "user-id"}}
	coordinators := uentity.Coordinators{{AdminID: "admin-id"}}
	producers := uentity.Producers{}

	tests := []struct {
		name      string
		setup     func(ctx context.Context, mocks *mocks)
		input     *messenger.NotifyNotificationInput
		expectErr error
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().Get(ctx, "notification-id").Return(notification, nil)
				mocks.user.EXPECT().GetAdmin(ctx, adminIn).Return(admin, nil)
				mocks.user.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Return(users, int64(1), nil)
				mocks.user.EXPECT().ListCoordinators(gomock.Any(), gomock.Any()).Return(coordinators, int64(1), nil)
				mocks.user.EXPECT().ListProducers(gomock.Any(), gomock.Any()).Return(producers, int64(0), nil)
				mocks.db.ReceivedQueue.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, queue *entity.ReceivedQueue) error {
						expect := map[entity.UserType]*entity.ReceivedQueue{
							entity.UserTypeUser: {
								ID:        queue.ID, // ignore
								EventType: entity.EventTypeNotification,
								UserType:  entity.UserTypeUser,
								UserIDs:   []string{"user-id"},
								Done:      false,
							},
							entity.UserTypeCoordinator: {
								ID:        queue.ID, // ignore
								EventType: entity.EventTypeNotification,
								UserType:  entity.UserTypeCoordinator,
								UserIDs:   []string{"admin-id"},
								Done:      false,
							},
							entity.UserTypeNone: {
								ID:        queue.ID, // ignore
								EventType: entity.EventTypeNotification,
								UserType:  entity.UserTypeNone,
								UserIDs:   nil,
								Done:      false,
							},
						}
						assert.Equal(t, expect[queue.UserType], queue)
						return nil
					}).Times(3)
				mocks.producer.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, b []byte) (string, error) {
						payload := &entity.WorkerPayload{}
						err := json.Unmarshal(b, payload)
						require.NoError(t, err)
						expect := map[entity.UserType]*entity.WorkerPayload{
							entity.UserTypeUser: {
								QueueID:   payload.QueueID, // ignore
								EventType: entity.EventTypeNotification,
								UserType:  entity.UserTypeUser,
								UserIDs:   []string{"user-id"},
								Message: &entity.MessageConfig{
									TemplateID:  entity.MessageTemplateIDNotificationLive,
									MessageType: entity.MessageTypeNotification,
									Title:       "お知らせ件名",
									Detail:      "お知らせ内容",
									ReceivedAt:  now.Local(),
								},
							},
							entity.UserTypeCoordinator: {
								QueueID:   payload.QueueID, // ignore
								EventType: entity.EventTypeNotification,
								UserType:  entity.UserTypeCoordinator,
								UserIDs:   []string{"admin-id"},
								Message: &entity.MessageConfig{
									TemplateID:  entity.MessageTemplateIDNotificationLive,
									MessageType: entity.MessageTypeNotification,
									Title:       "お知らせ件名",
									Detail:      "お知らせ内容",
									Link:        "http://admin.example.com/notifications/notification-id",
									ReceivedAt:  now.Local(),
								},
							},
							entity.UserTypeNone: {
								QueueID:   payload.QueueID, // ignore
								EventType: entity.EventTypeNotification,
								UserType:  entity.UserTypeNone,
								UserIDs:   nil,
								Report: &entity.ReportConfig{
									TemplateID:  entity.ReportTemplateIDNotification,
									Overview:    "お知らせ件名",
									Detail:      "お知らせ内容",
									Author:      "&. 管理者",
									Link:        "http://admin.example.com/notifications/notification-id",
									PublishedAt: now.Local(),
								},
							},
						}
						assert.Equal(t, expect[payload.UserType], payload)
						return "message-id", nil
					}).Times(3)
			},
			input: &messenger.NotifyNotificationInput{
				NotificationID: "notification-id",
			},
			expectErr: nil,
		},
		{
			name: "success to target none",
			setup: func(ctx context.Context, mocks *mocks) {
				notification := &entity.Notification{Targets: []entity.NotificationTarget{}, CreatedBy: "admin-id"}
				mocks.db.Notification.EXPECT().Get(ctx, "notification-id").Return(notification, nil)
				mocks.user.EXPECT().GetAdmin(ctx, adminIn).Return(admin, nil)
				mocks.db.ReceivedQueue.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				mocks.producer.EXPECT().SendMessage(gomock.Any(), gomock.Any()).Return("", nil)
			},
			input: &messenger.NotifyNotificationInput{
				NotificationID: "notification-id",
			},
			expectErr: nil,
		},
		{
			name:      "invalid argument",
			setup:     func(ctx context.Context, mocks *mocks) {},
			input:     &messenger.NotifyNotificationInput{},
			expectErr: exception.ErrInvalidArgument,
		},
		{
			name: "failed to get notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().Get(ctx, "notification-id").Return(nil, assert.AnError)
			},
			input: &messenger.NotifyNotificationInput{
				NotificationID: "notification-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to notify user notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().Get(ctx, "notification-id").Return(notification, nil)
				mocks.user.EXPECT().GetAdmin(ctx, adminIn).Return(nil, assert.AnError)
			},
			input: &messenger.NotifyNotificationInput{
				NotificationID: "notification-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to notify user notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().Get(ctx, "notification-id").Return(notification, nil)
				mocks.user.EXPECT().GetAdmin(ctx, adminIn).Return(admin, nil)
				mocks.user.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Return(nil, int64(0), assert.AnError)
				mocks.user.EXPECT().ListCoordinators(gomock.Any(), gomock.Any()).Return(uentity.Coordinators{}, int64(0), nil)
				mocks.user.EXPECT().ListProducers(gomock.Any(), gomock.Any()).Return(uentity.Producers{}, int64(0), nil)
			},
			input: &messenger.NotifyNotificationInput{
				NotificationID: "notification-id",
			},
			expectErr: exception.ErrInternal,
		},
		{
			name: "failed to notify admin notification",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.Notification.EXPECT().Get(ctx, "notification-id").Return(notification, nil)
				mocks.user.EXPECT().GetAdmin(ctx, adminIn).Return(admin, nil)
				mocks.user.EXPECT().ListUsers(gomock.Any(), gomock.Any()).Return(uentity.Users{}, int64(0), nil)
				mocks.user.EXPECT().ListCoordinators(gomock.Any(), gomock.Any()).Return(nil, int64(0), assert.AnError)
				mocks.user.EXPECT().ListProducers(gomock.Any(), gomock.Any()).Return(nil, int64(0), assert.AnError)
			},
			input: &messenger.NotifyNotificationInput{
				NotificationID: "notification-id",
			},
			expectErr: exception.ErrInternal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.NotifyNotification(ctx, tt.input)
			assert.ErrorIs(t, err, tt.expectErr)
		}, withNow(now)))
	}
}

func TestSendMessage(t *testing.T) {
	t.Parallel()
	queue := &entity.ReceivedQueue{
		ID:        "queue-id",
		EventType: entity.EventTypeUnknown,
		UserType:  entity.UserTypeUser,
		UserIDs:   []string{"user-id"},
	}
	tests := []struct {
		name    string
		setup   func(ctx context.Context, mocks *mocks)
		payload *entity.WorkerPayload
		hasErr  bool
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, queue).Return(nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email:     &entity.MailConfig{},
			},
			hasErr: false,
		},
		{
			name: "failed to create received queue",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, queue).Return(assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email:     &entity.MailConfig{},
			},
			hasErr: true,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, queue).Return(nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				UserType:  entity.UserTypeUser,
				UserIDs:   []string{"user-id"},
				Email:     &entity.MailConfig{},
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.sendMessage(ctx, tt.payload)
			assert.Equal(t, tt.hasErr, err != nil, err)
		}))
	}
}

func TestSendAllAdministrators(t *testing.T) {
	t.Parallel()

	in := &user.ListAdministratorsInput{
		Limit:  200,
		Offset: 0,
	}
	administrators := uentity.Administrators{
		{AdminID: "admin-id01"},
		{AdminID: "admin-id02"},
	}

	tests := []struct {
		name    string
		setup   func(ctx context.Context, mocks *mocks)
		payload *entity.WorkerPayload
		hasErr  bool
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().ListAdministrators(ctx, in).Return(administrators, int64(2), nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", nil)
				mocks.db.ReceivedQueue.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, queue *entity.ReceivedQueue) error {
						expect := &entity.ReceivedQueue{
							ID:        queue.ID, // ignore
							EventType: entity.EventTypeUnknown,
							UserType:  entity.UserTypeAdministrator,
							UserIDs:   []string{"admin-id01", "admin-id02"},
						}
						assert.Equal(t, expect, queue)
						return nil
					})
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: false,
		},
		{
			name: "success empty",
			setup: func(ctx context.Context, mocks *mocks) {
				administrators := uentity.Administrators{}
				mocks.user.EXPECT().ListAdministrators(ctx, in).Return(administrators, int64(0), nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: false,
		},
		{
			name: "failed to list administrators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().ListAdministrators(ctx, in).Return(nil, int64(0), assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: true,
		},
		{
			name: "failed to create received queue",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().ListAdministrators(ctx, in).Return(administrators, int64(2), nil)
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: true,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().ListAdministrators(ctx, in).Return(administrators, int64(2), nil)
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.sendAllAdministrators(ctx, tt.payload)
			assert.Equal(t, tt.hasErr, err != nil, err)
		}))
	}
}

func TestSendAllCoordinators(t *testing.T) {
	t.Parallel()

	in := &user.ListCoordinatorsInput{
		Limit:  200,
		Offset: 0,
	}
	coordinators := uentity.Coordinators{
		{AdminID: "admin-id01"},
		{AdminID: "admin-id02"},
	}

	tests := []struct {
		name    string
		setup   func(ctx context.Context, mocks *mocks)
		payload *entity.WorkerPayload
		hasErr  bool
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().ListCoordinators(ctx, in).Return(coordinators, int64(2), nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", nil)
				mocks.db.ReceivedQueue.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, queue *entity.ReceivedQueue) error {
						expect := &entity.ReceivedQueue{
							ID:        queue.ID, // ignore
							EventType: entity.EventTypeUnknown,
							UserType:  entity.UserTypeCoordinator,
							UserIDs:   []string{"admin-id01", "admin-id02"},
						}
						assert.Equal(t, expect, queue)
						return nil
					})
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: false,
		},
		{
			name: "success empty",
			setup: func(ctx context.Context, mocks *mocks) {
				coordinators := uentity.Coordinators{}
				mocks.user.EXPECT().ListCoordinators(ctx, in).Return(coordinators, int64(0), nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: false,
		},
		{
			name: "failed to list coordinators",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().ListCoordinators(ctx, in).Return(nil, int64(0), assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: true,
		},
		{
			name: "failed to create received queue",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().ListCoordinators(ctx, in).Return(coordinators, int64(2), nil)
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: true,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().ListCoordinators(ctx, in).Return(coordinators, int64(2), nil)
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.sendAllCoordinators(ctx, tt.payload)
			assert.Equal(t, tt.hasErr, err != nil, err)
		}))
	}
}

func TestSendAllProducers(t *testing.T) {
	t.Parallel()

	in := &user.ListProducersInput{
		Limit:  200,
		Offset: 0,
	}
	producers := uentity.Producers{
		{AdminID: "admin-id01"},
		{AdminID: "admin-id02"},
	}

	tests := []struct {
		name    string
		setup   func(ctx context.Context, mocks *mocks)
		payload *entity.WorkerPayload
		hasErr  bool
	}{
		{
			name: "success",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().ListProducers(ctx, in).Return(producers, int64(2), nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", nil)
				mocks.db.ReceivedQueue.EXPECT().
					Create(ctx, gomock.Any()).
					DoAndReturn(func(ctx context.Context, queue *entity.ReceivedQueue) error {
						expect := &entity.ReceivedQueue{
							ID:        queue.ID, // ignore
							EventType: entity.EventTypeUnknown,
							UserType:  entity.UserTypeProducer,
							UserIDs:   []string{"admin-id01", "admin-id02"},
						}
						assert.Equal(t, expect, queue)
						return nil
					})
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: false,
		},
		{
			name: "success empty",
			setup: func(ctx context.Context, mocks *mocks) {
				producers := uentity.Producers{}
				mocks.user.EXPECT().ListProducers(ctx, in).Return(producers, int64(0), nil)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: false,
		},
		{
			name: "failed to list producers",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().ListProducers(ctx, in).Return(nil, int64(0), assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: true,
		},
		{
			name: "failed to create received queue",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().ListProducers(ctx, in).Return(producers, int64(2), nil)
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, gomock.Any()).Return(assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: true,
		},
		{
			name: "failed to send message",
			setup: func(ctx context.Context, mocks *mocks) {
				mocks.user.EXPECT().ListProducers(ctx, in).Return(producers, int64(2), nil)
				mocks.db.ReceivedQueue.EXPECT().Create(ctx, gomock.Any()).Return(nil)
				mocks.producer.EXPECT().SendMessage(ctx, gomock.Any()).Return("", assert.AnError)
			},
			payload: &entity.WorkerPayload{
				QueueID:   "queue-id",
				EventType: entity.EventTypeUnknown,
				Email:     &entity.MailConfig{},
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, testService(tt.setup, func(ctx context.Context, t *testing.T, service *service) {
			err := service.sendAllProducers(ctx, tt.payload)
			assert.Equal(t, tt.hasErr, err != nil, err)
		}))
	}
}
