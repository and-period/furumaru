package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
)

func (h *handler) notificationRoutes(rg *gin.RouterGroup) {
	arg := rg.Use(h.authentication)
	arg.GET("", h.ListNotifications)
	arg.POST("", h.CreateNotification)
	arg.GET("/:notificationId", h.GetNotification)
	arg.PATCH("/:notificationId", h.UpdateNotifcation)
	arg.DELETE("/:notificationId", h.DeleteNotification)
}

func (h *handler) ListNotifications(ctx *gin.Context) {
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
	orders, err := h.newNotificationOrders(ctx)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	since, err := util.GetQueryInt64(ctx, "since", 0)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	until, err := util.GetQueryInt64(ctx, "until", 0)
	if err != nil {
		badRequest(ctx, err)
		return
	}
	onlyPublished, err := strconv.ParseBool(util.GetQuery(ctx, "onlyPublished", "false"))
	if err != nil {
		badRequest(ctx, err)
		return
	}

	in := &messenger.ListNotificationsInput{
		Limit:         limit,
		Offset:        offset,
		Since:         jst.ParseFromUnix(since),
		Until:         jst.ParseFromUnix(until),
		OnlyPublished: onlyPublished,
		Orders:        orders,
	}
	mnotifications, total, err := h.messenger.ListNotifications(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	notifications := service.NewNotifications(mnotifications)

	adminIn := &user.MultiGetAdminsInput{
		AdminIDs: notifications.AdminIDs(),
	}
	uadmins, err := h.user.MultiGetAdmins(ctx, adminIn)
	if err != nil {
		httpError(ctx, err)
		return
	}
	admins := service.NewAdmins(uadmins)

	notifications.Fill(admins.Map())

	res := &response.NotificationsResponse{
		Notifications: notifications.Response(),
		Total:         total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetNotification(ctx *gin.Context) {
	in := &messenger.GetNotificationInput{
		NotificationID: util.GetParam(ctx, "notificationId"),
	}
	mnotificaiton, err := h.messenger.GetNotification(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}
	notification := service.NewNotification(mnotificaiton)

	adminIn := &user.GetAdminInput{
		AdminID: notification.CreatedBy,
	}
	uadmin, err := h.user.GetAdmin(ctx, adminIn)
	if err != nil {
		httpError(ctx, err)
		return
	}
	admin := service.NewAdmin(uadmin)

	notification.Fill(admin)

	res := &response.NotificationResponse{
		Notification: notification.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newNotificationOrders(ctx *gin.Context) ([]*messenger.ListNotificationsOrder, error) {
	notifications := map[string]entity.NotificationOrderBy{
		"title":       entity.NotificationOrderByTitle,
		"public":      entity.NotificationOrderByPublic,
		"publishedAt": entity.NotificationOrderByPublishedAt,
	}
	params := util.GetOrders(ctx)
	res := make([]*messenger.ListNotificationsOrder, len(params))
	for i, p := range params {
		key, ok := notifications[p.Key]
		if !ok {
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", p.Key, errInvalidOrderkey)
		}
		res[i] = &messenger.ListNotificationsOrder{
			Key:        key,
			OrderByASC: p.Direction == util.OrderByASC,
		}
	}
	return res, nil
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

	publishedAt := jst.ParseFromUnix(req.PublishedAt)
	in := &messenger.CreateNotificationInput{
		CreatedBy:   getAdminID(ctx),
		Title:       req.Title,
		Body:        req.Body,
		Targets:     targets,
		Public:      req.Public,
		PublishedAt: publishedAt,
	}

	notification, err := h.messenger.CreateNotification(ctx, in)
	if err != nil {
		httpError(ctx, err)
		return
	}

	res := &response.NotificationResponse{
		Notification: service.NewNotification(notification).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) UpdateNotifcation(ctx *gin.Context) {
	req := &request.UpdateNotificationRequest{}
	if err := ctx.BindJSON(req); err != nil {
		badRequest(ctx, err)
		return
	}

	targets := make([]entity.TargetType, len(req.Targets))
	for i := range req.Targets {
		targets[i] = entity.TargetType(req.Targets[i])
	}
	publishedAt := jst.ParseFromUnix(req.PublishedAt)
	in := &messenger.UpdateNotificationInput{
		NotificationID: util.GetParam(ctx, "notificationId"),
		Title:          req.Title,
		Body:           req.Body,
		Targets:        targets,
		Public:         req.Public,
		PublishedAt:    publishedAt,
		UpdatedBy:      getAdminID(ctx),
	}
	if err := h.messenger.UpdateNotification(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (h *handler) DeleteNotification(ctx *gin.Context) {
	in := &messenger.DeleteNotificationInput{
		NotificationID: util.GetParam(ctx, "notificationId"),
	}
	if err := h.messenger.DeleteNotification(ctx, in); err != nil {
		httpError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
