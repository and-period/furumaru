package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/gateway"
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/service"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/backoff"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/rbac"
	"github.com/and-period/furumaru/api/pkg/sentry"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	sessionKey            = "session_id"
	sessionTTL            = 24 * 60 * 60 // 1 day
	defaultSyncInterval   = 5 * time.Minute
	defaultSyncMaxRetries = 3
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
type Params struct {
	WaitGroup *sync.WaitGroup
	User      user.Service
	Store     store.Service
	Messenger messenger.Service
	Media     media.Service
}

// @title               ふるマル API - 管理者向け
// @description         管理者向けのふるマルAPIです。
// @servers.url         https://admin-api.furumaru-stg.and-period.work
// @servers.description 検証環境
// @servers.url				  https://admin-api.furumaru.and-period.co.jp
// @servers.description 本番環境
// @securitydefinitions.cookieauth
// @securitydefinitions.bearerauth
type handler struct {
	appName        string
	env            string
	now            func() time.Time
	sentry         sentry.Client
	waitGroup      *sync.WaitGroup
	sharedGroup    *singleflight.Group
	enforcer       rbac.Enforcer
	user           user.Service
	store          store.Service
	messenger      messenger.Service
	media          media.Service
	syncMutex      *sync.Mutex
	syncInterval   time.Duration
	syncMaxRetries int64
}

type options struct {
	appName        string
	env            string
	sentry         sentry.Client
	syncInterval   time.Duration
	syncMaxRetries int64
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

func WithSyncInterval(interval time.Duration) Option {
	return func(opts *options) {
		opts.syncInterval = interval
	}
}

func WithSyncMaxRetries(retries int64) Option {
	return func(opts *options) {
		opts.syncMaxRetries = retries
	}
}

func NewHandler(params *Params, opts ...Option) gateway.Handler {
	dopts := &options{
		appName:        "admin-gateway",
		env:            "",
		sentry:         sentry.NewFixedMockClient(),
		syncInterval:   defaultSyncInterval,
		syncMaxRetries: defaultSyncMaxRetries,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &handler{
		appName:        dopts.appName,
		env:            dopts.env,
		now:            jst.Now,
		sentry:         dopts.sentry,
		waitGroup:      params.WaitGroup,
		sharedGroup:    &singleflight.Group{},
		enforcer:       nil,
		user:           params.User,
		store:          params.Store,
		messenger:      params.Messenger,
		media:          params.Media,
		syncMutex:      &sync.Mutex{},
		syncInterval:   dopts.syncInterval,
		syncMaxRetries: dopts.syncMaxRetries,
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
	h.shopRotues(v1)
	h.spotRoutes(v1)
	h.spotTypeRoutes(v1)
	h.threadRoutes(v1)
	h.topRoutes(v1)
	h.uploadRoutes(v1)
	h.userRoutes(v1)
	h.videoRoutes(v1)
	h.videoCommentRoutes(v1)

	// 認証が不要なエンドポイント
	h.guestBroadcastRoutes(v1)
}

/**
 * ###############################################
 * sync
 * ###############################################
 */
func (h *handler) Setup(ctx context.Context) error {
	if err := h.syncEnforcer(ctx); err != nil {
		return fmt.Errorf("handler: failed to sync enforcer: %w", err)
	}
	return nil
}

func (h *handler) Sync(ctx context.Context) error {
	ticker := time.NewTicker(h.syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := h.syncEnforcer(ctx); err != nil {
				slog.Error("Failed to sync enforcer", log.Error(err))
			}
		}
	}
}

func (h *handler) syncEnforcer(ctx context.Context) error {
	var rbacModel, rbacPolicy string
	retryFn := func() (err error) {
		rbacModel, rbacPolicy, err = h.user.GenerateAdminRole(ctx, &user.GenerateAdminRoleInput{})
		return err
	}
	retry := backoff.NewExponentialBackoff(h.syncMaxRetries)
	if err := backoff.Retry(ctx, retry, retryFn, backoff.WithRetryablel(exception.IsRetryable)); err != nil {
		return fmt.Errorf("handler: failed to generate admin role: %w", err)
	}
	enforcer, err := rbac.NewEnforcerFromString(rbacModel, rbacPolicy)
	if err != nil {
		return fmt.Errorf("handler: failed to new enforcer: %w", err)
	}
	h.syncMutex.Lock()
	h.enforcer = enforcer
	h.syncMutex.Unlock()
	return nil
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
			ID:        getAdminID(ctx),
			IPAddress: ctx.ClientIP(),
			Data:      map[string]string{"adminType": string(getAdminType(ctx))},
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

	auth, err := h.getAuth(ctx, token)
	if err != nil || auth.AdminID == "" {
		h.unauthorized(ctx, err)
		return
	}

	if err := h.setShop(ctx, auth); err != nil {
		h.httpError(ctx, err)
		return
	}

	setAuth(ctx, auth)

	// 認可情報の検証
	if err := h.enforce(ctx, auth); err != nil {
		h.httpError(ctx, err)
		return
	}

	ctx.Next()
}

func (h *handler) getAuth(ctx *gin.Context, token string) (*service.Auth, error) {
	in := &user.GetAdminAuthInput{
		AccessToken: token,
	}
	auth, err := h.user.GetAdminAuth(ctx, in)
	if err != nil {
		return nil, err
	}
	return service.NewAuth(auth), nil
}

func (h *handler) enforce(ctx *gin.Context, admin *service.Auth) error {
	if h.enforcer == nil {
		return nil
	}

	for _, groupID := range admin.GroupIDs {
		h.syncMutex.Lock()
		enforce, err := h.enforcer.Enforce(groupID, ctx.Request.URL.Path, ctx.Request.Method)
		h.syncMutex.Unlock()
		if err != nil {
			return fmt.Errorf("handler: failed to enforce: %w", err)
		}
		if enforce {
			return nil
		}
	}
	return fmt.Errorf("handler: you don't have the correct permissions: %w", exception.ErrForbidden)
}

func setAuth(ctx *gin.Context, auth *service.Auth) {
	if auth == nil {
		return
	}
	ctx.Request.Header.Set("Adminid", auth.AdminID)
	ctx.Request.Header.Set("Admintype", strconv.FormatInt(int64(auth.Type), 10))
}

func (h *handler) setShop(ctx *gin.Context, auth *service.Auth) error {
	if auth == nil {
		return nil
	}
	switch service.AdminType(auth.Type) {
	case service.AdminTypeAdministrator:
		return nil // 管理者は店舗を指定しない
	case service.AdminTypeCoordinator:
		shop, err := h.getShopByCoordinatorID(ctx, auth.AdminID)
		if err != nil {
			return err
		}
		ctx.Request.Header.Set("Shopid", shop.ID)
		return nil
	case service.AdminTypeProducer:
		in := &store.ListShopsInput{
			ProducerIDs: []string{auth.AdminID},
		}
		shops, _, err := h.store.ListShops(ctx, in)
		if err != nil {
			return err
		}
		ctx.Request.Header.Set("Shopids", strings.Join(shops.IDs(), ","))
		return nil
	default:
		return fmt.Errorf("handler: unknown admin role: %w", exception.ErrForbidden)
	}
}

func getAdminID(ctx *gin.Context) string {
	return ctx.GetHeader("Adminid")
}

func getAdminType(ctx *gin.Context) service.AdminType {
	role, _ := strconv.ParseInt(ctx.GetHeader("Admintype"), 10, 32)
	return service.AdminType(role)
}

func getShopID(ctx *gin.Context) string {
	return ctx.GetHeader("Shopid")
}

func getShopIDs(ctx *gin.Context) []string {
	return strings.Split(ctx.GetHeader("Shopids"), ",")
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
	switch getAdminType(ctx) {
	case service.AdminTypeAdministrator:
		return nil
	case service.AdminTypeCoordinator:
		if params == nil || params.coordinator == nil {
			return nil
		}
		if ok, err := params.coordinator(ctx); err != nil || ok {
			return err
		}
		return fmt.Errorf("handler: this coordinator is unauthenticated: %w", exception.ErrForbidden)
	case service.AdminTypeProducer:
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
