package handler

import (
	"context"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger"
	mentity "github.com/and-period/furumaru/api/internal/messenger/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
