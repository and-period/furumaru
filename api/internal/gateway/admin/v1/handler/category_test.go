package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/golang/mock/gomock"
)

func TestListCategories(t *testing.T) {
	t.Parallel()

	in := &store.ListCategoriesInput{
		Name:   "野菜",
		Limit:  20,
		Offset: 0,
	}
	categories := sentity.Categories{
		{
			ID:        "category-id01",
			Name:      "野菜",
			CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
		},
		{
			ID:        "category-id02",
			Name:      "野菜",
			CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
			UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
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
				mocks.store.EXPECT().ListCategories(gomock.Any(), in).Return(categories, int64(2), nil)
			},
			query: "?name=野菜",
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.CategoriesResponse{
					Categories: []*response.Category{
						{
							ID:        "category-id01",
							Name:      "野菜",
							CreatedAt: 1640962800,
							UpdatedAt: 1640962800,
						},
						{
							ID:        "category-id02",
							Name:      "野菜",
							CreatedAt: 1640962800,
							UpdatedAt: 1640962800,
						},
					},
					Total: 2,
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
			name: "failed to get categories",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().ListCategories(gomock.Any(), in).Return(nil, int64(0), errmock)
			},
			query: "?name=野菜",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/categories%s"
			path := fmt.Sprintf(format, tt.query)
			testGet(t, tt.setup, tt.expect, path)
		})
	}
}

func TestCreateCategory(t *testing.T) {
	t.Parallel()

	in := &store.CreateCategoryInput{
		Name: "野菜",
	}
	category := &sentity.Category{
		ID:        "category-id",
		Name:      "野菜",
		CreatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
		UpdatedAt: jst.Date(2022, 1, 1, 0, 0, 0, 0),
	}

	tests := []struct {
		name   string
		setup  func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		req    *request.CreateCategoryRequest
		expect *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().CreateCategory(gomock.Any(), in).Return(category, nil)
			},
			req: &request.CreateCategoryRequest{
				Name: "野菜",
			},
			expect: &testResponse{
				code: http.StatusOK,
				body: &response.CategoryResponse{
					Category: &response.Category{
						ID:        "category-id",
						Name:      "野菜",
						CreatedAt: 1640962800,
						UpdatedAt: 1640962800,
					},
				},
			},
		},
		{
			name: "failed to create category",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().CreateCategory(gomock.Any(), in).Return(nil, errmock)
			},
			req: &request.CreateCategoryRequest{
				Name: "野菜",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const path = "/v1/categories"
			testPost(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestUpdateCategory(t *testing.T) {
	t.Parallel()

	in := &store.UpdateCategoryInput{
		CategoryID: "category-id",
		Name:       "野菜",
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		categoryID string
		req        *request.UpdateCategoryRequest
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().UpdateCategory(gomock.Any(), in).Return(nil)
			},
			categoryID: "category-id",
			req: &request.UpdateCategoryRequest{
				Name: "野菜",
			},
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to update category",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().UpdateCategory(gomock.Any(), in).Return(errmock)
			},
			categoryID: "category-id",
			req: &request.UpdateCategoryRequest{
				Name: "野菜",
			},
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/categories/%s"
			path := fmt.Sprintf(format, tt.categoryID)
			testPatch(t, tt.setup, tt.expect, path, tt.req)
		})
	}
}

func TestDeleteCategory(t *testing.T) {
	t.Parallel()

	in := &store.DeleteCategoryInput{
		CategoryID: "category-id",
	}

	tests := []struct {
		name       string
		setup      func(t *testing.T, mocks *mocks, ctrl *gomock.Controller)
		categoryID string
		expect     *testResponse
	}{
		{
			name: "success",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().DeleteCategory(gomock.Any(), in).Return(nil)
			},
			categoryID: "category-id",
			expect: &testResponse{
				code: http.StatusNoContent,
			},
		},
		{
			name: "failed to delete category",
			setup: func(t *testing.T, mocks *mocks, ctrl *gomock.Controller) {
				mocks.store.EXPECT().DeleteCategory(gomock.Any(), in).Return(errmock)
			},
			categoryID: "category-id",
			expect: &testResponse{
				code: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			const format = "/v1/categories/%s"
			path := fmt.Sprintf(format, tt.categoryID)
			testDelete(t, tt.setup, tt.expect, path)
		})
	}
}
