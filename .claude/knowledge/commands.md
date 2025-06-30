# コマンド一覧

## 開発環境セットアップ

### 初回セットアップ
```bash
# リポジトリクローン
git clone <repository-url>
cd furumaru

# 環境設定ファイル作成
cp .env.temp .env
# .envファイルを編集（AWS認証情報、データベース設定等）

# 初回セットアップ（コンテナビルド、依存関係インストール、Swagger生成）
make setup

# 全サービス起動
make start
```

### 環境確認
```bash
# コンテナ状態確認
docker-compose ps

# ログ確認
make logs

# 特定サービスのログ確認
docker-compose logs -f user_web
docker-compose logs -f admin_web
docker-compose logs -f mysql
```

## サービス制御

### 起動・停止
```bash
# 全サービス起動
make start

# 購入者関連サービスのみ起動
make start-user

# 管理者関連サービスのみ起動
make start-admin

# フロントエンドのみ起動
make start-web

# APIのみ起動
make start-api

# 全サービス停止
make stop

# 全サービス停止・削除
make down

# 全サービス停止・削除（データも含む）
make remove
```

### 個別サービス操作
```bash
# 特定サービスの再起動
docker-compose restart user_web
docker-compose restart admin_web

# 特定サービスの停止
docker-compose stop mysql
docker-compose stop user_gateway

# 特定サービスの起動
docker-compose start mysql
docker-compose start user_gateway
```

## データベース操作

### マイグレーション
```bash
# データベースマイグレーション実行
make migrate

# マイグレーション状態確認
cd api/hack/database-migrate-mysql
go run ./main.go -db-host=mysql -db-port=3306 -dry-run
```

### データベース接続
```bash
# MySQLコンテナに接続
docker-compose exec mysql mysql -u root -p

# データベース一覧確認
SHOW DATABASES;

# 特定データベース使用
USE furumaru_gateway;
USE furumaru_user;
USE furumaru_store;
USE furumaru_media;
USE furumaru_messenger;
```

### テストデータ投入
```bash
# テストデータ生成・投入
cd api/hack/database-seeds
go run ./main.go

# 商品レビューダミーデータ作成
cd api/hack/store-create-product-reviews
go run ./main.go

# 管理者ユーザー作成
cd api/hack/user-create-admin
go run ./main.go
```

## API開発

### APIサービス操作
```bash
cd api

# 依存ライブラリインストール
make install

# モック生成（インターフェース変更時）
make mockgen

# コードフォーマット
make fmt-fix

# Lint実行・修正
make lint-fix

# テスト実行
make test

# 特定のサービスをビルド
make build SERVICE=gateway/admin
make build SERVICE=gateway/user

# 開発モードでサービス起動（ホットリロード）
make start-dev SERVICE=gateway/admin
make start-dev SERVICE=gateway/user
```

### Swagger（API仕様書）
```bash
# Swagger仕様書生成
make swagger

# Swagger UI起動
make start-swagger

# Swagger確認
# ユーザー向けAPI: http://localhost:8080
# 管理者向けAPI: http://localhost:8081
```

## フロントエンド開発

### 管理者ポータル
```bash
cd web/admin

# 開発サーバー起動
yarn dev

# ビルド
yarn build

# 本番サーバー起動
yarn start

# 静的サイト生成
yarn generate

# 型チェック
yarn typecheck

# Lint実行
yarn lint

# Lint自動修正
yarn format

# Service Worker ビルド
yarn sw:build
```

### 購入者ポータル
```bash
cd web/user

# 開発サーバー起動（HTTPS付き）
yarn dev

# ビルド
yarn build

# 本番ビルドプレビュー
yarn preview

# 本番サーバー起動
yarn start

# 静的サイト生成
yarn generate

# テスト実行
yarn test

# カバレッジ付きテスト実行
yarn coverage

# 型チェック
yarn typecheck

# Lint実行
yarn lint

# Lint自動修正
yarn format
```

## デバッグ・トラブルシューティング

### コンテナ・ネットワーク確認
```bash
# コンテナ状態確認
docker-compose ps

# リソース使用量確認
docker stats

# ネットワーク確認
docker network ls
docker network inspect furumaru_default

# ボリューム確認
docker volume ls
```

### ログ確認
```bash
# 全サービスのログ
make logs

# 特定サービスのログ（リアルタイム）
docker-compose logs -f user_web
docker-compose logs -f admin_web
docker-compose logs -f user_gateway
docker-compose logs -f admin_gateway
docker-compose logs -f mysql

# ログの保存
docker-compose logs user_web > user_web.log
```

### ポート確認
```bash
# 使用中ポート確認
lsof -i :3000   # user_web
lsof -i :3010   # admin_web
lsof -i :18000  # user_gateway
lsof -i :18010  # admin_gateway
lsof -i :3306   # mysql

# プロセス確認
ps aux | grep node
ps aux | grep go
```

### データベースデバッグ
```bash
# データベース接続テスト
docker-compose exec mysql mysql -u root -p -e "SELECT 1"

# テーブル確認
docker-compose exec mysql mysql -u root -p furumaru_store -e "SHOW TABLES"

# クエリ実行
docker-compose exec mysql mysql -u root -p furumaru_store -e "SELECT COUNT(*) FROM products"
```

## API テスト

### cURLでのAPIテスト
```bash
# 商品一覧取得
curl -X GET http://localhost:18000/v1/products

# 商品詳細取得
curl -X GET http://localhost:18000/v1/products/{product_id}

# ヘルスチェック
curl -X GET http://localhost:18000/health
curl -X GET http://localhost:18010/health
```

### 認証付きAPIテスト（JWT必要）
```bash
# JWT取得後のAPIテスト例
TOKEN="your_jwt_token_here"

curl -X GET \
  http://localhost:18000/v1/orders \
  -H "Authorization: Bearer $TOKEN"

curl -X POST \
  http://localhost:18000/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "products": [
      {"product_id": "prod_123", "quantity": 2}
    ]
  }'
```

## メンテナンス・クリーンアップ

### Docker環境クリーンアップ
```bash
# 停止コンテナ削除
docker container prune

# 未使用イメージ削除
docker image prune

# 未使用ボリューム削除
docker volume prune

# システム全体クリーンアップ（注意：実行前に確認）
docker system prune -a
```

### ログファイルクリーンアップ
```bash
# Dockerログのクリーンアップ
docker system prune --volumes

# 開発環境完全リセット
make down
docker system prune -a --volumes
make setup
```
