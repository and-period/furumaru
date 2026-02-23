# 実装パターン集

## Go言語 実装パターン

### エラーハンドリングパターン

#### カスタムエラータイプの定義
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}
```

#### エラーラッピング
```go
func (s *Service) ProcessOrder(ctx context.Context, orderID string) error {
    order, err := s.repository.GetOrder(ctx, orderID)
    if err != nil {
        return fmt.Errorf("failed to get order %s: %w", orderID, err)
    }
    // 処理続行...
}
```

### gRPCサービス実装パターン

#### 基本的なハンドラ構造
```go
func (h *Handler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
    // バリデーション
    if req.ProductId == "" {
        return nil, status.Error(codes.InvalidArgument, "product_id is required")
    }

    // サービス層呼び出し
    product, err := h.service.GetProduct(ctx, req.ProductId)
    if err != nil {
        return nil, h.handleError(err)
    }

    // レスポンス構築
    return &pb.GetProductResponse{
        Product: h.convertToProto(product),
    }, nil
}
```

### データベースアクセスパターン

#### リポジトリパターン
```go
type ProductRepository interface {
    Get(ctx context.Context, id string) (*entity.Product, error)
    List(ctx context.Context, params *ListParams) ([]*entity.Product, error)
    Create(ctx context.Context, product *entity.Product) error
    Update(ctx context.Context, product *entity.Product) error
    Delete(ctx context.Context, id string) error
}

type productRepository struct {
    db *gorm.DB
}

func (r *productRepository) Get(ctx context.Context, id string) (*entity.Product, error) {
    var product entity.Product
    if err := r.db.WithContext(ctx).Where("id = ?", id).First(&product).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrProductNotFound
        }
        return nil, err
    }
    return &product, nil
}
```

#### TiDB データベース実装パターン

本リポジトリのデータベース層は `database/tidb/` 配下に実装する。
各テーブルに対応する構造体とメソッドを定義し、`database.Xxx` インターフェースを実装する。

> 参照: `api/internal/store/database/tidb/category.go`

```go
package tidb

const categoryTable = "categories"

type category struct {
    db  *mysql.Client
    now func() time.Time
}

func NewCategory(db *mysql.Client) database.Category {
    return &category{
        db:  db,
        now: jst.Now,
    }
}

// List - フィルタリング・ソート・ページネーション付き一覧取得
func (c *category) List(
    ctx context.Context, params *database.ListCategoriesParams, fields ...string,
) (entity.Categories, error) {
    var categories entity.Categories
    p := listCategoriesParams(*params)
    stmt := c.db.Statement(ctx, c.db.DB, categoryTable, fields...)
    stmt = p.stmt(stmt)
    stmt = p.pagination(stmt)
    err := stmt.Find(&categories).Error
    return categories, dbError(err)
}

// Get - 単一取得
func (c *category) Get(ctx context.Context, categoryID string, fields ...string) (*entity.Category, error) {
    category, err := c.get(ctx, c.db.DB, categoryID, fields...)
    return category, dbError(err)
}

// Create - 作成（CreatedAt/UpdatedAt を now() で設定）
func (c *category) Create(ctx context.Context, category *entity.Category) error {
    now := c.now()
    category.CreatedAt, category.UpdatedAt = now, now
    err := c.db.DB.WithContext(ctx).Table(categoryTable).Create(&category).Error
    return dbError(err)
}

// Update - 更新（map[string]interface{} で部分更新）
func (c *category) Update(ctx context.Context, categoryID, name string) error {
    params := map[string]interface{}{
        "name":       name,
        "updated_at": c.now(),
    }
    stmt := c.db.DB.WithContext(ctx).
        Table(categoryTable).
        Where("id = ?", categoryID)
    err := stmt.Updates(params).Error
    return dbError(err)
}
```

**ポイント**:
- `dbError()` で GORM/MySQL エラーを `database.ErrXxx` に変換する
- `fields ...string` で取得カラムを選択可能にする
- `listXxxParams` 型で検索条件の `stmt()` / `pagination()` メソッドを提供する

### Presenter/Service ラッパーパターン

Gateway 層では、内部エンティティを API レスポンス DTO に変換する「サービスラッパー」を定義する。

> 参照: `api/internal/gateway/admin/v1/service/administrator.go`

#### 基本構造

```go
package service

import (
    "github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
    "github.com/and-period/furumaru/api/internal/user/entity"
)

// 単一エンティティのラッパー
type Administrator struct {
    types.Administrator  // Response DTO を埋め込み
}

// スライス型
type Administrators []*Administrator

// コンストラクタ: Entity → DTO 変換
func NewAdministrator(admin *entity.Administrator) *Administrator {
    return &Administrator{
        Administrator: types.Administrator{
            ID:        admin.ID,
            Status:    NewAdminStatus(admin.Status).Response(),
            Lastname:  admin.Lastname,
            Firstname: admin.Firstname,
            Email:     admin.Email,
            CreatedAt: admin.CreatedAt.Unix(),
            UpdatedAt: admin.UpdatedAt.Unix(),
        },
    }
}

// Response: DTO ポインタを返す
func (a *Administrator) Response() *types.Administrator {
    return &a.Administrator
}

// バッチ変換: Entity スライス → ラッパースライス
func NewAdministrators(admins entity.Administrators) Administrators {
    res := make(Administrators, len(admins))
    for i := range admins {
        res[i] = NewAdministrator(admins[i])
    }
    return res
}

// スライスの Response: []*types.Xxx を返す
func (as Administrators) Response() []*types.Administrator {
    res := make([]*types.Administrator, len(as))
    for i := range as {
        res[i] = as[i].Response()
    }
    return res
}

// Map: ID をキーとした map を返す（結合用）
func (as Administrators) Map() map[string]*Administrator {
    res := make(map[string]*Administrator, len(as))
    for _, a := range as {
        res[a.ID] = a
    }
    return res
}
```

#### enum 変換パターン

内部 enum と API enum の変換には専用の型と switch 文を使う。

> 参照: `api/internal/gateway/admin/v1/service/product.go`

```go
type ProductStatus types.ProductStatus

func NewProductStatus(status entity.ProductStatus) ProductStatus {
    switch status {
    case entity.ProductStatusPrivate:
        return ProductStatus(types.ProductStatusPrivate)
    case entity.ProductStatusForSale:
        return ProductStatus(types.ProductStatusForSale)
    default:
        return ProductStatus(types.ProductStatusUnknown)
    }
}

func (s ProductStatus) Response() types.ProductStatus {
    return types.ProductStatus(s)
}
```

### Gateway ハンドラーパターン

Gateway ハンドラーは Gin フレームワークを使い、以下の定型構造を持つ。

> 参照: `api/internal/gateway/admin/v1/handler/administrator.go`

#### ルーティング定義

```go
func (h *handler) administratorRoutes(rg *gin.RouterGroup) {
    r := rg.Group("/administrators", h.authentication)
    r.GET("", h.ListAdministrators)
    r.POST("", h.CreateAdministrator)
    r.GET("/:adminId", h.GetAdministrator)
    r.PATCH("/:adminId", h.UpdateAdministrator)
    r.DELETE("/:adminId", h.DeleteAdministrator)
}
```

#### 一覧取得ハンドラー

```go
// @Summary     一覧取得
// @Tags        Administrator
// @Router      /v1/administrators [get]
// @Param       limit query integer false "取得上限数" default(20)
// @Param       offset query integer false "取得開始位置" default(0)
// @Success     200 {object} types.AdministratorsResponse
func (h *handler) ListAdministrators(ctx *gin.Context) {
    const (
        defaultLimit  = 20
        defaultOffset = 0
    )

    limit, err := util.GetQueryInt64(ctx, "limit", defaultLimit)
    if err != nil {
        h.badRequest(ctx, err)
        return
    }
    offset, err := util.GetQueryInt64(ctx, "offset", defaultOffset)
    if err != nil {
        h.badRequest(ctx, err)
        return
    }

    in := &user.ListAdministratorsInput{
        Limit:  limit,
        Offset: offset,
    }
    admins, total, err := h.user.ListAdministrators(ctx, in)
    if err != nil {
        h.httpError(ctx, err)
        return
    }

    res := &types.AdministratorsResponse{
        Administrators: service.NewAdministrators(admins).Response(),
        Total:          total,
    }
    ctx.JSON(http.StatusOK, res)
}
```

#### 詳細取得ハンドラー

```go
func (h *handler) GetAdministrator(ctx *gin.Context) {
    in := &user.GetAdministratorInput{
        AdministratorID: util.GetParam(ctx, "adminId"),
    }
    admin, err := h.user.GetAdministrator(ctx, in)
    if err != nil {
        h.httpError(ctx, err)
        return
    }

    res := &types.AdministratorResponse{
        Administrator: service.NewAdministrator(admin).Response(),
    }
    ctx.JSON(http.StatusOK, res)
}
```

#### 作成ハンドラー

```go
func (h *handler) CreateAdministrator(ctx *gin.Context) {
    req := &types.CreateAdministratorRequest{}
    if err := ctx.BindJSON(req); err != nil {
        h.badRequest(ctx, err)
        return
    }

    in := &user.CreateAdministratorInput{
        Lastname:  req.Lastname,
        Firstname: req.Firstname,
        Email:     req.Email,
        // ... フィールドマッピング
    }
    admin, err := h.user.CreateAdministrator(ctx, in)
    if err != nil {
        h.httpError(ctx, err)
        return
    }

    res := &types.AdministratorResponse{
        Administrator: service.NewAdministrator(admin).Response(),
    }
    ctx.JSON(http.StatusOK, res)
}
```

**ハンドラーの共通パターン**:
1. リクエストのバインド / バリデーション（`ctx.BindJSON` または `util.GetQueryXxx`）
2. 入力構造体（`XxxInput`）の構築
3. サービス層メソッド呼び出し
4. エラー処理（`h.badRequest` / `h.httpError`）
5. サービスラッパー（`service.NewXxx`）でレスポンス変換
6. `ctx.JSON()` でレスポンス返却

### サービス層パターン

ビジネスロジック層はインターフェースを定義し、`service` パッケージで実装する。

> 参照: `api/internal/store/service/product.go`

```go
func (s *service) ListProducts(ctx context.Context, in *store.ListProductsInput) (entity.Products, int64, error) {
    // 1. バリデーション
    if err := s.validator.Struct(in); err != nil {
        return nil, 0, internalError(err)
    }

    // 2. 入力パラメータ変換
    params := &database.ListProductsParams{
        Name:   in.Name,
        Limit:  int(in.Limit),
        Offset: int(in.Offset),
    }

    // 3. 並行 DB 呼び出し（errgroup）
    var (
        products entity.Products
        total    int64
    )
    eg, ectx := errgroup.WithContext(ctx)
    eg.Go(func() (err error) {
        products, err = s.db.Product.List(ectx, params)
        return
    })
    eg.Go(func() (err error) {
        total, err = s.db.Product.Count(ectx, params)
        return
    })
    if err := eg.Wait(); err != nil {
        return nil, 0, internalError(err)
    }
    return products, total, nil
}
```

**ポイント**:
- `s.validator.Struct(in)` でバリデーション（`go-playground/validator`）
- `internalError(err)` で内部エラーを `exception.ErrXxx` に変換
- `errgroup` で独立した DB 呼び出しを並行実行

### JSONカラム処理パターン

#### mysql.JSONColumn[T] ジェネリック型
GORMのJSONカラムには `mysql.JSONColumn[T]` を使用する。各モジュールで独自のラッパー型を作成しない。

```go
import "github.com/and-period/furumaru/api/pkg/mysql"

// エンティティ定義
type Product struct {
    ID          string                         `gorm:"primaryKey"`
    Media       mysql.JSONColumn[[]Media]      `gorm:""`
    RecommendAt mysql.JSONColumn[[]time.Time]  `gorm:""`
}

// 値の作成
product := &Product{
    Media: mysql.NewJSONColumn([]Media{{URL: "https://example.com/img.jpg"}}),
}

// Update時のValue取得
val, _ := mysql.NewJSONColumn(media).Value()
updates["media"] = val
```

### testcontainers-go パターン

#### コンテナベースDBテスト
テストでは `mysql.NewContainerDB()` を使用してMySQL/TiDBコンテナを自動起動する。

```go
func TestMain(m *testing.M) {
    if !mysql.ShouldUseContainerDB() {
        // CI環境など既存DB使用時
        testDB, cleanup, _ = newTestDB(ctx)
    } else {
        // ローカル開発時: testcontainersでDB自動起動
        testDB, cleanup, _ = mysql.NewContainerDB(ctx,
            mysql.WithSchemaDir("path/to/schema"),
        )
    }
    defer cleanup()
    os.Exit(m.Run())
}
```

- `ShouldUseContainerDB()`: `DB_DRIVER` 環境変数が設定されていない場合に `true` を返す
- `WithSchemaDir()`: スキーマSQLファイルを自動実行してテーブルを初期化

### データベーステストパターン

データベース層のテストは実際の DB（testcontainers または外部 DB）に対して実行する。
テーブル駆動テスト + ヘルパー関数の組み合わせ。

> 参照: `api/internal/store/database/tidb/category_test.go`
> テンプレート詳細: `docs/knowledge/test-templates.md`

**構造**:
1. `TestMain` で DB 接続を初期化（パッケージに 1 つ）
2. `deleteAll()` でテスト前にデータクリーンアップ
3. `testXxx()` ヘルパーでテストデータ生成
4. テーブル駆動テスト（`setup` / `args` / `want` 構造）
5. 各テストケースで構造体を直接インスタンス化して実行

### サービス層テストパターン

サービス層テストでは `gomock` でデータベース等の依存をモック化し、`testService` ヘルパーで定型処理をラップする。

> 参照: `api/internal/store/service/service_test.go`, `api/internal/store/service/product_test.go`
> テンプレート詳細: `docs/knowledge/test-templates.md`

**構造**:
1. `mocks` 構造体に全依存のモックを集約
2. `testService()` ヘルパーで setup → サービス生成 → テスト実行をラップ
3. `testCaller` 型でテスト本体のシグネチャを統一
4. `withNow(now)` で時刻を固定

### モック生成パターン

モックは `go tool mockgen` で自動生成する。

```go
//go:generate go tool mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/store/$GOPACKAGE/$GOFILE
```

生成先: `api/mock/` 配下にパッケージ構造をミラーリング。

```bash
cd api && go generate ./internal/store/...
```

## Vue.js/Nuxt 実装パターン

### Composition API パターン

#### 基本的なコンポーネント構造
```vue
<template>
  <div>
    <h1>{{ title }}</h1>
    <ProductList
      :products="products"
      :loading="loading"
      @select="onProductSelect"
    />
  </div>
</template>

<script setup lang="ts">
interface Props {
  categoryId?: string
}

const props = withDefaults(defineProps<Props>(), {
  categoryId: ''
})

const emit = defineEmits<{
  productSelected: [product: Product]
}>()

// リアクティブデータ
const { data: products, pending: loading, error } = await useFetch<Product[]>(`/api/products`, {
  query: { categoryId: props.categoryId }
})

const title = computed(() => `Products (${products.value?.length || 0})`)

// メソッド
const onProductSelect = (product: Product) => {
  emit('productSelected', product)
}
</script>
```

### Pinia ストアパターン

#### 基本的なストア構造
```typescript
export const useProductStore = defineStore('product', () => {
  // State
  const products = ref<Product[]>([])
  const selectedProduct = ref<Product | null>(null)
  const loading = ref(false)

  // Getters
  const availableProducts = computed(() => 
    products.value.filter(p => p.status === 'available')
  )

  // Actions
  const fetchProducts = async (categoryId?: string) => {
    loading.value = true
    try {
      const { data } = await $fetch<{ products: Product[] }>('/api/products', {
        query: { categoryId }
      })
      products.value = data.products
    } catch (error) {
      console.error('Failed to fetch products:', error)
    } finally {
      loading.value = false
    }
  }

  const selectProduct = (product: Product) => {
    selectedProduct.value = product
  }

  return {
    // State
    products: readonly(products),
    selectedProduct: readonly(selectedProduct),
    loading: readonly(loading),
    // Getters
    availableProducts,
    // Actions
    fetchProducts,
    selectProduct
  }
})
```

## API設計パターン

### RESTful API設計
```
GET    /api/v1/products          # 商品一覧取得
GET    /api/v1/products/:id      # 商品詳細取得
POST   /api/v1/products          # 商品作成
PUT    /api/v1/products/:id      # 商品更新
DELETE /api/v1/products/:id      # 商品削除

# ネストしたリソース
GET    /api/v1/products/:id/reviews     # 商品レビュー一覧
POST   /api/v1/products/:id/reviews     # レビュー投稿
```

### レスポンス形式の統一
```json
{
  "data": {
    "products": [
      {
        "id": "prod_123",
        "name": "商品名",
        "price": 1000,
        "created_at": "2023-01-01T00:00:00Z"
      }
    ]
  },
  "meta": {
    "total": 100,
    "page": 1,
    "per_page": 20
  }
}
```

## パフォーマンス最適化

### データベース最適化
- 適切なインデックスの設定
- N+1クエリの回避（EagerLoadingの活用）
- ページネーションの実装
- 不要なデータの取得回避

### フロントエンド最適化
- 画像の最適化・遅延読み込み
- コンポーネントの遅延読み込み
- バンドルサイズの最適化
- キャッシュ戦略の実装

## セキュリティ実装

### 認証・認可パターン
```go
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractToken(r)
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        claims, err := h.verifyToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 入力値検証
```go
func ValidateCreateProductRequest(req *CreateProductRequest) error {
    if req.Name == "" {
        return errors.New("name is required")
    }
    if req.Price <= 0 {
        return errors.New("price must be positive")
    }
    if len(req.Description) > 1000 {
        return errors.New("description too long")
    }
    return nil
}
```
