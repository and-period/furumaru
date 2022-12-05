package handler

import (
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
				mocks.media.EXPECT().
					GenerateCoordinatorThumbnail(gomock.Any(), gomock.Any()).
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
		{
			name: "failed to generate",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.media.EXPECT().GenerateCoordinatorThumbnail(gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			field: "thumbnail",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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
				mocks.media.EXPECT().
					GenerateCoordinatorHeader(gomock.Any(), gomock.Any()).
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
		{
			name: "failed to generate",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.media.EXPECT().GenerateCoordinatorHeader(gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			field: "image",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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
				mocks.media.EXPECT().
					GenerateProducerThumbnail(gomock.Any(), gomock.Any()).
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
		{
			name: "failed to generate",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.media.EXPECT().GenerateProducerThumbnail(gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			field: "thumbnail",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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
				mocks.media.EXPECT().
					GenerateProducerHeader(gomock.Any(), gomock.Any()).
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
		{
			name: "failed to generate",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.media.EXPECT().GenerateProducerHeader(gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			field: "image",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/upload/producers/header"
			req := newMultipartRequest(t, http.MethodPost, path, tt.field)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}

func TestUploadProductImage(t *testing.T) {
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
				mocks.media.EXPECT().
					GenerateProductMediaImage(gomock.Any(), gomock.Any()).
					Return("https://and-period.jp/image.png", nil)
			},
			field: "image",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.UploadImageResponse{
					URL: "https://and-period.jp/image.png",
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
			const path = "/v1/upload/products/image"
			req := newMultipartRequest(t, http.MethodPost, path, tt.field)
			testHTTP(t, tt.setup, tt.expect, req, withRole(entity.AdminRoleCoordinator))
		})
	}
}

func TestUploadProductTypeIcon(t *testing.T) {
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
				mocks.media.EXPECT().
					GenerateProductTypeIcon(gomock.Any(), gomock.Any()).
					Return("https://and-period.jp/header.png", nil)
			},
			field: "icon",
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
			const path = "/v1/upload/product-types/icon"
			req := newMultipartRequest(t, http.MethodPost, path, tt.field)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}
