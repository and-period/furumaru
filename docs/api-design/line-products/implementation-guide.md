# LINE向け商品一覧API 実装ガイドライン

## 概要

このドキュメントは、LINE向け商品一覧APIの実装手順と考慮事項をまとめたものです。

## 実装手順

### 1. Protocol Buffers定義の追加

`api/proto/user/line.proto` を作成し、LINE向けのサービス定義を追加します。

```protobuf
syntax = "proto3";

package user.v1;

import "user/product.proto";

service LineService {
  rpc ListLineProducts(ListLineProductsRequest) returns (ListLineProductsResponse);
}

message ListLineProductsRequest {
  int64 limit = 1;
  int64 offset = 2;
  string category_id = 3;
  string producer_id = 4;
  int32 prefecture_code = 5;
  string sort = 6;
}

message ListLineProductsResponse {
  repeated LineProduct products = 1;
  Pagination pagination = 2;
  DisplaySettings display_settings = 3;
}
```

### 2. ハンドラー実装

#### ファイル構成

```
api/internal/gateway/user/v1/handler/
├── line_product.go      # LINEプロダクトハンドラー
└── line_product_test.go # テスト
```

#### 実装例（line_product.go）

```go
package handler

import (
    "net/http"
    
    "github.com/and-period/furumaru/api/internal/gateway/user/v1/response"
    "github.com/and-period/furumaru/api/internal/gateway/user/v1/service"
    "github.com/gin-gonic/gin"
)

type lineProductHandler struct {
    service service.LineProductService
}

func NewLineProductHandler(s service.LineProductService) *lineProductHandler {
    return &lineProductHandler{
        service: s,
    }
}

func (h *lineProductHandler) ListProducts(ctx *gin.Context) {
    // 1. パラメータ取得・バリデーション
    params := &listLineProductsParams{}
    if err := ctx.BindQuery(params); err != nil {
        badRequest(ctx, err)
        return
    }
    
    // 2. サービス層呼び出し
    products, pagination, err := h.service.ListLineProducts(ctx, params)
    if err != nil {
        httpError(ctx, err)
        return
    }
    
    // 3. レスポンス構築
    res := &response.LineProductsResponse{
        Products:        products,
        Pagination:      pagination,
        DisplaySettings: h.getDisplaySettings(),
    }
    
    ctx.JSON(http.StatusOK, res)
}
```

### 3. サービス層実装

#### ファイル構成

```
api/internal/gateway/user/v1/service/
├── line_product.go      # LINEプロダクトサービス
└── line_product_test.go # テスト
```

#### 実装のポイント

1. **データ取得の最適化**
   ```go
   // 並列でデータ取得
   eg, ectx := errgroup.WithContext(ctx)
   
   eg.Go(func() error {
       // 商品一覧取得
       return h.getProducts(ectx, params)
   })
   
   eg.Go(func() error {
       // 生産者情報取得
       return h.getProducers(ectx, productIDs)
   })
   ```

2. **レスポンスの軽量化**
   ```go
   func (s *lineProductService) toLineProduct(product *entity.Product) *response.LineProduct {
       return &response.LineProduct{
           ID:           product.ID,
           Name:         s.truncateName(product.Name, 40),
           Description:  s.summarizeDescription(product.Description, 100),
           ThumbnailURL: product.ThumbnailURL,
           // 必要最小限のフィールドのみ
       }
   }
   ```

3. **キャッシュの活用**
   ```go
   // Redisキャッシュの使用例
   cacheKey := fmt.Sprintf("line:products:%s", params.Hash())
   if cached := s.cache.Get(ctx, cacheKey); cached != nil {
       return cached.(*response.LineProductsResponse), nil
   }
   ```

### 4. レスポンス型定義

#### ファイル構成

```
api/internal/gateway/user/v1/response/
├── line_product.go      # LINE商品レスポンス型
└── line_product_test.go # テスト
```

#### 実装例

```go
package response

type LineProduct struct {
    ID             string       `json:"id"`
    Name           string       `json:"name"`
    Description    string       `json:"description"`
    ThumbnailURL   string       `json:"thumbnail_url"`
    Price          int64        `json:"price"`
    PriceText      string       `json:"price_text"`
    Producer       LineProducer `json:"producer"`
    Prefecture     string       `json:"prefecture"`
    Tags           []string     `json:"tags"`
    IsLimited      bool         `json:"is_limited"`
    IsOutOfStock   bool         `json:"is_out_of_stock"`
    LineURL        string       `json:"line_url"`
}
```

### 5. ルーティング設定

`api/internal/gateway/user/v1/handler/api.go` に追加：

```go
// LINE向けAPI（認証不要）
line := r.app.Group("/v1/line")
{
    line.GET("/products", h.lineProduct.ListProducts)
}
```

### 6. エンドポイントの設定

認証不要の公開APIとして実装するため、認証ミドルウェアを適用しません。

## パフォーマンス最適化

### 1. クエリ最適化

```go
// 必要なフィールドのみ取得
db.Select("id", "name", "description", "thumbnail_url", "price", "producer_id", "origin_prefecture_code").
   Where("public = ?", true).
   Where("deleted_at IS NULL").
   Limit(limit).
   Find(&products)
```

### 2. N+1問題の回避

```go
// プリロードを使用
db.Preload("Producer", "deleted_at IS NULL").
   Preload("ProductTags").
   Find(&products)
```

### 3. レスポンスキャッシュ

```yaml
# nginx設定例
location /v1/line/products {
    proxy_cache line_products_cache;
    proxy_cache_valid 200 5m;
    proxy_cache_key "$request_uri";
}
```

## テスト実装

### 1. 単体テスト

```go
func TestLineProductHandler_ListProducts(t *testing.T) {
    tests := []struct {
        name       string
        query      string
        setup      func(mocks)
        expect     func(*httptest.ResponseRecorder)
    }{
        {
            name:  "success",
            query: "?limit=10&sort=recommended",
            setup: func(m mocks) {
                m.service.EXPECT().
                    ListLineProducts(gomock.Any(), gomock.Any()).
                    Return(products, pagination, nil)
            },
            expect: func(rec *httptest.ResponseRecorder) {
                assert.Equal(t, http.StatusOK, rec.Code)
            },
        },
    }
}
```

### 2. 統合テスト

```go
func TestLineProductsIntegration(t *testing.T) {
    // Docker環境でのE2Eテスト
    ctx := context.Background()
    
    // テストデータ投入
    setupTestProducts(t, db)
    
    // API呼び出し
    res, err := client.GetLineProducts(ctx, &GetLineProductsRequest{
        Limit: 10,
    })
    
    assert.NoError(t, err)
    assert.Len(t, res.Products, 10)
}
```

## 監視とログ

### 1. メトリクス

```go
// Prometheusメトリクス
var (
    lineProductsRequests = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "line_products_requests_total",
            Help: "Total number of LINE products API requests",
        },
        []string{"status", "sort"},
    )
    
    lineProductsLatency = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "line_products_request_duration_seconds",
            Help:    "LINE products API request latency",
            Buckets: prometheus.DefBuckets,
        },
        []string{"status"},
    )
)
```

### 2. 構造化ログ

```go
logger.Info("LINE products request",
    zap.String("sort", params.Sort),
    zap.Int("limit", params.Limit),
    zap.Int("result_count", len(products)),
    zap.Duration("duration", time.Since(start)),
)
```

## セキュリティ考慮事項

1. **入力検証**
   ```go
   type listLineProductsParams struct {
       Limit  int    `form:"limit" binding:"min=1,max=20"`
       Offset int    `form:"offset" binding:"min=0"`
       Sort   string `form:"sort" binding:"omitempty,oneof=recommended new price_asc price_desc"`
   }
   ```

2. **CORS設定**
   ```go
   config := cors.DefaultConfig()
   config.AllowOrigins = []string{"https://line.me", "https://*.line-apps.com"}
   config.AllowMethods = []string{"GET", "OPTIONS"}
   ```

## デプロイメント

### 1. 環境変数

```env
# キャッシュ設定
REDIS_LINE_CACHE_URL=redis://cache:6379/2
LINE_PRODUCTS_CACHE_TTL=300
```

### 2. インフラ設定

```yaml
# docker-compose.override.yml
services:
  gateway-user:
    environment:
      - REDIS_LINE_CACHE_URL=${REDIS_LINE_CACHE_URL}
      - LINE_PRODUCTS_CACHE_TTL=${LINE_PRODUCTS_CACHE_TTL}
```

### 3. CI/CDパイプライン

```yaml
# .github/workflows/line-api.yml
- name: Run LINE API tests
  run: |
    cd api
    go test ./internal/gateway/user/v1/handler/line_*
    go test ./internal/gateway/user/v1/service/line_*
```

## トラブルシューティング

### よくある問題と対処法

1. **レスポンスが遅い**
   - キャッシュが効いているか確認
   - N+1クエリが発生していないか確認
   - インデックスが適切に設定されているか確認

2. **文字化けが発生する**
   - UTF-8エンコーディングを確認
   - 文字数カウントがバイト数でなく文字数になっているか確認

3. **CORSエラーが発生する**
   - AllowOriginsにクライアントのオリジンが含まれているか確認
   - OPTIONSメソッドが許可されているか確認

## 今後の改善案

1. **GraphQL対応**
   - より柔軟なデータ取得
   - オーバーフェッチの削減

2. **WebSocket対応**
   - リアルタイム在庫更新
   - 価格変更通知

3. **機械学習による推薦**
   - ユーザーの購買履歴分析
   - パーソナライズされた商品推薦