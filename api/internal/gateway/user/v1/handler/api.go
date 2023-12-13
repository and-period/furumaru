package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/sentry"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	sessionKey = "session_id"
	sessionTTL = 14 * 24 * 60 * 60 // 14days
)

var errNotFoundCart = errors.New("handler: not found cart")

/**
 * ###############################################
 * handler
 * ###############################################
 */
type Handler interface {
	Routes(rg *gin.RouterGroup) // エンドポイント一覧の定義
}

type Params struct {
	WaitGroup *sync.WaitGroup
	User      user.Service
	Store     store.Service
	Messenger messenger.Service
	Media     media.Service
}

type handler struct {
	appName     string
	env         string
	now         func() time.Time
	generateID  func() string
	logger      *zap.Logger
	sentry      sentry.Client
	waitGroup   *sync.WaitGroup
	sharedGroup *singleflight.Group
	user        user.Service
	store       store.Service
	messenger   messenger.Service
	media       media.Service
}

type options struct {
	appName string
	env     string
	logger  *zap.Logger
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

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithSentry(sentry sentry.Client) Option {
	return func(opts *options) {
		opts.sentry = sentry
	}
}

func NewHandler(params *Params, opts ...Option) Handler {
	dopts := &options{
		appName: "user-gateway",
		env:     "",
		logger:  zap.NewNop(),
		sentry:  sentry.NewFixedMockClient(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &handler{
		appName: dopts.appName,
		env:     dopts.env,
		now:     jst.Now,
		generateID: func() string {
			return uuid.Base58Encode(uuid.New())
		},
		logger:      dopts.logger,
		sentry:      dopts.sentry,
		waitGroup:   params.WaitGroup,
		sharedGroup: &singleflight.Group{},
		user:        params.User,
		store:       params.Store,
		messenger:   params.Messenger,
		media:       params.Media,
	}
}

/**
 * ###############################################
 * routes
 * ###############################################
 */
func (h *handler) Routes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	// 公開エンドポイント
	h.authRoutes(v1)
	h.coordinatorRoutes(v1)
	h.scheduleRoutes(v1)
	h.topRoutes(v1)
	h.postalCodeRoutes(v1)
	h.productRoutes(v1)
	h.statusRoutes(v1)
	// 要認証エンドポイント
	h.addressRoutes(v1)
	h.cartRoutes(v1)
	h.checkoutRoutes(v1)
}

/**
 * ###############################################
 * error handling
 * ###############################################
 */
func (h *handler) httpError(ctx *gin.Context, err error) {
	res, code := util.NewErrorResponse(err)
	h.reportError(ctx, err, res)
	h.filterResponse(res)
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

func (h *handler) filterResponse(res *util.ErrorResponse) {
	if res == nil || !strings.Contains(h.env, "prd") {
		return
	}
	// 本番環境の場合、エラーメッセージは返却しない
	res.Detail = ""
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
			ID:        getUserID(ctx),
			IPAddress: ctx.ClientIP(),
		}),
		sentry.WithTags(map[string]string{
			"app_name":   h.appName,
			"env":        h.env,
			"method":     ctx.Request.Method,
			"path":       ctx.FullPath(),
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
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}

	in := &user.GetUserAuthInput{AccessToken: token}
	auth, err := h.user.GetUserAuth(ctx, in)
	if err != nil || auth.UserID == "" {
		h.unauthorized(ctx, err)
		return
	}

	setAuth(ctx, auth.UserID)

	ctx.Next()
}

func (h *handler) getSessionID(ctx *gin.Context) string {
	sessionID, err := ctx.Cookie(sessionKey)
	if err == nil && sessionID != "" {
		return sessionID
	}
	// セッションIDが取得できない場合、新規IDを生成してCookieへ保存する
	sessionID = h.generateID()
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(sessionKey, sessionID, sessionTTL, "/", "", true, true)
	return sessionID
}

func setAuth(ctx *gin.Context, userID string) {
	if userID != "" {
		ctx.Request.Header.Set("userId", userID)
	}
}

func getUserID(ctx *gin.Context) string {
	return ctx.GetHeader("userId")
}
