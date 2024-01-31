package handler

import (
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	mentity "github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/gin-gonic/gin"
)

func (h *handler) messageRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/messages", h.authentication)

	r.GET("", h.ListMessages)
	r.GET("/:messageId", h.GetMessage)
}

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

	res := &response.MessagesResponse{
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

	res := &response.MessageResponse{
		Message: service.NewMessage(message).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
