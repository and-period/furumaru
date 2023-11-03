package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
)

func (h *handler) postalCodeRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/postal-codes")

	r.GET("/:postalCode", h.SearchPostalCode)
}

func (h *handler) SearchPostalCode(ctx *gin.Context) {
	in := &store.SearchPostalCodeInput{
		PostlCode: util.GetParam(ctx, "postalCode"),
	}
	postalCode, err := h.store.SearchPostalCode(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &response.PostalCodeResponse{
		PostalCode: service.NewPostalCode(postalCode).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
