package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/user/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/gin-gonic/gin"
)

// @tag.name        PostalCode
// @tag.description 郵便番号関連
func (h *handler) postalCodeRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/postal-codes")

	r.GET("/:postalCode", h.SearchPostalCode)
}

// @Summary     郵便番号検索
// @Description 郵便番号から住所情報を検索します。
// @Tags        PostalCode
// @Router      /postal-codes/{postalCode} [get]
// @Param       postalCode path string true "郵便番号（ハイフンなし7桁）"
// @Produce     json
// @Success     200 {object} types.PostalCodeResponse
// @Failure     404 {object} util.ErrorResponse "郵便番号が見つかりません"
func (h *handler) SearchPostalCode(ctx *gin.Context) {
	in := &store.SearchPostalCodeInput{
		PostlCode: util.GetParam(ctx, "postalCode"),
	}
	postalCode, err := h.store.SearchPostalCode(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	res := &types.PostalCodeResponse{
		PostalCode: service.NewPostalCode(postalCode).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
