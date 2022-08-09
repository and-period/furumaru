package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/gin-gonic/gin"
)

func (h *handler) promotionRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListPromotions)
	arg.POST("", h.CreatePromotion)
	arg.GET("/:promotionId", h.GetPromotion)
	arg.PATCH("/:promotionId", h.UpdatePromotion)
	arg.DELETE("/:promotionId", h.DeletePromotion)
}

func (h *handler) ListPromotions(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.PromotionsResponse{
		Promotions: []*response.Promotion{},
		Total:      0,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetPromotion(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.PromotionResponse{
		Promotion: &response.Promotion{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) CreatePromotion(ctx *gin.Context) {
	req := &request.CreatePromotionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	res := &response.PromotionResponse{
		Promotion: &response.Promotion{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdatePromotion(ctx *gin.Context) {
	req := &request.UpdatePromotionRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}
	// TODO: 詳細の実装
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) DeletePromotion(ctx *gin.Context) {
	// TODO: 詳細の実装
	ctx.JSON(http.StatusNoContent, gin.H{})
}
