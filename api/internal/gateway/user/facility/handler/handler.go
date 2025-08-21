package handler

import (
	"context"
	"sync"

	"github.com/and-period/furumaru/api/internal/gateway"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/pkg/sentry"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const userIDKey = "userId"

// @title               ふるマル API - 外部宿泊施設向け
// @description         外部宿泊施設向けのふるマルAPIです。（公開エンドポイントについては「ふるマルAPI - 購入者向け」を参照してください）
// @servers.url         https://api.furumaru-stg.and-period.work
// @servers.description 検証環境
// @servers.url				  https://api.furumaru.and-period.co.jp
// @servers.description 本番環境
// @securitydefinitions.bearerauth
type handler struct {
	appName   string
	env       string
	sentry    sentry.Client
	waitGroup *sync.WaitGroup
}

type Params struct {
	WaitGroup *sync.WaitGroup
}

type options struct {
	appName string
	env     string
	sentry  sentry.Client
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

func WithSentry(sentry sentry.Client) Option {
	return func(opts *options) {
		opts.sentry = sentry
	}
}

func NewHandler(params *Params, opts ...Option) gateway.Handler {
	dopts := &options{
		appName: "facility-gateway",
		env:     "",
		sentry:  sentry.NewFixedMockClient(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &handler{
		appName:   dopts.appName,
		env:       dopts.env,
		sentry:    dopts.sentry,
		waitGroup: params.WaitGroup,
	}
}

func (h *handler) Setup(ctx context.Context) error {
	return nil
}

func (h *handler) Sync(ctx context.Context) error {
	return nil
}

func (h *handler) Routes(rg *gin.RouterGroup) {
	g := rg.Group("/facilities/:facilityId")
	// 公開エンドポイント
	h.authRoutes(g)
	// 要認証エンドポイント
	h.authUserRoutes(g)
	h.cartRoutes(g)
	h.checkoutRoutes(g)
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
func (h *handler) authentication(ctx *gin.Context) {
	if err := h.setAuth(ctx); err != nil {
		h.unauthorized(ctx, err)
		return
	}
	ctx.Next()
}

func (h *handler) setAuth(ctx *gin.Context) error {
	if _, ok := ctx.Get(userIDKey); ok {
		return nil // すでに設定済みの場合はスキップする
	}
	// TODO: 詳細の実装
	return nil
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
