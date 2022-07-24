package handler

import (
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/gin-gonic/gin"
)

func (h *handler) messageRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.GET("", h.ListMessages)
	arg.GET("/:messageId", h.GetMessage)
}

func (h *handler) ListMessages(ctx *gin.Context) {
	const (
		defaultLimit  = 20
		defaultOffset = 0
	)

	limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
	if err != nil {
		badRequest(ctx, err)
		return
	}

	in := &messenger.ListMessagesInput{
		Limit:    limit,
		Offset:   offset,
		UserType: entity.UserTypeAdmin,
		UserID:   getAdminID(ctx),
	}
	messages, total, err := h.messenger.ListMessages(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.MessagesResponse{
		Messages: service.NewMessages(messages).Response(),
		Total:    total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetMessage(ctx *gin.Context) {
	in := &messenger.GetMessageInput{
		MessageID: util.GetParam(ctx, "messageId"),
		UserType:  entity.UserTypeAdmin,
		UserID:    getAdminID(ctx),
	}
	message, err := h.messenger.GetMessage(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.MessageResponse{
		Message: service.NewMessage(message).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}
