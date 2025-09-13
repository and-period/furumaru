# Web フロントエンドアーキテクチャ

## 📋 目次

- **[frontend-applications.md](./frontend-applications.md)** - フロントエンドアプリケーション詳細仕様
- **[components.md](./components.md)** - コンポーネント設計・分類・テスト戦略
- **[state-management.md](./state-management.md)** - Pinia状態管理・永続化・ストア連携
- **[api-integration.md](./api-integration.md)** - API連携・認証・エラーハンドリング

## 概要

Furumaruプロジェクトは、4つの独立したWebアプリケーションから構成されています：

- **admin**: 管理者ポータル（Nuxt 4 + Vuetify）
- **user**: 購入者ポータル（Nuxt 3 + Tailwind）
- **liff**: LINEアプリ（Nuxt 3 + LINE LIFF）
- **shared**: 共通コンポーネント（Vue 3 + Vite）

詳細な技術仕様とアーキテクチャについては、各リンク先のドキュメントを参照してください。