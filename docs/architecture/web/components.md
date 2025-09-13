# Webコンポーネント設計

## コンポーネント分類とルール

### 基本コンポーネント分類
```
components/
├── atoms/           # 単体要素（Button, Input等）
├── molecules/       # 複合要素（SearchBox, Card等）  
├── organisms/       # 機能単位（Header, ProductList等）
└── templates/       # ページ構造（Layout, Wrapper等）
```

### 命名規則
- **ファイル名**: `PascalCase.vue`
- **コンポーネント名**: `PascalCase`
- **Props**: `camelCase`
- **Events**: `kebab-case`

## 共通コンポーネント (`/web/shared`)

### 基本UI要素
```typescript
// 例: Button Component
<template>
  <button
    :class="buttonClasses"
    :disabled="disabled"
    @click="$emit('click', $event)"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'secondary' | 'danger'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
}

defineEmits<{
  click: [event: MouseEvent]
}>()
</script>
```

### データ表示コンポーネント
- **ProductCard**: 商品情報表示
- **UserAvatar**: ユーザーアバター
- **PriceDisplay**: 価格表示（通貨フォーマット）
- **StatusBadge**: ステータス表示

### 入力コンポーネント
- **FormInput**: バリデーション対応入力
- **FormSelect**: セレクトボックス
- **FormTextarea**: テキストエリア
- **ImageUploader**: 画像アップロード

## Admin専用コンポーネント

### Vuetify拡張
```vue
<!-- 管理画面特化のDataTable -->
<VDataTable
  :headers="headers"
  :items="items"
  :loading="loading"
  class="elevation-1"
/>
```

### 業務コンポーネント
- **AdminLayout**: 管理画面レイアウト
- **DataTable**: 一覧表示テーブル
- **FormDialog**: モーダルフォーム
- **ChartDisplay**: グラフ表示
- **RichEditor**: リッチテキストエディタ

## User専用コンポーネント

### Tailwind CSS活用
```vue
<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <!-- コンテンツ -->
  </div>
</template>
```

### EC特化コンポーネント
- **ProductGrid**: 商品グリッド表示
- **CartItem**: カートアイテム
- **OrderSummary**: 注文サマリー
- **ReviewForm**: レビュー投稿フォーム
- **LivePlayer**: ライブ配信プレイヤー

## コンポーネント設計原則

### 単一責任の原則
```vue
<!-- Good: 単一責任 -->
<ProductCard :product="product" />

<!-- Bad: 複数責任 -->
<ProductCardWithCartAndReview :product="product" />
```

### Props設計
```typescript
// Good: 明示的な型定義
interface Props {
  product: Product
  showPrice?: boolean
  size?: 'sm' | 'md' | 'lg'
}

// Bad: any型や曖昧な定義
interface Props {
  data: any
  options?: object
}
```

### Events設計
```vue
<script setup lang="ts">
// Good: 明示的なイベント定義
defineEmits<{
  'add-to-cart': [product: Product]
  'view-detail': [productId: string]
}>()
</script>
```

### スロット活用
```vue
<template>
  <BaseModal>
    <template #header>
      <slot name="header" />
    </template>
    <template #content>
      <slot />
    </template>
    <template #footer>
      <slot name="footer" />
    </template>
  </BaseModal>
</template>
```

## スタイリング戦略

### Admin (Vuetify)
```scss
// カスタマイズ用SCSS変数
$primary-color: #1976d2;
$secondary-color: #424242;

.v-btn--custom {
  @apply rounded-lg shadow-md;
}
```

### User (Tailwind CSS)
```vue
<template>
  <div class="container mx-auto px-4">
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- グリッドアイテム -->
    </div>
  </div>
</template>
```

### 共通スタイル
```css
/* 共通CSSカスタムプロパティ */
:root {
  --color-primary: #1976d2;
  --color-secondary: #424242;
  --spacing-unit: 0.25rem;
}
```

## テスト戦略

### コンポーネントテスト
```typescript
import { mount } from '@vue/test-utils'
import ProductCard from '@/components/ProductCard.vue'

describe('ProductCard', () => {
  test('renders product information correctly', () => {
    const product = {
      id: '1',
      name: 'Test Product',
      price: 1000
    }
    
    const wrapper = mount(ProductCard, {
      props: { product }
    })
    
    expect(wrapper.text()).toContain('Test Product')
    expect(wrapper.text()).toContain('¥1,000')
  })
})
```

### Storybook活用
```typescript
// ProductCard.stories.ts
export default {
  title: 'Components/ProductCard',
  component: ProductCard
}

export const Default = {
  args: {
    product: {
      id: '1',
      name: 'Sample Product',
      price: 1000,
      image: '/sample.jpg'
    }
  }
}
```

## パフォーマンス考慮

### 遅延ローディング
```vue
<script setup lang="ts">
// 重いコンポーネントの遅延読み込み
const RichEditor = defineAsyncComponent(
  () => import('@/components/RichEditor.vue')
)
</script>
```

### メモ化
```vue
<script setup lang="ts">
const expensiveCalculation = computed(() => {
  // 重い計算処理
  return heavyProcessing(props.data)
})
</script>
```

### 仮想スクロール
```vue
<!-- 大量データの表示 -->
<VirtualList
  :items="products"
  :item-height="200"
  #default="{ item }"
>
  <ProductCard :product="item" />
</VirtualList>
```

## アクセシビリティ

### キーボード対応
```vue
<template>
  <button
    @click="handleClick"
    @keydown.enter="handleClick"
    @keydown.space="handleClick"
  >
    Action Button
  </button>
</template>
```

### ARIA対応
```vue
<template>
  <div
    role="alert"
    :aria-live="type === 'error' ? 'assertive' : 'polite'"
  >
    {{ message }}
  </div>
</template>
```

### セマンティクス
```vue
<template>
  <article>
    <header>
      <h2>{{ product.name }}</h2>
    </header>
    <section>
      <p>{{ product.description }}</p>
    </section>
    <footer>
      <button type="button">Add to Cart</button>
    </footer>
  </article>
</template>
```