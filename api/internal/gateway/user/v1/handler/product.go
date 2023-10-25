package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/gin-gonic/gin"
)

func (h *handler) productRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.ListProducts)
	rg.GET("/:productId", h.GetProduct)
}

func (h *handler) ListProducts(ctx *gin.Context) {
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

func (h *handler) GetProduct(ctx *gin.Context) {
	res := &response.ProductResponse{
		ProductTags: []*response.ProductTag{},
	}
	ctx.JSON(http.StatusOK, res)
}
