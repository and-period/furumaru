package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/sentry"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	userIDKey  = "userId"
	sessionKey = "session_id"
	sessionTTL = 14 * 24 * 60 * 60 // 14days
)

var (
	errNotFoundCart  = errors.New("handler: not found cart")
	errNotFoundOrder = errors.New("handler: not found order")
)

/**
 * ###############################################
 * handler
 * ###############################################
 */
type Params struct {
	WaitGroup  *sync.WaitGroup
	UserWebURL *url.URL
	User       user.Service
	Store      store.Service
	Messenger  messenger.Service
	Media      media.Service
}

// @title               ふるマル API - 購入者向け
// @description         購入者向けのふるマルAPIです。
// @servers.url         https://api.furumaru-stg.and-period.work
// @servers.description 検証環境
// @servers.url				  https://api.furumaru.and-period.co.jp
// @servers.description 本番環境
// @securitydefinitions.cookieauth
// @securitydefinitions.bearerauth
type handler struct {
	appName          string
	env              string
	cookieBaseDomain string
	now              func() time.Time
	generateID       func() string
	sentry           sentry.Client
	waitGroup        *sync.WaitGroup
	sharedGroup      *singleflight.Group
	userWebURL       func() *url.URL
	user             user.Service
	store            store.Service
	messenger        messenger.Service
	media            media.Service
}

type options struct {
	appName          string
	env              string
	cookieBaseDomain string
	sentry           sentry.Client
}

type Option func(opts *options)

func WithAppName(name string) Option {
	return func(opts *options) {
		opts.appName = name
	}
}

func WithEnvironment(env string) Option {
	return func(opts *options) {
		opts.env = env
	}
}

func WithCookieBaseDomain(domain string) Option {
	return func(opts *options) {
		opts.cookieBaseDomain = domain
	}
}

func WithSentry(sentry sentry.Client) Option {
	return func(opts *options) {
		opts.sentry = sentry
	}
}

func NewHandler(params *Params, opts ...Option) gateway.Handler {
	dopts := &options{
		appName:          "user-gateway",
		env:              "",
		cookieBaseDomain: "",
		sentry:           sentry.NewFixedMockClient(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	userWebURL := func() *url.URL {
		url := *params.UserWebURL // copy
		return &url
	}
	return &handler{
		appName:          dopts.appName,
		env:              dopts.env,
		cookieBaseDomain: dopts.cookieBaseDomain,
		now:              jst.Now,
		generateID: func() string {
			return uuid.Base58Encode(uuid.New())
		},
		sentry:      dopts.sentry,
		waitGroup:   params.WaitGroup,
		sharedGroup: &singleflight.Group{},
		userWebURL:  userWebURL,
		user:        params.User,
		store:       params.Store,
		messenger:   params.Messenger,
		media:       params.Media,
	}
}

func (h *handler) Setup(ctx context.Context) error {
	return nil
}

func (h *handler) Sync(ctx context.Context) error {
	return nil
}

/**
 * ###############################################
 * routes
 * ###############################################
 */
func (h *handler) Routes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1", h.prerequest)
	// 公開エンドポイント
	h.authRoutes(v1)
	h.topRoutes(v1)
	h.scheduleRoutes(v1)
	h.productRoutes(v1)
	h.experienceRoutes(v1)
	h.coordinatorRoutes(v1)
	h.producerRoutes(v1)
	h.promotionRoutes(v1)
	h.postalCodeRoutes(v1)
	h.spotTypeRoutes(v1)
	h.statusRoutes(v1)
	h.videoRoutes(v1)
	// ゲスト用エンドポイント
	h.guestCheckoutRoutes(v1)
	h.guestLiveCommentRoutes(v1)
	h.guestVideoCommentRoutes(v1)
	// 要認証エンドポイント
	h.authUserRoutes(v1)
	h.addressRoutes(v1)
	h.cartRoutes(v1)
	h.checkoutRoutes(v1)
	h.experienceReviewRoutes(v1)
	h.liveCommentRoutes(v1)
	h.orderRoutes(v1)
	h.productReviewRoutes(v1)
	h.spotRoutes(v1)
	h.videoCommentRoutes(v1)
	h.uploadRoutes(v1)
}

/**
 * ###############################################
 * error handling
 * ###############################################
 */
func (h *handler) httpError(ctx *gin.Context, err error) {
	res, code := util.NewErrorResponse(err)
	h.reportError(ctx, err, res)
	ctx.JSON(code, res)
	ctx.Abort()
}

func (h *handler) badRequest(ctx *gin.Context, err error) {
	h.httpError(ctx, status.Error(codes.InvalidArgument, err.Error()))
}

func (h *handler) unauthorized(ctx *gin.Context, err error) {
	h.httpError(ctx, status.Error(codes.Unauthenticated, err.Error()))
}

func (h *handler) forbidden(ctx *gin.Context, err error) {
	h.httpError(ctx, status.Error(codes.PermissionDenied, err.Error()))
}

func (h *handler) notFound(ctx *gin.Context, err error) {
	h.httpError(ctx, status.Error(codes.NotFound, err.Error()))
}

func (h *handler) reportError(ctx *gin.Context, err error, res *util.ErrorResponse) {
	if h.sentry == nil || res.Status < 500 {
		return
	}
	opts := []sentry.ReportOption{
		sentry.WithLevel("error"),
		sentry.WithRequest(ctx.Request),
		sentry.WithFingerprint(
			ctx.Request.Method,
			ctx.FullPath(),
			res.GetDetail(),
		),
		sentry.WithUser(&sentry.User{
			ID:        h.getUserID(ctx),
			IPAddress: ctx.ClientIP(),
			Data:      map[string]string{"sessionId": h.getSessionID(ctx)},
		}),
		sentry.WithTags(map[string]string{
			"app_name":   h.appName,
			"env":        h.env,
			"method":     ctx.Request.Method,
			"path":       ctx.Request.URL.Path,
			"query":      ctx.Request.URL.RawQuery,
			"route":      ctx.FullPath(),
			"user_agent": ctx.Request.UserAgent(),
		}),
	}
	h.waitGroup.Add(1)
	go func(ctx context.Context, opts []sentry.ReportOption) {
		defer h.waitGroup.Done()
		h.sentry.ReportError(ctx, err, opts...)
	}(ctx, opts)
}

/**
 * ###############################################
 * other
 * ###############################################
 */
func (h *handler) prerequest(ctx *gin.Context) {
	h.setAuth(ctx) //nolint:errcheck
	ctx.Next()
}

func (h *handler) authentication(ctx *gin.Context) {
	if err := h.setAuth(ctx); err != nil {
		h.unauthorized(ctx, err)
		return
	}
	ctx.Next()
}

func (h *handler) createBroadcastViewerLog(ctx *gin.Context) {
	agent := ctx.Request.UserAgent()
	if agent == "node" {
		ctx.Next()
		return // サーバーサイドからのリクエストはスキップする
	}
	scheduleID := util.GetParam(ctx, "scheduleId")
	if scheduleID == "" {
		ctx.Next()
		return // 開催スケジュールIDがない場合、ライブ配信と関係ないエンドポイントとなるためスキップ
	}
	in := &media.CreateBroadcastViewerLogInput{
		ScheduleID: scheduleID,
		SessionID:  h.getSessionID(ctx),
		UserID:     h.getUserID(ctx),
		UserAgent:  ctx.Request.UserAgent(),
		ClientIP:   ctx.ClientIP(),
	}
	h.waitGroup.Add(1)
	go func() {
		defer h.waitGroup.Done()
		if err := h.media.CreateBroadcastViewerLog(context.Background(), in); err != nil {
			slog.Error("Failed to create broadcast viewer log", log.Error(err))
		}
	}()
	ctx.Next()
}

func (h *handler) createVideoViewerLog(ctx *gin.Context) {
	agent := ctx.Request.UserAgent()
	if agent == "node" {
		ctx.Next()
		return // サーバーサイドからのリクエストはスキップする
	}
	videoID := util.GetParam(ctx, "videoId")
	if videoID == "" {
		ctx.Next()
		return // オンデマンド動画IDがない場合、オンデマンド配信と関係ないエンドポイントとなるためスキップ
	}
	in := &media.CreateVideoViewerLogInput{
		VideoID:   videoID,
		SessionID: h.getSessionID(ctx),
		UserID:    h.getUserID(ctx),
		UserAgent: ctx.Request.UserAgent(),
		ClientIP:  ctx.ClientIP(),
	}
	h.waitGroup.Add(1)
	go func() {
		defer h.waitGroup.Done()
		if err := h.media.CreateVideoViewerLog(context.Background(), in); err != nil {
			slog.Error("Failed to create video viewer log", log.Error(err))
		}
	}()
	ctx.Next()
}

func (h *handler) setAuth(ctx *gin.Context) error {
	if _, ok := ctx.Get(userIDKey); ok {
		return nil // すでに設定済みの場合はスキップする
	}
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		return err
	}
	in := &user.GetUserAuthInput{AccessToken: token}
	auth, err := h.user.GetUserAuth(ctx, in)
	if err != nil || auth.UserID == "" {
		return err
	}
	ctx.Request.Header.Set("Userid", auth.UserID)
	return nil
}

func (h *handler) getSessionID(ctx *gin.Context) string {
	agent := ctx.Request.UserAgent()
	if agent == "node" {
		return "" // サーバーサイドからのリクエストはセッションIDを生成しない
	}
	sessionID, err := ctx.Cookie(sessionKey)
	if err == nil && sessionID != "" {
		return sessionID
	}
	// セッションIDが取得できない場合、新規IDを生成してCookieへ保存する
	sessionID = h.generateID()
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(sessionKey, sessionID, sessionTTL, "/", h.cookieBaseDomain, true, true)
	return sessionID
}

func (h *handler) getUserID(ctx *gin.Context) string {
	userID := ctx.GetHeader(userIDKey)
	if userID != "" {
		return userID
	}
	// ユーザーIDが取得できない場合、認証処理を行いヘッダーへ詰め直す
	h.setAuth(ctx) //nolint:errcheck
	return ctx.GetHeader(userIDKey)
}
