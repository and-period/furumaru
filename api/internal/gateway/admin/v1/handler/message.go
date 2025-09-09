package handler

import (
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	mentity "github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/gin-gonic/gin"
)

// @tag.name        Message
// @tag.description メッセージ関連
func (h *handler) messageRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/messages", h.authentication)

	r.GET("", h.ListMessages)
	r.GET("/:messageId", h.GetMessage)
}

// @Summary     メッセージ一覧取得
// @Description 管理者あてのメッセージ一覧を取得します。ソート順指定が可能です。
// @Tags        Message
// @Router      /v1/messages [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       orders query string false "ソート(type,-type,read,-read,receivedAt,-receivedAt)" example("-receivedAt")
// @Produce     json
// @Success     200 {object} types.MessagesResponse
func (h *handler) ListMessages(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	orders, err := h.newMessageOrders(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &messenger.ListMessagesInput{
		Limit:    limit,
		Offset:   offset,
		UserType: mentity.UserTypeAdmin,
		UserID:   getAdminID(ctx),
		Orders:   orders,
	}
	messages, total, err := h.messenger.ListMessages(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.MessagesResponse{
		Messages: service.NewMessages(messages).Response(),
		Total:    total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newMessageOrders(ctx *gin.Context) ([]*messenger.ListMessagesOrder, error) {
	messages := map[string]mentity.MessageOrderBy{
		"type":       mentity.MessageOrderByType,
		"read":       mentity.MessageOrderByRead,
		"receivedAt": mentity.MessageOrderByReceivedAt,
	}
	params := util.GetOrders(ctx)
	res := make([]*messenger.ListMessagesOrder, len(params))
	for i, p := range params {
		key, ok := messages[p.Key]
		if !ok {
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", p.Key, errInvalidOrderKey)
		}
		res[i] = &messenger.ListMessagesOrder{
			Key:        key,
			OrderByASC: p.Direction == util.OrderByASC,
		}
	}
	return res, nil
}

// @Summary     メッセージ取得
// @Description 指定されたメッセージの詳細情報を取得します。
// @Tags        Message
// @Router      /v1/messages/{messageId} [get]
// @Security    bearerauth
// @Param       messageId path string true "メッセージID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.MessageResponse
// @Failure     404 {object} util.ErrorResponse "メッセージが存在しない"
func (h *handler) GetMessage(ctx *gin.Context) {
	in := &messenger.GetMessageInput{
		MessageID: util.GetParam(ctx, "messageId"),
		UserType:  mentity.UserTypeAdmin,
		UserID:    getAdminID(ctx),
	}
	message, err := h.messenger.GetMessage(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.MessageResponse{
		Message: service.NewMessage(message).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
