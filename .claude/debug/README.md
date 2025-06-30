# デバッグ用README

このディレクトリには、Furumaruプロジェクトのデバッグに役立つファイルやスクリプトを配置します。

## デバッグ手順

### 1. 基本的なデバッグフロー

#### 問題の特定
1. **症状の確認**: どのような問題が発生しているか
2. **ログの確認**: エラーメッセージやスタックトレースの確認
3. **環境の確認**: 関連するサービスやコンテナの状態確認
4. **再現手順の特定**: 問題を再現できる最小限の手順

#### ログ収集スクリプト
```bash
#!/bin/bash
# logs-collect.sh - 全サービスのログを収集

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
LOG_DIR="debug_logs_$TIMESTAMP"

mkdir -p $LOG_DIR

echo "Collecting logs to $LOG_DIR/"

# Docker Compose logs
docker-compose logs --no-color user_web > $LOG_DIR/user_web.log
docker-compose logs --no-color admin_web > $LOG_DIR/admin_web.log
docker-compose logs --no-color user_gateway > $LOG_DIR/user_gateway.log
docker-compose logs --no-color admin_gateway > $LOG_DIR/admin_gateway.log
docker-compose logs --no-color mysql > $LOG_DIR/mysql.log

# System info
docker-compose ps > $LOG_DIR/containers_status.log
docker stats --no-stream > $LOG_DIR/containers_stats.log
docker network ls > $LOG_DIR/networks.log

echo "Logs collected in $LOG_DIR/"
```

### 2. データベースデバッグ

#### 接続テスト
```bash
# MySQL接続テスト
docker-compose exec mysql mysql -u root -proot -e "SELECT 'Connection OK' as status"

# データベース一覧
docker-compose exec mysql mysql -u root -proot -e "SHOW DATABASES"

# テーブル状態確認
docker-compose exec mysql mysql -u root -proot furumaru_store -e "SHOW TABLE STATUS"
```

#### よく使うクエリ
```sql
-- ユーザー数確認
SELECT COUNT(*) as user_count FROM users;

-- 商品数確認
SELECT COUNT(*) as product_count FROM products;

-- 注文状態別件数
SELECT status, COUNT(*) as count FROM orders GROUP BY status;

-- 最新の注文
SELECT * FROM orders ORDER BY created_at DESC LIMIT 10;

-- エラーログ確認（ある場合）
SELECT * FROM error_logs ORDER BY created_at DESC LIMIT 100;
```

### 3. APIデバッグ

#### ヘルスチェック
```bash
# 基本的なヘルスチェック
curl -v http://localhost:18000/health
curl -v http://localhost:18010/health

# レスポンス時間測定
time curl -s http://localhost:18000/v1/products > /dev/null
```

#### API詳細デバッグ
```bash
# 詳細なリクエスト・レスポンス確認
curl -v -X GET \
  http://localhost:18000/v1/products \
  -H "Accept: application/json" \
  -H "User-Agent: Debug/1.0"

# 認証付きAPIテスト
TOKEN="your_jwt_token"
curl -v -X GET \
  http://localhost:18000/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Accept: application/json"
```

### 4. フロントエンドデバッグ

#### ブラウザでのデバッグ
- **開発者ツール**: ネットワークタブでAPI通信を確認
- **コンソール**: JavaScriptエラーの確認
- **Vue DevTools**: Vue.jsコンポーネントの状態確認

#### サーバーサイドデバッグ
```bash
# Nuxt開発モードでの詳細ログ
cd web/user
DEBUG=nuxt:* yarn dev

# 型チェック
yarn typecheck

# ビルドエラーの確認
yarn build
```

### 5. 性能デバッグ

#### レスポンス時間測定
```bash
# APIエンドポイントの性能測定
ab -n 100 -c 10 http://localhost:18000/v1/products

# または、より詳細な測定
wrk -t12 -c400 -d30s http://localhost:18000/v1/products
```

#### データベースの性能分析
```sql
-- スロークエリの確認
SHOW VARIABLES LIKE 'slow_query_log%';
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;

-- プロセス一覧
SHOW PROCESSLIST;

-- インデックスの使用状況
SHOW INDEX FROM products;
EXPLAIN SELECT * FROM products WHERE category_id = 1;
```

## デバッグツール

### 1. ログ解析
- **grep**: 特定のエラーの検索
- **awk**: ログの集計・分析
- **jq**: JSON形式のレスポンス分析

### 2. ネットワーク解析
- **curl**: HTTP/HTTPS リクエストテスト
- **wget**: ファイルダウンロードテスト
- **nc**: ポート接続テスト

### 3. コンテナ・システム
- **docker logs**: コンテナログの確認
- **docker exec**: コンテナ内でのコマンド実行
- **docker stats**: リソース使用状況の監視

## トラブルシューティングチェックリスト

### 環境確認
- [ ] 必要なポートが開いている（3000, 3010, 18000, 18010, 3306）
- [ ] 環境変数が正しく設定されている（.env ファイル）
- [ ] Dockerコンテナが全て起動している
- [ ] データベースマイグレーションが完了している

### ログ確認
- [ ] エラーメッセージの確認
- [ ] スタックトレースの分析
- [ ] 関連するサービスのログ確認

### 設定確認
- [ ] データベース接続設定
- [ ] AWS認証設定
- [ ] CORS設定（フロントエンド）
- [ ] プロキシ設定（開発環境）

### 依存関係確認
- [ ] 依存サービスの起動順序
- [ ] ライブラリバージョンの整合性
- [ ] 外部サービスの可用性（AWS等）

## 緊急時対応

### 1. サービス全停止・再起動
```bash
make down
make start
```

### 2. データベースリセット
```bash
# 注意: 開発環境のみ実施
make remove  # データも削除される
make setup   # 再セットアップ
```

### 3. ログ・状態の保存
```bash
# 問題調査用にログと状態を保存
mkdir emergency_debug_$(date +%Y%m%d_%H%M%S)
cd emergency_debug_*

# 各種ログを収集
docker-compose logs > all_services.log
docker ps -a > containers.log
docker images > images.log
docker volume ls > volumes.log
docker network ls > networks.log

# 設定ファイルをコピー
cp ../.env env_backup
cp ../docker-compose.yml compose_backup.yml
```

これらの手順を参考に、問題の特定と解決を行ってください。
