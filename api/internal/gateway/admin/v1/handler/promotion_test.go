package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/golang/mock/gomock"
)

func TestListPromotions(t *testing.T) {
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
				body: &response.PromotionsResponse{
					Promotions: []*response.Promotion{},
					Total:      0,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/promotions%s"
			path := fmt.Sprintf(format, tt.query)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestGetPromotion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		promotionID string
		expect      *testResponse
	}{
		{
			name:        "success",
			setup:       func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			promotionID: "promotion-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.PromotionResponse{
					Promotion: &response.Promotion{},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/promotions/%s"
			path := fmt.Sprintf(format, tt.promotionID)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestCreatePromotion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreatePromotionRequest
		expect *testResponse
	}{
		{
			name:  "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			req: &request.CreatePromotionRequest{
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				PublishedAt:  1640962800,
				DiscountType: 2,
				DiscountRate: 10,
				Code:         "excode01",
				StartAt:      1640962800,
				EndAt:        1640962800,
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.PromotionResponse{
					Promotion: &response.Promotion{},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/promotions"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdatePromotion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		promotionID string
		req         *request.UpdatePromotionRequest
		expect      *testResponse
	}{
		{
			name:        "success",
			setup:       func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			promotionID: "promotion-id",
			req: &request.UpdatePromotionRequest{
				Title:        "プロモーションタイトル",
				Description:  "プロモーションの詳細です。",
				Public:       true,
				PublishedAt:  1640962800,
				DiscountType: 2,
				DiscountRate: 10,
				Code:         "excode01",
				StartAt:      1640962800,
				EndAt:        1640962800,
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/promotions/%s"
			path := fmt.Sprintf(format, tt.promotionID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestDeletePromotion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setup       func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		promotionID string
		expect      *testResponse
	}{
		{
			name:        "success",
			setup:       func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			promotionID: "promotion-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/promotions/%s"
			path := fmt.Sprintf(format, tt.promotionID)
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}
