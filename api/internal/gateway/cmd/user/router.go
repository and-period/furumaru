package cmd

import (
	"bytes"
	"net/http"

	"github.com/and-period/furumaru/api/pkg/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type wrapResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *wrapResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func newRouter(reg *registry, opts ...gin.HandlerFunc) *gin.Engine {
	opts = append(opts, cors.NewGinMiddleware())
	opts = append(opts, gzip.Gzip(gzip.DefaultCompression))
	opts = append(opts, notifyError(reg.logger))
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

func recoveryWithWriter() gin.HandlerFunc {
	recovery := func(ctx *gin.Context, err interface{}) {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	return gin.CustomRecovery(recovery)
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
