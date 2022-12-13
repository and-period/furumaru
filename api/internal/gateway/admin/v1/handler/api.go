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
	h.authRoutes(v1.Group("/auth"))
	h.administratorRoutes(v1.Group("/administrators"))
	h.coordinatorRoutes(v1.Group("/coordinators"))
	h.producerRoutes(v1.Group("/producers"))
	h.relatedProducerRoutes(v1.Group("/coordinators/:coordinatorId/producers"))
	h.categoryRoutes(v1.Group("/categories"))
	h.productTypeRoutes(v1.Group("/categories/:categoryId/product-types"))
	h.shippingRoutes(v1.Group("/shippings"))
	h.productRoutes(v1.Group("/products"))
	h.promotionRoutes(v1.Group("/promotions"))
	h.orderRoutes(v1.Group("/orders"))
	h.notificationRoutes(v1.Group("/notifications"))
	h.contactRoutes(v1.Group("/contacts"))
	h.messageRoutes(v1.Group("/messages"))
	h.scheduleRoutes(v1.Group("/schedules"))
	h.liveRoutes(v1.Group("/schedules/:scheduleId/lives"))
	h.userRoutes(v1.Group("/users"))
	h.postalCodeRoutes(v1.Group("/postal-codes"))
	v1.GET("/categories/-/product-types", h.authentication, h.ListProductTypes)
	h.uploadRoutes(v1.Group("/upload"))
}

/**
 * ###############################################
 * error handling
 * ###############################################
 */
func httpError(ctx *gin.Context, err error) {
	res, code := util.NewErrorResponse(err)
	ctx.JSON(code, res)
	ctx.Abort()
}

func badRequest(ctx *gin.Context, err error) {
	httpError(ctx, status.Error(codes.InvalidArgument, err.Error()))
}

func unauthorized(ctx *gin.Context, err error) {
	httpError(ctx, status.Error(codes.Unauthenticated, err.Error()))
}

func forbidden(ctx *gin.Context, err error) {
	httpError(ctx, status.Error(codes.PermissionDenied, err.Error()))
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
		unauthorized(ctx, err)
		return
	}

	in := &user.GetAdminAuthInput{AccessToken: token}
	auth, err := h.user.GetAdminAuth(ctx, in)
	if err != nil || auth.AdminID == "" {
		unauthorized(ctx, err)
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
		httpError(ctx, status.Error(codes.Internal, err.Error()))
		return
	}
	if !enforce {
		forbidden(ctx, errors.New("handler: you don't have the correct permissions"))
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
