# Web状態管理アーキテクチャ

## Pinia状態管理戦略

### ストア分割方針
```
stores/
├── auth.ts         # 認証状態・権限管理
├── cart.ts         # カート状態・セッション管理  
├── user.ts         # ユーザー情報・プロフィール
├── product.ts      # 商品データ・検索状態
├── order.ts        # 注文履歴・進行状況
├── media.ts        # 動画・ライブ配信状態
├── ui.ts           # UI状態・モーダル管理
└── notification.ts # 通知・アラート管理
```

## 認証ストア (`auth.ts`)

### 基本構造
```typescript
interface AuthState {
  user: User | null
  token: string | null
  refreshToken: string | null
  isAuthenticated: boolean
  permissions: Permission[]
  loginProvider: 'cognito' | 'google' | 'line' | null
}

export const useAuthStore = defineStore('auth', () => {
  const state = reactive<AuthState>({
    user: null,
    token: null,
    refreshToken: null,
    isAuthenticated: false,
    permissions: [],
    loginProvider: null
  })

  // Getters
  const hasPermission = computed(() => 
    (permission: string) => state.permissions.some(p => p.name === permission)
  )

  // Actions
  const login = async (credentials: LoginCredentials) => {
    try {
      const { user, token, refreshToken } = await authApi.login(credentials)
      state.user = user
      state.token = token
      state.refreshToken = refreshToken
      state.isAuthenticated = true
    } catch (error) {
      throw new AuthError('Login failed', error)
    }
  }

  const logout = async () => {
    await authApi.logout()
    $reset()
  }

  const refreshAuthToken = async () => {
    if (!state.refreshToken) throw new Error('No refresh token')
    const { token } = await authApi.refresh(state.refreshToken)
    state.token = token
  }

  return {
    ...toRefs(state),
    hasPermission,
    login,
    logout,
    refreshAuthToken
  }
})
```

### 認証フロー管理
```typescript
// middleware/auth.ts
export default defineNuxtRouteMiddleware((to, from) => {
  const { isAuthenticated } = useAuthStore()
  
  if (!isAuthenticated) {
    return navigateTo('/login', { 
      query: { redirect: to.fullPath } 
    })
  }
})
```

## カートストア (`cart.ts`)

### セッション対応
```typescript
interface CartItem {
  id: string
  productId: string
  quantity: number
  price: number
  options?: Record<string, any>
}

interface CartState {
  items: CartItem[]
  sessionId: string | null
  guestId: string | null
  lastUpdated: Date
}

export const useCartStore = defineStore('cart', () => {
  const state = reactive<CartState>({
    items: [],
    sessionId: null,
    guestId: null,
    lastUpdated: new Date()
  })

  // Getters
  const totalAmount = computed(() => 
    state.items.reduce((sum, item) => sum + (item.price * item.quantity), 0)
  )

  const itemCount = computed(() =>
    state.items.reduce((sum, item) => sum + item.quantity, 0)
  )

  // Actions
  const addItem = async (product: Product, quantity = 1) => {
    const existingItem = state.items.find(item => item.productId === product.id)
    
    if (existingItem) {
      existingItem.quantity += quantity
    } else {
      state.items.push({
        id: generateId(),
        productId: product.id,
        quantity,
        price: product.price
      })
    }
    
    await syncWithServer()
  }

  const syncWithServer = async () => {
    const { isAuthenticated } = useAuthStore()
    
    if (isAuthenticated) {
      await cartApi.syncUserCart(state.items)
    } else {
      await cartApi.syncGuestCart(state.guestId, state.items)
    }
    
    state.lastUpdated = new Date()
  }

  return {
    ...toRefs(state),
    totalAmount,
    itemCount,
    addItem,
    syncWithServer
  }
}, {
  persist: {
    key: 'furumaru-cart',
    storage: persistedState.localStorage,
    pick: ['items', 'guestId'] // sessionIdは永続化しない
  }
})
```

## 商品ストア (`product.ts`)

### 検索・フィルタ機能
```typescript
interface ProductFilters {
  category?: string[]
  priceRange?: [number, number]
  producer?: string[]
  prefecture?: string[]
  inStock?: boolean
}

interface ProductState {
  products: Product[]
  searchQuery: string
  filters: ProductFilters
  sortBy: 'price' | 'name' | 'created_at'
  sortOrder: 'asc' | 'desc'
  loading: boolean
  pagination: {
    page: number
    limit: number
    total: number
  }
}

export const useProductStore = defineStore('product', () => {
  const state = reactive<ProductState>({
    products: [],
    searchQuery: '',
    filters: {},
    sortBy: 'created_at',
    sortOrder: 'desc',
    loading: false,
    pagination: {
      page: 1,
      limit: 20,
      total: 0
    }
  })

  // Getters
  const filteredProducts = computed(() => {
    let result = state.products

    // 検索クエリフィルタ
    if (state.searchQuery) {
      result = result.filter(product => 
        product.name.includes(state.searchQuery) ||
        product.description.includes(state.searchQuery)
      )
    }

    // カテゴリフィルタ
    if (state.filters.category?.length) {
      result = result.filter(product =>
        state.filters.category!.includes(product.categoryId)
      )
    }

    return result
  })

  // Actions
  const searchProducts = async (query: string) => {
    state.loading = true
    state.searchQuery = query
    
    try {
      const response = await productApi.search({
        query,
        filters: state.filters,
        sort: `${state.sortBy}:${state.sortOrder}`,
        page: state.pagination.page,
        limit: state.pagination.limit
      })
      
      state.products = response.products
      state.pagination.total = response.total
    } finally {
      state.loading = false
    }
  }

  return {
    ...toRefs(state),
    filteredProducts,
    searchProducts
  }
})
```

## UI状態ストア (`ui.ts`)

### モーダル・通知管理
```typescript
interface UIState {
  modals: Record<string, boolean>
  loading: Record<string, boolean>
  notifications: Notification[]
  theme: 'light' | 'dark'
  sidebar: {
    open: boolean
    collapsed: boolean
  }
}

export const useUIStore = defineStore('ui', () => {
  const state = reactive<UIState>({
    modals: {},
    loading: {},
    notifications: [],
    theme: 'light',
    sidebar: {
      open: false,
      collapsed: false
    }
  })

  // Actions
  const showModal = (modalId: string) => {
    state.modals[modalId] = true
  }

  const hideModal = (modalId: string) => {
    state.modals[modalId] = false
  }

  const setLoading = (key: string, loading: boolean) => {
    state.loading[key] = loading
  }

  const addNotification = (notification: Omit<Notification, 'id'>) => {
    const id = generateId()
    state.notifications.push({ ...notification, id })
    
    if (notification.autoClose !== false) {
      setTimeout(() => {
        removeNotification(id)
      }, notification.duration || 5000)
    }
  }

  const removeNotification = (id: string) => {
    const index = state.notifications.findIndex(n => n.id === id)
    if (index > -1) {
      state.notifications.splice(index, 1)
    }
  }

  return {
    ...toRefs(state),
    showModal,
    hideModal,
    setLoading,
    addNotification,
    removeNotification
  }
}, {
  persist: {
    pick: ['theme', 'sidebar']
  }
})
```

## ストア間連携

### Composition API活用
```typescript
// composables/useCheckout.ts
export const useCheckout = () => {
  const cart = useCartStore()
  const auth = useAuthStore()
  const ui = useUIStore()

  const proceedToCheckout = async () => {
    if (!auth.isAuthenticated) {
      ui.showModal('login')
      return
    }

    if (cart.itemCount === 0) {
      ui.addNotification({
        type: 'error',
        message: 'カートが空です'
      })
      return
    }

    ui.setLoading('checkout', true)
    
    try {
      const order = await checkoutApi.create({
        items: cart.items,
        userId: auth.user?.id
      })
      
      await navigateTo(`/checkout/${order.id}`)
    } catch (error) {
      ui.addNotification({
        type: 'error',
        message: 'チェックアウトに失敗しました'
      })
    } finally {
      ui.setLoading('checkout', false)
    }
  }

  return {
    proceedToCheckout
  }
}
```

## 永続化戦略

### アプリ別永続化設定
```typescript
// Admin (非永続化）
export const useAdminAuthStore = defineStore('admin-auth', () => {
  // セッションのみ、永続化しない
})

// User (永続化）
export const useUserStore = defineStore('user', () => {
  // ...
}, {
  persist: {
    storage: persistedState.localStorage,
    pick: ['preferences', 'recentlyViewed']
  }
})

// LIFF (制限的永続化）  
export const useLiffStore = defineStore('liff', () => {
  // ...
}, {
  persist: {
    storage: persistedState.sessionStorage, // sessionStorageを使用
    pick: ['tempData'] // 最小限のデータのみ
  }
})
```

## エラーハンドリング

### 統一エラー管理
```typescript
// stores/error.ts
export const useErrorStore = defineStore('error', () => {
  const errors = ref<AppError[]>([])

  const handleError = (error: unknown) => {
    const appError = normalizeError(error)
    errors.value.push(appError)
    
    // Sentry等への送信
    if (appError.level === 'error') {
      reportError(appError)
    }
  }

  const clearErrors = () => {
    errors.value = []
  }

  return {
    errors: readonly(errors),
    handleError,
    clearErrors
  }
})
```

## テスト戦略

### ストアテスト
```typescript
import { createPinia, setActivePinia } from 'pinia'
import { useCartStore } from '@/stores/cart'

describe('CartStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  test('adds item to cart', async () => {
    const cart = useCartStore()
    const product = { id: '1', name: 'Test', price: 1000 }
    
    await cart.addItem(product, 2)
    
    expect(cart.itemCount).toBe(2)
    expect(cart.totalAmount).toBe(2000)
  })
})
```

## パフォーマンス最適化

### 状態正規化
```typescript
// 正規化されたストア構造
interface NormalizedState {
  products: {
    byId: Record<string, Product>
    allIds: string[]
  }
  categories: {
    byId: Record<string, Category>
    allIds: string[]
  }
}

// セレクター関数
const getProductById = (id: string) => state.products.byId[id]
const getAllProducts = () => state.products.allIds.map(getProductById)
```

### 遅延更新
```typescript
import { debounce } from 'lodash-es'

const debouncedSync = debounce(async () => {
  await syncWithServer()
}, 1000)

const updateQuantity = (itemId: string, quantity: number) => {
  const item = state.items.find(i => i.id === itemId)
  if (item) {
    item.quantity = quantity
    debouncedSync()
  }
}
```