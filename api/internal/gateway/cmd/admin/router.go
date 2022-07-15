package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/and-period/furumaru/api/pkg/cors"
	"github.com/gin-contrib/gzip"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

func newRouter(reg *registry, logger *zap.Logger) *gin.Engine {
	opts := make([]gin.HandlerFunc, 0)
	opts = append(opts, accessLogger(logger))
	opts = append(opts, cors.NewGinMiddleware())
	opts = append(opts, gzip.Gzip(gzip.DefaultCompression))
	opts = append(opts, notifyError(logger, reg))
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

func notifyError(logger *zap.Logger, reg *registry) gin.HandlerFunc {
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

		if reg.line == nil {
			return
		}
		const service = "admin-gateway"
		altText := fmt.Sprintf("[ふるマルアラート] %s", service)
		components := []linebot.FlexComponent{
			newAlertContent("service", service),
			newAlertContent("env", reg.env),
			newAlertContent("status", strconv.FormatInt(int64(status), 10)),
			newAlertContent("method", method),
			newAlertContent("path", path),
			newAlertContent("detail", res),
		}
		_ = reg.line.PushMessage(ctx, newAlertMessage(altText, components))
	}
}

func newAlertMessage(altText string, components []linebot.FlexComponent) *linebot.FlexMessage {
	container := &linebot.BubbleContainer{
		Type:      linebot.FlexContainerTypeBubble,
		Direction: linebot.FlexBubbleDirectionTypeLTR,
		Header: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "ふるマル APIアラート",
					Size:   linebot.FlexTextSizeTypeXl,
					Weight: linebot.FlexTextWeightTypeBold,
					Color:  "#FFFFFF", // white
				},
			},
		},
		Body: &linebot.BoxComponent{
			Type:     linebot.FlexComponentTypeBox,
			Layout:   linebot.FlexBoxLayoutTypeVertical,
			Contents: components,
		},
		Size: linebot.FlexBubbleSizeTypeGiga,
		Styles: &linebot.BubbleStyle{
			Header: &linebot.BlockStyle{
				BackgroundColor: "#F44336", // red
			},
		},
	}
	return linebot.NewFlexMessage(altText, container)
}

func newAlertContent(field, value string) *linebot.BoxComponent {
	const (
		fcolor = "#aaaaaa" // gray
		vcolor = "#666666" // black
	)
	var (
		fflex = 1
		vflex = 4
	)
	return &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeBaseline,
		Contents: []linebot.FlexComponent{
			&linebot.TextComponent{
				Type:  linebot.FlexComponentTypeText,
				Text:  field,
				Size:  linebot.FlexTextSizeTypeSm,
				Color: fcolor,
				Flex:  &fflex,
			},
			&linebot.TextComponent{
				Type:  linebot.FlexComponentTypeText,
				Text:  value,
				Wrap:  true,
				Size:  linebot.FlexTextSizeTypeSm,
				Color: vcolor,
				Flex:  &vflex,
			},
		},
	}
}
