package handler

import (
	"fmt"
	"net/http"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/messenger"
	mentity "github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// @tag.name        Notification
// @tag.description 通知関連
func (h *handler) notificationRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/notifications", h.authentication)

	r.GET("", h.ListNotifications)
	r.POST("", h.CreateNotification)
	r.GET("/:notificationId", h.GetNotification)
	r.PATCH("/:notificationId", h.UpdateNotifcation)
	r.DELETE("/:notificationId", h.DeleteNotification)
}

// @Summary     通知一覧取得
// @Description 通知の一覧を取得します。期間や配信日時でのフィルタリング、ソート順指定が可能です。
// @Tags        Notification
// @Router      /v1/notifications [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20) example(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0) example(0)
// @Param       since query integer false "検索開始日時（unixtime）" example("1640962800")
// @Param       until query integer false "検索終了日時（unixtime）" example("1640962800")
// @Param       orders query string false "ソート(title,-title,publishedAt,-publishedAt)" example("-publishedAt")
// @Produce     json
// @Success     200 {object} types.NotificationsResponse
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
		res := &types.NotificationsResponse{
			Notifications: []*types.Notification{},
			Admins:        []*types.Admin{},
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

	res := &types.NotificationsResponse{
		Notifications: snotifications.Response(),
		Admins:        admins.Response(),
		Total:         total,
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     通知取得
// @Description 指定された通知の詳細情報を取得します。
// @Tags        Notification
// @Router      /v1/notifications/{notificationId} [get]
// @Security    bearerauth
// @Param       notificationId path string true "通知ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     200 {object} types.NotificationResponse
// @Failure     404 {object} util.ErrorResponse "通知が存在しない"
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

	res := &types.NotificationResponse{
		Notification: snotification.Response(),
		Admin:        admin.Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *handler) newNotificationOrders(ctx *gin.Context) ([]*messenger.ListNotificationsOrder, error) {
	notifications := map[string]messenger.ListNotificationsOrderKey{
		"title":       messenger.ListNotificationsOrderByTitle,
		"publishedAt": messenger.ListNotificationsOrderByPublishedAt,
	}
	params := util.GetOrders(ctx)
	if len(params) == 0 {
		res := []*messenger.ListNotificationsOrder{
			{Key: messenger.ListNotificationsOrderByPublishedAt, OrderByASC: false},
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

// @Summary     通知登録
// @Description 新しい通知を登録します。配信対象、配信日時などを指定できます。
// @Tags        Notification
// @Router      /v1/notifications [post]
// @Security    bearerauth
// @Accept      json
// @Param       request body types.CreateNotificationRequest true "通知情報"
// @Produce     json
// @Success     200 {object} types.NotificationResponse
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
func (h *handler) CreateNotification(ctx *gin.Context) {
	req := &types.CreateNotificationRequest{}
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

	res := &types.NotificationResponse{
		Notification: service.NewNotification(notification).Response(),
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary     通知更新
// @Description 通知の情報を更新します。
// @Tags        Notification
// @Router      /v1/notifications/{notificationId} [patch]
// @Security    bearerauth
// @Param       notificationId path string true "通知ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Accept      json
// @Param       request body types.UpdateNotificationRequest true "通知情報"
// @Produce     json
// @Success     204
// @Failure     400 {object} util.ErrorResponse "バリデーションエラー"
// @Failure     404 {object} util.ErrorResponse "通知が存在しない"
func (h *handler) UpdateNotifcation(ctx *gin.Context) {
	req := &types.UpdateNotificationRequest{}
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

// @Summary     通知削除
// @Description 通知を削除します。
// @Tags        Notification
// @Router      /v1/notifications/{notificationId} [delete]
// @Security    bearerauth
// @Param       notificationId path string true "通知ID" example("kSByoE6FetnPs5Byk3a9Zx")
// @Produce     json
// @Success     204
// @Failure     404 {object} util.ErrorResponse "通知が存在しない"
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
