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
	AuthRoutes(rg *gin.RouterGroup)   // 認証済みでアクセス可能なエンドポイント一覧
	NoAuthRoutes(rg *gin.RouterGroup) // 未認証でもアクセス可能なエンドポイント一覧
	Authentication() gin.HandlerFunc  // 認証情報の検証
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
func (h *apiV1Handler) AuthRoutes(rg *gin.RouterGroup) {}

func (h *apiV1Handler) NoAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/v1/auth", h.SignIn)
	rg.DELETE("/v1/auth", h.SignOut)
	rg.POST("/v1/auth/refresh-token", h.RefreshAuthToken)
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
func (h *apiV1Handler) Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: User Serviceの実装完了後に修正
		// token, err := util.GetAuthToken(ctx)
		// if err != nil {
		// 	unauthorized(ctx, err)
		// 	return
		// }

		// in := &user.GetUserAuthRequest{AccessToken: token}
		// out, err := h.user.GetUserAuth(ctx, in)
		// if err != nil || out.UserID == "" {
		// 	unauthorized(ctx, err)
		// 	return
		// }

		setAuth(ctx, "")

		ctx.Next()
	}
}

func setAuth(ctx *gin.Context, userID string) {
	if userID != "" {
		ctx.Request.Header.Set("userId", userID)
	}
}

// func getUserID(ctx *gin.Context) string {
// 	return ctx.GetHeader("userId")
// }
