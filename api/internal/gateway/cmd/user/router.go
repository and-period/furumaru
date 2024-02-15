package user

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/pkg/cors"
	"github.com/and-period/furumaru/api/pkg/jst"
	sentrygin "github.com/getsentry/sentry-go/gin"
	ginzip "github.com/gin-contrib/gzip"
	ginpprof "github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (a *app) newRouter() *gin.Engine {
	opts := make([]gin.HandlerFunc, 0)
	opts = append(opts, nrgin.Middleware(a.newRelic))
	opts = append(opts, sentrygin.New(sentrygin.Options{}))
	opts = append(opts, a.accessLogger())
	opts = append(opts, cors.NewGinMiddleware())
	opts = append(opts, ginzip.Gzip(ginzip.DefaultCompression))
	opts = append(opts, ginzap.RecoveryWithZap(a.logger, true))

	rt := gin.New()
	rt.Use(opts...)

	a.v1.Routes(rt.Group(""))

	// other routes
	rt.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "ok")
	})
	rt.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, "not found")
	})
	ginpprof.Register(rt)

	return rt
}

type wrapResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *wrapResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *wrapResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (w *wrapResponseWriter) errorResponse() (*util.ErrorResponse, error) {
	r, err := gzip.NewReader(w.body)
	if err != nil {
		return nil, err
	}
	var res *util.ErrorResponse
	return res, json.NewDecoder(r).Decode(&res)
}

func (a *app) accessLogger() gin.HandlerFunc {
	skipPaths := map[string]bool{
		"/health": true,
	}

	return func(ctx *gin.Context) {
		var req []byte
		if a.enableDebugMode(ctx) {
			req, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(req))
		}

		w := &wrapResponseWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBufferString(""),
		}
		ctx.Writer = w

		start := jst.Now()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		ctx.Next()

		if _, ok := skipPaths[path]; ok {
			return
		}

		end := jst.Now()
		status := ctx.Writer.Status()

		fields := []zapcore.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", ctx.Request.URL.RawQuery),
			zap.String("ip", ctx.ClientIP()),
			zap.String("userAgent", ctx.Request.UserAgent()),
			zap.Int64("latency", end.Sub(start).Milliseconds()),
			zap.String("requestedAt", start.Format(time.RFC3339Nano)),
			zap.String("responsedAt", end.Format(time.RFC3339Nano)),
			zap.String("userId", ctx.GetHeader("userId")),
		}
		if a.enableDebugMode(ctx) {
			str := strings.ReplaceAll(bytes.NewBuffer(req).String(), "\n", "")
			fields = append(fields, zap.String("request", str))
		}

		// ~ 399
		if status < 400 {
			a.logger.Info(path, fields...)
			return
		}

		res, err := w.errorResponse()
		if err != nil {
			a.logger.Error("Failed to parse http response", zap.Error(err))
		}
		fields = append(fields, zap.Any("response", res))

		// 400 ~ 499
		if status < 500 {
			a.logger.Warn(path, fields...)
			return
		}

		// 500 ~
		fields = append(fields, zap.Strings("errors", ctx.Errors.Errors()))
		a.logger.Warn(path, fields...)

		if a.slack == nil {
			return
		}

		details, _ := json.Marshal(res)
		params := &alertMessageParams{
			title:   "ふるまる APIアラート",
			appName: a.AppName,
			env:     a.Environment,
			status:  int64(status),
			method:  method,
			path:    path,
			details: string(details),
		}
		msg := newAlertMessage(params)
		if err := a.slack.SendMessage(ctx, msg); err != nil {
			a.logger.Error("Failed to alert message", zap.Error(err))
		}
	}
}

func (a *app) enableDebugMode(ctx *gin.Context) bool {
	const contentType = "application/json"
	if !a.debugMode {
		return false
	}
	return strings.Contains(ctx.ContentType(), contentType)
}

type alertMessageParams struct {
	title   string
	appName string
	env     string
	status  int64
	method  string
	path    string
	details string
}

func newAlertMessage(params *alertMessageParams) slack.MsgOption {
	attachment := slack.Attachment{
		Title: params.title,
		Color: string(slack.StyleDanger),
		Fields: []slack.AttachmentField{
			{Title: "service", Value: params.appName, Short: true},
			{Title: "environment", Value: params.env, Short: true},
			{Title: "method", Value: params.method, Short: true},
			{Title: "path", Value: params.path, Short: true},
			{Title: "status", Value: strconv.FormatInt(params.status, 10), Short: false},
			{Title: "details", Value: params.details, Short: false},
		},
	}
	return slack.MsgOptionAttachments(attachment)
}
