package handler

import (
	"context"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/sentry"
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
	Store     store.Service
}

type handler struct {
	appName   string
	env       string
	now       func() time.Time
	logger    *zap.Logger
	sentry    sentry.Client
	waitGroup *sync.WaitGroup
	store     store.Service
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
		appName: "komoju-gateway",
		env:     "",
		logger:  zap.NewNop(),
		sentry:  sentry.NewFixedMockClient(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &handler{
		appName:   dopts.appName,
		env:       dopts.env,
		now:       jst.Now,
		logger:    dopts.logger,
		sentry:    dopts.sentry,
		waitGroup: params.WaitGroup,
		store:     params.Store,
	}
}

/**
 * ###############################################
 * routes
 * ###############################################
 */
func (h *handler) Routes(rg *gin.RouterGroup) {
	komoju := rg.Group("/komoju")
	komoju.POST("/webhooks", h.Event)
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
		sentry.WithTags(map[string]string{
			"app_name":   h.appName,
			"env":        h.env,
			"method":     ctx.Request.Method,
			"path":       ctx.FullPath(),
			"user_agent": ctx.Request.UserAgent(),
		}),
	}
	h.waitGroup.Add(1)
	go func(ctx context.Context, opts []sentry.ReportOption) {
		defer h.waitGroup.Done()
		h.sentry.ReportError(ctx, err, opts...)
	}(ctx, opts)
}
