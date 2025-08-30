package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/gin-gonic/gin"
)

// @tag.name        ContactRead
// @tag.description お問い合わせ既読関連
func (h *handler) contactReadRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/contact-reads", h.authentication)

	r.POST("", h.CreateContactRead)
}

// @Summary     お問い合わせ既読登録
// @Description お問い合わせの既読状態を登録します。
// @Tags        ContactRead
// @Router      /v1/contact-reads [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body request.CreateContactReadRequest true "既読情報"
// @Produce     json
// @Success     200 {object} response.ContactReadResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateContactRead(ctx *gin.Context) {
	req := &request.CreateContactReadRequest{}
	if err := ctx.BindJSON(req); err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &messenger.CreateContactReadInput{
		ContactID: req.ContactID,
		UserID:    req.UserID,
		UserType:  entity.ContactUserType(req.UserType),
	}
	scontactRead, err := h.messenger.CreateContactRead(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	contactRead := service.NewContactRead(scontactRead)

	res := &response.ContactReadResponse{
		ContactRead: contactRead.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
