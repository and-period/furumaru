package handler

import (
	"context"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/sentry"
	pkgstripe "github.com/and-period/furumaru/api/pkg/stripe"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Params struct {
	WaitGroup *sync.WaitGroup
	Store     store.Service
	Receiver  pkgstripe.Receiver
}

type handler struct {
	appName   string
	env       string
	now       func() time.Time
	sentry    sentry.Client
	waitGroup *sync.WaitGroup
	store     store.Service
	receiver  pkgstripe.Receiver
}

type options struct {
	appName string
	env     string
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

func WithSentry(sentry sentry.Client) Option {
	return func(opts *options) {
		opts.sentry = sentry
	}
}

func NewHandler(params *Params, opts ...Option) gateway.Handler {
	dopts := &options{
		appName: "stripe-gateway",
		env:     "",
		sentry:  sentry.NewFixedMockClient(),
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &handler{
		appName:   dopts.appName,
		env:       dopts.env,
		now:       jst.Now,
		sentry:    dopts.sentry,
		waitGroup: params.WaitGroup,
		store:     params.Store,
		receiver:  params.Receiver,
	}
}

func (h *handler) Setup(ctx context.Context) error {
	return nil
}

func (h *handler) Sync(ctx context.Context) error {
	return nil
}

/**
 * ###############################################
 * routes
 * ###############################################
 */
func (h *handler) Routes(rg *gin.RouterGroup) {
	stripe := rg.Group("/stripe")
	stripe.POST("/webhooks", h.Event)
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
