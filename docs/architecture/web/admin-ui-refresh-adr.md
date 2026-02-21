# 管理者Web UI/UX 刷新 - 技術選定方針設計（ADR）

本ドキュメントは、ふるマル管理者Web の UI/UX 刷新における主要な技術選定の判断根拠を Architecture Decision Record（ADR）形式で記録する。

---

## ADR-001: UI コンポーネントライブラリの選定

### ステータス: 決定済み

### コンテキスト

現行の管理者Webは Vuetify 3.10.4（MD3 Blueprint）を採用している。UI/UX 刷新にあたり、以下の選択肢を検討した。

### 検討した選択肢

| 選択肢 | 概要 | メリット | デメリット |
|--------|------|---------|----------|
| **A. Vuetify 3.x 継続** | 現行ライブラリを維持しカスタマイズ | 移行コスト最小、既存コード資産活用 | MD デザインの制約、バンドルサイズ |
| B. Nuxt UI v3 | Nuxt公式UIライブラリ（Reka UI + Tailwind v4） | Nuxt統合、モダンアーキテクチャ、125+コンポーネント | 全面書き換え必要、エコシステムが若い |
| C. PrimeVue v4 | エンタープライズ向けリッチUI | データテーブルが強力、Unstyled モード | デザインがやや硬質、畑テーマとの親和性が低い |
| D. shadcn-vue | コピー&ペースト方式のコンポーネント集 | 完全なコード所有、自由なカスタマイズ | 保守負荷が高い、自前でのアクセシビリティ保証 |

### 決定

**A. Vuetify 3.x を継続し、段階的にカスタマイズを強化する**

### 理由

1. **移行コストの最小化**: 79個のコンポーネント + 31個のストアが Vuetify に依存。全面移行は数ヶ月規模の作業
2. **既存資産の活用**: 現行の Vuetify コンポーネント設定（`defaults`）を畑テーマ向けに拡張するだけで大きな効果
3. **成熟したアクセシビリティ**: Vuetify は WAI-ARIA を組み込み済み。独自実装より信頼性が高い
4. **MD3 Blueprint**: Material Design 3 のベースラインがあり、カスタマイズの起点として優秀
5. **Vuetify 4.0 への備え**: Vuetify 4.0（現在alpha）は既存コードとの互換性を重視。将来の移行パスが明確

### トレードオフ

- Material Design のビジュアル言語からの完全な脱却は困難
- バンドルサイズの課題は `unplugin-vuetify` で対処が必要
- 一部の先進的なUI（ヘッドレスコンポーネント等）は Vuetify の制約を受ける

### 将来の見直し条件

- Vuetify 4.0 安定版リリース時に再評価
- Nuxt UI v3 のエコシステムが十分成熟した場合（2027年以降を見込む）

---

## ADR-002: スタイリング方針の選定

### ステータス: 決定済み

### コンテキスト

現行は SCSS + Vuetify Variables で構成。デザイントークンの体系的な管理ができていない。

### 検討した選択肢

| 選択肢 | 概要 |
|--------|------|
| **A. SCSS + CSS Custom Properties + Vuetify Variables** | 現行ベースに CSS Custom Properties によるデザイントークン層を追加 |
| B. Tailwind CSS v4 + Vuetify | Tailwind のユーティリティクラスを Vuetify と併用 |
| C. UnoCSS + Vuetify | アトミック CSS エンジンを併用 |

### 決定

**A. SCSS + CSS Custom Properties + Vuetify Variables の3層構成**

### 理由

1. **段階的移行が可能**: 既存の SCSS をそのまま活用しつつ、デザイントークンを上に積む
2. **Vuetify との親和性**: Vuetify の `$vuetify` SCSS 変数と CSS Custom Properties の両方でテーマ制御
3. **ランタイムコスト0**: CSS Custom Properties はブラウザネイティブ。JSバンドルに影響しない
4. **ダークモード対応**: `:root` と `.dark` でトークン値を切り替えるだけ
5. **Tailwind 導入の複雑性回避**: Vuetify と Tailwind を同居させると CSS 優先度の競合が頻発

### 実装方針

```
assets/styles/tokens/
  _colors.scss      → CSS Custom Properties としてカラートークン定義
  _spacing.scss     → スペーシングトークン定義
  _typography.scss  → タイポグラフィトークン定義
  _shadows.scss     → シャドウトークン定義
```

```scss
// _colors.scss
:root {
  --color-primary: #7cb342;
  --color-primary-light: #dcedc8;
  --color-accent: #ffb300;
  // ...
}

// Vuetify テーマとの連携
// plugins/vuetify.ts で同じ値を参照
```

---

## ADR-003: フォーム入力バリアントの選定

### ステータス: 決定済み

### コンテキスト

現行は Vuetify のテキスト入力に `variant: 'underlined'` を採用。高齢者のユーザビリティ観点で見直しが必要。

### 検討した選択肢

| バリアント | 特徴 | 高齢者向き度 |
|-----------|------|------------|
| underlined（現行） | 下線のみ。ミニマル | △ 入力領域が不明確 |
| **outlined（推奨）** | 全周囲に枠線。入力域が明確 | ◎ 領域認識が容易 |
| filled | 背景色付き。存在感あり | ○ だが画面が重くなりがち |
| solo | 独立した入力ボックス | ○ ただし Vuetify 3 では非推奨 |

### 決定

**`outlined` を全テキスト入力のデフォルトとする**

### 理由

1. **入力領域の明確化**: 枠線で囲まれた入力欄は、どこに入力すべきかが一目瞭然
2. **高齢者の認知負荷軽減**: 下線のみでは「ここが入力欄」と認識しにくい
3. **エラー状態の明確化**: 枠線全体が赤くなるため、エラー箇所が目立つ
4. **Web標準との一貫性**: 一般的なWebフォームの入力欄と同じ見た目で、学習コストが低い

### 移行方針

```typescript
// plugins/vuetify.ts の defaults を変更
defaults: {
  VTextField: { variant: 'outlined' },
  VTextarea: { variant: 'outlined' },
  VSelect: { variant: 'outlined' },
  VAutocomplete: { variant: 'outlined' },
  VCombobox: { variant: 'outlined' },
}
```

---

## ADR-004: Vuetify コンポーネントのインポート方式

### ステータス: 決定済み

### コンテキスト

現行は `import * from 'vuetify/components'` で全コンポーネントをバンドル。バンドルサイズの肥大化が懸念。

### 検討した選択肢

| 方式 | バンドルサイズ | 開発体験 |
|------|-------------|---------|
| **全量インポート（現行）** | 大（未使用コンポーネント含む） | 簡単（何でも使える） |
| 手動個別インポート | 最小 | 面倒（新コンポーネント追加時にimport追加） |
| **unplugin-vuetify** | 使用分のみ（自動） | 簡単（全量と同じ開発体験） |

### 決定

**`unplugin-vuetify`（vite-plugin-vuetify）を導入し自動 Tree-shaking を行う**

### 理由

1. **バンドルサイズ削減**: 使用しないコンポーネントを自動で除外
2. **開発体験の維持**: コード上は全量インポートと同じ書き方で利用可能
3. **Labs コンポーネントも対応**: 実験的コンポーネントも自動 Tree-shaking
4. **Vuetify 公式推奨**: Vuetify ドキュメントで推奨されている方式

### 実装方針

```typescript
// nuxt.config.ts
export default defineNuxtConfig({
  build: {
    transpile: ['vuetify'],
  },
  vite: {
    plugins: [
      vuetify({ autoImport: true }),
    ],
  },
})
```

---

## ADR-005: 状態管理パターンの選定

### ステータス: 決定済み

### コンテキスト

現行は Pinia 2.3.1 + Options API スタイルのストア定義。32個のストアファイルが存在。Pinia プラグインで API クライアントファクトリ（`this.apiMethodName()`）とエラーハンドラ（`this.errorHandler()`）を注入するパターンを採用しており、ボイラープレートが多い。

### 検討した選択肢

| パターン | 概要 |
|---------|------|
| Options API ストア（現行） | `state`, `getters`, `actions` のオブジェクト定義 |
| **Composition API ストア** | `ref`, `computed`, 関数によるセットアップ関数定義 |

### 決定

**新規ストアは Composition API スタイルで作成。既存ストアは改修時に段階的に移行**

### 理由

1. **Composition API との一貫性**: コンポーネントが `<script setup>` を使用しているため、ストアも同じパラダイムに統一
2. **TypeScript 推論の向上**: セットアップ関数内の型推論がより正確
3. **コンポーザブルとの共有**: ストアとコンポーザブルで同じパターンを使用でき、ロジックの共有が容易
4. **Vue/Pinia 公式推奨**: 2026年現在のベストプラクティスとして Composition API スタイルが推奨

### 移行方針

```typescript
// 新規作成: Composition API スタイル
export const useProductStore = defineStore('product', () => {
  const items = ref<Product[]>([])
  const loading = ref(false)

  const activeItems = computed(() =>
    items.value.filter(p => p.status === 'active')
  )

  async function fetchList(params: ListParams) { ... }

  return { items, loading, activeItems, fetchList }
})

// 既存: 改修対象になった時点で順次移行
```

---

## ADR-006: コンポーネントのドメイン別整理方針

### ステータス: 決定済み

### コンテキスト

現行の `components/` ディレクトリは Atomic Design の各層（atoms/molecules/organisms/templates）にフラットに配置。テンプレートだけで47ファイルあり、ファイル探索が困難。

### 検討した選択肢

| 方式 | 概要 |
|------|------|
| フラット（現行） | 各層にファイルを直置き |
| **ドメイン別サブディレクトリ** | 各層の中をドメイン（商品、注文等）で分類 |
| ドメインファースト | ドメインの中にatomic 層を配置（`products/atoms/`, `products/templates/`） |

### 決定

**Atomic Design の各層内にドメイン別サブディレクトリを設ける**

### 理由

1. **既存構造との互換性**: Atomic Design の層構造はそのまま維持。破壊的変更を避ける
2. **探索性の向上**: 「商品のテンプレート」→ `templates/products/` と直感的に辿れる
3. **Nuxt の auto-import との親和性**: サブディレクトリ化しても Nuxt の auto-import は動作
4. **段階的移行が可能**: 改修対象のファイルから順次移動

### 具体的なディレクトリ構成

```
components/
├── atoms/              # 層構造は変更なし
├── molecules/
│   ├── forms/          # フォーム系
│   ├── lists/          # リスト系
│   └── feedback/       # フィードバック系
├── organisms/
│   ├── products/       # 商品ドメイン
│   ├── orders/         # 注文ドメイン
│   ├── schedules/      # スケジュールドメイン
│   ├── videos/         # 動画ドメイン
│   └── coordinators/   # コーディネータードメイン
└── templates/
    ├── products/       # 商品ページテンプレート
    ├── orders/         # 注文ページテンプレート
    ├── schedules/      # スケジュールページテンプレート
    ├── auth/           # 認証ページテンプレート
    └── settings/       # 設定ページテンプレート
```

---

## ADR-007: アクセシビリティ基準の選定

### ステータス: 決定済み

### コンテキスト

ふるマルの管理者は農家（高齢者含む）が主要ユーザー。アクセシビリティへの体系的な取り組みが必要。

### 検討した選択肢

| 基準 | 概要 | 状態 |
|------|------|------|
| **WCAG 2.2 Level AA** | W3C 現行標準（2023年10月公開） | 安定・推奨 |
| WCAG 2.2 Level AAA | 最も厳格な基準 | 過剰。全ページ適用は非現実的 |
| WCAG 3.0 | 次世代標準（Working Draft） | 未完成。2028年以降の策定見込み |
| JIS X 8341-3 | 日本工業規格のアクセシビリティ指針 | WCAG 2.x と整合 |

### 決定

**WCAG 2.2 Level AA を準拠基準とする**

### 理由

1. **国際標準**: W3C の最新勧告。日本の JIS X 8341-3 とも整合
2. **AA レベルで十分な実用性**: AAA は一部のコンテンツでは達成困難
3. **ユーザー特性への適合**: WCAG 2.2 は高齢者・認知障害者向けの新基準（2.5.7 ドラッグ操作の代替、2.5.8 ターゲットサイズ等）を追加
4. **WCAG 3.0 は時期尚早**: 2028年以降の策定見込みで、現時点での準拠は不可能
5. **Vuetify の既存対応**: Vuetify はWAI-ARIA対応済み。追加対応が少なく済む

### 重点対応項目（WCAG 2.2 新基準）

| 基準 | 内容 | 管理者Webでの対応 |
|------|------|-----------------|
| 2.4.11 Focus Not Obscured | フォーカスが隠れない | スティッキーヘッダーの下にフォーカスが隠れないよう `scroll-margin-top` を設定 |
| 2.5.7 Dragging Movements | ドラッグの代替手段 | SortableProductList に上下ボタンを追加 |
| 2.5.8 Target Size | 最小 24×24px | 全インタラクティブ要素を 44×44px 以上に |
| 3.3.7 Redundant Entry | 再入力の回避 | フォームの自動補完、住所入力の郵便番号連携 |
| 3.3.8 Accessible Authentication | 認知テスト不要の認証 | CAPTCHAは使用しない。ソーシャルログイン推奨 |

---

## ADR-008: フォントの選定

### ステータス: 決定済み

### コンテキスト

現行は BIZ UDPGothic を使用。ウェイトのバリエーションが限定的（Regular, Bold のみ）。

### 決定

**BIZ UDPGothic を主フォントとして維持。必要に応じて Noto Sans JP でウェイトを補完**

### 理由

1. **UD書体の優位性**: BIZ UDPGothic はモリサワが開発したユニバーサルデザイン書体。高齢者・弱視者の可読性に最適化
2. **日本語Web標準**: Google Fonts で提供。追加コストなし
3. **既存との互換**: フォント変更による既存レイアウトへの影響を回避
4. **Noto Sans JP は補完用**: BIZ UDPGothic が対応しない Medium (500)、Semibold (600) ウェイトが必要な場合のフォールバック

```scss
--font-family-primary: "BIZ UDPGothic", "Noto Sans JP", sans-serif;
```

---

## ADR-009: Nuxt バージョンアップ方針

### ステータス: 決定済み

### コンテキスト

現行は Nuxt 4.1.2（`compatibilityVersion: 4` 設定済み）。2026年2月現在の最新安定版は Nuxt 4.3.x。

### 決定

**Nuxt 4.3.x にアップグレードする。Nuxt 5 は 2026年後半以降に評価**

### 理由

1. **マイナーアップデート**: 4.1.2 → 4.3.x は破壊的変更が少ない
2. **パフォーマンス改善**: フック最適化、ルーティング改善等が含まれる
3. **互換性設定済み**: `future: { compatibilityVersion: 4 }` により v4 の新ディレクトリ構造に対応済み
4. **Nuxt 3 のサポート終了**: 2026年7月に Nuxt 3 のメンテナンスが終了。v4 系の最新を維持する必要性

### アップグレード対象パッケージ

| パッケージ | 現行 | 目標 | 備考 |
|-----------|------|------|------|
| nuxt | 4.1.2 | 4.3.x | フレームワーク本体 |
| @nuxt/devtools | 2.6.5 | 3.2.x | 開発ツール |
| vue | 3.5.21 | 3.5.x 最新 | 安定版維持 |
| vuetify | 3.10.4 | 3.10.x 最新 | パッチ更新 |
| vite | 7.1.11 | 7.x 最新 | パッチ更新 |
| typescript | 5.9.3 | 5.9.x 最新 | パッチ更新 |
| @sentry/vue | 7.120.4 | 8.x | メジャーアップデート（別途計画） |

### 監視対象（将来の採用候補）

| 技術 | 現状 | 採用見込み |
|------|------|-----------|
| Vue 3.6 (Vapor Mode) | Beta 6（2026年2月） | 安定版リリース後に評価（2026年後半） |
| Vite 8 | Beta | 安定版リリース後に評価 |
| Vuetify 4.0 | Alpha | 安定版リリース後に評価 |
| Nuxt 5 | 開発中 | 2027年以降 |

---

## ADR-010: 重量級ライブラリの遅延読み込み方針

### ステータス: 決定済み

### コンテキスト

Chart.js、ECharts、TipTap、HLS.js など、特定ページでのみ使用する重量級ライブラリがバンドルに含まれている。

### 決定

**特定ページ専用のライブラリは `defineAsyncComponent` および動的 `import()` で遅延読み込みする**

### 対象ライブラリと適用方針

| ライブラリ | サイズ（概算） | 使用ページ | 遅延読み込み方式 |
|-----------|-------------|-----------|---------------|
| echarts | ~800KB | ダッシュボード、分析画面 | 動的import + defineAsyncComponent |
| chart.js + vue-chart-3 | ~200KB | ダッシュボード | 動的import + defineAsyncComponent |
| @tiptap/* | ~300KB | コンテンツ編集画面 | 動的import + defineAsyncComponent |
| hls.js | ~200KB | ライブ配信画面 | 動的import |
| sortablejs | ~40KB | 商品並べ替え画面 | 動的import |

### 実装例

```typescript
// TipTap エディタの遅延読み込み
const TiptapEditor = defineAsyncComponent(() =>
  import('~/components/TiptapEditor.vue')
)

// ECharts の遅延読み込み
const ScheduleAnalytics = defineAsyncComponent({
  loader: () => import('~/components/organisms/schedules/ScheduleAnalytics.vue'),
  loadingComponent: LoadingSpinner,
  delay: 200,
})
```

---

## ADR-011: CSS モダン機能の採用方針

### ステータス: 決定済み

### コンテキスト

CSS Container Queries、`:has()`、`@layer` 等のモダンCSS機能が主要ブラウザで安定サポートされている。

### 決定

**以下のモダンCSS機能を段階的に採用する**

| 機能 | ブラウザサポート | 採用方針 |
|------|--------------|---------|
| CSS Custom Properties | 99%+ | 即時採用（デザイントークン基盤） |
| `:has()` | 95%+ | 即時採用（既に一部使用中） |
| Container Queries | 95%+ | 段階的採用（再利用コンポーネントから） |
| `@layer` | 95%+ | 段階的採用（Vuetify との優先度制御） |
| CSS Nesting | 90%+ | 新規SCSSから段階的に移行 |
| `color-mix()` | 90%+ | デザイントークンのバリエーション生成に |

### 理由

1. **管理者WebはSPA**: 対象ブラウザが限定的（Chrome/Edge/Safari最新版）
2. **ランタイムコスト0**: CSS機能はJSバンドルに影響しない
3. **既存使用実績**: `:has()` は `main.scss` で既に使用中
4. **Container Queries の実用性**: コンポーネントの再利用性向上に直結

### Container Queries の適用候補

```scss
// 例: 商品カードが一覧・詳細・ダッシュボードで異なるレイアウトを取る場合
.product-card-container {
  container-type: inline-size;
}

@container (min-width: 400px) {
  .product-card { /* 横並びレイアウト */ }
}

@container (max-width: 399px) {
  .product-card { /* 縦積みレイアウト */ }
}
```

---

## 決定の要約

| ADR | 決定事項 | 採用理由の要点 |
|-----|---------|-------------|
| 001 | Vuetify 3.x 継続 | 移行コスト最小、既存資産活用 |
| 002 | SCSS + CSS Custom Properties | 段階的移行、ランタイムコスト0 |
| 003 | フォーム outlined バリアント | 高齢者の入力領域認識向上 |
| 004 | unplugin-vuetify 導入 | バンドルサイズ削減、開発体験維持 |
| 005 | Composition API ストア | コンポーネントとの一貫性、型推論向上 |
| 006 | ドメイン別サブディレクトリ | 探索性向上、段階的移行可能 |
| 007 | WCAG 2.2 Level AA | 国際標準、高齢者配慮の新基準対応 |
| 008 | BIZ UDPGothic 維持 | UD書体の可読性、既存互換 |
| 009 | Nuxt 4.3.x アップグレード | マイナー更新、パフォーマンス改善 |
| 010 | 重量級ライブラリの遅延読み込み | バンドルサイズ削減、初期ロード高速化 |
| 011 | モダンCSS段階的採用 | ブラウザサポート十分、コスト0 |

---

## 関連ドキュメント

- [全体設計（アーキテクチャ）](./admin-ui-refresh.md)
- [デザインガイドライン](./admin-design-guidelines.md)
- [改修設計書](../../spec/internal/20260221_admin-ui-ux-refresh.md)
- [既存の設計決定](../design-decisions.md)
