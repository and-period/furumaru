# 技術スタックルール

## バックエンド

### 言語・フレームワーク
- **言語**: Go 1.26
- **HTTPフレームワーク**: Gin
- **内部通信**: 直接関数呼び出し
- **ORマッパー**: GORM

### データベース
- **メインDB**: MySQL 8.0 / TiDB（本番環境）
- **キャッシュ**: 現在は未使用（必要に応じてRedis導入を検討）
- **マイグレーション**: 自作ツール（`api/hack/database-migrate-mysql`）

### メッセージング・キュー
- **キュー**: AWS SQS
- **通知**: Firebase Cloud Messaging (FCM)

### 認証・認可
- **認証基盤**: AWS Cognito
- **トークン形式**: JWT
- **認可**: RBAC（Role-Based Access Control）

## フロントエンド

### フレームワーク
- **フレームワーク**: 
  - 管理者ポータル: Nuxt 4
  - 購入者ポータル: Nuxt 3
- **UIライブラリ**: Vue 3（Composition API）
- **状態管理**: Pinia

### スタイリング
- **管理者ポータル**: Vuetify 3（Material Design）
- **購入者ポータル**: Tailwind CSS

### ビルドツール
- **パッケージマネージャー**: Yarn
- **ビルドツール**: Vite（Nuxt内蔵）

## インフラストラクチャ

### コンテナ
- **コンテナ**: Docker
- **オーケストレーション（開発）**: Docker Compose
- **オーケストレーション（本番）**: AWS ECS/Fargate（推測）

### AWS サービス
- **ストレージ**: S3
- **動画配信**: MediaLive, MediaConvert
- **サーバーレス**: Lambda
- **API Gateway**: 未使用（独自ゲートウェイ実装）

### モニタリング・ログ
- **エラー監視**: Sentry
- **APM**: New Relic
- **ログ**: CloudWatch Logs（推測）

## 開発ツール

### コード品質
- **Linter（Go）**: golangci-lint
- **Linter（JS/TS）**: ESLint
- **フォーマッター（Go）**: gofmt, goimports
- **フォーマッター（JS/TS）**: Prettier

### テスト
- **Go**: testing標準パッケージ + testify
- **モック**: go.uber.org/mock (gomock)
- **DBテスト**: testcontainers-go
- **フロントエンド**: Vitest

### API仕様
- **仕様記述**: OpenAPI 3.0 (Swagger)
- **コード生成**: oapi-codegen

## バージョン管理ルール

### 依存関係の更新
- セキュリティアップデートは即座に適用
- メジャーバージョンアップは慎重に検討
- 依存関係は定期的に見直し

### サポートバージョン
- Go: 最新の2バージョンをサポート
- Node.js: LTS版を使用
- データベース: 最新の安定版を使用

## 技術選定の基準

1. **実績と安定性**: 本番環境での実績がある
2. **コミュニティ**: 活発なコミュニティが存在
3. **保守性**: 長期的な保守が可能
4. **パフォーマンス**: 要求性能を満たす
5. **セキュリティ**: セキュリティアップデートが提供される