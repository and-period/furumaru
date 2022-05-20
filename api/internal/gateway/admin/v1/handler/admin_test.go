package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/user"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
)

func TestListAdmins(t *testing.T) {
	t.Parallel()

	in := &user.ListAdminsInput{
		Limit:  20,
		Offset: 0,
		Roles:  []int32{int32(uentity.AdminRoleAdministrator), int32(uentity.AdminRoleProducer)},
	}
	Admins := uentity.Admins{
		{
			ID:            "admin-id01",
			Lastname:      "&.",
			Firstname:     "管理者",
			LastnameKana:  "あんどどっと",
			FirstnameKana: "かんりしゃ",
			Email:         "test-admin01@and-period.jp",
			Role:          uentity.AdminRoleAdministrator,
			ThumbnailURL:  "https://and-period.jp",
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
			Role:          uentity.AdminRoleProducer,
			ThumbnailURL:  "https://and-period.jp",
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
				mocks.user.EXPECT().ListAdmins(gomock.Any(), in).Return(Admins, nil)
			},
			query: "?roles=1,2",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AdminsResponse{
					Admins: []*response.Admin{
						{
							ID:            "admin-id01",
							Lastname:      "&.",
							Firstname:     "管理者",
							LastnameKana:  "あんどどっと",
							FirstnameKana: "かんりしゃ",
							Email:         "test-admin01@and-period.jp",
							Role:          int32(service.AdminRoleAdministrator),
							ThumbnailURL:  "https://and-period.jp",
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
							Role:          int32(service.AdminRoleProducer),
							ThumbnailURL:  "https://and-period.jp",
							CreatedAt:     1640962800,
							UpdatedAt:     1640962800,
						},
					},
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
			name: "invalid roles",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
			},
			query: "?roles=a,b",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get admins",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().ListAdmins(gomock.Any(), in).Return(nil, errmock)
			},
			query: "?roles=1,2",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const prefix = "/v1/admins"
			path := fmt.Sprintf("%s%s", prefix, tt.query)
			req := newHTTPRequest(t, http.MethodGet, path, nil)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestGetAdmin(t *testing.T) {
	t.Parallel()

	in := &user.GetAdminInput{
		AdminID: "admin-id",
	}
	admin := &uentity.Admin{
		ID:            "admin-id",
		Lastname:      "&.",
		Firstname:     "管理者",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "かんりしゃ",
		Email:         "test-admin01@and-period.jp",
		Role:          uentity.AdminRoleAdministrator,
		ThumbnailURL:  "https://and-period.jp",
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
				mocks.user.EXPECT().GetAdmin(gomock.Any(), in).Return(admin, nil)
			},
			adminID: "admin-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AdminResponse{
					Admin: &response.Admin{
						ID:            "admin-id",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						Role:          int32(service.AdminRoleAdministrator),
						ThumbnailURL:  "https://and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
		},
		{
			name: "failed to get admin",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().GetAdmin(gomock.Any(), in).Return(nil, errmock)
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
			t.Parallel()
			const prefix = "/v1/admins"
			path := fmt.Sprintf("%s/%s", prefix, tt.adminID)
			req := newHTTPRequest(t, http.MethodGet, path, nil)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestCreateAdmin(t *testing.T) {
	t.Parallel()

	in := &user.CreateAdminInput{
		Lastname:      "&.",
		Firstname:     "管理者",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "かんりしゃ",
		Email:         "test-admin01@and-period.jp",
		Role:          int32(uentity.AdminRoleAdministrator),
	}
	admin := &uentity.Admin{
		ID:            "admin-id",
		Lastname:      "&.",
		Firstname:     "管理者",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "かんりしゃ",
		Email:         "test-admin01@and-period.jp",
		Role:          uentity.AdminRoleAdministrator,
		ThumbnailURL:  "https://and-period.jp",
		CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreateAdminRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().CreateAdmin(gomock.Any(), in).Return(admin, nil)
			},
			req: &request.CreateAdminRequest{
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-admin01@and-period.jp",
				Role:          int32(service.AdminRoleAdministrator),
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AdminResponse{
					Admin: &response.Admin{
						ID:            "admin-id",
						Lastname:      "&.",
						Firstname:     "管理者",
						LastnameKana:  "あんどどっと",
						FirstnameKana: "かんりしゃ",
						Email:         "test-admin01@and-period.jp",
						Role:          int32(service.AdminRoleAdministrator),
						ThumbnailURL:  "https://and-period.jp",
						CreatedAt:     1640962800,
						UpdatedAt:     1640962800,
					},
				},
			},
		},
		{
			name: "failed to create admin",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.user.EXPECT().CreateAdmin(gomock.Any(), in).Return(nil, errmock)
			},
			req: &request.CreateAdminRequest{
				Lastname:      "&.",
				Firstname:     "管理者",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "かんりしゃ",
				Email:         "test-admin01@and-period.jp",
				Role:          int32(service.AdminRoleAdministrator),
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const path = "/v1/admins"
			req := newHTTPRequest(t, http.MethodPost, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestGetUserMe(t *testing.T) {
	t.Parallel()

	admin := &uentity.Admin{
		ID:            "admin-id",
		Lastname:      "&.",
		Firstname:     "管理者",
		LastnameKana:  "あんどどっと",
		FirstnameKana: "かんりしゃ",
		Email:         "test-admin01@and-period.jp",
		Role:          uentity.AdminRoleAdministrator,
		ThumbnailURL:  "https://and-period.jp",
		CreatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:     jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.GetAdminInput{AdminID: idmock}
				mocks.user.EXPECT().GetAdmin(gomock.Any(), in).Return(admin, nil)
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.AdminMeResponse{
					ID:            "admin-id",
					Lastname:      "&.",
					Firstname:     "管理者",
					LastnameKana:  "あんどどっと",
					FirstnameKana: "かんりしゃ",
					Email:         "test-admin01@and-period.jp",
					Role:          int32(service.AdminRoleAdministrator),
					ThumbnailURL:  "https://and-period.jp",
				},
			},
		},
		{
			name: "failed to get admin",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.GetAdminInput{AdminID: idmock}
				mocks.user.EXPECT().GetAdmin(gomock.Any(), in).Return(nil, errmock)
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const path = "/v1/admins/me"
			req := newHTTPRequest(t, http.MethodGet, path, nil)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestUpdateAdminEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.UpdateAdminEmailRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.UpdateAdminEmailInput{
					AccessToken: tokenmock,
					Email:       "test-user@and-period.jp",
				}
				mocks.user.EXPECT().UpdateAdminEmail(gomock.Any(), in).Return(nil)
			},
			req: &request.UpdateAdminEmailRequest{
				Email: "test-user@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update user email",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.UpdateAdminEmailInput{
					AccessToken: tokenmock,
					Email:       "test-user@and-period.jp",
				}
				mocks.user.EXPECT().UpdateAdminEmail(gomock.Any(), in).Return(errmock)
			},
			req: &request.UpdateAdminEmailRequest{
				Email: "test-user@and-period.jp",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const path = "/v1/admins/me/email"
			req := newHTTPRequest(t, http.MethodPatch, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestVerifyAdminEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.VerifyAdminEmailRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.VerifyAdminEmailInput{
					AccessToken: tokenmock,
					VerifyCode:  "123456",
				}
				mocks.user.EXPECT().VerifyAdminEmail(gomock.Any(), in).Return(nil)
			},
			req: &request.VerifyAdminEmailRequest{
				VerifyCode: "123456",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to veirify user email",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.VerifyAdminEmailInput{
					AccessToken: tokenmock,
					VerifyCode:  "123456",
				}
				mocks.user.EXPECT().VerifyAdminEmail(gomock.Any(), in).Return(errmock)
			},
			req: &request.VerifyAdminEmailRequest{
				VerifyCode: "123456",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const path = "/v1/admins/me/email/verified"
			req := newHTTPRequest(t, http.MethodPost, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestUpdateAdminPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.UpdateAdminPasswordRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.UpdateAdminPasswordInput{
					AccessToken:          tokenmock,
					OldPassword:          "!Qaz2wsx",
					NewPassword:          "!Qaz3edc",
					PasswordConfirmation: "!Qaz3edc",
				}
				mocks.user.EXPECT().UpdateAdminPassword(gomock.Any(), in).Return(nil)
			},
			req: &request.UpdateAdminPasswordRequest{
				OldPassword:          "!Qaz2wsx",
				NewPassword:          "!Qaz3edc",
				PasswordConfirmation: "!Qaz3edc",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update user password",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				in := &user.UpdateAdminPasswordInput{
					AccessToken:          tokenmock,
					OldPassword:          "!Qaz2wsx",
					NewPassword:          "!Qaz3edc",
					PasswordConfirmation: "!Qaz3edc",
				}
				mocks.user.EXPECT().UpdateAdminPassword(gomock.Any(), in).Return(errmock)
			},
			req: &request.UpdateAdminPasswordRequest{
				OldPassword:          "!Qaz2wsx",
				NewPassword:          "!Qaz3edc",
				PasswordConfirmation: "!Qaz3edc",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const path = "/v1/admins/me/password"
			req := newHTTPRequest(t, http.MethodPatch, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}
