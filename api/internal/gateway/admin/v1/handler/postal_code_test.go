package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSearchPostalCode(t *testing.T) {
	t.Parallel()

	in := &store.SearchPostalCodeInput{
		PostlCode: "1000014",
	}
	potalCode := &entity.PostalCode{
		PostalCode:     "1000014",
		PrefectureCode: "tokyo",
		Prefecture:     "東京都",
		City:           "千代田区",
		Town:           "永田町",
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		postalCode string
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().SearchPostalCode(gomock.Any(), in).Return(potalCode, nil)
			},
			postalCode: "1000014",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.PostalCodeResponse{
					PostalCode: &response.PostalCode{
						PostalCode:     "1000014",
						PrefectureCode: "tokyo",
						Prefecture:     "東京都",
						City:           "千代田区",
						Town:           "永田町",
					},
				},
			},
		},
		{
			name: "failed to search postal code",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().SearchPostalCode(gomock.Any(), in).Return(nil, assert.AnError)
			},
			postalCode: "1000014",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/postal-codes/%s"
			path := fmt.Sprintf(format, tt.postalCode)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}
