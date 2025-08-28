# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Claude Code 作業ログ・知識管理ルール

### 作業ログの記録
Claude Code での作業時は以下のディレクトリに適切に情報を記録すること：

1. **`.claude/debug/`** - 作業過程とデバッグ情報
   - 作業セッションの詳細ログ
   - エラー発生時の調査・解決過程
   - 試行錯誤の記録
   - ファイル命名: `YYYY-MM-DD_session.md`, `debug_YYYYMMDD_issue-name.md`

2. **`.claude/knowledge/`** - ドメイン知識の蓄積
   - アーキテクチャの理解・分析結果
   - ビジネスロジックの解析
   - エンティティ・API仕様の整理
   - 設計パターンやベストプラクティス
   - ファイル命名: `component-name.md`, `analysis_YYYYMMDD_topic.md`

3. **`docs/`** - 正式なプロジェクトドキュメント
   - 外部向けドキュメント
   - 運用マニュアル
   - アーキテクチャドキュメント

### /load-rules コマンド
`/load-rules` コマンドが実行された場合、以下のディレクトリから学習すること：
- `docs/` ディレクトリ配下の全てのマークダウンファイル
- `.claude/knowledge/` ディレクトリ配下の知識ファイル
- プロジェクトの最新の理解を更新し、作業に反映させる

### 記録のタイミング
- **作業開始時**: セッションログを `.claude/debug/` に作成
- **新しい理解を得た時**: 即座に `.claude/knowledge/` に記録
- **問題発生・解決時**: デバッグ過程を `.claude/debug/` に記録
- **作業完了時**: 重要な知識を整理し適切なディレクトリに保存

## プロジェクト概要

**Furumaru (全国ふるさとマルシェ「ふるマル」)** は、地域・地方の特産品を扱う日本のECマーケットプレイスプラットフォームで、ライブコマース機能を備えています。複数のマイクロサービスとフロントエンドアプリケーションを含むモノレポ構成です。

## アーキテクチャ

### バックエンドサービス (`/api`)
- **言語**: Go 1.24.2
- **アーキテクチャ**: ドメイン駆動設計によるマイクロサービス
- **サービス構成**:
  - `gateway/admin`: 管理者向けAPIゲートウェイ
  - `gateway/user`: 購入者向けAPIゲートウェイ
  - `media`: 動画ストリーミング、配信、メディアコンテンツ
  - `messenger`: 通知とメッセージング
  - `store`: ECコア機能（商品、注文、決済、配送）
  - `user`: ユーザー管理と認証
- **通信方式**: サービス間はgRPC、外部向けはREST API
- **データベース**: MySQL/TiDB（ドメインごとに分離）
- **認証**: AWS Cognito
- **決済**: KomojuとStripe

### フロントエンドアプリケーション (`/web`)
- **管理者ポータル** (`/web/admin`): Nuxt 4 + Vue 3 + Vuetify (ポート3010)
- **購入者ポータル** (`/web/user`): Nuxt 3 + Vue 3 + Tailwind CSS (ポート3000)
- **状態管理**: Pinia
- **リアルタイム機能**: HLS.jsによる動画ストリーミング

### インフラストラクチャ
- **コンテナ**: DockerとDocker Compose
- **クラウド**: AWS (S3、Cognito、MediaLive、SQS、Lambda)
- **APIドキュメント**: `/docs/swagger`にOpenAPI/Swagger仕様

## 開発コマンド

### 初期セットアップ
```bash
make setup          # 初回セットアップ：コンテナビルド、依存関係インストール、Swagger生成
cp .env.temp .env   # 環境ファイルを作成し、AWS認証情報を設定
```

### サービスの起動
```bash
make start          # 全サービスを起動
make start-user     # 購入者向けサービスのみ起動
make start-admin    # 管理者向けサービスのみ起動
make migrate        # データベースマイグレーションを実行
```

### API開発 (`/api`)
```bash
cd api
make mockgen        # インターフェース変更後にモックを生成
make test           # カバレッジ付きで全テストを実行
make lint-fix       # Lintエラーを修正
make fmt-fix        # フォーマットを修正
make start-dev SERVICE=gateway/admin  # ホットリロード付きでサービスを起動
```

### フロントエンド開発
```bash
# 管理者ポータル
cd web/admin
yarn dev            # 開発サーバーを起動
yarn typecheck      # 型チェックを実行
yarn lint           # Lintを実行
yarn format         # Lintエラーを自動修正

# 購入者ポータル
cd web/user
yarn dev            # HTTPS付き開発サーバーを起動
yarn test           # Vitestテストを実行
yarn coverage       # カバレッジ付きでテストを実行
yarn typecheck      # 型チェックを実行
yarn lint           # Lintを実行
yarn format         # Lintエラーを自動修正
```

### テスト
```bash
# APIテスト
cd api && make test

# フロントエンドテスト
cd web/user && yarn test
cd web/admin && yarn test  # テストが存在する場合
```

### よく使うタスク
```bash
make swagger        # APIドキュメントを再生成
make logs           # コンテナログを表示
make down           # 全コンテナを停止・削除
```

## 高レベルアーキテクチャパターン

### サービス通信フロー
1. **フロントエンド** → **ゲートウェイ** → **内部サービス** → **データベース**
2. ゲートウェイが認証、リクエスト検証、ルーティングを処理
3. 内部サービスがビジネスロジックとデータベース操作を担当
4. サービス間はProtocol Buffersを使用したgRPCで通信

### ドメイン構造
各サービスは以下のパターンに従います：
```
/api/internal/{service}/
├── database/          # データベースモデルとクエリ
├── entity/            # ドメインエンティティと値オブジェクト
├── service/           # ビジネスロジック層
├── api/               # gRPCハンドラー
└── *.proto            # サービス定義
```

### フロントエンド構造
両方のWebアプリはNuxt 3の規約に従います：
```
/web/{app}/
├── components/        # Vueコンポーネント
├── composables/       # Composition APIユーティリティ
├── pages/             # ルートベースのページ
├── stores/            # Piniaストア
├── middleware/        # ルートミドルウェア
└── plugins/           # Vueプラグイン
```

### 重要なアーキテクチャ決定事項

1. **マルチデータベース戦略**: 各ドメインが独自のデータベースを持ち、サービスの独立性を確保
2. **ゲートウェイパターン**: セキュリティと一貫性のため、全ての外部リクエストはゲートウェイを経由
3. **イベント駆動**: 非同期処理（注文処理、通知）にAWS SQSを使用
4. **CQRS的な分離**: 最適化のため、読み取りモデルと書き込みモデルが異なることが多い
5. **エンティティファースト設計**: リッチなエンティティオブジェクトによる強力なドメインモデリング

### 重要な規約

1. **エラーハンドリング**: 全サービスで標準化されたエラーコードとメッセージを使用
2. **認証**: AWS CognitoからのJWTトークン、ゲートウェイレベルで検証
3. **テスト**: 実装ファイルと並んで単体テストを配置（例：`service.go` → `service_test.go`）
4. **モック生成**: インターフェースにはテスト用の対応するモックが存在
5. **APIバージョニング**: 現在v1、URLパスでのバージョニングの可能性あり

### よくある注意点

1. **データベースマイグレーション**: 変更をpullした後は必ず`make migrate`を実行
2. **環境変数**: 多くの機能で`.env`にAWS認証情報が必要
3. **ポートの競合**: ポート3000、3010、18000、18010、3306が利用可能であることを確認
4. **サービスの依存関係**: 一部のサービスは他のサービスの起動が必要（docker-compose.yamlを確認）
5. **モックの再生成**: インターフェース変更後は`make mockgen`を実行

### Git/PR運用ルール

1. **mainブランチ保護**: mainブランチには直接pushしない
2. **フィーチャーブランチ**: 新機能や修正は必ずフィーチャーブランチで作業
3. **PR必須**: mainブランチへの変更は必ずPull Requestを経由
4. **自動化**: `/push`コマンド実行時は自動的に新しいブランチを作成し、PRを作成

### テストアプローチ

- **API**: モック化された依存関係を使用したテーブル駆動テスト
- **フロントエンド**: Vitestによるコンポーネントテスト
- **統合テスト**: Docker Compose環境が本番環境を模倣
- **カバレッジ**: 特にサービス層で高いカバレッジを目指す

### デバッグに便利なコマンド

```bash
# サービスログの確認
docker-compose logs -f {service_name}

# データベースアクセス
docker-compose exec mysql mysql -u root -p

# curlでAPIテスト
curl -X GET http://localhost:18000/v1/products

# テストデータの生成
cd api/hack/database-seeds && go run ./main.go
```
