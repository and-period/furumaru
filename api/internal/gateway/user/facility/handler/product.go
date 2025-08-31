package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/facility/response"
	"github.com/gin-gonic/gin"
)

// @tag.name        Product
// @tag.description 商品関連
func (h *handler) productRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/products", h.authentication)

	r.GET("", h.ListProducts)
	r.GET("/:productId", h.GetProduct)
}

// @Summary     商品一覧取得
// @Description 商品の一覧を取得します。
// @Tags        Product
// @Router      /facilities/{facilityId}/products [get]
// @Param       facilityId path string true "施設ID"
// @Param       limit query int64 false "取得件数" default(20)
// @Param       offset query int64 false "取得開始位置" default(0)
// @Produce     json
// @Success     200 {object} response.ProductsResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) ListProducts(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.ProductsResponse{
		Products:     []*response.Product{},
		Coordinators: []*response.Coordinator{},
		Producers:    []*response.Producer{},
		Categories:   []*response.Category{},
		ProductTypes: []*response.ProductType{},
		ProductTags:  []*response.ProductTag{},
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     商品詳細取得
// @Description 商品の詳細を取得します。
// @Tags        Product
// @Router      /facilities/{facilityId}/products/{productId} [get]
// @Param       facilityId path string true "施設ID"
// @Param       productId path string true "商品ID"
// @Produce     json
// @Success     200 {object} response.ProductResponse
// @Failure     404 {object} util.ErrorResponse "商品が見つからない"
func (h *handler) GetProduct(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.ProductResponse{
		ProductTags: []*response.ProductTag{},
	}
	ctx.JSON(http.StatusOK, res)
}
