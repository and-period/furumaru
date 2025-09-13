# Web API連携アーキテクチャ

## API連携の基本構造

### レイヤー構成
```
Vue Component -> Composable -> API Client -> HTTP Gateway -> gRPC Service
     |              |            |              |              |
   UI Logic     Business      HTTP通信         認証・変換      ビジネス
               Logic                                           ロジック
```

## API Client層

### 基本クライアント構造
```typescript
// lib/api/client.ts
interface ApiClientConfig {
  baseURL: string
  timeout: number
  headers?: Record<string, string>
}

class ApiClient {
  private axios: AxiosInstance
  
  constructor(config: ApiClientConfig) {
    this.axios = axios.create({
      baseURL: config.baseURL,
      timeout: config.timeout,
      headers: {
        'Content-Type': 'application/json',
        ...config.headers
      }
    })
    
    this.setupInterceptors()
  }

  private setupInterceptors() {
    // リクエストインターセプター
    this.axios.interceptors.request.use(
      (config) => {
        const { token } = useAuthStore()
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }
        return config
      },
      (error) => Promise.reject(error)
    )

    // レスポンスインターセプター
    this.axios.interceptors.response.use(
      (response) => response,
      async (error) => {
        if (error.response?.status === 401) {
          await this.handleUnauthorized()
        }
        return Promise.reject(error)
      }
    )
  }

  private async handleUnauthorized() {
    const { refreshAuthToken, logout } = useAuthStore()
    
    try {
      await refreshAuthToken()
    } catch {
      await logout()
      await navigateTo('/login')
    }
  }
}
```

### サービス別APIクライアント
```typescript
// lib/api/product.ts
interface ProductSearchParams {
  query?: string
  category?: string[]
  priceRange?: [number, number]
  page?: number
  limit?: number
}

interface ProductSearchResponse {
  products: Product[]
  total: number
  page: number
  limit: number
}

export class ProductApiClient {
  constructor(private client: ApiClient) {}

  async search(params: ProductSearchParams): Promise<ProductSearchResponse> {
    const response = await this.client.get('/products/search', { params })
    return response.data
  }

  async getById(id: string): Promise<Product> {
    const response = await this.client.get(`/products/${id}`)
    return response.data
  }

  async getReviews(productId: string): Promise<ProductReview[]> {
    const response = await this.client.get(`/products/${productId}/reviews`)
    return response.data
  }
}
```

## Composable層

### 基本Composable構造
```typescript
// composables/useProducts.ts
export const useProducts = () => {
  const productStore = useProductStore()
  const { addNotification } = useUIStore()
  
  const loading = ref(false)
  const error = ref<Error | null>(null)

  const searchProducts = async (params: ProductSearchParams) => {
    loading.value = true
    error.value = null
    
    try {
      await productStore.searchProducts(params)
    } catch (err) {
      error.value = err as Error
      addNotification({
        type: 'error',
        message: '商品の検索に失敗しました'
      })
    } finally {
      loading.value = false
    }
  }

  const getProduct = async (id: string) => {
    try {
      const product = await productApi.getById(id)
      return product
    } catch (err) {
      addNotification({
        type: 'error',
        message: '商品の取得に失敗しました'
      })
      throw err
    }
  }

  return {
    // State
    products: computed(() => productStore.products),
    loading: readonly(loading),
    error: readonly(error),
    
    // Actions
    searchProducts,
    getProduct
  }
}
```

### リアルタイム対応Composable
```typescript
// composables/useLiveStream.ts
export const useLiveStream = (streamId: string) => {
  const stream = ref<LiveStream | null>(null)
  const comments = ref<Comment[]>([])
  const connected = ref(false)
  
  let ws: WebSocket | null = null

  const connect = () => {
    const wsUrl = `${config.WS_BASE_URL}/live/${streamId}`
    ws = new WebSocket(wsUrl)
    
    ws.onopen = () => {
      connected.value = true
    }
    
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      
      switch (data.type) {
        case 'comment':
          comments.value.push(data.comment)
          break
        case 'stream_update':
          stream.value = { ...stream.value, ...data.stream }
          break
      }
    }
    
    ws.onclose = () => {
      connected.value = false
    }
  }

  const sendComment = (message: string) => {
    if (ws && connected.value) {
      ws.send(JSON.stringify({
        type: 'comment',
        message,
        timestamp: new Date().toISOString()
      }))
    }
  }

  const disconnect = () => {
    if (ws) {
      ws.close()
      ws = null
    }
  }

  onMounted(connect)
  onUnmounted(disconnect)

  return {
    stream: readonly(stream),
    comments: readonly(comments),
    connected: readonly(connected),
    sendComment
  }
}
```

## 認証連携

### 認証方式別実装
```typescript
// composables/useAuth.ts
export const useAuth = () => {
  const authStore = useAuthStore()
  const route = useRoute()

  // Cognito認証
  const loginWithCognito = async (email: string, password: string) => {
    try {
      const response = await authApi.loginCognito({ email, password })
      authStore.setAuth(response)
      
      const redirect = route.query.redirect as string
      await navigateTo(redirect || '/')
    } catch (error) {
      throw new AuthError('Cognito認証に失敗しました')
    }
  }

  // Google OAuth
  const loginWithGoogle = async () => {
    const googleAuthUrl = `${config.API_BASE_URL}/auth/google/redirect`
    window.location.href = googleAuthUrl
  }

  // LINE OAuth  
  const loginWithLine = async () => {
    const lineAuthUrl = `${config.API_BASE_URL}/auth/line/redirect`
    window.location.href = lineAuthUrl
  }

  // LIFF認証
  const initializeLiff = async () => {
    if (!process.client) return
    
    const liff = await import('@line/liff')
    await liff.default.init({ liffId: config.LIFF_ID })
    
    if (liff.default.isLoggedIn()) {
      const profile = await liff.default.getProfile()
      authStore.setLiffAuth(profile)
    }
  }

  return {
    isAuthenticated: computed(() => authStore.isAuthenticated),
    user: computed(() => authStore.user),
    loginWithCognito,
    loginWithGoogle,
    loginWithLine,
    initializeLiff
  }
}
```

## エラーハンドリング

### 統一エラー処理
```typescript
// lib/api/errors.ts
export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public code?: string,
    public details?: any
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export const handleApiError = (error: any): ApiError => {
  if (error.response) {
    const { status, data } = error.response
    return new ApiError(
      data.message || 'APIエラーが発生しました',
      status,
      data.code,
      data.details
    )
  }
  
  if (error.request) {
    return new ApiError('ネットワークエラーが発生しました', 0)
  }
  
  return new ApiError(error.message || '不明なエラーが発生しました', 0)
}
```

### リトライ機能
```typescript
// lib/api/retry.ts
interface RetryConfig {
  maxRetries: number
  delay: number
  backoff: 'linear' | 'exponential'
}

export const withRetry = async <T>(
  fn: () => Promise<T>,
  config: RetryConfig
): Promise<T> => {
  let lastError: Error
  
  for (let attempt = 0; attempt <= config.maxRetries; attempt++) {
    try {
      return await fn()
    } catch (error) {
      lastError = error as Error
      
      if (attempt === config.maxRetries) {
        break
      }
      
      const delay = config.backoff === 'exponential' 
        ? config.delay * Math.pow(2, attempt)
        : config.delay
        
      await new Promise(resolve => setTimeout(resolve, delay))
    }
  }
  
  throw lastError!
}
```

## 型安全性

### API型定義の自動生成
```typescript
// types/api/generated.ts (OpenAPIから自動生成)
export interface Product {
  id: string
  name: string
  description: string
  price: number
  categoryId: string
  producerId: string
  images: ProductImage[]
  status: 'active' | 'inactive' | 'sold_out'
  createdAt: string
  updatedAt: string
}

export interface ProductSearchRequest {
  query?: string
  categoryIds?: string[]
  producerIds?: string[]
  priceMin?: number
  priceMax?: number
  page?: number
  limit?: number
}
```

### Runtime型検証
```typescript
// lib/api/validation.ts
import { z } from 'zod'

const ProductSchema = z.object({
  id: z.string(),
  name: z.string(),
  price: z.number().positive(),
  status: z.enum(['active', 'inactive', 'sold_out'])
})

export const validateProduct = (data: unknown): Product => {
  return ProductSchema.parse(data)
}
```

## キャッシュ戦略

### SWR実装
```typescript
// composables/useSWR.ts
interface SWROptions<T> {
  revalidateOnFocus?: boolean
  revalidateOnReconnect?: boolean
  refreshInterval?: number
  initialData?: T
}

export const useSWR = <T>(
  key: string,
  fetcher: () => Promise<T>,
  options: SWROptions<T> = {}
) => {
  const data = ref<T | undefined>(options.initialData)
  const error = ref<Error | null>(null)
  const loading = ref(false)

  const mutate = async (newData?: T) => {
    if (newData) {
      data.value = newData
    } else {
      await fetchData()
    }
  }

  const fetchData = async () => {
    loading.value = true
    error.value = null
    
    try {
      const result = await fetcher()
      data.value = result
    } catch (err) {
      error.value = err as Error
    } finally {
      loading.value = false
    }
  }

  // 初回読み込み
  onMounted(fetchData)
  
  // フォーカス時の再検証
  if (options.revalidateOnFocus) {
    useEventListener('focus', fetchData)
  }

  return {
    data: readonly(data),
    error: readonly(error),
    loading: readonly(loading),
    mutate
  }
}
```

## パフォーマンス最適化

### Request Deduplication
```typescript
// lib/api/deduplication.ts
const pendingRequests = new Map<string, Promise<any>>()

export const dedupRequest = <T>(
  key: string,
  fetcher: () => Promise<T>
): Promise<T> => {
  if (pendingRequests.has(key)) {
    return pendingRequests.get(key)!
  }
  
  const promise = fetcher()
    .finally(() => {
      pendingRequests.delete(key)
    })
  
  pendingRequests.set(key, promise)
  return promise
}
```

### バッチリクエスト
```typescript
// composables/useBatchLoader.ts
export const useBatchLoader = <K, V>(
  batchFn: (keys: K[]) => Promise<V[]>,
  options: { batchSize?: number; delay?: number } = {}
) => {
  const { batchSize = 10, delay = 10 } = options
  const queue: Array<{ key: K; resolve: (value: V) => void; reject: (error: any) => void }> = []

  const processBatch = async () => {
    if (queue.length === 0) return
    
    const batch = queue.splice(0, batchSize)
    const keys = batch.map(item => item.key)
    
    try {
      const results = await batchFn(keys)
      batch.forEach((item, index) => {
        item.resolve(results[index])
      })
    } catch (error) {
      batch.forEach(item => item.reject(error))
    }
  }

  const load = (key: K): Promise<V> => {
    return new Promise((resolve, reject) => {
      queue.push({ key, resolve, reject })
      
      // 遅延実行でバッチ処理
      setTimeout(processBatch, delay)
    })
  }

  return { load }
}
```

## リアルタイム通信

### WebSocket管理
```typescript
// composables/useWebSocket.ts
export const useWebSocket = (url: string) => {
  const socket = ref<WebSocket | null>(null)
  const connected = ref(false)
  const error = ref<Event | null>(null)

  const connect = () => {
    socket.value = new WebSocket(url)
    
    socket.value.onopen = () => {
      connected.value = true
      error.value = null
    }
    
    socket.value.onerror = (event) => {
      error.value = event
    }
    
    socket.value.onclose = () => {
      connected.value = false
    }
  }

  const send = (data: any) => {
    if (socket.value?.readyState === WebSocket.OPEN) {
      socket.value.send(JSON.stringify(data))
    }
  }

  const disconnect = () => {
    socket.value?.close()
  }

  onMounted(connect)
  onUnmounted(disconnect)

  return {
    connected: readonly(connected),
    error: readonly(error),
    send,
    disconnect
  }
}
```

## テスト戦略

### API Client テスト
```typescript
// tests/api/product.test.ts
import { describe, test, expect, vi } from 'vitest'
import { ProductApiClient } from '@/lib/api/product'

describe('ProductApiClient', () => {
  test('searches products successfully', async () => {
    const mockClient = {
      get: vi.fn().mockResolvedValue({
        data: {
          products: [{ id: '1', name: 'Test Product' }],
          total: 1
        }
      })
    }
    
    const client = new ProductApiClient(mockClient as any)
    const result = await client.search({ query: 'test' })
    
    expect(result.products).toHaveLength(1)
    expect(mockClient.get).toHaveBeenCalledWith(
      '/products/search',
      { params: { query: 'test' } }
    )
  })
})
```

### Composable テスト
```typescript
// tests/composables/useProducts.test.ts
import { describe, test, expect } from 'vitest'
import { useProducts } from '@/composables/useProducts'

describe('useProducts', () => {
  test('handles search correctly', async () => {
    const { searchProducts, loading, products } = useProducts()
    
    expect(loading.value).toBe(false)
    
    await searchProducts({ query: 'test' })
    
    expect(products.value).toBeDefined()
  })
})
```