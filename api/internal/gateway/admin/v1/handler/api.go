package handler

import (
	"errors"
	"fmt"
	"strconv"
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
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errInvalidOrderkey   = errors.New("handler: invalid order key")
	errInvalidFileFormat = errors.New("handler: invalid file format")
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
	now         func() time.Time
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	sharedGroup *singleflight.Group
	enforcer    rbac.Enforcer
	user        user.Service
	store       store.Service
	messenger   messenger.Service
	media       media.Service
}

type options struct {
	logger *zap.Logger
}

type Option func(opts *options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func NewHandler(params *Params, opts ...Option) Handler {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &handler{
		now:         jst.Now,
		logger:      dopts.logger,
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
	h.liveRoutes(v1)
	h.messageRoutes(v1)
	h.notificationRoutes(v1)
	h.orderRoutes(v1)
	h.postalCodeRoutes(v1)
	h.producerRoutes(v1)
	h.productRoutes(v1)
	h.productTagRoutes(v1)
	h.productTypeRoutes(v1)
	h.promotionRoutes(v1)
	h.relatedProducerRoutes(v1)
	h.scheduleRoutes(v1)
	h.shippingRoutes(v1)
	h.threadRoutes(v1)
	h.uploadRoutes(v1)
	h.userRoutes(v1)
}

/**
 * ###############################################
 * error handling
 * ###############################################
 */
func (h *handler) httpError(ctx *gin.Context, err error) {
	res, code := util.NewErrorResponse(err)
	if code >= 500 {
		h.logger.Error("Internal server error", zap.Error(err), zap.Any("request", ctx.Request))
	}
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
		fmt.Println("debug", err)
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
