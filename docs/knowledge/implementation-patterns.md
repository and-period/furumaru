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