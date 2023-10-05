package handler

import (
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler interface {
	Routes(rg *gin.RouterGroup) // エンドポイント一覧の定義
}

type Params struct {
	WaitGroup *sync.WaitGroup
}

type handler struct {
	now       func() time.Time
	logger    *zap.Logger
	waitGroup *sync.WaitGroup
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
		now:       jst.Now,
		logger:    dopts.logger,
		waitGroup: params.WaitGroup,
	}
}

func (h *handler) Routes(rg *gin.RouterGroup) {
	komoju := rg.Group("/komoju")
	komoju.POST("/webhooks", h.Event)
}

func httpError(ctx *gin.Context, err error) {
	res, code := util.NewErrorResponse(err)
	ctx.JSON(code, res)
	ctx.Abort()
}

func badRequest(ctx *gin.Context, err error) {
	httpError(ctx, status.Error(codes.InvalidArgument, err.Error()))
}
