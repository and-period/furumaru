package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/gin-gonic/gin"
)

func (h *handler) cartRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.GetCart)
	rg.POST("/-/items", h.AddCartItem)
	rg.DELETE("/:number/items/:productId", h.RemoveCartItem)
}

func (h *handler) GetCart(ctx *gin.Context) {
	res := &response.CartResponse{
		Carts:        []*response.Cart{},
		Coordinators: []*response.Coordinator{},
		Products:     []*response.Product{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) AddCartItem(ctx *gin.Context) {
	req := &request.AddCartItemRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) RemoveCartItem(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, gin.H{})
}
