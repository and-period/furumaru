# トラブルシューティング

## 起動・接続エラー

### 1. コンテナが起動しない

#### 症状
```bash
make start
# エラー: Cannot start service xxx: port is already allocated
```

#### 原因と対処法
**原因**: ポートが既に使用されている

**対処法**:
```bash
# 使用中のポートを確認
lsof -i :3000   # user_web
lsof -i :3010   # admin_web
lsof -i :18000  # user_gateway
lsof -i :18010  # admin_gateway
lsof -i :3306   # mysql

# プロセス停止
kill -9 <PID>

# または、既存のコンテナを停止
docker-compose down
make start
```

### 2. データベース接続エラー

#### 症状
```bash
Error: dial tcp 127.0.0.1:3306: connect: connection refused
```

#### 原因と対処法
**原因**: MySQLコンテナが起動していない、または接続設定が間違っている

**対処法**:
```bash
# MySQLコンテナの状態確認
docker-compose ps mysql

# MySQLコンテナが停止している場合
docker-compose start mysql

# ログでエラー詳細確認
docker-compose logs mysql

# 接続テスト
docker-compose exec mysql mysql -u root -p
```

### 3. 環境変数エラー

#### 症状
```bash
Error: required environment variable XXX is not set
```

#### 原因と対処法
**原因**: 必要な環境変数が設定されていない

**対処法**:
```bash
# .envファイルの存在確認
ls -la .env

# .envファイルが存在しない場合
cp .env.temp .env

# .envファイルを編集して必要な値を設定
# 特に以下の項目を確認：
# - AWS_ACCESS_KEY_ID
# - AWS_SECRET_ACCESS_KEY
# - DATABASE_URL
```

## ビルド・デプロイエラー

### 1. Dockerイメージビルド失敗

#### 症状
```bash
docker build failed: no space left on device
```

#### 対処法
```bash
# Docker環境のクリーンアップ
docker system prune -a --volumes

# 不要なイメージを削除
docker images -q --filter "dangling=true" | xargs docker rmi

# 再ビルド
make build
```

### 2. マイグレーション失敗

#### 症状
```bash
Error: migration failed: duplicate column name
```

#### 対処法
```bash
# マイグレーション状態確認
cd api/hack/database-migrate-mysql
go run ./main.go -db-host=mysql -db-port=3306 -dry-run

# 問題のあるマイグレーションをスキップ
# またはデータベースをリセット（開発環境のみ）
docker-compose down
docker volume rm furumaru_mysql_data
make migrate
```

## パフォーマンス問題

### 1. APIレスポンスが遅い

#### 調査方法
```bash
# スロークエリログの確認
docker-compose exec mysql mysql -u root -p -e "SHOW VARIABLES LIKE 'slow_query_log%'"

# プロセスリストの確認
docker-compose exec mysql mysql -u root -p -e "SHOW PROCESSLIST"

# インデックスの確認
docker-compose exec mysql mysql -u root -p furumaru_store -e "SHOW INDEX FROM products"
```

#### 改善策
- N+1クエリの解消（Preloadの使用）
- 適切なインデックスの追加
- クエリの最適化
- キャッシュの導入検討

### 2. メモリ不足

#### 症状
```bash
fatal error: runtime: out of memory
```

#### 対処法
```bash
# リソース使用状況確認
docker stats

# コンテナのメモリ制限を増やす
# docker-compose.ymlで設定：
services:
  service_name:
    mem_limit: 2g
    mem_reservation: 1g
```

## 開発環境固有の問題

### 1. ホットリロードが効かない

#### 原因と対処法
**原因**: ファイルシステムの監視が機能していない

**対処法**:
```bash
# 開発サーバーを再起動
cd web/user
yarn dev

# またはDockerコンテナを再起動
docker-compose restart user_web
```

### 2. CORS エラー

#### 症状
```
Access to XMLHttpRequest has been blocked by CORS policy
```

#### 対処法
- 開発環境の設定確認（nuxt.config.tsのproxy設定）
- APIゲートウェイのCORS設定確認
- 正しいAPIエンドポイントを使用しているか確認

## デバッグTips

### ログの確認方法
```bash
# 全サービスのログ
make logs

# 特定サービスのログ（リアルタイム）
docker-compose logs -f user_gateway

# ログをファイルに保存
docker-compose logs user_web > debug.log 2>&1
```

### データベースの直接確認
```bash
# データベース接続
docker-compose exec mysql mysql -u root -p

# テーブル構造確認
USE furumaru_store;
DESCRIBE products;

# データ確認
SELECT * FROM products LIMIT 10;
```

### APIのデバッグ
```bash
# curlでAPIテスト（詳細表示）
curl -v -X GET http://localhost:18000/v1/products

# レスポンスヘッダーの確認
curl -I http://localhost:18000/health

# JSONレスポンスの整形
curl -s http://localhost:18000/v1/products | jq .
```