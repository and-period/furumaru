# 管理者Web UI/UX 刷新 - 全体設計（アーキテクチャ）

## 概要

本ドキュメントは、ふるマル管理者Web（`web/admin`）の UI/UX 刷新プロジェクトの全体アーキテクチャを定義する。農家を主要ユーザーとする管理画面を、2026年現在のベストプラクティスに基づき、アクセシビリティとユーザビリティを重視した「畑をテーマにした柔らかいデザイン」へ刷新する。

## 現状アーキテクチャ

### 技術スタック（現行）

| カテゴリ | 技術 | バージョン |
|---------|------|-----------|
| フレームワーク | Nuxt | 4.1.2 |
| UI ライブラリ | Vue.js | 3.5.21 |
| コンポーネント | Vuetify | 3.10.4 (MD3 Blueprint) |
| ビルドツール | Vite | 7.1.11 |
| 状態管理 | Pinia | 2.3.1 |
| スタイリング | SCSS + Vuetify Variables |
| フォント | BIZ UDPGothic |
| レンダリング | SPA（SSR無効） |
| TypeScript | 5.9.3 |

### ディレクトリ構成（現行）

```
web/admin/src/
├── assets/               # SCSS スタイル（main.scss, variables.scss）
├── components/
│   ├── atoms/            # 3 コンポーネント（AppLogoWithTitle, AppTitle, FileUploadFiled）
│   ├── molecules/        # 11 コンポーネント（フォーム部品、ソート可能リスト等）
│   ├── organisms/        # 18 コンポーネント（カレンダー、分析、動画関連等）
│   └── templates/        # 52 コンポーネント（各ページのテンプレート、一部20KB超の大規模ファイル）
├── composables/          # 4 composable（通知、プロモーション関連）
├── constants/            # 12 定数ファイル（都道府県、市区町村、支払い等。cities.ts は 346KB）
├── layouts/              # 3 レイアウト（default, auth, error）
├── lib/                  # 共通ユーティリティ
│   ├── externals/        # 外部サービス連携
│   ├── formatter/        # データフォーマッター
│   ├── helpers/          # ヘルパー関数（画像リサイズ、都道府県等）
│   ├── hooks/            # 再利用可能 composable（useAlert, usePagination, useSearchAddress）
│   ├── prefectures/      # 都道府県データ
│   └── validations/      # Vuelidate カスタムバリデーション（かな、電話番号、郵便番号等）
├── middleware/           # 2 ミドルウェア（auth, notification）
├── pages/                # 約30 ページディレクトリ
├── plugins/              # 6 プラグイン（vuetify, firebase, sentry, api-client, api-error-handler等）
├── store/                # 32 ストアファイル（Options API スタイル + Pinia プラグインで API クライアント注入）
└── types/                # TypeScript 型定義
    ├── api/v1/           # OpenAPI 自動生成の API 型
    ├── enum.ts           # カスタム列挙型
    ├── exception.ts      # カスタム例外クラス
    ├── plugins/          # プラグイン型
    ├── props/            # コンポーネント Props 型
    └── validations/      # バリデーション型
```

### 現行の課題

1. **Vuetify コンポーネントの全量インポート**: `import * from 'vuetify/components'` によるバンドルサイズの肥大化
2. **コンポーネントの機能分割不足**: テンプレートコンポーネント（52個）が肥大化（一部20KB超）。ページロジックとUIが密結合
3. **composable の分散と不足**: `composables/` に4ファイル、`lib/hooks/` に3ファイルが分散。ビジネスロジックがコンポーネントやストア内に散在
4. **ストアの Options API パターン**: 32個のストアが Options API + Pinia プラグイン注入（`this.apiMethodName()`）パターンで記述。ボイラープレートが多い
5. **デザイントークンの不在**: カラーやスペーシングがVuetifyのデフォルト値に依存
6. **アクセシビリティ未対応**: WCAG準拠の体系的な取り組みがない
7. **テストの不足**: フロントエンドのユニットテスト・E2Eテストが未整備（テストファイルが存在しない）
8. **定数ファイルの肥大化**: `constants/cities.ts` が 346KB。都道府県・市区町村データのインライン定義

---

## 刷新後アーキテクチャ

### 技術スタック（目標）

| カテゴリ | 技術 | バージョン | 変更理由 |
|---------|------|-----------|---------|
| フレームワーク | Nuxt | 4.3.x | 最新安定版へ更新 |
| UI ライブラリ | Vue.js | 3.5.x | 安定版維持（3.6 Vapor Mode は監視） |
| コンポーネント | Vuetify | 3.10.x → 段階的に最新へ | 既存資産を活用しつつ最適化 |
| ビルドツール | Vite | 7.x | 現行維持（Rolldown ベース） |
| 状態管理 | Pinia | 2.3.x | 現行維持 |
| スタイリング | SCSS + CSS Custom Properties + Vuetify Variables |
| フォント | BIZ UDPGothic + Noto Sans JP（ウェイト補完） |
| アイコン | @mdi/js（維持） |
| DevTools | @nuxt/devtools 3.2.x | 最新版へ更新 |
| アクセシビリティ | WCAG 2.2 Level AA 準拠 |

### アーキテクチャ方針

```
┌─────────────────────────────────────────────────────────────────┐
│                      Design System Layer                        │
│  ┌──────────┐  ┌──────────┐  ┌────────────┐  ┌──────────────┐ │
│  │ Design   │  │ Theme    │  │ Typography │  │ Accessibility│ │
│  │ Tokens   │  │ Config   │  │ System     │  │ Primitives   │ │
│  └──────────┘  └──────────┘  └────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                      Component Layer                            │
│  ┌──────────┐  ┌──────────┐  ┌────────────┐  ┌──────────────┐ │
│  │ Atoms    │  │Molecules │  │ Organisms  │  │ Templates    │ │
│  │ (基礎)   │  │(部品)    │  │(機能ブロック)│ │(ページ構成)   │ │
│  └──────────┘  └──────────┘  └────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                      Business Logic Layer                       │
│  ┌──────────┐  ┌──────────┐  ┌────────────┐  ┌──────────────┐ │
│  │Composable│  │ Pinia    │  │ API Client │  │ Validation   │ │
│  │ (hooks)  │  │ Stores   │  │ (axios)    │  │ (vuelidate)  │ │
│  └──────────┘  └──────────┘  └────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                      Infrastructure Layer                       │
│  ┌──────────┐  ┌──────────┐  ┌────────────┐  ┌──────────────┐ │
│  │ Nuxt     │  │ Vite     │  │ Firebase   │  │ Sentry       │ │
│  │ Runtime  │  │ Build    │  │ Auth/FCM   │  │ Monitoring   │ │
│  └──────────┘  └──────────┘  └────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

### ディレクトリ構成（刷新後）

```
web/admin/src/
├── assets/
│   ├── styles/
│   │   ├── main.scss             # グローバルスタイル
│   │   ├── variables.scss        # Vuetify 変数オーバーライド
│   │   └── tokens/
│   │       ├── _colors.scss      # カラートークン定義
│   │       ├── _spacing.scss     # スペーシングトークン
│   │       ├── _typography.scss  # タイポグラフィトークン
│   │       └── _shadows.scss     # シャドウ・エレベーション
│   └── images/                   # 静的画像アセット
│
├── components/
│   ├── atoms/                    # 基礎コンポーネント
│   │   ├── AppButton.vue         # ボタン（サイズ・バリアント対応）
│   │   ├── AppIcon.vue           # アイコンラッパー
│   │   ├── AppBadge.vue          # バッジ・ステータス表示
│   │   ├── AppLogoWithTitle.vue  # ロゴ（既存）
│   │   └── AppTitle.vue          # タイトル（既存）
│   │
│   ├── molecules/                # 複合コンポーネント
│   │   ├── forms/                # フォーム系コンポーネント
│   │   │   ├── AddressForm.vue
│   │   │   ├── ImageSelectForm.vue
│   │   │   └── ShippingRateForm.vue
│   │   ├── lists/                # リスト系コンポーネント
│   │   │   ├── SortableProductList.vue
│   │   │   └── SortableProductThumbnail.vue
│   │   └── feedback/             # フィードバック系
│   │       ├── AppAlert.vue
│   │       └── AppConfirmDialog.vue
│   │
│   ├── organisms/                # 機能ブロック
│   │   ├── products/             # 商品関連
│   │   │   └── ProductTypeList.vue
│   │   ├── schedules/            # スケジュール関連
│   │   │   ├── ScheduleAnalytics.vue
│   │   │   ├── ScheduleShow.vue
│   │   │   └── ScheduleStreaming.vue
│   │   ├── videos/               # 動画関連
│   │   │   ├── VideoAnalytics.vue
│   │   │   ├── VideoForm.vue
│   │   │   └── VideoPreview.vue
│   │   ├── orders/               # 注文関連
│   │   │   └── OrderExperienceDetails.vue
│   │   └── coordinators/         # コーディネーター関連
│   │       ├── CoordinatorShow.vue
│   │       └── CoordinatorShipping.vue
│   │
│   └── templates/                # ページテンプレート（既存構成を維持）
│       ├── products/             # 商品ページテンプレート
│       ├── orders/               # 注文ページテンプレート
│       ├── schedules/            # スケジュールページテンプレート
│       ├── auth/                 # 認証ページテンプレート
│       └── settings/             # 設定ページテンプレート
│
├── composables/                  # ビジネスロジック（拡充。既存 lib/hooks/ を統合）
│   ├── useAuth.ts                # 認証ロジック
│   ├── useApiClient.ts           # API 呼び出しラッパー
│   ├── useAlert.ts               # アラート表示（lib/hooks/ から移動）
│   ├── useFormValidation.ts      # フォームバリデーション共通
│   ├── usePagination.ts          # ページネーション共通（lib/hooks/ から移動）
│   ├── useSearchAddress.ts       # 住所検索（lib/hooks/ から移動）
│   ├── useFileUpload.ts          # ファイルアップロード共通
│   ├── useNotificationDisplay.ts # 通知表示（既存）
│   ├── useNotificationForm.ts    # 通知フォーム（既存）
│   ├── usePromotionForm.ts       # プロモーションフォーム（既存）
│   ├── usePromotionValidation.ts # プロモーションバリデーション（既存）
│   └── useAccessibility.ts       # アクセシビリティヘルパー
│
├── layouts/                      # レイアウト（既存構成を維持）
├── middleware/                   # ミドルウェア（既存構成を維持）
├── pages/                        # ページ（既存構成を維持）
├── plugins/                      # プラグイン（既存構成を維持）
├── store/                        # Pinia ストア（既存構成を維持、段階的に Composition API 移行）
├── lib/                          # 共通ユーティリティ（既存構成を維持）
│   ├── externals/                # 外部サービス連携
│   ├── formatter/                # データフォーマッター
│   ├── helpers/                  # ヘルパー関数
│   ├── prefectures/              # 都道府県データ
│   └── validations/              # Vuelidate カスタムバリデーション
├── constants/                    # 定数定義（既存構成を維持）
└── types/                        # TypeScript 型定義（OpenAPI 自動生成含む）
```

---

## デザインシステム設計

### デザインテーマ: 「畑（Hatake）」

ふるマルの管理者画面は、農家の日常に寄り添う温かみのあるデザインを採用する。

#### テーマコンセプト

```
┌─────────────────────────────────────────────────────┐
│                   畑テーマの構成要素                    │
├──────────────┬──────────────────────────────────────┤
│ 色彩         │ 大地の茶、若葉の緑、収穫の橙、空の青     │
│ 形状         │ 角丸（柔らかさ）、カード型レイアウト       │
│ タイポグラフィ │ BIZ UDPGothic（UD書体 = 読みやすさ）    │
│ アイコン      │ Material Design Icons（丸みアウトライン）│
│ イラスト      │ 農産物・季節をモチーフにした装飾（最小限） │
│ トーン       │ 温かく親しみやすい、専門用語を避ける       │
└──────────────┴──────────────────────────────────────┘
```

### カラーシステム

3層のデザイントークンアーキテクチャを採用する。

```
Layer 1: Primitive Tokens（原色定義）
    ↓
Layer 2: Semantic Tokens（意味的な色）
    ↓
Layer 3: Component Tokens（コンポーネント固有の色）
```

#### Primitive Tokens（原色パレット）

```scss
// 畑テーマ カラーパレット
$hatake-green-50:  #f1f8e9;   // 若葉のかすかな緑
$hatake-green-100: #dcedc8;   // 新芽の淡い緑
$hatake-green-200: #c5e1a5;   // 春の若葉
$hatake-green-300: #aed581;   // 成長する緑
$hatake-green-400: #9ccc65;   // 畑の緑
$hatake-green-500: #8bc34a;   // 生命力ある緑
$hatake-green-600: #7cb342;   // 深い若葉（プライマリ候補）
$hatake-green-700: #689f38;   // 森の入口（現行primary: lightGreen.darken2）
$hatake-green-800: #558b2f;   // 深い森の緑
$hatake-green-900: #33691e;   // 濃い常緑

$hatake-earth-50:  #efebe9;   // 乾いた土
$hatake-earth-100: #d7ccc8;   // 明るい土
$hatake-earth-200: #bcaaa4;   // 砂の色
$hatake-earth-300: #a1887f;   // 土の色
$hatake-earth-400: #8d6e63;   // 肥沃な土
$hatake-earth-500: #795548;   // 耕された土
$hatake-earth-600: #6d4c41;   // 深い土

$hatake-harvest-50:  #fff8e1;  // 淡い陽光
$hatake-harvest-100: #ffecb3;  // 朝焼け
$hatake-harvest-200: #ffe082;  // 麦の穂
$hatake-harvest-300: #ffd54f;  // 収穫の金
$hatake-harvest-400: #ffca28;  // 実りの黄金
$hatake-harvest-500: #ffc107;  // 豊穣
$hatake-harvest-600: #ffb300;  // 夕暮れ（現行accent: amber.darken1）

$hatake-sky-50:  #e0f7fa;     // 朝の空
$hatake-sky-100: #b2ebf2;     // 晴天
$hatake-sky-200: #80deea;     // 澄んだ空
$hatake-sky-300: #4dd0e1;     // 清流
$hatake-sky-400: #26c6da;     // 湧き水（現行info: teal.lighten1）
```

#### Semantic Tokens（意味的カラー）

```scss
// ライトテーマ
:root {
  // ブランドカラー
  --color-primary:        #{$hatake-green-600};     // メインアクション、ナビゲーション
  --color-primary-light:  #{$hatake-green-100};     // ホバー背景、選択状態
  --color-primary-dark:   #{$hatake-green-800};     // アクティブ状態

  // アクセントカラー
  --color-accent:         #{$hatake-harvest-600};   // 注目ポイント、CTA
  --color-accent-light:   #{$hatake-harvest-100};   // アクセント背景

  // サーフェス
  --color-surface:        #ffffff;                  // カード・パネル背景
  --color-surface-dim:    #{$hatake-green-50};      // ページ背景（薄緑で柔らかさ）
  --color-surface-warm:   #{$hatake-earth-50};      // ウォームサーフェス

  // テキスト
  --color-text-primary:   #1a1a1a;                  // 本文テキスト
  --color-text-secondary: #616161;                  // 補助テキスト
  --color-text-hint:      #9e9e9e;                  // ヒントテキスト

  // セマンティック
  --color-success:        #66bb6a;                  // 成功・完了
  --color-warning:        #{$hatake-harvest-500};   // 警告・注意
  --color-error:          #ef5350;                  // エラー・危険
  --color-info:           #{$hatake-sky-400};       // 情報・ヘルプ
}
```

### タイポグラフィシステム

```scss
// フォントファミリー
--font-family-primary: "BIZ UDPGothic", "Noto Sans JP", sans-serif;

// フォントサイズスケール（高齢者に配慮した大きめベース）
--font-size-xs:    0.8125rem;   // 13px - キャプション
--font-size-sm:    0.875rem;    // 14px - 補助テキスト
--font-size-base:  1rem;        // 16px - 本文ベース
--font-size-md:    1.125rem;    // 18px - 強調本文
--font-size-lg:    1.25rem;     // 20px - 小見出し
--font-size-xl:    1.5rem;      // 24px - 見出し
--font-size-2xl:   1.875rem;    // 30px - 大見出し
--font-size-3xl:   2.25rem;     // 36px - ページタイトル

// 行の高さ
--line-height-tight:   1.25;
--line-height-normal:  1.625;    // 日本語テキストに適した行間
--line-height-relaxed: 1.875;    // 読みやすさ重視
```

### スペーシングシステム

```scss
// 4px ベースのスケール
--space-1:  0.25rem;   // 4px
--space-2:  0.5rem;    // 8px
--space-3:  0.75rem;   // 12px
--space-4:  1rem;      // 16px
--space-5:  1.25rem;   // 20px
--space-6:  1.5rem;    // 24px
--space-8:  2rem;      // 32px
--space-10: 2.5rem;    // 40px
--space-12: 3rem;      // 48px
--space-16: 4rem;      // 64px

// コンポーネントスペーシング
--spacing-card-padding:    var(--space-6);
--spacing-section-gap:     var(--space-8);
--spacing-form-gap:        var(--space-5);
--spacing-inline-gap:      var(--space-3);
```

### シャドウ・エレベーション

```scss
// 柔らかいシャドウ（畑テーマ）
--shadow-sm:  0 1px 3px rgba(0, 0, 0, 0.06), 0 1px 2px rgba(0, 0, 0, 0.04);
--shadow-md:  0 4px 6px rgba(0, 0, 0, 0.05), 0 2px 4px rgba(0, 0, 0, 0.03);
--shadow-lg:  0 10px 15px rgba(0, 0, 0, 0.05), 0 4px 6px rgba(0, 0, 0, 0.03);
--shadow-xl:  0 20px 25px rgba(0, 0, 0, 0.05), 0 8px 10px rgba(0, 0, 0, 0.02);

// カード用シャドウ（温かみのある影）
--shadow-card:       0 2px 8px rgba(104, 159, 56, 0.08);
--shadow-card-hover: 0 4px 16px rgba(104, 159, 56, 0.12);
```

### 角丸（Border Radius）

```scss
// 丸みのあるデザイン（柔らかさの表現）
--radius-sm:   6px;     // ボタン、インプット
--radius-md:   10px;    // カード、パネル
--radius-lg:   16px;    // モーダル、ダイアログ
--radius-xl:   24px;    // 特殊要素
--radius-full: 9999px;  // アバター、バッジ
```

---

## アクセシビリティ設計

### WCAG 2.2 Level AA 準拠方針

| カテゴリ | 要件 | 対応方針 |
|---------|------|---------|
| コントラスト比 | 通常テキスト 4.5:1 以上、大テキスト 3:1 以上 | デザイントークンで保証 |
| タッチターゲット | 最小 44×44px（推奨 48×48px） | Vuetify defaults でサイズ制御 |
| フォーカス表示 | 全操作可能要素に視認可能なフォーカスインジケータ | CSS でカスタムフォーカスリング |
| キーボード操作 | 全機能がキーボードのみで操作可能 | Vuetify の既存対応 + 追加対応 |
| ドラッグ代替 | ドラッグ操作に代替手段を提供 | 並べ替えボタンの追加 |
| エラー表示 | 色だけでなくアイコン・テキストで表示 | エラーコンポーネントの統一 |
| 一貫したヘルプ | ヘルプ機能を一貫した場所に配置 | ヘッダーにヘルプボタン常設 |

### 高齢者・低ITリテラシーユーザー向け配慮

1. **大きなクリックターゲット**: ボタン最小高さ 48px、十分なパディング
2. **明確なラベリング**: アイコン単独での操作を避け、必ずテキストラベルを併記
3. **段階的な操作**: 複雑な操作はウィザード形式で段階的に案内
4. **確認ダイアログ**: 破壊的操作（削除等）は必ず確認ステップを設ける
5. **平易な日本語**: 専門用語を避け、分かりやすい表現を使用
6. **色+アイコン+テキスト**: 状態表示は必ず3要素で伝達（色覚多様性への配慮）
7. **進捗表示**: 処理中は常にプログレスインジケータを表示
8. **パンくずリスト**: 現在位置を常に表示し、迷子にならない設計

---

## コンポーネントアーキテクチャ

### Atomic Design の拡張

```
┌────────────────────────────────────────────────────────────────┐
│ Templates（ページテンプレート）                                   │
│  - ページ全体のレイアウトと構成                                    │
│  - ドメイン別にサブディレクトリで整理                               │
│  - 例: ProductEdit, OrderShow, ScheduleNew                     │
├────────────────────────────────────────────────────────────────┤
│ Organisms（機能ブロック）                                        │
│  - 独立した機能を持つUIブロック                                    │
│  - ドメイン別にサブディレクトリで整理                               │
│  - 例: ScheduleAnalytics, VideoForm, CoordinatorShow           │
├────────────────────────────────────────────────────────────────┤
│ Molecules（複合コンポーネント）                                    │
│  - 機能別にサブディレクトリで整理（forms/, lists/, feedback/）      │
│  - 複数のAtomを組み合わせた再利用可能なUI部品                       │
│  - 例: AddressForm, SortableProductList, ImageSelectForm       │
├────────────────────────────────────────────────────────────────┤
│ Atoms（基礎コンポーネント）                                       │
│  - Vuetify コンポーネントのラッパーまたは独自の最小単位             │
│  - デザイントークンを適用した基礎部品                               │
│  - 例: AppButton, AppBadge, AppIcon                            │
└────────────────────────────────────────────────────────────────┘
```

### コンポーネント設計原則

1. **単一責務**: 1コンポーネント = 1機能
2. **Props-down, Events-up**: データは親から子へ、イベントは子から親へ
3. **Composition API**: `<script setup lang="ts">` を全コンポーネントで統一
4. **型安全**: Props/Emits は TypeScript で厳密に型定義
5. **スコープドスタイル**: `<style scoped>` + デザイントークン参照
6. **テスト可能性**: ビジネスロジックは composable に分離

### Composable 設計方針

```typescript
// composable の命名規則
useXxx()        // 汎用composable
useXxxForm()    // フォームロジック
useXxxList()    // リスト操作ロジック
useXxxApi()     // API通信ロジック

// composable 設計パターン
export function useProductForm(initialData?: Product) {
  // State
  const formData = ref<ProductFormData>({ ... })
  const loading = ref(false)
  const errors = ref<Record<string, string>>({})

  // Computed
  const isValid = computed(() => Object.keys(errors.value).length === 0)
  const isDirty = computed(() => /* 変更検知ロジック */)

  // Actions
  async function submit() { ... }
  function reset() { ... }
  function validate() { ... }

  return {
    formData,
    loading,
    errors,
    isValid,
    isDirty,
    submit,
    reset,
    validate,
  }
}
```

---

## 状態管理設計

### Pinia ストア設計方針

```
store/
├── auth.ts           # 認証状態（セッション、トークン）
├── common.ts         # グローバルUI状態（ローディング、通知）
├── admin.ts          # 管理者データ
├── product.ts        # 商品データ
├── order.ts          # 注文データ
├── schedule.ts       # スケジュールデータ
├── ...               # 各ドメインストア
└── index.ts          # ストアエクスポート
```

### ストア設計パターン

```typescript
// Composition API スタイル（推奨）
export const useProductStore = defineStore('product', () => {
  // State
  const items = ref<Product[]>([])
  const current = ref<Product | null>(null)
  const loading = ref(false)
  const total = ref(0)

  // Getters
  const activeItems = computed(() =>
    items.value.filter(item => item.status === 'active')
  )

  // Actions
  async function fetchList(params: ListParams) {
    loading.value = true
    try {
      const response = await apiClient.products.list(params)
      items.value = response.data.items
      total.value = response.data.total
    } finally {
      loading.value = false
    }
  }

  async function fetchById(id: string) {
    loading.value = true
    try {
      current.value = await apiClient.products.get(id)
    } finally {
      loading.value = false
    }
  }

  return { items, current, loading, total, activeItems, fetchList, fetchById }
})
```

---

## パフォーマンス最適化方針

### バンドルサイズ削減

1. **Vuetify Tree-shaking**: `unplugin-vuetify` を導入し、使用コンポーネントのみをバンドル
2. **動的インポート**: Chart.js, ECharts, TipTap, HLS.js を遅延読み込み
3. **ルートベースコード分割**: Nuxt のページ自動分割を活用

### レンダリング最適化

1. **仮想スクロール**: 大量データのリスト表示に `v-virtual-scroll` を活用
2. **画像最適化**: `@nuxt/image` による最適化（将来検討）
3. **コンポーネントの遅延読み込み**: `defineAsyncComponent` の活用

---

## 関連ドキュメント

- [デザインガイドライン](./admin-design-guidelines.md) - 畑テーマのデザインガイドライン詳細
- [技術選定 ADR](./admin-ui-refresh-adr.md) - 技術選定の判断根拠
- [改修設計書](../../spec/internal/20260221_admin-ui-ux-refresh.md) - 段階的な改修手順と設計詳細
- [フロントエンドアプリケーション構成](./frontend-applications.md) - 全Webアプリ横断のアーキテクチャ
