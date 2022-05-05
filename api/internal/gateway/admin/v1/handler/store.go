package handler

import (
	"net/http"

	"github.com/and-period/marche/api/internal/gateway/admin/v1/response"
	"github.com/gin-gonic/gin"
)

func (h *apiV1Handler) storeRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListStores)
	arg.GET("/:storeId", h.GetStore)
}

func (h *apiV1Handler) ListStores(ctx *gin.Context) {
	// mock
	res := &response.StoresResponse{
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
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *apiV1Handler) GetStore(ctx *gin.Context) {
	// mock
	res := &response.StoreResponse{
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
	}
	ctx.JSON(http.StatusOK, res)
}
