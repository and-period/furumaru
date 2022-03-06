package log

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

func NewGinMiddleware(params *Params) (gin.HandlerFunc, error) {
	logger, err := NewLogger(params)
	if err != nil {
		return nil, err
	}
	return ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        false,
		SkipPaths:  []string{"/health"},
	}), nil
}
