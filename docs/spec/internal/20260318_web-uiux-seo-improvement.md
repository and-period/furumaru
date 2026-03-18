# Web UI/UX・SEO 総合改善

| 項目 |  |
|----|--|
| 機能 | web/user・web/liff・web/shared の UI/UX 改善および SEO 最適化 |

## 仕様

ユーザー体験の向上と検索エンジン最適化を目的に、2026 年現在のベストプラクティスに沿った改善を行う。
対象は `web/user`（ユーザー向け EC サイト）、`web/liff`（LINE ミニアプリ）、`web/shared`（共通コンポーネントライブラリ）の 3 プロジェクト。

## 設計概要

改善領域を以下の 5 カテゴリに分類し、優先度順に対応する。

1. **SEO 最適化**（web/user）— 検索エンジンからの流入増加
2. **アクセシビリティ改善**（全体）— WCAG 2.2 AA 準拠
3. **パフォーマンス最適化**（全体）— Core Web Vitals 改善
4. **UI/UX デザイン改善**（全体）— ユーザー体験の質的向上
5. **共通基盤の強化**（web/shared）— 一貫性と保守性の向上

## 設計詳細

---

### 1. SEO 最適化（web/user）

#### 1-1. 構造化データ（JSON-LD）の導入

現状、構造化データが一切存在しない。以下の schema.org マークアップを各ページに追加する。

| ページ | schema.org Type | 対象データ |
|--------|----------------|-----------|
| トップページ | `Organization`, `WebSite`, `SearchAction` | サイト情報、サイト内検索 |
| 商品一覧 | `ItemList`, `CollectionPage` | 商品リスト |
| 商品詳細 | `Product`, `Offer`, `AggregateRating`, `Review` | 価格・在庫・レビュー |
| 体験詳細 | `Event`, `Offer` | 体験情報・価格・日程 |
| マルシェ詳細 | `Event`, `Place` | イベント情報・開催場所 |
| ブログ記事 | `Article`, `BreadcrumbList` | 記事内容・パンくず |
| コーディネーター | `Person`, `Organization` | 生産者・団体情報 |

```vue
<!-- 商品詳細ページでの実装例 -->
<script setup>
useHead({
  script: [
    {
      type: 'application/ld+json',
      innerHTML: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'Product',
        name: product.name,
        image: product.media.map(m => m.url),
        description: product.description,
        offers: {
          '@type': 'Offer',
          price: product.price,
          priceCurrency: 'JPY',
          availability: product.inventory > 0
            ? 'https://schema.org/InStock'
            : 'https://schema.org/OutOfStock',
        },
        aggregateRating: {
          '@type': 'AggregateRating',
          ratingValue: product.averageRating,
          reviewCount: product.reviewCount,
        },
      }),
    },
  ],
})
</script>
```

#### 1-2. メタタグの動的生成

現状、多くのページで固定タイトルのみ設定されている。各ページに適切な動的メタタグを設定する。

```typescript
// 商品詳細ページの例
useSeoMeta({
  title: `${product.name} | ふるマル`,
  description: product.description?.slice(0, 120),
  ogTitle: product.name,
  ogDescription: product.description?.slice(0, 120),
  ogImage: product.thumbnailUrl,
  ogType: 'product',
  ogUrl: `https://www.furumaru.and-period.co.jp/items/${product.id}`,
  twitterCard: 'summary_large_image',
  twitterTitle: product.name,
  twitterDescription: product.description?.slice(0, 120),
  twitterImage: product.thumbnailUrl,
})
```

**対象ページと動的メタタグ一覧:**

| ページ | title | description | og:image |
|--------|-------|-------------|----------|
| 商品詳細 | `{商品名} \| ふるマル` | 商品説明文（120文字） | 商品サムネイル |
| 体験詳細 | `{体験名} \| ふるマル` | 体験説明文（120文字） | 体験サムネイル |
| マルシェ詳細 | `{マルシェ名} \| ふるマル` | マルシェ説明文 | マルシェ画像 |
| コーディネーター | `{名前} \| ふるマル` | コーディネーター紹介文 | プロフィール画像 |
| カテゴリ一覧 | `{カテゴリ名}の商品一覧 \| ふるマル` | カテゴリ説明文 | カテゴリ画像 |
| ブログ記事 | `{記事タイトル} \| ふるマル` | 記事冒頭120文字 | 記事アイキャッチ |

#### 1-3. canonical URL の設定

ページごとに正規 URL を明示し、重複コンテンツを防止する。

```typescript
// nuxt.config.ts に追加
app: {
  head: {
    link: [
      { rel: 'canonical', href: 'https://www.furumaru.and-period.co.jp' },
    ],
  },
},
```

各ページで動的に `useHead` で canonical を上書きする。

#### 1-4. hreflang タグの追加

i18n で日本語・英語をサポートしているため、hreflang を設定する。

```typescript
// nuxt.config.ts の i18n 設定に追加
i18n: {
  locales: [
    { code: 'ja', language: 'ja-JP' },
    { code: 'en', language: 'en-US' },
  ],
  // Nuxt i18n モジュールが自動で hreflang を生成
}
```

#### 1-5. サイトマップの最適化

`@nuxtjs/sitemap` モジュールは導入済みだが、動的ルートの設定が不足している。

```typescript
// nuxt.config.ts
sitemap: {
  sources: ['/api/__sitemap__/urls'],
  defaults: {
    changefreq: 'daily',
    priority: 0.8,
    lastmod: new Date().toISOString(),
  },
  exclude: [
    '/account/**',
    '/signin',
    '/signup',
    '/verify',
    '/auth/**',
    '/checkout/**',
    '/cart',
  ],
},
```

サーバーサイドで動的 URL（商品・体験・マルシェ・ブログ）を生成するエンドポイントを追加する。

#### 1-6. robots.txt の最適化

```typescript
// nuxt.config.ts
robots: {
  groups: [
    {
      userAgent: '*',
      allow: '/',
      disallow: ['/account/', '/checkout/', '/cart', '/auth/', '/verify', '/signin', '/signup'],
    },
  ],
  sitemap: 'https://www.furumaru.and-period.co.jp/sitemap.xml',
},
```

#### 1-7. パンくずリストの実装

全ページにパンくずリストを表示し、`BreadcrumbList` 構造化データも同時に出力する。

```vue
<!-- components/organisms/TheBreadcrumb.vue -->
<template>
  <nav aria-label="パンくずリスト">
    <ol class="flex items-center gap-1 text-sm text-gray-500">
      <li v-for="(crumb, i) in crumbs" :key="crumb.path">
        <NuxtLink v-if="i < crumbs.length - 1" :to="crumb.path" class="hover:underline">
          {{ crumb.label }}
        </NuxtLink>
        <span v-else aria-current="page">{{ crumb.label }}</span>
        <span v-if="i < crumbs.length - 1" aria-hidden="true">/</span>
      </li>
    </ol>
  </nav>
</template>
```

#### 1-8. セマンティック HTML の強化

| 現状 | 改善後 | 対象 |
|------|--------|------|
| `<div>` でセクション区切り | `<main>`, `<section>`, `<article>`, `<aside>` | 全ページ |
| 見出しレベルの飛び | `<h1>` → `<h2>` → `<h3>` の正しい階層 | 全ページ |
| リンクリスト | `<nav>` で囲む | ヘッダー・フッター |
| 商品一覧 | `<ul>` + `<li>` でリスト化 | 一覧ページ |

---

### 2. アクセシビリティ改善（WCAG 2.2 AA 準拠）

#### 2-1. ARIA 属性の追加

**web/user — 対応必要箇所:**

| コンポーネント | 問題 | 対応 |
|--------------|------|------|
| アイコンボタン（30箇所以上） | `aria-label` 未設定 | 操作内容を表す `aria-label` を追加 |
| ドロップダウンメニュー | `aria-expanded` 未設定 | 開閉状態を `aria-expanded` で通知 |
| カートメニュー | `aria-expanded`, `aria-controls` 未設定 | 展開状態とターゲットを関連付け |
| スナックバー通知 | `aria-live` 未設定 | `aria-live="polite"` を追加 |
| モーダルダイアログ | フォーカストラップ未実装 | `focus-trap` を実装 |
| 画像（商品サムネイル等） | `alt` 未設定 | 商品名等を `alt` に設定 |

**web/liff — 対応必要箇所:**

| コンポーネント | 問題 | 対応 |
|--------------|------|------|
| 注文ステータスバッジ | 色のみで区別 | テキストラベルを併記 |
| ハンバーガーメニュー | キーボードナビゲーション不可 | フォーカス管理を追加 |
| 商品画像 | 一部 `alt` 未設定 | 商品名を `alt` に設定 |
| カートパネル（固定下部） | フォーカストラップなし | 展開時にフォーカスを管理 |

**web/shared — 対応必要箇所:**

| コンポーネント | 問題 | 対応 |
|--------------|------|------|
| FmTextInput | `aria-invalid`, `aria-describedby` 未設定 | エラー状態とメッセージを関連付け |
| FmProductItem | 動画要素に `aria-label` なし | 動画の説明を追加 |
| FmProductDetail | 動画に字幕なし | 代替テキストまたはキャプションを追加 |
| FmCreditCardForm | `autocomplete` 属性未設定 | `cc-number`, `cc-name` 等を設定 |
| 全エラーメッセージ | `role="alert"` 未設定 | 動的エラーに `role="alert"` を追加 |

#### 2-2. キーボードナビゲーション

```typescript
// スキップリンクの追加（default.vue レイアウト）
<a href="#main-content" class="sr-only focus:not-sr-only focus:absolute focus:z-50 focus:p-4 focus:bg-white">
  メインコンテンツへスキップ
</a>

// メインコンテンツに id を付与
<main id="main-content" tabindex="-1">
  <slot />
</main>
```

- すべてのインタラクティブ要素に `focus-visible` スタイルを設定
- モーダル・ドロップダウンにフォーカストラップを実装
- Escape キーでモーダル・メニューを閉じる

#### 2-3. カラーコントラスト

| 要素 | 現状 | 問題 | 対応 |
|------|------|------|------|
| `typography: #707070` on `base: #F9F6EA` | コントラスト比 4.57:1 | WCAG AA を満たす（4.5:1 以上）がギリギリ | `#5a5a5a`（6.37:1）に変更し余裕を持たせる |
| `orange: #F48D26` on white | コントラスト比 2.42:1 | WCAG AA 不適合 | ボタンテキストは白以外に変更、または背景色を調整 |
| プレースホルダーテキスト | 薄いグレー | 読みにくい | コントラスト比 3:1 以上に調整 |

---

### 3. パフォーマンス最適化

#### 3-1. 画像最適化（web/user）

`@nuxt/image` は導入済みだが、最適化設定が不足している。

```vue
<!-- 現状: 最適化なし -->
<img :src="product.thumbnail" />

<!-- 改善後: レスポンシブ画像 + 遅延読み込み + WebP -->
<NuxtImg
  :src="product.thumbnail"
  :alt="product.name"
  width="400"
  height="300"
  sizes="(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 33vw"
  format="webp"
  loading="lazy"
  placeholder
/>
```

**対応項目:**
- すべての `<img>` を `<NuxtImg>` / `<NuxtPicture>` に置き換え
- `sizes` 属性でレスポンシブ画像を提供
- `format="webp"` で次世代フォーマットを使用
- Above the fold 以外は `loading="lazy"` を設定
- LCP 要素（ヒーロー画像等）には `loading="eager"` + `fetchpriority="high"` を設定
- プレースホルダー（blur-up）を有効化

#### 3-2. 動画の遅延読み込み

```vue
<!-- ヒーロー動画: IntersectionObserver で遅延読み込み -->
<script setup>
const videoRef = ref<HTMLElement | null>(null)
const isVideoVisible = ref(false)

onMounted(() => {
  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          isVideoVisible.value = true
          observer.disconnect()
        }
      })
    },
    { threshold: 0.25 },
  )

  if (videoRef.value) {
    observer.observe(videoRef.value)
  }
})
</script>

<template>
  <div ref="videoRef">
    <video
      v-if="isVideoVisible"
      :src="videoUrl"
      autoplay
      muted
      loop
      playsinline
      preload="none"
    />
  </div>
</template>
```

#### 3-3. スケルトン UI の導入

現状、データ取得中に空白表示またはスピナーのみ。スケルトン UI でレイアウトシフトを防ぐ。

**対象ページ:**

| ページ | スケルトン内容 |
|--------|-------------|
| 商品一覧 | カードグリッド（画像 + テキスト行 × 2） |
| 商品詳細 | 画像ギャラリー + テキストブロック |
| 注文履歴 | リストアイテム × 5 |
| トップページ | ヒーローエリア + カルーセル |

```vue
<!-- components/atoms/TheSkeleton.vue -->
<template>
  <div
    :class="[
      'animate-pulse rounded bg-gray-200',
      props.class,
    ]"
    :style="{ width, height }"
  />
</template>
```

#### 3-4. コード分割とプリフェッチ

```typescript
// nuxt.config.ts
experimental: {
  componentIslands: true, // アイランドアーキテクチャ
},
routeRules: {
  '/': { prerender: true },                     // トップページをプリレンダリング
  '/items/**': { swr: 3600 },                   // 商品ページを 1 時間キャッシュ
  '/experiences/**': { swr: 3600 },             // 体験ページを 1 時間キャッシュ
  '/about': { prerender: true },                // 静的ページ
  '/account/**': { ssr: false },                // アカウントページはCSR
},
```

#### 3-5. サードパーティスクリプトの最適化

現状、GTM・Meta Pixel・New Relic・Clarity が同時に読み込まれている。

```typescript
// nuxt.config.ts — パフォーマンス影響の大きいスクリプトを遅延読み込み
app: {
  head: {
    script: [
      // New Relic, Clarity: requestIdleCallback または afterInteractive で読み込み
      { src: 'clarity-script', tagPosition: 'bodyClose', defer: true },
    ],
  },
},
```

---

### 4. UI/UX デザイン改善

#### 4-1. web/user 改善項目

| # | カテゴリ | 現状の問題 | 改善案 | 優先度 |
|---|---------|-----------|--------|-------|
| 1 | ローディング | データ取得中に空白表示 | スケルトン UI + プログレスインジケーター | 高 |
| 2 | エラーハンドリング | 汎用的なエラーメッセージのみ | コンテキストに応じたエラー表示 + リトライボタン | 高 |
| 3 | 検索機能 | ナビゲーションで検索がコメントアウト | サイト内検索を実装（商品・体験・ブログ横断） | 高 |
| 4 | 商品一覧 | 基本的なグリッド表示のみ | フィルタリング・ソート機能の追加 | 中 |
| 5 | 画像ギャラリー | 基本的な表示 | ピンチズーム・スワイプ対応のギャラリー | 中 |
| 6 | フォーム UX | 送信時のみバリデーション | リアルタイムバリデーション + インラインエラー | 中 |
| 7 | ページ遷移 | 即時切り替え | `<ViewTransition>` API によるスムーズな遷移 | 低 |
| 8 | ダークモード | 未対応 | Tailwind dark モード + システム設定連動 | 低 |
| 9 | オフライン対応 | 未対応 | Service Worker でオフラインフォールバック | 低 |
| 10 | LINE バナー | 常に固定表示（閉じても再表示） | 表示状態を永続化、スクロール位置に応じた表示制御 | 中 |

#### 4-2. web/liff 改善項目

| # | カテゴリ | 現状の問題 | 改善案 | 優先度 |
|---|---------|-----------|--------|-------|
| 1 | 数量選択 | 1 個ずつしか追加できない | 商品カードに数量ピッカーを追加 | 高 |
| 2 | 在庫表示 | カートで在庫変動の警告なし | カートアイテムに在庫状態を反映 | 高 |
| 3 | カラー一貫性 | `#F48D26` がハードコード（テーマ定義と不一致） | Tailwind テーマカラーに統一 | 中 |
| 4 | 注文ステータス | ハードコードされたステータスコード（1-5） | 定数化 + ステータスラベルの一元管理 | 中 |
| 5 | 注文一覧 | 50 件固定でページネーションなし | 無限スクロールまたはページネーション | 中 |
| 6 | チェックイン日時 | `datetime-local` のネイティブ UI | カスタム日時ピッカー（UX 向上） | 低 |
| 7 | クーポン UX | 適用後のコード表示が分かりにくい | 適用状態の明確なフィードバック | 中 |
| 8 | プルトゥリフレッシュ | 未対応 | リスト画面でプルトゥリフレッシュを追加 | 低 |

#### 4-3. トップページ改善（web/user）

```text
現状のレイアウト:
┌─────────────────────────┐
│ ヒーロー動画（全幅）       │ ← アニメーションタイトル（3秒交互）
├─────────────────────────┤
│ CTA ボタン × 2           │ ← 商品一覧・体験一覧
├─────────────────────────┤
│ ライブ配信カルーセル       │ ← 横スクロール
├─────────────────────────┤
│ アーカイブカルーセル       │ ← 横スクロール
├─────────────────────────┤
│ 動画カルーセル            │ ← 横スクロール
├─────────────────────────┤
│ About セクション          │
└─────────────────────────┘

改善案:
┌─────────────────────────┐
│ ヒーロー動画（全幅）       │ ← 動画遅延読み込み + poster 画像
│ + 検索バー               │ ← サイト内検索の入口
├─────────────────────────┤
│ CTA ボタン × 3           │ ← 商品・体験・マルシェ
├─────────────────────────┤
│ おすすめ商品              │ ← パーソナライズ表示
├─────────────────────────┤
│ 旬の特集                 │ ← 季節に応じたコンテンツ
├─────────────────────────┤
│ ライブ配信 / アーカイブ    │ ← タブ切り替え（カルーセル統合）
├─────────────────────────┤
│ 生産者紹介               │ ← コーディネーターカルーセル
├─────────────────────────┤
│ レビュー・お客様の声       │ ← 社会的証明
├─────────────────────────┤
│ About + FAQ              │
└─────────────────────────┘
```

#### 4-4. 共通 UI 改善

**トースト通知の統一:**
```text
現状: スナックバー（DOM 直接操作）+ アラートコンポーネント（混在）
改善: 統一されたトースト通知システム（Headless UI パターン）
  - 成功: 緑背景 + チェックアイコン
  - エラー: 赤背景 + エラーアイコン
  - 情報: 青背景 + 情報アイコン
  - 自動消去（5秒）+ 手動閉じ
  - aria-live="polite" で読み上げ対応
```

**空状態の改善:**
```text
現状: 空リストでテキストのみ表示
改善: イラスト + メッセージ + CTA ボタン
  例: 「注文履歴がありません」→ イラスト + 「商品を探す」ボタン
```

---

### 5. 共通基盤の強化（web/shared）

#### 5-1. コンポーネントの拡充

現在 5 コンポーネントのみ。以下を追加してデザインの一貫性を向上させる。

| コンポーネント | 用途 | 利用箇所 |
|--------------|------|---------|
| `FmSkeleton` | スケルトンローディング | 全ページ |
| `FmToast` | トースト通知 | 全ページ |
| `FmBadge` | ステータスバッジ | 注文一覧、商品一覧 |
| `FmModal` | モーダルダイアログ | カート、確認ダイアログ |
| `FmPagination` | ページネーション | 一覧ページ |
| `FmEmptyState` | 空状態表示 | 一覧ページ |
| `FmBreadcrumb` | パンくずリスト | 全ページ（user） |
| `FmImageGallery` | 画像ギャラリー | 商品詳細 |
| `FmQuantityPicker` | 数量選択 | カート、商品詳細 |
| `FmRating` | 星評価表示 | 商品詳細、レビュー |

#### 5-2. デザイントークンの統一

web/user と web/liff でカラー定義が異なっている（OKLch vs HEX、値の差異）。

```css
/* shared/src/assets/design-tokens.css — 統一されたデザイントークン */
@theme {
  /* カラーパレット */
  --color-main: #604c3f;
  --color-orange: #d97a38;
  --color-green: #7cb342;
  --color-base: #f9f6ea;
  --color-typography: #5a5a5a;  /* コントラスト改善: #707070 → #5a5a5a */
  --color-error: #f44336;
  --color-success: #66bb6a;
  --color-placeholder: #a8a19a;

  /* スペーシング */
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-md: 16px;
  --spacing-lg: 24px;
  --spacing-xl: 32px;

  /* ボーダー */
  --radius-sm: 4px;
  --radius-md: 8px;
  --radius-lg: 16px;
  --radius-full: 9999px;

  /* シャドウ */
  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.1);
  --shadow-lg: 0 10px 15px rgba(0, 0, 0, 0.1);
}
```

#### 5-3. FmTextInput のアクセシビリティ強化

```vue
<!-- 改善後 -->
<template>
  <div>
    <label :for="id">{{ label }}</label>
    <input
      :id="id"
      v-model="modelValue"
      :type="inputType"
      :aria-invalid="!!errorMessage"
      :aria-describedby="errorMessage ? `${id}-error` : undefined"
      :autocomplete="autocompleteAttr"
    />
    <p v-if="errorMessage" :id="`${id}-error`" role="alert" class="text-error text-sm">
      {{ errorMessage }}
    </p>
  </div>
</template>
```

#### 5-4. FmCreditCardForm の autocomplete 対応

```html
<input name="cc-number" autocomplete="cc-number" inputmode="numeric" />
<input name="cc-name" autocomplete="cc-name" />
<input name="cc-exp-month" autocomplete="cc-exp-month" />
<input name="cc-exp-year" autocomplete="cc-exp-year" />
<input name="cc-csc" autocomplete="cc-csc" inputmode="numeric" />
```

#### 5-5. テストカバレッジの向上

```text
現状: Storybook ストーリーのみ（ユニットテストなし）
改善:
  - 各コンポーネントに Vitest ユニットテストを追加
  - アクセシビリティテスト（@storybook/addon-a11y を "todo" → "error" に変更）
  - ビジュアルリグレッションテスト（Chromatic 統合済み → CI パイプラインに組み込み）
```

---

## チェックリスト

### 実装開始前
* [ ] Core Web Vitals のベースライン計測（Lighthouse CI — LCP, INP, CLS の現在値を記録）
* [ ] 現在のバンドルサイズの計測（web/user, web/liff, web/shared）
* [ ] ページごとの初期ロードサイズの記録
* [ ] アクセシビリティ監査の実施（axe-core）
* [ ] 現在の検索エンジンインデックス状況の確認（Google Search Console）
* [ ] デザイントークンの統一方針について合意
* [ ] 優先度の確定（SEO → アクセシビリティ → パフォーマンス → UI/UX の順を推奨）

### Phase 1: SEO 最適化（高優先度）
* [ ] JSON-LD 構造化データの実装（商品詳細・体験詳細・トップページ）
* [ ] 動的メタタグの実装（全対象ページ）
* [ ] canonical URL の設定
* [ ] hreflang タグの追加
* [ ] サイトマップの動的ルート対応
* [ ] robots.txt の最適化
* [ ] パンくずリストの実装
* [ ] セマンティック HTML の見直し

### Phase 2: アクセシビリティ改善（高優先度）
* [ ] 全アイコンボタンに `aria-label` を追加
* [ ] フォーム要素に `aria-invalid`, `aria-describedby` を追加
* [ ] スキップリンクの実装
* [ ] フォーカスマネジメントの実装（モーダル・メニュー）
* [ ] カラーコントラストの修正
* [ ] 画像の `alt` テキスト設定
* [ ] `autocomplete` 属性の追加（クレジットカードフォーム）
* [ ] Storybook a11y addon を "error" モードに変更

### Phase 3: パフォーマンス最適化（中優先度）
* [ ] `<img>` → `<NuxtImg>` への置き換え
* [ ] スケルトン UI の実装
* [ ] ルートルールの設定（SWR キャッシュ・プリレンダリング）
* [ ] サードパーティスクリプトの遅延読み込み
* [ ] 動画の遅延読み込み

### Phase 4: UI/UX 改善（段階的）
* [ ] サイト内検索の実装
* [ ] 商品一覧のフィルタリング・ソート機能
* [ ] トースト通知の統一
* [ ] 空状態コンポーネントの実装
* [ ] 数量ピッカーの実装（liff）
* [ ] 注文一覧のページネーション（liff）
* [ ] リアルタイムフォームバリデーション

### Phase 5: 共通基盤強化
* [ ] デザイントークンの統一（shared）
* [ ] 新規共通コンポーネントの追加
* [ ] ユニットテストの追加
* [ ] ビジュアルリグレッションテストの CI 統合

### 動作確認
* [ ] Lighthouse スコア: Performance ≥ 90, Accessibility ≥ 95, SEO ≥ 95
* [ ] Google Rich Results Test でリッチリザルト表示を確認
* [ ] axe-core でアクセシビリティ違反 0 件
* [ ] 全ページでキーボードナビゲーションが機能することを確認
* [ ] モバイル（iOS Safari, Android Chrome）での動作確認
* [ ] i18n 切り替え時の hreflang 正常動作

## リリース時確認事項

### リリース順

Phase ごとに段階的にリリースする。各 Phase は独立してリリース可能。

1. Phase 1（SEO） — Web のみ
2. Phase 2（アクセシビリティ） — Web のみ
3. Phase 3（パフォーマンス） — Web のみ
4. Phase 4（UI/UX） — Web のみ
5. Phase 5（共通基盤） — shared → user / liff の順

### リリース制御

特になし。各改善は段階的に適用可能。

### インフラ設定

- サイトマップ動的生成用の API エンドポイント追加（API 側）
- CDN キャッシュルールの見直し（CloudFront）
- Google Search Console での構造化データ検証

### パフォーマンスチェック

- リリース前後で Lighthouse CI スコアを比較
- Core Web Vitals（LCP, FID, CLS）の回帰がないことを確認
- バンドルサイズの増加量を確認（+10% 以内を目標）

### その他

- Google Search Console でインデックス状況をモニタリング
- リッチリザルトの表示状況を 2 週間追跡
- ユーザーからのアクセシビリティフィードバック収集体制を整備

## 関連リンク

- [Web Vitals](https://web.dev/vitals/)
- [WCAG 2.2](https://www.w3.org/TR/WCAG22/)
- [schema.org](https://schema.org/)
- [Nuxt SEO](https://nuxtseo.com/)
- [Nuxt Image](https://image.nuxt.com/)
