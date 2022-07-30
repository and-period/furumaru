package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
)

func TestListCoordinator(t *testing.T) {
	t.Parallel()

	in := &user.ListCoordinatorsInput{
		Limit:  20,
		Offset: 0,
	}
	coordinators := uentity.Coordinators{
		{
			ID:               "coordinator-id01",
			Lastname:         "&.",
			Firstname:        "管理者",
			LastnameKana:     "あんどどっと",
			FirstnameKana:    "かんりしゃ",
			CompanyName:      "&.株式会社",
			StoreName:        "&.農園",
			ThumbnailURL:     "https://and-period.jp/thumbnail.png",
			HeaderURL:        "https://and-period.jp/header.png",
			TwitterAccount:   "twitter-id",
			InstagramAccount: "instagram-id",
			FacebookAccount:  "facebook-id",
			Email:            "test-coordinator@and-period.jp",
			PhoneNumber:      "+819012345678",
			PostalCode:       "1000014",
			Prefecture:       "東京都",
			City:             "千代田区",
			CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		{
			ID:               "coordinator-id02",
			Lastname:         "&.",
			Firstname:        "管理者",
			LastnameKana:     "あんどどっと",
			FirstnameKana:    "かんりしゃ",
			CompanyName:      "&.株式会社",
			StoreName:        "&.農園",
			ThumbnailURL:     "https://and-period.jp/thumbnail.png",
			HeaderURL:        "https://and-period.jp/header.png",
			TwitterAccount:   "twitter-id",
			InstagramAccount: "instagram-id",
			FacebookAccount:  "facebook-id",
			Email:            "test-coordinator@and-period.jp",
			PhoneNumber:      "+819012345678",
			PostalCode:       "1000014",
			Prefecture:       "東京都",
			City:             "千代田区",
			CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
				mocks.user.EXPECT().ListCoordinators(gomock.Any(), in).Return(coordinators, int64(2), nil)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.CoordinatorsResponse{
					Coordinators: []*response.Coordinator{
						{
							ID:               "coordinator-id01",
							Lastname:         "&.",
							Firstname:        "管理者",
							LastnameKana:     "あんどどっと",
							FirstnameKana:    "かんりしゃ",
							CompanyName:      "&.株式会社",
							StoreName:        "&.農園",
							ThumbnailURL:     "https://and-period.jp/thumbnail.png",
							HeaderURL:        "https://and-period.jp/header.png",
							TwitterAccount:   "twitter-id",
							InstagramAccount: "instagram-id",
							FacebookAccount:  "facebook-id",
							Email:            "test-coordinator@and-period.jp",
							PhoneNumber:      "+819012345678",
							PostalCode:       "1000014",
							Prefecture:       "東京都",
							City:             "千代田区",
							CreatedAt:        1640962800,
							UpdatedAt:        1640962800,
						},
						{
							ID:               "coordinator-id02",
							Lastname:         "&.",
							Firstname:        "管理者",
							LastnameKana:     "あんどどっと",
							FirstnameKana:    "かんりしゃ",
							CompanyName:      "&.株式会社",
							StoreName:        "&.農園",
							ThumbnailURL:     "https://and-period.jp/thumbnail.png",
							HeaderURL:        "https://and-period.jp/header.png",
							TwitterAccount:   "twitter-id",
							InstagramAccount: "instagram-id",
							FacebookAccount:  "facebook-id",
							Email:            "test-coordinator@and-period.jp",
							PhoneNumber:      "+819012345678",
							PostalCode:       "1000014",
							Prefecture:       "東京都",
							City:             "千代田区",
							CreatedAt:        1640962800,
							UpdatedAt:        1640962800,
						},
					},
					Total: 2,
				},
			},
		},
		{
			name: "invalid limit",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
			},
			query: "?limit=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "invalid offset",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
			},
			query: "?offset=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get coordinators",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListCoordinators(gomock.Any(), in).Return(nil, int64(0), errmock)
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
			const format = "/v1/coordinators%s"
			path := fmt.Sprintf(format, tt.query)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestGetCoordinator(t *testing.T) {
	t.Parallel()

	in := &user.GetCoordinatorInput{
		CoordinatorID: "coordinator-id",
	}
	coordinator := &uentity.Coordinator{
		ID:               "coordinator-id",
		Lastname:         "&.",
		Firstname:        "管理者",
		LastnameKana:     "あんどどっと",
		FirstnameKana:    "かんりしゃ",
		CompanyName:      "&.株式会社",
		StoreName:        "&.農園",
		ThumbnailURL:     "https://and-period.jp/thumbnail.png",
		HeaderURL:        "https://and-period.jp/header.png",
		TwitterAccount:   "twitter-id",
		InstagramAccount: "instagram-id",
		FacebookAccount:  "facebook-id",
		Email:            "test-coordinator@and-period.jp",
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
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), in).Return(coordinator, nil)
			},
			coordinatorID: "coordinator-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.CoordinatorResponse{
					Coordinator: &response.Coordinator{
						ID:               "coordinator-id",
						Lastname:         "&.",
						Firstname:        "管理者",
						LastnameKana:     "あんどどっと",
						FirstnameKana:    "かんりしゃ",
						CompanyName:      "&.株式会社",
						StoreName:        "&.農園",
						ThumbnailURL:     "https://and-period.jp/thumbnail.png",
						HeaderURL:        "https://and-period.jp/header.png",
						TwitterAccount:   "twitter-id",
						InstagramAccount: "instagram-id",
						FacebookAccount:  "facebook-id",
						Email:            "test-coordinator@and-period.jp",
						PhoneNumber:      "+819012345678",
						PostalCode:       "1000014",
						Prefecture:       "東京都",
						City:             "千代田区",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
		},
		{
			name: "failed to get coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetCoordinator(gomock.Any(), in).Return(nil, errmock)
			},
			coordinatorID: "coordinator-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/coordinators/%s"
			path := fmt.Sprintf(format, tt.coordinatorID)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestCreateCoordinator(t *testing.T) {
	t.Parallel()

	in := &user.CreateCoordinatorInput{
		Lastname:         "&.",
		Firstname:        "生産者",
		LastnameKana:     "あんどどっと",
		FirstnameKana:    "せいさんしゃ",
		CompanyName:      "&.株式会社",
		StoreName:        "&.農園",
		ThumbnailURL:     "https://and-period.jp/thumbnail.png",
		HeaderURL:        "https://and-period.jp/header.png",
		TwitterAccount:   "twitter-id",
		InstagramAccount: "instagram-id",
		FacebookAccount:  "facebook-id",
		Email:            "test-coordinator@and-period.jp",
		PhoneNumber:      "+819012345678",
		PostalCode:       "1000014",
		Prefecture:       "東京都",
		City:             "千代田区",
		AddressLine1:     "永田町1-7-1",
		AddressLine2:     "",
	}
	coordinator := &uentity.Coordinator{
		ID:               "coordinator-id",
		Lastname:         "&.",
		Firstname:        "管理者",
		LastnameKana:     "あんどどっと",
		FirstnameKana:    "かんりしゃ",
		CompanyName:      "&.株式会社",
		StoreName:        "&.農園",
		ThumbnailURL:     "https://and-period.jp/thumbnail.png",
		HeaderURL:        "https://and-period.jp/header.png",
		TwitterAccount:   "twitter-id",
		InstagramAccount: "instagram-id",
		FacebookAccount:  "facebook-id",
		Email:            "test-coordinator@and-period.jp",
		PhoneNumber:      "+819012345678",
		PostalCode:       "1000014",
		Prefecture:       "東京都",
		City:             "千代田区",
		AddressLine1:     "永田町1-7-1",
		AddressLine2:     "",
		CreatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreateCoordinatorRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().CreateCoordinator(gomock.Any(), in).Return(coordinator, nil)
			},
			req: &request.CreateCoordinatorRequest{
				Lastname:         "&.",
				Firstname:        "生産者",
				LastnameKana:     "あんどどっと",
				FirstnameKana:    "せいさんしゃ",
				CompanyName:      "&.株式会社",
				StoreName:        "&.農園",
				ThumbnailURL:     "https://and-period.jp/thumbnail.png",
				HeaderURL:        "https://and-period.jp/header.png",
				TwitterAccount:   "twitter-id",
				InstagramAccount: "instagram-id",
				FacebookAccount:  "facebook-id",
				Email:            "test-coordinator@and-period.jp",
				PhoneNumber:      "+819012345678",
				PostalCode:       "1000014",
				Prefecture:       "東京都",
				City:             "千代田区",
				AddressLine1:     "永田町1-7-1",
				AddressLine2:     "",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.CoordinatorResponse{
					Coordinator: &response.Coordinator{
						ID:               "coordinator-id",
						Lastname:         "&.",
						Firstname:        "管理者",
						LastnameKana:     "あんどどっと",
						FirstnameKana:    "かんりしゃ",
						CompanyName:      "&.株式会社",
						StoreName:        "&.農園",
						ThumbnailURL:     "https://and-period.jp/thumbnail.png",
						HeaderURL:        "https://and-period.jp/header.png",
						TwitterAccount:   "twitter-id",
						InstagramAccount: "instagram-id",
						FacebookAccount:  "facebook-id",
						Email:            "test-coordinator@and-period.jp",
						PhoneNumber:      "+819012345678",
						PostalCode:       "1000014",
						Prefecture:       "東京都",
						City:             "千代田区",
						AddressLine1:     "永田町1-7-1",
						AddressLine2:     "",
						CreatedAt:        1640962800,
						UpdatedAt:        1640962800,
					},
				},
			},
		},
		{
			name: "failed to create coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().CreateCoordinator(gomock.Any(), in).Return(nil, errmock)
			},
			req: &request.CreateCoordinatorRequest{
				Lastname:         "&.",
				Firstname:        "生産者",
				LastnameKana:     "あんどどっと",
				FirstnameKana:    "せいさんしゃ",
				CompanyName:      "&.株式会社",
				StoreName:        "&.農園",
				ThumbnailURL:     "https://and-period.jp/thumbnail.png",
				HeaderURL:        "https://and-period.jp/header.png",
				TwitterAccount:   "twitter-id",
				InstagramAccount: "instagram-id",
				FacebookAccount:  "facebook-id",
				Email:            "test-coordinator@and-period.jp",
				PhoneNumber:      "+819012345678",
				PostalCode:       "1000014",
				Prefecture:       "東京都",
				City:             "千代田区",
				AddressLine1:     "永田町1-7-1",
				AddressLine2:     "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/coordinators"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateCoordinator(t *testing.T) {
	t.Parallel()

	in := &user.UpdateCoordinatorInput{
		CoordinatorID: "coordinator-id",
		Lastname:      "&.",
		Firstname:     "生産者",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "せいさんしゃ",
		StoreName:     "&.農園",
		ThumbnailURL:  "https://and-period.jp/thumbnail.png",
		HeaderURL:     "https://and-period.jp/header.png",
		PhoneNumber:   "+819012345678",
		PostalCode:    "1000014",
		Prefecture:    "東京都",
		City:          "千代田区",
		AddressLine1:  "永田町1-7-1",
		AddressLine2:  "",
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		coordinatorID string
		req           *request.UpdateCoordinatorRequest
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().UpdateCoordinator(gomock.Any(), in).Return(nil)
			},
			coordinatorID: "coordinator-id",
			req: &request.UpdateCoordinatorRequest{
				Lastname:      "&.",
				Firstname:     "生産者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "せいさんしゃ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				HeaderURL:     "https://and-period.jp/header.png",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update coordinator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().UpdateCoordinator(gomock.Any(), in).Return(errmock)
			},
			coordinatorID: "coordinator-id",
			req: &request.UpdateCoordinatorRequest{
				Lastname:      "&.",
				Firstname:     "生産者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "せいさんしゃ",
				StoreName:     "&.農園",
				ThumbnailURL:  "https://and-period.jp/thumbnail.png",
				HeaderURL:     "https://and-period.jp/header.png",
				PhoneNumber:   "+819012345678",
				PostalCode:    "1000014",
				Prefecture:    "東京都",
				City:          "千代田区",
				AddressLine1:  "永田町1-7-1",
				AddressLine2:  "",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/coordinators/%s"
			path := fmt.Sprintf(format, tt.coordinatorID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateCoordinatorEmail(t *testing.T) {
	t.Parallel()

	in := &user.UpdateCoordinatorEmailInput{
		CoordinatorID: "coordinator-id",
		Email:         "test-producer@and-period.jp",
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		coordinatorID string
		req           *request.UpdateCoordinatorEmailRequest
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().UpdateCoordinatorEmail(gomock.Any(), in).Return(nil)
			},
			coordinatorID: "coordinator-id",
			req: &request.UpdateCoordinatorEmailRequest{
				Email: "test-producer@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update coordinator email",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().UpdateCoordinatorEmail(gomock.Any(), in).Return(errmock)
			},
			coordinatorID: "coordinator-id",
			req: &request.UpdateCoordinatorEmailRequest{
				Email: "test-producer@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/coordinators/%s/email"
			path := fmt.Sprintf(format, tt.coordinatorID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestResetCoordinatorPassword(t *testing.T) {
	t.Parallel()

	in := &user.ResetCoordinatorPasswordInput{
		CoordinatorID: "coordinator-id",
	}

	tests := []struct {
		name          string
		setup         func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		coordinatorID string
		expect        *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ResetCoordinatorPassword(gomock.Any(), in).Return(nil)
			},
			coordinatorID: "coordinator-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to reset coordinator password",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ResetCoordinatorPassword(gomock.Any(), in).Return(errmock)
			},
			coordinatorID: "coordinator-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/coordinators/%s/password"
			path := fmt.Sprintf(format, tt.coordinatorID)
			testPatch(t, tt.setup, tt.expect, path, nil)
		})
	}
}
