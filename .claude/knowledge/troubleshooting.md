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
# - JWT_SECRET
```

## フロントエンドエラー

### 1. yarn devが起動しない

#### 症状
```bash
cd web/user
yarn dev
# Error: Cannot find module '@nuxt/xxx'
```

#### 原因と対処法
**原因**: 依存関係がインストールされていない

**対処法**:
```bash
# 依存関係の再インストール
yarn install

# node_modulesをクリアして再インストール
rm -rf node_modules yarn.lock
yarn install

# 全体のセットアップからやり直し
cd ../../
make install
```

### 2. 型エラー

#### 症状
```bash
yarn typecheck
# Type 'xxx' is not assignable to type 'yyy'
```

#### 原因と対処法
**原因**: TypeScript型定義の不整合

**対処法**:
```bash
# 型定義ファイルの確認
ls -la types/

# Nuxtの型生成
yarn dev --prepare

# 型チェック実行
yarn typecheck
```

### 3. API接続エラー

#### 症状
フロントエンドでAPIエラーが発生

#### 原因と対処法
**原因**: APIサーバーが起動していない、またはプロキシ設定の問題

**対処法**:
```bash
# APIサーバーの状態確認
curl http://localhost:18000/health
curl http://localhost:18010/health

# APIサーバーが停止している場合
make start-api

# ネットワーク接続確認
docker network ls
docker network inspect furumaru_default
```

## APIサーバーエラー

### 1. gRPCエラー

#### 症状
```bash
rpc error: code = Unavailable desc = connection error
```

#### 原因と対処法
**原因**: 内部サービス間の通信エラー

**対処法**:
```bash
# 内部サービスの状態確認
docker-compose ps

# 依存サービスの起動順序を確認
# gatewayサービスより先に内部サービスが起動している必要がある

# サービス再起動
docker-compose restart user_gateway
docker-compose restart admin_gateway
```

### 2. データベースマイグレーションエラー

#### 症状
```bash
make migrate
# Error: migration failed
```

#### 原因と対処法
**原因**: マイグレーションファイルに問題がある、または既存データとの整合性問題

**対処法**:
```bash
# マイグレーション状態確認
cd api/hack/database-migrate-mysql
go run ./main.go -db-host=mysql -db-port=3306 -dry-run

# データベースバックアップ
docker-compose exec mysql mysqldump -u root -p --all-databases > backup.sql

# 問題のあるマイグレーションをロールバック
# （プロジェクト固有の手順に従う）

# 段階的にマイグレーション実行
# 一度に全てではなく、一つずつ確認しながら実行
```

## パフォーマンス問題

### 1. 応答速度が遅い

#### 症状
APIの応答が遅い、またはフロントエンドの描画が遅い

#### 原因調査
```bash
# データベースの負荷確認
docker-compose exec mysql mysql -u root -p -e "SHOW PROCESSLIST"

# スロークエリの確認
docker-compose exec mysql mysql -u root -p -e "SHOW VARIABLES LIKE 'slow_query_log'"

# システムリソース確認
docker stats

# APIの応答時間測定
time curl -X GET http://localhost:18000/v1/products
```

#### 対処法
- データベースクエリの最適化
- インデックスの追加
- N+1クエリの解消（EagerLoadingの活用）
- 不要なデータ取得の削除

### 2. メモリ不足

#### 症状
```bash
docker stats
# メモリ使用率が高い
```

#### 対処法
```bash
# 不要なコンテナの停止
docker-compose stop <service_name>

# メモリ制限の設定（docker-compose.yml）
services:
  user_web:
    mem_limit: 512m
    memswap_limit: 512m
```

## 開発環境固有の問題

### 1. ホットリロードが効かない

#### 症状
コードを変更してもブラウザに反映されない

#### 対処法
```bash
# フロントエンド
cd web/user
yarn dev --host

# APIサーバー（開発モード）
cd api
make start-dev SERVICE=gateway/user
```

### 2. Docker Composeでファイル変更が反映されない

#### 症状
ローカルファイルの変更がコンテナ内に反映されない

#### 対処法
```bash
# ボリュームマウントの確認
docker-compose config

# コンテナ再起動
docker-compose restart <service_name>

# 完全な再ビルド
docker-compose down
docker-compose up --build
```

## よくある設定問題

### 1. AWS認証エラー

#### 症状
```bash
Error: NoCredentialsError: Unable to locate credentials
```

#### 対処法
```bash
# .envファイルの確認
grep AWS .env

# AWS認証情報の設定
AWS_ACCESS_KEY_ID=your_access_key
AWS_SECRET_ACCESS_KEY=your_secret_key
AWS_REGION=ap-northeast-1
```

### 2. CORS エラー

#### 症状
フロントエンドでCORSエラーが発生

#### 対処法
APIサーバーのCORS設定を確認し、必要に応じて許可オリジンを追加

## 緊急時の対応

### 1. 完全な環境リセット
```bash
# 全てのコンテナとデータを削除
make remove

# Dockerのクリーンアップ
docker system prune -a --volumes

# 再セットアップ
make setup
```

### 2. ログの収集
```bash
# 全サービスのログを収集
mkdir -p logs
docker-compose logs user_web > logs/user_web.log
docker-compose logs admin_web > logs/admin_web.log
docker-compose logs user_gateway > logs/user_gateway.log
docker-compose logs admin_gateway > logs/admin_gateway.log
docker-compose logs mysql > logs/mysql.log
```

### 3. 問題の報告
問題が解決しない場合、以下の情報を整理して報告：
- 発生した症状
- 実行したコマンド
- エラーメッセージ
- 環境情報（OS、Dockerバージョンなど）
- ログファイル