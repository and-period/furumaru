package log

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func NewGinMiddleware(opts ...Option) (gin.HandlerFunc, error) {
	logger, err := NewLogger(opts...)
	if err != nil {
		return nil, err
	}
	return ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        false,
		SkipPaths:  []string{}, // ヘルスチェックのデバッグのため、一時的にすべて出力
		// SkipPaths:  []string{"/health"},
	}), nil
}
