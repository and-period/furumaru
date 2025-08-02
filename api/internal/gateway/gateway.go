package gateway

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Routes(rg *gin.RouterGroup)      // エンドポイント一覧の定義
	Setup(ctx context.Context) error // 初期化処理
	Sync(ctx context.Context) error  // 定期的な同期処理
}
