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

func (h *handler) NotificationRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication())
	arg.POST("", h.CreateNotification)
}

func (h *handler) CreateNotification(ctx *gin.Context) {
	req := &request.CreateNotificationRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}
	targets := make([]entity.TargetType, len(req.Targets))
	for i := range req.Targets {
		targets[i] = entity.TargetType(req.Targets[i])
	}
	// targetsに0があれば、1, 2, 3を詰める
	if isTargetAll(targets) {
		targets = []entity.TargetType{
			entity.PostTargetUsers,
			entity.PostTargetProducers,
			entity.PostTargetCoordinators,
		}
	}

	in := &messenger.CreateNotificationInput{
		CreatedBy:   getAdminID(ctx),
		Title:       req.Title,
		Body:        req.Body,
		Targets:     targets,
		Public:      req.Public,
		PublishedAt: req.PublishedAt,
	}

	notification, err := h.messenger.CreateNotification(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.NotificationResponse{
		Notifiation: service.NewNotification(notification).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func isTargetAll(arr []entity.TargetType) bool {
	for _, v := range arr {
		if v == entity.PostTargetAll {
			return true
		}
	}
	return false
}
