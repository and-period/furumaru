package handler

import (
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/request"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	mentity "github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *handler) notificationRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/notifications", h.authentication)

	r.GET("", h.ListNotifications)
	r.POST("", h.CreateNotification)
	r.GET("/:notificationId", h.GetNotification)
	r.PATCH("/:notificationId", h.UpdateNotifcation)
	r.DELETE("/:notificationId", h.DeleteNotification)
}

func (h *handler) ListNotifications(ctx *gin.Context) {
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
	orders, err := h.newNotificationOrders(ctx)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	since, err := util.GetQueryInt64(ctx, "since", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}
	until, err := util.GetQueryInt64(ctx, "until", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	in := &messenger.ListNotificationsInput{
		Limit:  limit,
		Offset: offset,
		Since:  jst.ParseFromUnix(since),
		Until:  jst.ParseFromUnix(until),
		Orders: orders,
	}
	notifications, total, err := h.messenger.ListNotifications(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}
	if len(notifications) == 0 {
		res := &response.NotificationsResponse{
			Notifications: []*response.Notification{},
			Admins:        []*response.Admin{},
		}
		ctx.JSON(http.StatusOK, res)
		return
	}

	var (
		admins     service.Admins
		promotions service.Promotions
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		admins, err = h.multiGetAdmins(ectx, notifications.AdminIDs())
		return
	})
	eg.Go(func() (err error) {
		promotions, err = h.multiGetPromotions(ectx, notifications.PromotionIDs())
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	snotifications := service.NewNotifications(notifications)
	snotifications.Fill(promotions.Map())

	res := &response.NotificationsResponse{
		Notifications: snotifications.Response(),
		Admins:        admins.Response(),
		Total:         total,
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) GetNotification(ctx *gin.Context) {
	in := &messenger.GetNotificationInput{
		NotificationID: util.GetParam(ctx, "notificationId"),
	}
	notification, err := h.messenger.GetNotification(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	var (
		admin     *service.Admin
		promotion *service.Promotion
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		admin, err = h.getAdmin(ectx, notification.CreatedBy)
		return
	})
	eg.Go(func() (err error) {
		if notification.Type != mentity.NotificationTypePromotion {
			return
		}
		promotion, err = h.getPromotion(ectx, notification.PromotionID)
		return
	})
	if err := eg.Wait(); err != nil {
		h.httpError(ctx, err)
		return
	}

	snotification := service.NewNotification(notification)
	snotification.Fill(promotion)

	res := &response.NotificationResponse{
		Notification: snotification.Response(),
		Admin:        admin.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newNotificationOrders(ctx *gin.Context) ([]*messenger.ListNotificationsOrder, error) {
	notifications := map[string]mentity.NotificationOrderBy{
		"title":       mentity.NotificationOrderByTitle,
		"publishedAt": mentity.NotificationOrderByPublishedAt,
	}
	params := util.GetOrders(ctx)
	if len(params) == 0 {
		res := []*messenger.ListNotificationsOrder{
			{Key: mentity.NotificationOrderByPublishedAt, OrderByASC: false},
		}
		return res, nil
	}
	res := make([]*messenger.ListNotificationsOrder, len(params))
	for i, p := range params {
		key, ok := notifications[p.Key]
		if !ok {
			return nil, fmt.Errorf("handler: unknown order key. key=%s: %w", p.Key, errInvalidOrderKey)
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
		h.badRequest(ctx, err)
		return
	}
	targets := make([]mentity.NotificationTarget, len(req.Targets))
	for i := range req.Targets {
		targets[i] = mentity.NotificationTarget(req.Targets[i])
	}

	publishedAt := jst.ParseFromUnix(req.PublishedAt)
	in := &messenger.CreateNotificationInput{
		Type:        mentity.NotificationType(req.Type),
		Title:       req.Title,
		Body:        req.Body,
		Note:        req.Note,
		Targets:     targets,
		PublishedAt: publishedAt,
		CreatedBy:   getAdminID(ctx),
		PromotionID: req.PromotionID,
	}

	notification, err := h.messenger.CreateNotification(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
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
		h.badRequest(ctx, err)
		return
	}
	targets := make([]mentity.NotificationTarget, len(req.Targets))
	for i := range req.Targets {
		targets[i] = mentity.NotificationTarget(req.Targets[i])
	}

	publishedAt := jst.ParseFromUnix(req.PublishedAt)
	in := &messenger.UpdateNotificationInput{
		NotificationID: util.GetParam(ctx, "notificationId"),
		Title:          req.Title,
		Body:           req.Body,
		Note:           req.Note,
		Targets:        targets,
		PublishedAt:    publishedAt,
		UpdatedBy:      getAdminID(ctx),
	}
	if err := h.messenger.UpdateNotification(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *handler) DeleteNotification(ctx *gin.Context) {
	in := &messenger.DeleteNotificationInput{
		NotificationID: util.GetParam(ctx, "notificationId"),
	}
	if err := h.messenger.DeleteNotification(ctx, in); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
