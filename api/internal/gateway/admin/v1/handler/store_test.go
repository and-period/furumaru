package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/marche/api/internal/gateway/admin/v1/request"
	"github.com/and-period/marche/api/internal/gateway/admin/v1/response"
	sentity "github.com/and-period/marche/api/internal/store/entity"
	store "github.com/and-period/marche/api/internal/store/service"
	uentity "github.com/and-period/marche/api/internal/user/entity"
	user "github.com/and-period/marche/api/internal/user/service"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/golang/mock/gomock"
)

func TestListStores(t *testing.T) {
	t.Parallel()

	in := &store.ListStoresInput{
		Limit:  20,
		Offset: 0,
	}
	stores := sentity.Stores{
		{
			ID:           1,
			Name:         "&.農園",
			ThumbnailURL: "https://and-period.jp",
			CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		{
			ID:           2,
			Name:         "&.水産",
			ThumbnailURL: "https://and-period.jp",
			CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
				mocks.store.EXPECT().ListStores(gomock.Any(), in).Return(stores, nil)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.StoresResponse{
					Stores: []*response.Store{
						{
							ID:           1,
							Name:         "&.農園",
							ThumbnailURL: "https://and-period.jp",
							CreatedAt:    1640962800,
							UpdatedAt:    1640962800,
						},
						{
							ID:           2,
							Name:         "&.水産",
							ThumbnailURL: "https://and-period.jp",
							CreatedAt:    1640962800,
							UpdatedAt:    1640962800,
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
			name: "failed to get stores",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListStores(gomock.Any(), in).Return(nil, errmock)
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
			t.Parallel()
			const prefix = "/v1/stores"
			path := fmt.Sprintf("%s%s", prefix, tt.query)
			req := newHTTPRequest(t, http.MethodGet, path, nil)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestGetStore(t *testing.T) {
	t.Parallel()

	storeIn := &store.GetStoreInput{
		StoreID: 1,
	}
	staffsIn := &store.ListStaffsByStoreIDInput{
		StoreID: 1,
	}
	adminsIn := &user.MultiGetAdminsInput{
		AdminIDs: []string{"kSByoE6FetnPs5Byk3a9Zx", "kSByoE6FetnPs5Byk3a9Za"},
	}
	s := &sentity.Store{
		ID:           1,
		Name:         "&.農園",
		ThumbnailURL: "https://and-period.jp",
		CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	staffs := sentity.Staffs{
		{
			StoreID: 1,
			UserID:  "kSByoE6FetnPs5Byk3a9Zx",
			Role:    sentity.StoreRoleAdministrator,
		},
		{
			StoreID: 1,
			UserID:  "kSByoE6FetnPs5Byk3a9Za",
			Role:    sentity.StoreRoleEditor,
		},
	}
	admins := uentity.Admins{
		{
			ID:        "kSByoE6FetnPs5Byk3a9Zx",
			Lastname:  "&.",
			Firstname: "スタッフ1",
			Email:     "test-user01@and-period.jp",
		},
		{
			ID:        "kSByoE6FetnPs5Byk3a9Za",
			Lastname:  "&.",
			Firstname: "スタッフ2",
			Email:     "test-user02@and-period.jp",
		},
	}

	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		storeID string
		expect  *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetStore(gomock.Any(), storeIn).Return(s, nil)
				mocks.store.EXPECT().ListStaffsByStoreID(gomock.Any(), staffsIn).Return(staffs, nil)
				mocks.user.EXPECT().MultiGetAdmins(gomock.Any(), adminsIn).Return(admins, nil)
			},
			storeID: "1",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.StoreResponse{
					Store: &response.Store{
						ID:           1,
						Name:         "&.農園",
						ThumbnailURL: "https://and-period.jp",
						Staffs: []*response.Staff{
							{
								ID:    "kSByoE6FetnPs5Byk3a9Zx",
								Name:  "&. スタッフ1",
								Email: "test-user01@and-period.jp",
								Role:  1,
							},
							{
								ID:    "kSByoE6FetnPs5Byk3a9Za",
								Name:  "&. スタッフ2",
								Email: "test-user02@and-period.jp",
								Role:  2,
							},
						},
						CreatedAt: 1640962800,
						UpdatedAt: 1640962800,
					},
				},
			},
		},
		{
			name: "success staff empty",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetStore(gomock.Any(), storeIn).Return(s, nil)
				mocks.store.EXPECT().ListStaffsByStoreID(gomock.Any(), staffsIn).Return(sentity.Staffs{}, nil)
			},
			storeID: "1",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.StoreResponse{
					Store: &response.Store{
						ID:           1,
						Name:         "&.農園",
						ThumbnailURL: "https://and-period.jp",
						Staffs:       []*response.Staff{},
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
		},
		{
			name:    "invalid store id",
			setup:   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			storeID: "a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to get store",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetStore(gomock.Any(), storeIn).Return(nil, errmock)
				mocks.store.EXPECT().ListStaffsByStoreID(gomock.Any(), staffsIn).Return(staffs, nil)
				mocks.user.EXPECT().MultiGetAdmins(gomock.Any(), adminsIn).Return(admins, nil)
			},
			storeID: "1",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get staffs",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetStore(gomock.Any(), storeIn).Return(s, nil)
				mocks.store.EXPECT().ListStaffsByStoreID(gomock.Any(), staffsIn).Return(nil, errmock)
			},
			storeID: "1",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to get admins",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetStore(gomock.Any(), storeIn).Return(s, nil)
				mocks.store.EXPECT().ListStaffsByStoreID(gomock.Any(), staffsIn).Return(staffs, nil)
				mocks.user.EXPECT().MultiGetAdmins(gomock.Any(), adminsIn).Return(nil, errmock)
			},
			storeID: "1",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const prefix = "/v1/stores"
			path := fmt.Sprintf("%s/%s", prefix, tt.storeID)
			req := newHTTPRequest(t, http.MethodGet, path, nil)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestCreateStore(t *testing.T) {
	t.Parallel()

	in := &store.CreateStoreInput{
		Name: "&.農園",
	}
	store := &sentity.Store{
		ID:           1,
		Name:         "&.農園",
		ThumbnailURL: "https://and-period.jp/thumbnail.png",
		CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreateStoreRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().CreateStore(gomock.Any(), in).Return(store, nil)
			},
			req: &request.CreateStoreRequest{
				Name: "&.農園",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.StoreResponse{
					Store: &response.Store{
						ID:           1,
						Name:         "&.農園",
						ThumbnailURL: "https://and-period.jp/thumbnail.png",
						Staffs:       []*response.Staff{},
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
		},
		{
			name: "failed to create store",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().CreateStore(gomock.Any(), in).Return(nil, errmock)
			},
			req: &request.CreateStoreRequest{
				Name: "&.農園",
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
			const path = "/v1/stores"
			req := newHTTPRequest(t, http.MethodPost, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestUpdateStore(t *testing.T) {
	t.Parallel()

	in := &store.UpdateStoreInput{
		StoreID:      1,
		Name:         "&.農園",
		ThumbnailURL: "https://and-period.jp/thumbnail.png",
	}

	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		storeID string
		req     *request.UpdateStoreRequest
		expect  *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().UpdateStore(gomock.Any(), in).Return(nil)
			},
			storeID: "1",
			req: &request.UpdateStoreRequest{
				Name:         "&.農園",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name:    "invalid store id",
			setup:   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			storeID: "a",
			req: &request.UpdateStoreRequest{
				Name:         "&.農園",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
			},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to update store",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().UpdateStore(gomock.Any(), in).Return(errmock)
			},
			storeID: "1",
			req: &request.UpdateStoreRequest{
				Name:         "&.農園",
				ThumbnailURL: "https://and-period.jp/thumbnail.png",
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
			const prefix = "/v1/stores"
			path := fmt.Sprintf("%s/%s", prefix, tt.storeID)
			req := newHTTPRequest(t, http.MethodPatch, path, tt.req)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}
