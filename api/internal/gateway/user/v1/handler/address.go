package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/gin-gonic/gin"
)

func (h *handler) addressRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("", h.ListAddresses)
	arg.GET("/:addressId", h.GetAddress)
}

func (h *handler) ListAddresses(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.AddressesResponse{
		Addresses: []*response.Address{},
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetAddress(ctx *gin.Context) {
	// TODO: 詳細の実装
	res := &response.AddressResponse{}
	ctx.JSON(http.StatusOK, res)
}
