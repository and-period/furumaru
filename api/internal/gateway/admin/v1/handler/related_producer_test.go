package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFilterAccessRelateProducers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		options       []testOption
		coordinatorID string
		expect        int
	}{
		{
			name:          "administrator success",
			setup:         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			options:       []testOption{withRole(uentity.AdminRoleAdministrator)},
			coordinatorID: "coordinator-id",
			expect:        http.StatusOK,
		},
		{
			name:          "coordinator success",
			setup:         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			options:       []testOption{withRole(uentity.AdminRoleCoordinator), withAdminID("coordinator-id")},
			coordinatorID: "coordinator-id",
			expect:        http.StatusOK,
		},
		{
			name:          "coordinator forbidden",
			setup:         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			options:       []testOption{withRole(uentity.AdminRoleCoordinator)},
			coordinatorID: "coordinator-id",
			expect:        http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const (
				route  = "/coordinators/:coordinatorId/producers/:producerId"
				format = "/coordinators/%s/producers/producer-id"
			)
			path := fmt.Sprintf(format, tt.coordinatorID)
			testMiddleware(t, tt.setup, route, path, tt.expect, func(h *handler) gin.HandlerFunc {
				return h.filterAccessRelatedProducer
			}, tt.options...)
		})
	}
}

func TestListRelateProducerss(t *testing.T) {
	t.Parallel()

	in := &user.ListProducersInput{
		CoordinatorID: "coordinator-id",
		Limit:         20,
		Offset:        0,
	}
	producers := uentity.Producers{
		{
			Admin: uentity.Admin{
				ID:            "producer-id01",
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-producer@and-period.jp",
			},
			AdminID:       "producer-id01",
			CoordinatorID: "coordinator-id",
			StoreName:     "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
			HeaderURL:     "https://and-period.jp/header.png",
			PhoneNumber:   "+819012345678",
			PostalCode:    "1000014",
			Prefecture:    "東京都",
			City:          "千代田区",
			CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		{
			Admin: uentity.Admin{
				ID:            "producer-id02",
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-producer@and-period.jp",
			},
			AdminID:       "producer-id02",
			CoordinatorID: "coordinator-id",
			StoreName:     "&.農園",
			ThumbnailURL:  "https://and-period.jp/thumbnail.png",
			HeaderURL:     "https://and-period.jp/header.png",
			PhoneNumber:   "+819012345678",
			PostalCode:    "1000014",
			Prefecture:    "東京都",
			City:          "千代田区",
			CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		options       []testOption
		coordinatorID string
		query         string
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListProducers(gomock.Any(), in).Return(producers, int64(2), nil)
			},
			options:       []testOption{withRole(uentity.AdminRoleAdministrator)},
			coordinatorID: "coordinator-id",
			query:         "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ProducersResponse{
					Producers: []*response.Producer{
						{
							ID:            "producer-id01",
							CoordinatorID: "coordinator-id",
							Lastname:      "&.",
							Firstname:     "管理者",
							LastnameKana:  "あんどどっと",
							FirstnameKana: "かんりしゃ",
							StoreName:     "&.農園",
							ThumbnailURL:  "https://and-period.jp/thumbnail.png",
							HeaderURL:     "https://and-period.jp/header.png",
							Email:         "test-producer@and-period.jp",
							PhoneNumber:   "+819012345678",
							PostalCode:    "1000014",
							Prefecture:    "東京都",
							City:          "千代田区",
							CreatedAt:     1640962800,
							UpdatedAt:     1640962800,
						},
						{
							ID:            "producer-id02",
							CoordinatorID: "coordinator-id",
							Lastname:      "&.",
							Firstname:     "管理者",
							LastnameKana:  "あんどどっと",
							FirstnameKana: "かんりしゃ",
							StoreName:     "&.農園",
							ThumbnailURL:  "https://and-period.jp/thumbnail.png",
							HeaderURL:     "https://and-period.jp/header.png",
							Email:         "test-producer@and-period.jp",
							PhoneNumber:   "+819012345678",
							PostalCode:    "1000014",
							Prefecture:    "東京都",
							City:          "千代田区",
							CreatedAt:     1640962800,
							UpdatedAt:     1640962800,
						},
					},
					Total: 2,
				},
			},
		},
		{
			name:          "invalid limit",
			setup:         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			coordinatorID: "coordinator-id",
			query:         "?limit=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:          "invalid offset",
			setup:         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			coordinatorID: "coordinator-id",
			query:         "?offset=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get producers",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListProducers(gomock.Any(), in).Return(nil, int64(0), errmock)
			},
			coordinatorID: "coordinator-id",
			query:         "",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/coordinators/%s/producers%s"
			path := fmt.Sprintf(format, tt.coordinatorID, tt.query)
			testGet(t, tt.setup, tt.expect, path, tt.options...)
		})
	}
}

func TestRelateProducers(t *testing.T) {
	t.Parallel()

	producerIn := &user.RelateProducersInput{
		CoordinatorID: "coordinator-id",
		ProducerIDs:   []string{"producer-id"},
	}
	coordinatorIn := &user.GetCoordinatorInput{
		CoordinatorID: "coordinator-id",
	}
	coordinator := &uentity.Coordinator{
		Admin: uentity.Admin{
			ID:            "coordinator-id",
			Lastname:      "&.",
			Firstname:     "管理者",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "かんりしゃ",
			Email:         "test-coordinator@and-period.jp",
		},
		AdminID:          "coordinator-id",
		CompanyName:      "&.株式会社",
		StoreName:        "&.農園",
		ThumbnailURL:     "https://and-period.jp/thumbnail.png",
		HeaderURL:        "https://and-period.jp/header.png",
		TwitterAccount:   "twitter-id",
		InstagramAccount: "instagram-id",
		FacebookAccount:  "facebook-id",
		PhoneNumber:      "+819012345678",
		PostalCode:       "1000014",
		Prefecture:       "東京都",
		City:             "千代田区",
		CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		options       []testOption
		coordinatorID string
		req           *request.RelateProducersRequest
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().RelateProducers(gomock.Any(), producerIn).Return(nil)
			},
			options:       []testOption{withRole(entity.AdminRoleAdministrator)},
			coordinatorID: "coordinator-id",
			req: &request.RelateProducersRequest{
				ProducerIDs: []string{"producer-id"},
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to get coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, assert.AnError)
			},
			coordinatorID: "coordinator-id",
			req: &request.RelateProducersRequest{
				ProducerIDs: []string{"producer-id"},
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to related producer",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().RelateProducers(gomock.Any(), producerIn).Return(assert.AnError)
			},
			coordinatorID: "coordinator-id",
			req: &request.RelateProducersRequest{
				ProducerIDs: []string{"producer-id"},
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/coordinators/%s/producers"
			path := fmt.Sprintf(format, tt.coordinatorID)
			testPost(t, tt.setup, tt.expect, path, tt.req, tt.options...)
		})
	}
}

func TestUnrelateProducer(t *testing.T) {
	t.Parallel()

	producerIn := &user.UnrelateProducerInput{
		ProducerID: "producer-id",
	}
	coordinatorIn := &user.GetCoordinatorInput{
		CoordinatorID: "coordinator-id",
	}
	coordinator := &uentity.Coordinator{
		Admin: uentity.Admin{
			ID:            "coordinator-id",
			Lastname:      "&.",
			Firstname:     "管理者",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "かんりしゃ",
			Email:         "test-coordinator@and-period.jp",
		},
		AdminID:          "coordinator-id",
		CompanyName:      "&.株式会社",
		StoreName:        "&.農園",
		ThumbnailURL:     "https://and-period.jp/thumbnail.png",
		HeaderURL:        "https://and-period.jp/header.png",
		TwitterAccount:   "twitter-id",
		InstagramAccount: "instagram-id",
		FacebookAccount:  "facebook-id",
		PhoneNumber:      "+819012345678",
		PostalCode:       "1000014",
		Prefecture:       "東京都",
		City:             "千代田区",
		CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		coordinatorID string
		producerID    string
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().UnrelateProducer(gomock.Any(), producerIn).Return(nil)
			},
			coordinatorID: "coordinator-id",
			producerID:    "producer-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to get coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(nil, assert.AnError)
			},
			coordinatorID: "coordinator-id",
			producerID:    "producer-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to unrelated producer",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), coordinatorIn).Return(coordinator, nil)
				mocks.user.EXPECT().UnrelateProducer(gomock.Any(), producerIn).Return(errmock)
			},
			coordinatorID: "coordinator-id",
			producerID:    "producer-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/coordinators/%s/producers/%s"
			path := fmt.Sprintf(format, tt.coordinatorID, tt.producerID)
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}
