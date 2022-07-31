package handler

import (
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger"
	mentity "github.com/and-period/furumaru/api/internal/messenger/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
)

func TestCreateNotification(t *testing.T) {
	t.Parallel()
	date := jst.Date(2022, 1, 1, 0, 0, 0, 0)

	in := &messenger.CreateNotificationInput{
		CreatedBy: idmock,
		Title:     "キャベツ祭り開催",
		Body:      "旬のキャベツを大安売り",
		Targets: []mentity.TargetType{
			mentity.PostTargetUsers,
			mentity.PostTargetProducers,
		},
		Public:      true,
		PublishedAt: date,
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
		PublishedAt: date,
		CreatedAt:   date,
		UpdatedAt:   date,
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
				mocks.messenger.EXPECT().CreateNotification(gomock.Any(), in).Return(notification, nil)
			},
			req: &request.CreateNotificationRequest{
				Title: "キャベツ祭り開催",
				Body:  "旬のキャベツ大安売り",
				Targets: []request.TargetType{
					request.PostTargetUsers,
					request.PostTargetCoordinators,
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
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/notification"
			testPost(t, tt.setup, tt.expect, path, tt.req, withRole(uentity.AdminRoleCoordinator))
		})
	}
}
