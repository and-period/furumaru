package handler

import (
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/jst"
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
	now         func() time.Time
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	sharedGroup *singleflight.Group
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
	h.topRoutes(v1.Group("/top"))
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
func (h *handler) authentication(ctx *gin.Context) {
	token, err := util.GetAuthToken(ctx)
	if err != nil {
		unauthorized(ctx, err)
		return
	}

	in := &user.GetUserAuthInput{AccessToken: token}
	auth, err := h.user.GetUserAuth(ctx, in)
	if err != nil || auth.UserID == "" {
		unauthorized(ctx, err)
		return
	}

	setAuth(ctx, auth.UserID)

	ctx.Next()
}

func setAuth(ctx *gin.Context, userID string) {
	if userID != "" {
		ctx.Request.Header.Set("userId", userID)
	}
}

func getUserID(ctx *gin.Context) string {
	return ctx.GetHeader("userId")
}
