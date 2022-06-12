package handler

import (
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/golang/mock/gomock"
)

func TestUploadCoordinatorThumbnail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		field  string
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.storage.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("https://and-period.jp/thumbnail.png", nil)
			},
			field: "thumbnail",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.UploadImageResponse{
					URL: "https://and-period.jp/thumbnail.png",
				},
			},
		},
		{
			name:  "invalid field",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			field: "",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const path = "/v1/upload/coordinators/thumbnail"
			req := newMultipartRequest(t, http.MethodPost, path, tt.field)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestUploadCoordinatorHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		field  string
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.storage.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("https://and-period.jp/header.png", nil)
			},
			field: "image",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.UploadImageResponse{
					URL: "https://and-period.jp/header.png",
				},
			},
		},
		{
			name:  "invalid field",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			field: "",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const path = "/v1/upload/coordinators/header"
			req := newMultipartRequest(t, http.MethodPost, path, tt.field)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestUploadProducerThumbnail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		field  string
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.storage.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("https://and-period.jp/thumbnail.png", nil)
			},
			field: "thumbnail",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.UploadImageResponse{
					URL: "https://and-period.jp/thumbnail.png",
				},
			},
		},
		{
			name:  "invalid field",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			field: "",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const path = "/v1/upload/producers/thumbnail"
			req := newMultipartRequest(t, http.MethodPost, path, tt.field)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestUploadProducerHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		field  string
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.storage.EXPECT().
					Upload(gomock.Any(), gomock.Any(), gomock.Any()).
					Return("https://and-period.jp/header.png", nil)
			},
			field: "image",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.UploadImageResponse{
					URL: "https://and-period.jp/header.png",
				},
			},
		},
		{
			name:  "invalid field",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			field: "",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const path = "/v1/upload/producers/header"
			req := newMultipartRequest(t, http.MethodPost, path, tt.field)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}
