package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/audit"
	"github.com/gin-gonic/gin"
)

// auditRouteInfo はルートからリソース情報を取得するための定義
type auditRouteInfo struct {
	resourceType string // リソース種別
	idParam      string // リソースIDのパスパラメータ名
}

// auditRouteMap はGinのルートパターンからリソース情報へのマッピング
var auditRouteMap = map[string]auditRouteInfo{
	// 管理者
	"/v1/administrators":                   {resourceType: "administrator", idParam: ""},
	"/v1/administrators/:adminId":          {resourceType: "administrator", idParam: "adminId"},
	"/v1/administrators/:adminId/email":    {resourceType: "administrator", idParam: "adminId"},
	"/v1/administrators/:adminId/password": {resourceType: "administrator", idParam: "adminId"},
	// コーディネータ
	"/v1/coordinators": {resourceType: "coordinator", idParam: ""},
	"/v1/coordinators/:coordinatorId": {
		resourceType: "coordinator",
		idParam:      "coordinatorId",
	},
	"/v1/coordinators/:coordinatorId/email": {
		resourceType: "coordinator",
		idParam:      "coordinatorId",
	},
	"/v1/coordinators/:coordinatorId/password": {
		resourceType: "coordinator",
		idParam:      "coordinatorId",
	},
	// 生産者
	"/v1/producers":             {resourceType: "producer", idParam: ""},
	"/v1/producers/:producerId": {resourceType: "producer", idParam: "producerId"},
	// 商品
	"/v1/products":            {resourceType: "product", idParam: ""},
	"/v1/products/:productId": {resourceType: "product", idParam: "productId"},
	// 商品タイプ
	"/v1/product-types":                {resourceType: "product_type", idParam: ""},
	"/v1/product-types/:productTypeId": {resourceType: "product_type", idParam: "productTypeId"},
	// 商品タグ
	"/v1/product-tags":               {resourceType: "product_tag", idParam: ""},
	"/v1/product-tags/:productTagId": {resourceType: "product_tag", idParam: "productTagId"},
	// カテゴリ
	"/v1/categories":             {resourceType: "category", idParam: ""},
	"/v1/categories/:categoryId": {resourceType: "category", idParam: "categoryId"},
	// 注文
	"/v1/orders/-/export":                             {resourceType: "order", idParam: ""},
	"/v1/orders/:orderId/draft":                       {resourceType: "order", idParam: "orderId"},
	"/v1/orders/:orderId/capture":                     {resourceType: "order", idParam: "orderId"},
	"/v1/orders/:orderId/complete":                    {resourceType: "order", idParam: "orderId"},
	"/v1/orders/:orderId/cancel":                      {resourceType: "order", idParam: "orderId"},
	"/v1/orders/:orderId/refund":                      {resourceType: "order", idParam: "orderId"},
	"/v1/orders/:orderId/fulfillments/:fulfillmentId": {resourceType: "order", idParam: "orderId"},
	// 配送
	"/v1/shippings":                             {resourceType: "shipping", idParam: ""},
	"/v1/shippings/default":                     {resourceType: "shipping", idParam: ""},
	"/v1/shippings/:shippingId":                 {resourceType: "shipping", idParam: "shippingId"},
	"/v1/shippings/:shippingId/activation":      {resourceType: "shipping", idParam: "shippingId"},
	"/v1/coordinators/:coordinatorId/shippings": {resourceType: "shipping", idParam: ""},
	// プロモーション
	"/v1/promotions":              {resourceType: "promotion", idParam: ""},
	"/v1/promotions/:promotionId": {resourceType: "promotion", idParam: "promotionId"},
	// スケジュール
	"/v1/schedules":                      {resourceType: "schedule", idParam: ""},
	"/v1/schedules/:scheduleId":          {resourceType: "schedule", idParam: "scheduleId"},
	"/v1/schedules/:scheduleId/approval": {resourceType: "schedule", idParam: "scheduleId"},
	"/v1/schedules/:scheduleId/publish":  {resourceType: "schedule", idParam: "scheduleId"},
	// ライブ
	"/v1/lives":         {resourceType: "live", idParam: ""},
	"/v1/lives/:liveId": {resourceType: "live", idParam: "liveId"},
	// ライブコメント
	"/v1/lives/:liveId/comments/:commentId": {resourceType: "live_comment", idParam: "commentId"},
	// 配信
	"/v1/schedules/:scheduleId/broadcasts": {
		resourceType: "broadcast",
		idParam:      "scheduleId",
	},
	"/v1/schedules/:scheduleId/broadcasts/archive-video": {
		resourceType: "broadcast",
		idParam:      "scheduleId",
	},
	"/v1/schedules/:scheduleId/broadcasts/static-image": {
		resourceType: "broadcast",
		idParam:      "scheduleId",
	},
	"/v1/schedules/:scheduleId/broadcasts/rtmp": {
		resourceType: "broadcast",
		idParam:      "scheduleId",
	},
	"/v1/schedules/:scheduleId/broadcasts/mp4": {
		resourceType: "broadcast",
		idParam:      "scheduleId",
	},
	"/v1/schedules/:scheduleId/broadcasts/youtube/auth": {
		resourceType: "broadcast",
		idParam:      "scheduleId",
	},
	// 動画
	"/v1/videos":          {resourceType: "video", idParam: ""},
	"/v1/videos/:videoId": {resourceType: "video", idParam: "videoId"},
	// 動画コメント
	"/v1/videos/:videoId/comments/:commentId": {
		resourceType: "video_comment",
		idParam:      "commentId",
	},
	// 通知
	"/v1/notifications":                 {resourceType: "notification", idParam: ""},
	"/v1/notifications/:notificationId": {resourceType: "notification", idParam: "notificationId"},
	// メッセージ
	"/v1/messages":            {resourceType: "message", idParam: ""},
	"/v1/messages/:messageId": {resourceType: "message", idParam: "messageId"},
	// 問い合わせ
	"/v1/contacts":                  {resourceType: "contact", idParam: ""},
	"/v1/contacts/:contactId":       {resourceType: "contact", idParam: "contactId"},
	"/v1/contacts/:contactId/reads": {resourceType: "contact", idParam: "contactId"},
	// スレッド
	"/v1/contacts/:contactId/threads":           {resourceType: "thread", idParam: ""},
	"/v1/contacts/:contactId/threads/:threadId": {resourceType: "thread", idParam: "threadId"},
	// 店舗
	"/v1/shops/:shopId": {resourceType: "shop", idParam: "shopId"},
	// 体験
	"/v1/experiences":               {resourceType: "experience", idParam: ""},
	"/v1/experiences/:experienceId": {resourceType: "experience", idParam: "experienceId"},
	// 体験タイプ
	"/v1/experience-types": {resourceType: "experience_type", idParam: ""},
	"/v1/experience-types/:experienceTypeId": {
		resourceType: "experience_type",
		idParam:      "experienceTypeId",
	},
	// スポット
	"/v1/spots":                  {resourceType: "spot", idParam: ""},
	"/v1/spots/:spotId":          {resourceType: "spot", idParam: "spotId"},
	"/v1/spots/:spotId/approval": {resourceType: "spot", idParam: "spotId"},
	// スポットタイプ
	"/v1/spot-types":             {resourceType: "spot_type", idParam: ""},
	"/v1/spot-types/:spotTypeId": {resourceType: "spot_type", idParam: "spotTypeId"},
	// アップロード
	"/v1/upload/coordinators/thumbnail":       {resourceType: "upload", idParam: ""},
	"/v1/upload/coordinators/header":          {resourceType: "upload", idParam: ""},
	"/v1/upload/coordinators/promotion-video": {resourceType: "upload", idParam: ""},
	"/v1/upload/coordinators/bonus-video":     {resourceType: "upload", idParam: ""},
	"/v1/upload/experiences/image":            {resourceType: "upload", idParam: ""},
	"/v1/upload/experiences/video":            {resourceType: "upload", idParam: ""},
	"/v1/upload/experiences/promotion-video":  {resourceType: "upload", idParam: ""},
	"/v1/upload/producers/thumbnail":          {resourceType: "upload", idParam: ""},
	"/v1/upload/producers/header":             {resourceType: "upload", idParam: ""},
	"/v1/upload/producers/promotion-video":    {resourceType: "upload", idParam: ""},
	"/v1/upload/producers/bonus-video":        {resourceType: "upload", idParam: ""},
	"/v1/upload/products/image":               {resourceType: "upload", idParam: ""},
	"/v1/upload/products/video":               {resourceType: "upload", idParam: ""},
	"/v1/upload/product-types/icon":           {resourceType: "upload", idParam: ""},
	"/v1/upload/schedules/thumbnail":          {resourceType: "upload", idParam: ""},
	"/v1/upload/schedules/image":              {resourceType: "upload", idParam: ""},
	"/v1/upload/schedules/opening-video":      {resourceType: "upload", idParam: ""},
	"/v1/upload/schedules/:scheduleId/broadcasts/archive": {
		resourceType: "upload",
		idParam:      "scheduleId",
	},
	"/v1/upload/schedules/-/broadcasts/live": {resourceType: "upload", idParam: ""},
	"/v1/upload/videos/thumbnail":            {resourceType: "upload", idParam: ""},
	"/v1/upload/videos/file":                 {resourceType: "upload", idParam: ""},
	"/v1/upload/spots/thumbnail":             {resourceType: "upload", idParam: ""},
	// ユーザー
	"/v1/users/:userId": {resourceType: "user", idParam: "userId"},
	// 決済システム
	"/v1/payment-systems/:methodType": {resourceType: "payment_system", idParam: "methodType"},
	// 認証 (自身)
	"/v1/auth/email":       {resourceType: "auth", idParam: ""},
	"/v1/auth/password":    {resourceType: "auth", idParam: ""},
	"/v1/auth/device":      {resourceType: "auth", idParam: ""},
	"/v1/auth/google":      {resourceType: "auth", idParam: ""},
	"/v1/auth/line":        {resourceType: "auth", idParam: ""},
	"/v1/auth/coordinator": {resourceType: "auth", idParam: ""},
	// プロダクトレビュー
	"/v1/product-reviews": {resourceType: "product_review", idParam: ""},
}

// auditMiddleware は監査ログ記録用のGinミドルウェア
func auditMiddleware(w *audit.Writer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// GET, OPTIONS, HEADリクエストは記録しない
		if !isAuditableMethod(ctx.Request.Method) {
			ctx.Next()
			return
		}

		start := time.Now()

		ctx.Next()

		// レスポンス後に監査ログを記録
		route := ctx.FullPath()
		if route == "" {
			return // ルートが見つからない場合はスキップ
		}

		adminID := getAdminID(ctx)
		adminType := getAdminType(ctx)

		info, ok := auditRouteMap[route]
		if !ok {
			info = auditRouteInfo{resourceType: "unknown", idParam: ""}
		}

		var resourceID string
		if info.idParam != "" {
			resourceID = ctx.Param(info.idParam)
		}

		status := ctx.Writer.Status()
		duration := time.Since(start)

		log := entity.NewAuditLog(&entity.NewAuditLogParams{
			AdminID:      adminID,
			AdminType:    adminTypeToEntity(adminType),
			Action:       httpMethodToAction(ctx.Request.Method, route),
			ResourceType: info.resourceType,
			ResourceID:   resourceID,
			Result:       httpStatusToResult(status),
			HttpMethod:   ctx.Request.Method,
			HttpPath:     ctx.Request.URL.Path,
			HttpRoute:    route,
			HttpStatus:   status,
			ClientIP:     ctx.ClientIP(),
			UserAgent:    ctx.Request.UserAgent(),
			RequestID:    ctx.GetHeader("X-Request-Id"),
			DurationMs:   int(duration.Milliseconds()),
		})

		w.Send(log)
	}
}

func isAuditableMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	}

	return false
}

func httpMethodToAction(method, route string) entity.AuditLogAction {
	// 特殊なルートの判定
	if route == "/v1/orders/-/export" {
		return entity.AuditLogActionExport
	}
	// アップロード系
	if len(route) > 10 && route[:10] == "/v1/upload" {
		return entity.AuditLogActionUpload
	}

	switch method {
	case http.MethodPost:
		return entity.AuditLogActionCreate
	case http.MethodPut, http.MethodPatch:
		return entity.AuditLogActionUpdate
	case http.MethodDelete:
		return entity.AuditLogActionDelete
	default:
		return entity.AuditLogActionUnknown
	}
}

func httpStatusToResult(status int) entity.AuditLogResult {
	switch {
	case status >= 200 && status < 300:
		return entity.AuditLogResultSuccess
	case status == http.StatusForbidden:
		return entity.AuditLogResultDenied
	case status >= 400 && status < 500:
		return entity.AuditLogResultFailure
	case status >= 500:
		return entity.AuditLogResultError
	default:
		return entity.AuditLogResultUnknown
	}
}

func adminTypeToEntity(t service.AdminType) entity.AdminType {
	return entity.AdminType(int32(t))
}

// recordAuditLog はハンドラ内から明示的にaudit logを記録するヘルパー
func (h *handler) recordAuditLog(
	ctx *gin.Context,
	action entity.AuditLogAction,
	resourceType string,
	resourceID string,
	result entity.AuditLogResult,
	adminID string,
	adminType entity.AdminType,
) {
	if h.auditWriter == nil {
		return
	}

	log := entity.NewAuditLog(&entity.NewAuditLogParams{
		AdminID:      adminID,
		AdminType:    adminType,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Result:       result,
		HttpMethod:   ctx.Request.Method,
		HttpPath:     ctx.Request.URL.Path,
		HttpRoute:    ctx.FullPath(),
		HttpStatus:   ctx.Writer.Status(),
		ClientIP:     ctx.ClientIP(),
		UserAgent:    ctx.Request.UserAgent(),
		RequestID:    ctx.GetHeader("X-Request-Id"),
	})
	h.auditWriter.Send(log)
}

// @tag.name        AuditLog
// @tag.description 監査ログ関連
func (h *handler) auditLogRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/audit-logs", h.authentication)
	r.GET("", h.ListAuditLogs)
}

// @Summary     監査ログ一覧取得
// @Description 監査ログの一覧を取得します。Administrator のみアクセス可能です。
// @Tags        AuditLog
// @Router      /v1/audit-logs [get]
// @Security    bearerauth
// @Param       limit query integer false "取得上限数(max:200)" default(20)
// @Param       offset query integer false "取得開始位置(min:0)" default(0)
// @Param       adminId query string false "管理者IDでフィルタ"
// @Param       resourceType query string false "リソース種別でフィルタ"
// @Param       action query integer false "操作種別でフィルタ"
// @Param       startAt query integer false "開始日時(unix timestamp)"
// @Param       endAt query integer false "終了日時(unix timestamp)"
// @Produce     json
// @Success     200 {object} types.AuditLogsResponse
// @Failure     401 {object} util.ErrorResponse "認証エラー"
// @Failure     403 {object} util.ErrorResponse "権限エラー"
func (h *handler) ListAuditLogs(ctx *gin.Context) {
	// Administrator のみ許可
	if getAdminType(ctx).Response() != types.AdminTypeAdministrator {
		h.forbidden(ctx, errNotAdministrator)
		return
	}

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

	action, err := util.GetQueryInt32(ctx, "action", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	startAt, err := util.GetQueryInt64(ctx, "startAt", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	endAt, err := util.GetQueryInt64(ctx, "endAt", 0)
	if err != nil {
		h.badRequest(ctx, err)
		return
	}

	var startTime, endTime time.Time
	if startAt > 0 {
		startTime = time.Unix(startAt, 0)
	}

	if endAt > 0 {
		endTime = time.Unix(endAt, 0)
	}

	in := &user.ListAuditLogsInput{
		AdminID:      util.GetQuery(ctx, "adminId", ""),
		ResourceType: util.GetQuery(ctx, "resourceType", ""),
		Action:       entity.AuditLogAction(action),
		StartAt:      startTime,
		EndAt:        endTime,
		Limit:        limit,
		Offset:       offset,
	}

	logs, total, err := h.user.ListAuditLogs(ctx, in)
	if err != nil {
		h.httpError(ctx, err)
		return
	}

	res := &types.AuditLogsResponse{
		AuditLogs: newAuditLogsResponse(logs),
		Total:     total,
	}
	ctx.JSON(http.StatusOK, res)
}

func newAuditLogsResponse(logs entity.AuditLogs) []*types.AuditLog {
	res := make([]*types.AuditLog, len(logs))
	for i, log := range logs {
		res[i] = newAuditLogResponse(log)
	}

	return res
}

func newAuditLogResponse(log *entity.AuditLog) *types.AuditLog {
	return &types.AuditLog{
		ID:            log.ID,
		CreatedAt:     log.CreatedAt.Unix(),
		AdminID:       log.AdminID,
		AdminType:     int32(log.AdminType),
		Action:        int32(log.Action),
		ResourceType:  log.ResourceType,
		ResourceID:    log.ResourceID,
		Result:        int32(log.Result),
		ResultDetail:  log.ResultDetail,
		HttpMethod:    log.HttpMethod,
		HttpPath:      log.HttpPath,
		HttpRoute:     log.HttpRoute,
		HttpStatus:    log.HttpStatus,
		ClientIP:      log.ClientIP,
		UserAgent:     log.UserAgent,
		RequestID:     log.RequestID,
		DurationMs:    log.DurationMs,
		ChangedFields: log.ChangedFields,
	}
}

var errNotAdministrator = errors.New("handler: this user is not administrator")
