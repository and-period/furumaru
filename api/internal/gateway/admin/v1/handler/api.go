package handler

import (
	"sync"
	"time"

	"github.com/and-period/marche/api/internal/gateway/util"
	user "github.com/and-period/marche/api/internal/user/service"
	"github.com/and-period/marche/api/pkg/jst"
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
	WaitGroup   *sync.WaitGroup
	UserService user.UserService
}

type apiV1Handler struct {
	now         func() time.Time
	logger      *zap.Logger
	sharedGroup *singleflight.Group
	waitGroup   *sync.WaitGroup
	user        user.UserService
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
		user:      params.UserService,
	}
}

/**
 * ###############################################
 * routes
 * ###############################################
 */
func (h *apiV1Handler) Routes(rg *gin.RouterGroup) {}

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

// func badRequest(ctx *gin.Context, err error) {
// 	httpError(ctx, status.Error(codes.InvalidArgument, err.Error()))
// }

func unauthorized(ctx *gin.Context, err error) {
	httpError(ctx, status.Error(codes.Unauthenticated, err.Error()))
}

/**
 * ###############################################
 * other
 * ###############################################
 */
func (h *apiV1Handler) authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := util.GetAuthToken(ctx)
		if err != nil {
			unauthorized(ctx, err)
			return
		}

		// TODO: 管理者情報取得処理の追加

		ctx.Next()
	}
}

func setAuth(ctx *gin.Context, userID string) {
	if userID != "" {
		ctx.Request.Header.Set("userId", userID)
	}
}

func getUserID(ctx *gin.Context) string {
	return ctx.GetHeader("userId")
}
