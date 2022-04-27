package handler

import (
	"sync"
	"time"

	"github.com/and-period/marche/api/internal/gateway/util"
	"github.com/and-period/marche/api/pkg/jst"
	"github.com/and-period/marche/api/proto/user"
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
	Logger      *zap.Logger
	WaitGroup   *sync.WaitGroup
	UserService user.UserServiceClient
}

type apiV1Handler struct {
	now         func() time.Time
	logger      *zap.Logger
	sharedGroup *singleflight.Group
	waitGroup   *sync.WaitGroup
	user        user.UserServiceClient
}

func NewAPIV1Handler(params *Params) APIV1Handler {
	return &apiV1Handler{
		now:       jst.Now,
		logger:    params.Logger,
		waitGroup: params.WaitGroup,
		user:      params.UserService,
	}
}

/**
 * ###############################################
 * routes
 * ###############################################
 */
func (h *apiV1Handler) Routes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	h.authRoutes(v1.Group("/auth"))
	h.userRoutes(v1.Group("/users"))
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

/**
 * ###############################################
 * other
 * ###############################################
 */
func (h *apiV1Handler) authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := util.GetAuthToken(ctx)
		if err != nil {
			unauthorized(ctx, err)
			return
		}

		in := &user.GetUserAuthRequest{AccessToken: token}
		out, err := h.user.GetUserAuth(ctx, in)
		if err != nil || out.Auth.UserId == "" {
			unauthorized(ctx, err)
			return
		}

		setAuth(ctx, out.Auth.UserId)

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
