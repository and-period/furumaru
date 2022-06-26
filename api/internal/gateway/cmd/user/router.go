package cmd

import (
	"bytes"
	"net/http"
	"time"

	"github.com/and-period/furumaru/api/pkg/cors"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func newRouter(reg *registry, logger *zap.Logger) *gin.Engine {
	opts := make([]gin.HandlerFunc, 0)
	opts = append(opts, accessLogger(logger))
	opts = append(opts, cors.NewGinMiddleware())
	opts = append(opts, gzip.Gzip(gzip.DefaultCompression))
	opts = append(opts, notifyError(logger))
	opts = append(opts, recoveryWithWriter())

	rt := gin.New()
	rt.Use(opts...)

	reg.v1.Routes(rt.Group(""))

	// other routes
	rt.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "ok")
	})
	rt.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, "not found")
	})

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

func recoveryWithWriter() gin.HandlerFunc {
	recovery := func(ctx *gin.Context, err interface{}) {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	return gin.CustomRecovery(recovery)
}

func accessLogger(logger *zap.Logger) gin.HandlerFunc {
	return ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        false,
		SkipPaths:  []string{"/health"},
	})
}

func notifyError(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		w := &wrapResponseWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBufferString(""),
		}
		ctx.Writer = w
		ctx.Next()

		status := ctx.Writer.Status()
		if status < 500 {
			return
		}

		path := ctx.Request.URL.String()
		method := ctx.Request.Method
		res := w.body.String()

		logger.Error("Internal Server Error",
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("response", res),
		)
	}
}
