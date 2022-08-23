package handler

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger"
	mentity "github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListNotifications(t *testing.T) {
	t.Parallel()
	var date int64 = 1640962800

	notificationsIn := &messenger.ListNotificationsInput{
		Limit:         20,
		Offset:        0,
		Since:         jst.ParseFromUnix(date),
		Until:         jst.ParseFromUnix(date),
		OnlyPublished: false,
		Orders:        []*messenger.ListNotificationsOrder{},
	}
	notifications := mentity.Notifications{
		{
			ID:          "notification-id",
			CreatedBy:   "admin-id",
			CreatorName: "&. 管理者",
			UpdatedBy:   "admin-id",
			Title:       "キャベツ祭り開催",
			Body:        "旬のキャベツを大安売り",
			Targets: []mentity.TargetType{
				mentity.PostTargetUsers,
				mentity.PostTargetProducers,
			},
			Public:      false,
			PublishedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}
	adminsIn := &user.MultiGetAdminsInput{
		AdminIDs: []string{"admin-id"},
	}
	admins := uentity.Admins{
		{
			ID:            "admin-id",
			Role:          uentity.AdminRoleAdministrator,
			Lastname:      "&.",
			Firstname:     "管理者",
			LastnameKana:  "あんどぴりおど",
			FirstnameKana: "かんりしゃ",
			Email:         "test-admin@and-period.jp",
			CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		query  string
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().ListNotifications(gomock.Any(), notificationsIn).Return(notifications, int64(1), nil)
				mocks.user.EXPECT().MultiGetAdmins(gomock.Any(), adminsIn).Return(admins, nil)
			},
			query: "?since=1640962800&until=1640962800",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.NotificationsResponse{
					Notifications: []*response.Notification{
						{
							ID:          "notification-id",
							CreatedBy:   "admin-id",
							CreatorName: "&. 管理者",
							UpdatedBy:   "admin-id",
							Title:       "キャベツ祭り開催",
							Body:        "旬のキャベツを大安売り",
							Targets: []response.TargetType{
								response.PostTargetUsers,
								response.PostTargetProducers,
							},
							PublishedAt: 1640962800,
							CreatedAt:   1640962800,
							UpdatedAt:   1640962800,
						},
					},
					Total: 1,
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
			query: "?orders=body",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to list notifications",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().ListNotifications(gomock.Any(), notificationsIn).Return(nil, int64(0), errmock)
			},
			query: "?since=1640962800&until=1640962800",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to multi get admins",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().ListNotifications(gomock.Any(), notificationsIn).Return(notifications, int64(1), nil)
				mocks.user.EXPECT().MultiGetAdmins(gomock.Any(), adminsIn).Return(nil, errmock)
			},
			query: "?since=1640962800&until=1640962800",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/notifications%s"
			path := fmt.Sprintf(format, tt.query)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestCreateNotification(t *testing.T) {
	t.Parallel()
	var date int64 = 1640962800

	in := &messenger.CreateNotificationInput{
		CreatedBy: idmock,
		Title:     "キャベツ祭り開催",
		Body:      "旬のキャベツを大安売り",
		Targets: []mentity.TargetType{
			mentity.PostTargetUsers,
			mentity.PostTargetProducers,
		},
		Public:      true,
		PublishedAt: jst.ParseFromUnix(date),
	}
	notification := &mentity.Notification{
		ID:          "notification-id",
		CreatedBy:   idmock,
		CreatorName: "登録者",
		UpdatedBy:   idmock,
		Title:       "キャベツ祭り開催",
		Body:        "旬のキャベツを大安売り",
		Targets: []mentity.TargetType{
			mentity.PostTargetUsers,
			mentity.PostTargetProducers,
		},
		Public:      true,
		PublishedAt: jst.ParseFromUnix(date),
		CreatedAt:   jst.ParseFromUnix(date),
		UpdatedAt:   jst.ParseFromUnix(date),
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreateNotificationRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().
					CreateNotification(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, actual *messenger.CreateNotificationInput) (*mentity.Notification, error) {
						assert.Equal(t, in, actual)
						return notification, nil
					})
			},
			req: &request.CreateNotificationRequest{
				Title: "キャベツ祭り開催",
				Body:  "旬のキャベツを大安売り",
				Targets: []request.TargetType{
					request.PostTargetUsers,
					request.PostTargetProducers,
				},
				PublishedAt: date,
				Public:      true,
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.NotificationResponse{
					Notification: &response.Notification{
						ID:          "notification-id",
						CreatedBy:   idmock,
						CreatorName: "登録者",
						UpdatedBy:   idmock,
						Title:       "キャベツ祭り開催",
						Body:        "旬のキャベツを大安売り",
						Targets: []response.TargetType{
							response.PostTargetUsers,
							response.PostTargetProducers,
						},
						PublishedAt: 1640962800,
						Public:      true,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
			},
		},
		{
			name: "failed to create notification",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.messenger.EXPECT().CreateNotification(gomock.Any(), in).Return(nil, errmock)
			},
			req: &request.CreateNotificationRequest{
				Title: "キャベツ祭り開催",
				Body:  "旬のキャベツを大安売り",
				Targets: []request.TargetType{
					request.PostTargetUsers,
					request.PostTargetProducers,
				},
				PublishedAt: date,
				Public:      true,
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/notifications"
			testPost(t, tt.setup, tt.expect, path, tt.req, withRole(uentity.AdminRoleAdministrator))
		})
	}
}
