package handler

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/and-period/marche/api/internal/gateway/admin/v1/service"
	"github.com/and-period/marche/api/internal/gateway/util"
	store "github.com/and-period/marche/api/internal/store/service"
	uentity "github.com/and-period/marche/api/internal/user/entity"
	user "github.com/and-period/marche/api/internal/user/service"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/and-period/marche/api/pkg/rbac"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
 * ###############################################
 * handler
 * ###############################################
 */
type APIV1Handler interface {
	Routes(rg *gin.RouterGroup) // エンドポイント一覧の定義
}

type Params struct {
	WaitGroup    *sync.WaitGroup
	Enforcer     rbac.Enforcer
	UserService  user.UserService
	StoreService store.StoreService
}

type apiV1Handler struct {
	now         func() time.Time
	logger      *zap.Logger
	sharedGroup *singleflight.Group
	waitGroup   *sync.WaitGroup
	enforcer    rbac.Enforcer
	user        user.UserService
	store       store.StoreService
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

func NewAPIV1Handler(params *Params, opts ...Option) APIV1Handler {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &apiV1Handler{
		now:       jst.Now,
		logger:    dopts.logger,
		waitGroup: params.WaitGroup,
		enforcer:  params.Enforcer,
		user:      params.UserService,
		store:     params.StoreService,
	}
}

/**
 * ###############################################
 * routes
 * ###############################################
 */
func (h *apiV1Handler) Routes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	h.storeRoutes(v1.Group("/stores"))
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
func (h *apiV1Handler) authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 認証情報の検証
		_, err := util.GetAuthToken(ctx)
		if err != nil {
			unauthorized(ctx, err)
			return
		}

		// TODO: 管理者情報取得処理の追加
		auth := &uentity.AdminAuth{Role: uentity.AdminRoleAdministrator}
		role := service.NewAdminRole(auth.Role)

		setAuth(ctx, auth.AdminID, role)

		// 認可情報の検証
		if h.enforcer == nil {
			ctx.Next()
			return
		}

		enforce, err := h.enforcer.Enforce(role.String(), ctx.Request.URL.Path, ctx.Request.Method)
		if err != nil {
			httpError(ctx, status.Error(codes.Internal, err.Error()))
			return
		}
		if !enforce {
			forbidden(ctx, errors.New("handler: you don't have the correct permissions"))
			return
		}

		ctx.Next()
	}
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

// func getRole(ctx *gin.Context) service.AdminRole {
// 	role, _ := strconv.ParseInt(ctx.GetHeader("role"), 10, 64)
// 	return service.AdminRole(role)
// }
