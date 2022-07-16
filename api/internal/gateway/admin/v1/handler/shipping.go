package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/gin-gonic/gin"
)

func (h *handler) shippingRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListShippings)
	arg.POST("", h.CreateShipping)
	arg.GET("/:shippingId", h.GetShipping)
	arg.PATCH("/:shippingId", h.UpdateShipping)
	arg.DELETE("/:shippingId", h.DeleteShipping)
}

func (h *handler) ListShippings(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.ShippingsResponse{
		Shippings: []*response.Shipping{},
		Total:     0,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetShipping(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.ShippingResponse{
		Shipping: &response.Shipping{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreateShipping(ctx *gin.Context) {
	req := &request.CreateShippingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	res := &response.ShippingResponse{
		Shipping: &response.Shipping{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateShipping(ctx *gin.Context) {
	req := &request.UpdateShippingRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) DeleteShipping(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.JSON(http.StatusNoContent, gin.H{})
}
