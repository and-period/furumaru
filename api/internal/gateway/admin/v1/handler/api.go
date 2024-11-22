package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/rbac"
	"github.com/and-period/furumaru/api/pkg/sentry"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	sessionKey = "session_id"
	sessionTTL = 24 * 60 * 60 // 1 day
)

var (
	errInvalidOrderKey = errors.New("handler: invalid order key")
	errInvalidSession  = errors.New("handler: invalid session")
)

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
	Enforcer  rbac.Enforcer
	User      user.Service
	Store     store.Service
	Messenger messenger.Service
	Media     media.Service
}

type handler struct {
	appName     string
	env         string
	now         func() time.Time
	logger      *zap.Logger
	sentry      sentry.Client
	waitGroup   *sync.WaitGroup
	sharedGroup *singleflight.Group
	enforcer    rbac.Enforcer
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
		appName: "admin-gateway",
		env:     "",
		logger:  zap.NewNop(),
		sentry:  sentry.NewFixedMockClient(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &handler{
		appName:     dopts.appName,
		env:         dopts.env,
		now:         jst.Now,
		logger:      dopts.logger,
		sentry:      dopts.sentry,
		waitGroup:   params.WaitGroup,
		sharedGroup: &singleflight.Group{},
		enforcer:    params.Enforcer,
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
	h.administratorRoutes(v1)
	h.authRoutes(v1)
	h.broadcastRoutes(v1)
	h.categoryRoutes(v1)
	h.contactRoutes(v1)
	h.contactCategoryRoutes(v1)
	h.contactReadRoutes(v1)
	h.coordinatorRoutes(v1)
	h.experienceRoutes(v1)
	h.experienceTypeRoutes(v1)
	h.liveRoutes(v1)
	h.liveCommentRoutes(v1)
	h.messageRoutes(v1)
	h.notificationRoutes(v1)
	h.orderRoutes(v1)
	h.paymentSystemRoutes(v1)
	h.postalCodeRoutes(v1)
	h.producerRoutes(v1)
	h.productRoutes(v1)
	h.productTagRoutes(v1)
	h.productTypeRoutes(v1)
	h.promotionRoutes(v1)
	h.relatedProducerRoutes(v1)
	h.scheduleRoutes(v1)
	h.shippingRoutes(v1)
	h.spotRoutes(v1)
	h.spotTypeRoutes(v1)
	h.threadRoutes(v1)
	h.uploadRoutes(v1)
	h.userRoutes(v1)
	h.videoRoutes(v1)
	h.videoCommentRoutes(v1)

	// 認証が不要なエンドポイント
	h.guestBroadcastRoutes(v1)
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
			ID:        getAdminID(ctx),
			IPAddress: ctx.ClientIP(),
			Data:      map[string]string{"role": string(getRole(ctx))},
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
	// 認証情報の検証
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		h.unauthorized(ctx, err)
		return
	}

	in := &user.GetAdminAuthInput{AccessToken: token}
	auth, err := h.user.GetAdminAuth(ctx, in)
	if err != nil || auth.AdminID == "" {
		h.unauthorized(ctx, err)
		return
	}
	role := service.NewAdminRole(auth.Role)

	setAuth(ctx, auth.AdminID, role)

	// 認可情報の検証
	if h.enforcer == nil {
		ctx.Next()
		return
	}

	enforce, err := h.enforcer.Enforce(role.String(), ctx.Request.URL.Path, ctx.Request.Method)
	if err != nil {
		h.httpError(ctx, status.Error(codes.Internal, err.Error()))
		return
	}
	if !enforce {
		h.forbidden(ctx, errors.New("handler: you don't have the correct permissions"))
		return
	}

	ctx.Next()
}

func setAuth(ctx *gin.Context, adminID string, role service.AdminRole) {
	if adminID != "" {
		ctx.Request.Header.Set("adminId", adminID)
		ctx.Request.Header.Set("role", strconv.FormatInt(int64(role), 10))
	}
}

func getAdminID(ctx *gin.Context) string {
	return ctx.GetHeader("adminId")
}

func getRole(ctx *gin.Context) service.AdminRole {
	role, _ := strconv.ParseInt(ctx.GetHeader("role"), 10, 64)
	return service.AdminRole(role)
}

func currentAdmin(ctx *gin.Context, adminID string) bool {
	return getAdminID(ctx) == adminID
}

func (h *handler) getSessionID(ctx *gin.Context) (string, error) {
	sessionID, err := ctx.Cookie(sessionKey)
	if err != nil || sessionID == "" {
		return "", errInvalidSession
	}
	return sessionID, nil
}

func (h *handler) setSessionID(ctx *gin.Context, sessionID string) {
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(sessionKey, sessionID, sessionTTL, "/", "", true, true)
}

type filterAccessParams struct {
	coordinator func(ctx *gin.Context) (bool, error)
	producer    func(ctx *gin.Context) (bool, error)
}

func filterAccess(ctx *gin.Context, params *filterAccessParams) error {
	switch getRole(ctx) {
	case service.AdminRoleAdministrator:
		return nil
	case service.AdminRoleCoordinator:
		if params == nil || params.coordinator == nil {
			return nil
		}
		if ok, err := params.coordinator(ctx); err != nil || ok {
			return err
		}
		return fmt.Errorf("handler: this coordinator is unauthenticated: %w", exception.ErrForbidden)
	case service.AdminRoleProducer:
		if params == nil || params.producer == nil {
			return nil
		}
		if ok, err := params.producer(ctx); err != nil || ok {
			return err
		}
		return fmt.Errorf("handler: this producer is unauthenticated: %w", exception.ErrForbidden)
	default:
		return fmt.Errorf("handler: unknown admin role: %w", exception.ErrForbidden)
	}
}
