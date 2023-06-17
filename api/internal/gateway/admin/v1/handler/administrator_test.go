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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListAdministrators(t *testing.T) {
	t.Parallel()

	in := &user.ListAdministratorsInput{
		Limit:  20,
		Offset: 0,
	}
	admins := uentity.Administrators{
		{
			Admin: uentity.Admin{
				ID:            "admin-id01",
				Role:          entity.AdminRoleAdministrator,
				Status:        entity.AdminStatusActivated,
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-admin01@and-period.jp",
			},
			AdminID:     "admin-id01",
			PhoneNumber: "+819012345678",
			CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		{
			Admin: uentity.Admin{
				ID:            "admin-id02",
				Role:          entity.AdminRoleAdministrator,
				Status:        entity.AdminStatusActivated,
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
				Email:         "test-admin02@and-period.jp",
			},
			AdminID:     "admin-id02",
			PhoneNumber: "+819012345678",
			CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
							Status:        entity.AdminStatusActivated,
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
							Status:        entity.AdminStatusActivated,
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
				mocks.user.EXPECT().ListAdministrators(gomock.Any(), in).Return(nil, int64(0), assert.AnError)
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
			const format = "/v1/administrators%s"
			path := fmt.Sprintf(format, tt.query)
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
		Admin: uentity.Admin{
			ID:            "admin-id",
			Role:          entity.AdminRoleAdministrator,
			Status:        entity.AdminStatusActivated,
			Lastname:      "&.",
			Firstname:     "管理者",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "かんりしゃ",
			Email:         "test-admin01@and-period.jp",
		},
		AdminID:     "admin-id",
		PhoneNumber: "+819012345678",
		CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
						Status:        entity.AdminStatusActivated,
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
				mocks.user.EXPECT().GetAdministrator(gomock.Any(), in).Return(nil, assert.AnError)
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
			const format = "/v1/administrators/%s"
			path := fmt.Sprintf(format, tt.adminID)
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
		Admin: uentity.Admin{
			ID:            "admin-id",
			Lastname:      "&.",
			Firstname:     "管理者",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "かんりしゃ",
			Email:         "test-admin01@and-period.jp",
		},
		AdminID:     "admin-id",
		PhoneNumber: "+819012345678",
		CreatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:   jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
				mocks.user.EXPECT().CreateAdministrator(gomock.Any(), in).Return(nil, assert.AnError)
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

func TestUpdateAdministrator(t *testing.T) {
	t.Parallel()

	in := &user.UpdateAdministratorInput{
		AdministratorID: "administrator-id",
		Lastname:        "&.",
		Firstname:       "管理者",
		LastnameKana:    "あんどどっと",
		FirstnameKana:   "かんりしゃ",
		PhoneNumber:     "+819012345678",
	}

	tests := []struct {
		name            string
		setup           func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req             *request.UpdateAdministratorRequest
		administratorID string
		expect          *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().UpdateAdministrator(gomock.Any(), in).Return(nil)
			},
			administratorID: "administrator-id",
			req: &request.UpdateAdministratorRequest{
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				PhoneNumber:   "+819012345678",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update administrator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().UpdateAdministrator(gomock.Any(), in).Return(assert.AnError)
			},
			administratorID: "administrator-id",
			req: &request.UpdateAdministratorRequest{
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
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
			const format = "/v1/administrators/%s"
			path := fmt.Sprintf(format, tt.administratorID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateAdministratorEmail(t *testing.T) {
	t.Parallel()

	in := &user.UpdateAdministratorEmailInput{
		AdministratorID: "administrator-id",
		Email:           "test-admin01@and-period.jp",
	}

	tests := []struct {
		name            string
		setup           func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req             *request.UpdateAdministratorEmailRequest
		administratorID string
		expect          *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().UpdateAdministratorEmail(gomock.Any(), in).Return(nil)
			},
			administratorID: "administrator-id",
			req: &request.UpdateAdministratorEmailRequest{
				Email: "test-admin01@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update administartor email",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().UpdateAdministratorEmail(gomock.Any(), in).Return(assert.AnError)
			},
			administratorID: "administrator-id",
			req: &request.UpdateAdministratorEmailRequest{
				Email: "test-admin01@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/administrators/%s/email"
			path := fmt.Sprintf(format, tt.administratorID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestResetAdministratorPassword(t *testing.T) {
	t.Parallel()

	in := &user.ResetAdministratorPasswordInput{
		AdministratorID: "administrator-id",
	}

	tests := []struct {
		name            string
		setup           func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		administratorID string
		expect          *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ResetAdministratorPassword(gomock.Any(), in).Return(nil)
			},
			administratorID: "administrator-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update reset administrator password",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ResetAdministratorPassword(gomock.Any(), in).Return(assert.AnError)
			},
			administratorID: "administrator-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/administrators/%s/password"
			path := fmt.Sprintf(format, tt.administratorID)
			testPatch(t, tt.setup, tt.expect, path, nil)
		})
	}
}

func TestDeleteAdministrator(t *testing.T) {
	t.Parallel()

	in := &user.DeleteAdministratorInput{
		AdministratorID: "administrator-id",
	}

	tests := []struct {
		name            string
		setup           func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		administratorID string
		expect          *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().DeleteAdministrator(gomock.Any(), in).Return(nil)
			},
			administratorID: "administrator-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to delete administrator",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().DeleteAdministrator(gomock.Any(), in).Return(assert.AnError)
			},
			administratorID: "administrator-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/administrators/%s"
			path := fmt.Sprintf(format, tt.administratorID)
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}
