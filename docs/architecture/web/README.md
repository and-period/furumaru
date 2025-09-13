# Web フロントエンドアーキテクチャ

## 📋 目次

- **[components.md](./components.md)** - コンポーネント設計・分類・テスト戦略
- **[state-management.md](./state-management.md)** - Pinia状態管理・永続化・ストア連携
- **[api-integration.md](./api-integration.md)** - API連携・認証・エラーハンドリング

## アプリケーション構成

### `/web/admin` - 管理者ポータル
- **Framework**: Nuxt 4 (Compatibility Mode)
- **UI**: Vuetify 3 + Material Design
- **State**: Pinia
- **Rendering**: SPA (SSR無効)
- **目的**: 管理者向け管理ツール

**特徴:**
- リッチエディタ (TipTap) 搭載
- Firebase連携 (認証・プッシュ通知)
- チャート・グラフ表示 (Chart.js, ECharts)
- Sentry エラー監視

### `/web/user` - 購入者ポータル
- **Framework**: Nuxt 3
- **UI**: Tailwind CSS
- **State**: Pinia + Persisted State
- **Rendering**: SSR有効
- **目的**: ECサイト・ライブコマース

**特徴:**
- 多言語対応 (i18n)
- SEO対応 (SSR + Meta Tags)
- Google Maps連携
- 動画配信 (HLS.js)
- microCMS連携

### `/web/liff` - LINEミニアプリ
- **Framework**: Nuxt 3
- **UI**: Tailwind CSS
- **State**: Pinia
- **Rendering**: SPA
- **目的**: LINE内購買体験・チャット統合

**特徴:**
- LINE LIFF SDK v2.27+ 連携
- 軽量・高速起動
- ユーザー認証連携
- チャット機能統合

### `/web/shared` - 共通コンポーネントライブラリ
- **Framework**: Vue 3 + TypeScript
- **Build**: Vite
- **Documentation**: Storybook
- **Testing**: Vitest
- **目的**: デザインシステム・再利用可能コンポーネント

**特徴:**
- モノレポ対応（admin/user/liff共通）
- コンポーネントカタログ
- 型安全性保証
- テスト駆動開発

## アーキテクチャ原則

### レイヤー構造
```
├── pages/          # ページコンポーネント (ルーティング)
├── layouts/        # レイアウトテンプレート
├── components/     # 再利用可能コンポーネント
├── composables/    # Vue Composition 関数
├── stores/         # Pinia状態管理
├── middleware/     # ルートミドルウェア
├── plugins/        # プラグイン初期化
├── types/          # TypeScript型定義
└── constants/      # 定数・設定
```

### API連携アーキテクチャ
```
Vue Component -> Composable -> API Client -> Gateway
     |              |            |            |
   UI State      ビジネス      HTTP通信     認証・変換
                 ロジック
```

### 認証方式
- **Admin**: AWS Cognito JWT + Bearer Token
- **User**: Cookie Session + Bearer Token 両対応
- **LIFF**: LINE OAuth + セッション管理

## 状態管理戦略

### Piniaストア分割方針
```
stores/
├── auth.ts         # 認証状態
├── cart.ts         # カート状態
├── user.ts         # ユーザー情報
├── product.ts      # 商品データ
└── ui.ts           # UI状態
```

### 状態永続化
- **User App**: `@pinia-plugin-persistedstate` 使用
- **Admin App**: セッションベース (非永続)
- **LIFF App**: localStorage活用

## パフォーマンス最適化

### ビルド最適化
- **Code Splitting**: 自動ページ分割
- **Tree Shaking**: 不要コード除去
- **Bundle Analysis**: vite-bundle-analyzer

### 画像最適化
- **User App**: Nuxt Image + CloudFront連携
- **Admin App**: 手動最適化

### キャッシュ戦略
- **SSG**: 静的ページ事前生成
- **ISR**: 差分更新
- **Browser Cache**: 適切なCache-Control

## 開発・運用

### 開発環境
```bash
# 各アプリ個別起動
cd web/admin && yarn dev
cd web/user && yarn dev
cd web/liff && yarn dev

# 共通ライブラリ開発
cd web/shared && yarn storybook
```

### ビルド・デプロイ
- **Development**: Docker Compose
- **Production**: AWS S3 + CloudFront
- **CI/CD**: GitHub Actions

### 品質管理
- **ESLint**: コード品質
- **Prettier**: フォーマット
- **TypeScript**: 型安全性
- **Vitest**: 単体テスト
- **Stylelint**: CSS品質 (admin のみ)

## セキュリティ対策

### XSS対策
- Vue.js自動エスケープ
- Content Security Policy
- DOMPurify (リッチエディタ)

### 認証セキュリティ
- JWT検証
- CSRF対策
- Secure Cookie設定

### データ保護
- 機密情報の環境変数管理
- ログ出力制限
- エラー情報の適切な隠蔽