package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
)

func TestListShippings(t *testing.T) {
	t.Parallel()

	in := &store.ListShippingsInput{
		Limit:  20,
		Offset: 0,
	}
	shippings := entity.Shippings{
		{
			ID:   "shipping-id",
			Name: "デフォルト配送設定",
			Box60Rates: entity.ShippingRates{
				{Number: 1, Name: "東京都", Price: 0, Prefectures: []int64{13}},
			},
			Box60Refrigerated:  500,
			Box60Frozen:        800,
			Box80Rates:         entity.ShippingRates{},
			Box80Refrigerated:  500,
			Box80Frozen:        800,
			Box100Rates:        entity.ShippingRates{},
			Box100Refrigerated: 500,
			Box100Frozen:       800,
			HasFreeShipping:    true,
			FreeShippingRates:  3000,
			CreatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
				mocks.store.EXPECT().ListShippings(gomock.Any(), in).Return(shippings, int64(1), nil)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ShippingsResponse{
					Shippings: []*response.Shipping{
						{
							ID:   "shipping-id",
							Name: "デフォルト配送設定",
							Box60Rates: []*response.ShippingRate{
								{Number: 1, Name: "東京都", Price: 0, Prefectures: []string{"tokyo"}},
							},
							Box60Refrigerated:  500,
							Box60Frozen:        800,
							Box80Rates:         []*response.ShippingRate{},
							Box80Refrigerated:  500,
							Box80Frozen:        800,
							Box100Rates:        []*response.ShippingRate{},
							Box100Refrigerated: 500,
							Box100Frozen:       800,
							HasFreeShipping:    true,
							FreeShippingRates:  3000,
							CreatedAt:          1640962800,
							UpdatedAt:          1640962800,
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
			name: "failed to list shippings",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListShippings(gomock.Any(), in).Return(nil, int64(0), errmock)
			},
			query: "",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to new shippings",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				shippings := entity.Shippings{{Box60Rates: []*entity.ShippingRate{{Prefectures: []int64{0}}}}}
				mocks.store.EXPECT().ListShippings(gomock.Any(), in).Return(shippings, int64(1), nil)
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
			const format = "/v1/shippings%s"
			path := fmt.Sprintf(format, tt.query)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestGetShipping(t *testing.T) {
	t.Parallel()

	in := &store.GetShippingInput{
		ShippingID: "shipping-id",
	}
	shipping := &entity.Shipping{
		ID:   "shipping-id",
		Name: "デフォルト配送設定",
		Box60Rates: entity.ShippingRates{
			{Number: 1, Name: "東京都", Price: 0, Prefectures: []int64{13}},
		},
		Box60Refrigerated:  500,
		Box60Frozen:        800,
		Box80Rates:         entity.ShippingRates{},
		Box80Refrigerated:  500,
		Box80Frozen:        800,
		Box100Rates:        entity.ShippingRates{},
		Box100Refrigerated: 500,
		Box100Frozen:       800,
		HasFreeShipping:    true,
		FreeShippingRates:  3000,
		CreatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		shippingID string
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetShipping(gomock.Any(), in).Return(shipping, nil)
			},
			shippingID: "shipping-id",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ShippingResponse{
					Shipping: &response.Shipping{
						ID:   "shipping-id",
						Name: "デフォルト配送設定",
						Box60Rates: []*response.ShippingRate{
							{Number: 1, Name: "東京都", Price: 0, Prefectures: []string{"tokyo"}},
						},
						Box60Refrigerated:  500,
						Box60Frozen:        800,
						Box80Rates:         []*response.ShippingRate{},
						Box80Refrigerated:  500,
						Box80Frozen:        800,
						Box100Rates:        []*response.ShippingRate{},
						Box100Refrigerated: 500,
						Box100Frozen:       800,
						HasFreeShipping:    true,
						FreeShippingRates:  3000,
						CreatedAt:          1640962800,
						UpdatedAt:          1640962800,
					},
				},
			},
		},
		{
			name: "failed to get shipping",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().GetShipping(gomock.Any(), in).Return(nil, errmock)
			},
			shippingID: "shipping-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to new shipping",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				shipping := &entity.Shipping{Box60Rates: []*entity.ShippingRate{{Prefectures: []int64{0}}}}
				mocks.store.EXPECT().GetShipping(gomock.Any(), in).Return(shipping, nil)
			},
			shippingID: "shipping-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/shippings/%s"
			path := fmt.Sprintf(format, tt.shippingID)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestCreateShipping(t *testing.T) {
	t.Parallel()

	in := &store.CreateShippingInput{
		Name: "デフォルト配送設定",
		Box60Rates: []*store.CreateShippingRate{
			{Name: "東京都", Price: 0, Prefectures: []int64{13}},
		},
		Box60Refrigerated:  500,
		Box60Frozen:        800,
		Box80Rates:         []*store.CreateShippingRate{},
		Box80Refrigerated:  500,
		Box80Frozen:        800,
		Box100Rates:        []*store.CreateShippingRate{},
		Box100Refrigerated: 500,
		Box100Frozen:       800,
		HasFreeShipping:    true,
		FreeShippingRates:  3000,
	}
	shipping := &entity.Shipping{
		ID:   "shipping-id",
		Name: "デフォルト配送設定",
		Box60Rates: entity.ShippingRates{
			{Number: 1, Name: "東京都", Price: 0, Prefectures: []int64{13}},
		},
		Box60Refrigerated:  500,
		Box60Frozen:        800,
		Box80Rates:         entity.ShippingRates{},
		Box80Refrigerated:  500,
		Box80Frozen:        800,
		Box100Rates:        entity.ShippingRates{},
		Box100Refrigerated: 500,
		Box100Frozen:       800,
		HasFreeShipping:    true,
		FreeShippingRates:  3000,
		CreatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt:          jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreateShippingRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().CreateShipping(gomock.Any(), in).Return(shipping, nil)
			},
			req: &request.CreateShippingRequest{
				Name: "デフォルト配送設定",
				Box60Rates: []*request.CreateShippingRate{
					{Name: "東京都", Price: 0, Prefectures: []string{"tokyo"}},
				},
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         []*request.CreateShippingRate{},
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        []*request.CreateShippingRate{},
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.ShippingResponse{
					Shipping: &response.Shipping{
						ID:   "shipping-id",
						Name: "デフォルト配送設定",
						Box60Rates: []*response.ShippingRate{
							{Number: 1, Name: "東京都", Price: 0, Prefectures: []string{"tokyo"}},
						},
						Box60Refrigerated:  500,
						Box60Frozen:        800,
						Box80Rates:         []*response.ShippingRate{},
						Box80Refrigerated:  500,
						Box80Frozen:        800,
						Box100Rates:        []*response.ShippingRate{},
						Box100Refrigerated: 500,
						Box100Frozen:       800,
						HasFreeShipping:    true,
						FreeShippingRates:  3000,
						CreatedAt:          1640962800,
						UpdatedAt:          1640962800,
					},
				},
			},
		},
		{
			name:  "failed to new box 60 rates",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			req: &request.CreateShippingRequest{
				Box60Rates: []*request.CreateShippingRate{
					{Name: "東京都", Price: 0, Prefectures: []string{""}},
				},
			},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:  "failed to new box 80 rates",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			req: &request.CreateShippingRequest{
				Box80Rates: []*request.CreateShippingRate{
					{Name: "東京都", Price: 0, Prefectures: []string{""}},
				},
			},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:  "failed to new box 100 rates",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			req: &request.CreateShippingRequest{
				Box100Rates: []*request.CreateShippingRate{
					{Name: "東京都", Price: 0, Prefectures: []string{""}},
				},
			},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to create shipping",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().CreateShipping(gomock.Any(), in).Return(nil, errmock)
			},
			req: &request.CreateShippingRequest{
				Name: "デフォルト配送設定",
				Box60Rates: []*request.CreateShippingRate{
					{Name: "東京都", Price: 0, Prefectures: []string{"tokyo"}},
				},
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         []*request.CreateShippingRate{},
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        []*request.CreateShippingRate{},
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
		{
			name: "failed to new shipping",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				shipping := &entity.Shipping{Box60Rates: []*entity.ShippingRate{{Prefectures: []int64{0}}}}
				mocks.store.EXPECT().CreateShipping(gomock.Any(), in).Return(shipping, nil)
			},
			req: &request.CreateShippingRequest{
				Name: "デフォルト配送設定",
				Box60Rates: []*request.CreateShippingRate{
					{Name: "東京都", Price: 0, Prefectures: []string{"tokyo"}},
				},
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         []*request.CreateShippingRate{},
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        []*request.CreateShippingRate{},
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/shippings"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateShipping(t *testing.T) {
	t.Parallel()

	in := &store.UpdateShippingInput{
		ShippingID: "shipping-id",
		Name:       "デフォルト配送設定",
		Box60Rates: []*store.UpdateShippingRate{
			{Name: "東京都", Price: 0, Prefectures: []int64{13}},
		},
		Box60Refrigerated:  500,
		Box60Frozen:        800,
		Box80Rates:         []*store.UpdateShippingRate{},
		Box80Refrigerated:  500,
		Box80Frozen:        800,
		Box100Rates:        []*store.UpdateShippingRate{},
		Box100Refrigerated: 500,
		Box100Frozen:       800,
		HasFreeShipping:    true,
		FreeShippingRates:  3000,
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		shippingID string
		req        *request.UpdateShippingRequest
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().UpdateShipping(gomock.Any(), in).Return(nil)
			},
			shippingID: "shipping-id",
			req: &request.UpdateShippingRequest{
				Name: "デフォルト配送設定",
				Box60Rates: []*request.UpdateShippingRate{
					{Name: "東京都", Price: 0, Prefectures: []string{"tokyo"}},
				},
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         []*request.UpdateShippingRate{},
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        []*request.UpdateShippingRate{},
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name:       "failed to new box 60 rates",
			setup:      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			shippingID: "shipping-id",
			req: &request.UpdateShippingRequest{
				Box60Rates: []*request.UpdateShippingRate{
					{Name: "東京都", Price: 0, Prefectures: []string{""}},
				},
			},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:       "failed to new box 80 rates",
			setup:      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			shippingID: "shipping-id",
			req: &request.UpdateShippingRequest{
				Box80Rates: []*request.UpdateShippingRate{
					{Name: "東京都", Price: 0, Prefectures: []string{""}},
				},
			},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name:       "failed to new box 100 rates",
			setup:      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {},
			shippingID: "shipping-id",
			req: &request.UpdateShippingRequest{
				Box100Rates: []*request.UpdateShippingRate{
					{Name: "東京都", Price: 0, Prefectures: []string{""}},
				},
			},
			expect: &testResponse{
				code: http.StatusBadRequest,
			},
		},
		{
			name: "failed to update shipping",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().UpdateShipping(gomock.Any(), in).Return(errmock)
			},
			shippingID: "shipping-id",
			req: &request.UpdateShippingRequest{
				Name: "デフォルト配送設定",
				Box60Rates: []*request.UpdateShippingRate{
					{Name: "東京都", Price: 0, Prefectures: []string{"tokyo"}},
				},
				Box60Refrigerated:  500,
				Box60Frozen:        800,
				Box80Rates:         []*request.UpdateShippingRate{},
				Box80Refrigerated:  500,
				Box80Frozen:        800,
				Box100Rates:        []*request.UpdateShippingRate{},
				Box100Refrigerated: 500,
				Box100Frozen:       800,
				HasFreeShipping:    true,
				FreeShippingRates:  3000,
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/shippings/%s"
			path := fmt.Sprintf(format, tt.shippingID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestDeleteShipping(t *testing.T) {
	t.Parallel()

	in := &store.DeleteShippingInput{
		ShippingID: "shipping-id",
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		shippingID string
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().DeleteShipping(gomock.Any(), in).Return(nil)
			},
			shippingID: "shipping-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to delete shipping",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().DeleteShipping(gomock.Any(), in).Return(errmock)
			},
			shippingID: "shipping-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/shippings/%s"
			path := fmt.Sprintf(format, tt.shippingID)
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}
