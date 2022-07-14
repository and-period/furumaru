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

func TestListAdministrators(t *testing.T) {
	t.Parallel()

	in := &user.ListAdministratorsInput{
		Limit:  20,
		Offset: 0,
	}
	admins := uentity.Administrators{
		{
			ID:            "admin-id01",
			Lastname:      "&.",
			Firstname:     "管理者",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "かんりしゃ",
			Email:         "test-admin01@and-period.jp",
			PhoneNumber:   "+819012345678",
			CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		{
			ID:            "admin-id02",
			Lastname:      "&.",
			Firstname:     "スタッフ",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "すたっふ",
			Email:         "test-admin02@and-period.jp",
			PhoneNumber:   "+819012345678",
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
				mocks.user.EXPECT().ListAdministrators(gomock.Any(), in).Return(admins, int64(2), nil)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AdministratorsResponse{
					Administrators: []*response.Administrator{
						{
							ID:            "admin-id01",
							Lastname:      "&.",
							Firstname:     "管理者",
							LastnameKana:  "あんどどっと",
							FirstnameKana: "かんりしゃ",
							Email:         "test-admin01@and-period.jp",
							PhoneNumber:   "+819012345678",
							CreatedAt:     1640962800,
							UpdatedAt:     1640962800,
						},
						{
							ID:            "admin-id02",
							Lastname:      "&.",
							Firstname:     "スタッフ",
							LastnameKana:  "あんどどっと",
							FirstnameKana: "すたっふ",
							Email:         "test-admin02@and-period.jp",
							PhoneNumber:   "+819012345678",
							CreatedAt:     1640962800,
							UpdatedAt:     1640962800,
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
			name: "failed to get administrators",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListAdministrators(gomock.Any(), in).Return(nil, int64(0), errmock)
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
			const prefix = "/v1/administrators"
			path := fmt.Sprintf("%s%s", prefix, tt.query)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestGetAdministrator(t *testing.T) {
	t.Parallel()

	in := &user.GetAdministratorInput{
		AdministratorID: "admin-id",
	}
	admin := &uentity.Administrator{
		ID:            "admin-id",
		Lastname:      "&.",
		Firstname:     "管理者",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "かんりしゃ",
		Email:         "test-admin01@and-period.jp",
		PhoneNumber:   "+819012345678",
		CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		adminID string
		expect  *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetAdministrator(gomock.Any(), in).Return(admin, nil)
			},
			adminID: "admin-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AdministratorResponse{
					Administrator: &response.Administrator{
						ID:            "admin-id",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
		},
		{
			name: "failed to get admin",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetAdministrator(gomock.Any(), in).Return(nil, errmock)
			},
			adminID: "admin-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const prefix = "/v1/administrators"
			path := fmt.Sprintf("%s/%s", prefix, tt.adminID)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestCreateAdministrator(t *testing.T) {
	t.Parallel()

	in := &user.CreateAdministratorInput{
		Lastname:      "&.",
		Firstname:     "管理者",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "かんりしゃ",
		Email:         "test-admin01@and-period.jp",
		PhoneNumber:   "+819012345678",
	}
	admin := &uentity.Administrator{
		ID:            "admin-id",
		Lastname:      "&.",
		Firstname:     "管理者",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "かんりしゃ",
		Email:         "test-admin01@and-period.jp",
		PhoneNumber:   "+819012345678",
		CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreateAdministratorRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().CreateAdministrator(gomock.Any(), in).Return(admin, nil)
			},
			req: &request.CreateAdministratorRequest{
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-admin01@and-period.jp",
				PhoneNumber:   "+819012345678",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AdministratorResponse{
					Administrator: &response.Administrator{
						ID:            "admin-id",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						PhoneNumber:   "+819012345678",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
		},
		{
			name: "failed to create admin",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().CreateAdministrator(gomock.Any(), in).Return(nil, errmock)
			},
			req: &request.CreateAdministratorRequest{
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-admin01@and-period.jp",
				PhoneNumber:   "+819012345678",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/administrators"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}
