package handler

import (
	"net/http"
	"testing"

	"github.com/and-period/marche/api/internal/gateway/admin/v1/response"
	"github.com/golang/mock/gomock"
)

func TestUploadStoreThumbnail(t *testing.T) {
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
				body: &response.UploaderResponse{
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
			const path = "/v1/upload/stores/thumbnail"
			req := newMultipartRequest(t, http.MethodPost, path, tt.field)
			testHTTP(t, tt.setup, tt.expect, req)
		})
	}
}
