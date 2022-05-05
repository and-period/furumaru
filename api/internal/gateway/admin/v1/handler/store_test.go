package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/marche/api/internal/gateway/admin/v1/response"
	"github.com/golang/mock/gomock"
)

func TestListStores(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		query  string
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
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

	tests := []struct {
		name    string
		setup   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		storeID string
		expect  *testResponse
	}{
		{
			name:    "success",
			setup:   func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			storeID: "store-id",
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
								Name:  "&.スタッフ1",
								Email: "test-user01@and-period.jp",
								Role:  1,
							},
							{
								ID:    "kSByoE6FetnPs5Byk3a9Za",
								Name:  "&.スタッフ2",
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
