package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListPromotions(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	in := &store.ListPromotionsInput{
		Limit:  20,
		Offset: 0,
		Orders: []*store.ListPromotionsOrder{
			{Key: sentity.PromotionOrderByPublishedAt, OrderByASC: false},
			{Key: sentity.PromotionOrderByPublic, OrderByASC: false},
		},
	}
	promotions := sentity.Promotions{
		{
			ID:           "promotion-id",
			Title:        "夏の採れたて野菜マルシェを開催!!",
			Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
			Public:       true,
			PublishedAt:  now,
			DiscountType: sentity.DiscountTypeFreeShipping,
			DiscountRate: 0,
			Code:         "code0001",
			CodeType:     sentity.PromotionCodeTypeOnce,
			StartAt:      now,
			EndAt:        now.AddDate(0, 1, 0),
			CreatedAt:    now,
			UpdatedAt:    now,
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
				mocks.store.EXPECT().ListPromotions(gomock.Any(), in).Return(promotions, int64(1), nil)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.PromotionsResponse{
					Promotions: []*response.Promotion{
						{
							ID:           "promotion-id",
							Title:        "夏の採れたて野菜マルシェを開催!!",
							Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
							Public:       true,
							PublishedAt:  1640962800,
							DiscountType: int32(service.DiscountTypeFreeShipping),
							DiscountRate: 0,
							Code:         "code0001",
							StartAt:      1640962800,
							EndAt:        1643641200,
							CreatedAt:    1640962800,
							UpdatedAt:    1640962800,
						},
					},
					Total: 1,
				},
			},
		},
		{
			name:  "invalid limit",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?limit=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:  "invalid offset",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?offset=a",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:  "invalid orders",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			query: "?orders=public,other",
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failedo to list promotions",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListPromotions(gomock.Any(), in).Return(nil, int64(0), assert.AnError)
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
			const format = "/v1/promotions%s"
			path := fmt.Sprintf(format, tt.query)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestGetPromotion(t *testing.T) {
	t.Parallel()

	now := jst.Date(2022, 1, 1, 0, 0, 0, 0)
	in := &store.GetPromotionInput{
		PromotionID: "promotion-id",
	}
	promotion := &sentity.Promotion{
		ID:           "promotion-id",
		Title:        "夏の採れたて野菜マルシェを開催!!",
		Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
		Public:       true,
		PublishedAt:  now,
		DiscountType: sentity.DiscountTypeFreeShipping,
		DiscountRate: 0,
		Code:         "code0001",
		CodeType:     sentity.PromotionCodeTypeOnce,
		StartAt:      now,
		EndAt:        now.AddDate(0, 1, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	tests := []struct {
		name        string
		setup       func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		promotionID string
		expect      *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetPromotion(gomock.Any(), in).Return(promotion, nil)
			},
			promotionID: "promotion-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.PromotionResponse{
					Promotion: &response.Promotion{
						ID:           "promotion-id",
						Title:        "夏の採れたて野菜マルシェを開催!!",
						Description:  "採れたての夏野菜を紹介するマルシェを開催ます!!",
						Public:       true,
						PublishedAt:  1640962800,
						DiscountType: int32(service.DiscountTypeFreeShipping),
						DiscountRate: 0,
						Code:         "code0001",
						StartAt:      1640962800,
						EndAt:        1643641200,
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
		},
		{
			name: "failed to get promotion",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetPromotion(gomock.Any(), in).Return(nil, assert.AnError)
			},
			promotionID: "promotion-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
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

	in := &store.CreatePromotionInput{
		Title:        "プロモーションタイトル",
		Description:  "プロモーションの詳細です。",
		Public:       true,
		PublishedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
		DiscountType: sentity.DiscountTypeRate,
		DiscountRate: 10,
		Code:         "excode01",
		CodeType:     sentity.PromotionCodeTypeAlways,
		StartAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
		EndAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}
	promotion := &sentity.Promotion{
		ID:           "promotion-id",
		Title:        "プロモーションタイトル",
		Description:  "プロモーションの詳細です。",
		Public:       true,
		PublishedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
		DiscountType: sentity.DiscountTypeRate,
		DiscountRate: 10,
		Code:         "code0001",
		CodeType:     sentity.PromotionCodeTypeAlways,
		StartAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
		EndAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
		CreatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:    jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreatePromotionRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().CreatePromotion(gomock.Any(), in).Return(promotion, nil)
			},
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
					Promotion: &response.Promotion{
						ID:           "promotion-id",
						Title:        "プロモーションタイトル",
						Description:  "プロモーションの詳細です。",
						Public:       true,
						PublishedAt:  1640962800,
						DiscountType: int32(service.DiscountTypeRate),
						DiscountRate: 10,
						Code:         "code0001",
						StartAt:      1640962800,
						EndAt:        1640962800,
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
		},
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().CreatePromotion(gomock.Any(), in).Return(nil, assert.AnError)
			},
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
				code: http.StatusInternalServerError,
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

	in := &store.UpdatePromotionInput{
		PromotionID:  "promotion-id",
		Title:        "プロモーションタイトル",
		Description:  "プロモーションの詳細です。",
		Public:       true,
		PublishedAt:  jst.Date(2022, 1, 1, 0, 0, 0, 0),
		DiscountType: sentity.DiscountTypeRate,
		DiscountRate: 10,
		Code:         "excode01",
		CodeType:     sentity.PromotionCodeTypeAlways,
		StartAt:      jst.Date(2022, 1, 1, 0, 0, 0, 0),
		EndAt:        jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name        string
		setup       func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		promotionID string
		req         *request.UpdatePromotionRequest
		expect      *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().UpdatePromotion(gomock.Any(), in).Return(nil)
			},
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
		{
			name: "failed to update promotion",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().UpdatePromotion(gomock.Any(), in).Return(assert.AnError)
			},
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
				code: http.StatusInternalServerError,
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

	in := &store.DeletePromotionInput{
		PromotionID: "promotion-id",
	}

	tests := []struct {
		name        string
		setup       func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		promotionID string
		expect      *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().DeletePromotion(gomock.Any(), in).Return(nil)
			},
			promotionID: "promotion-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to delete promotion",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().DeletePromotion(gomock.Any(), in).Return(assert.AnError)
			},
			promotionID: "promotion-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
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
