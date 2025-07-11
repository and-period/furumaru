# アーキテクチャ設計ルール

## 基本原則

### 1. マイクロサービスアーキテクチャ
- 各サービスは単一の責務を持つ
- サービス間の依存は最小限に保つ
- データベースはサービスごとに分離する

### 2. レイヤードアーキテクチャ
各サービス内部は以下の層で構成する：
- **API層** (`api/`): gRPCハンドラー、リクエスト/レスポンス処理
- **サービス層** (`service/`): ビジネスロジック
- **データベース層** (`database/`): データアクセス、リポジトリ
- **エンティティ層** (`entity/`): ドメインモデル、値オブジェクト

### 3. 通信パターン
- **外部 → ゲートウェイ**: REST API (HTTP/HTTPS)
- **ゲートウェイ → 内部サービス**: gRPC
- **内部サービス間**: gRPC
- **非同期処理**: AWS SQS

## サービス設計ルール

### 1. ゲートウェイの責務
- 認証・認可の処理
- リクエストの検証
- 複数サービスの集約
- レスポンスの整形

### 2. 内部サービスの責務
- ビジネスロジックの実装
- データの永続化
- ドメインルールの適用
- イベントの発行

### 3. エラーハンドリング
- 各層で適切なエラーハンドリングを実装
- エラーは構造化して上位層に伝播
- クライアントには適切なHTTPステータスコードを返す

## データ設計ルール

### 1. データベース分離
- サービスごとに独立したデータベースを持つ
- 他サービスのデータベースに直接アクセスしない
- 必要な場合はAPIを通じてデータを取得

### 2. トランザクション
- トランザクションはサービス内で完結させる
- 分散トランザクションは避ける
- 最終的整合性を許容する設計にする

### 3. イベント駆動
- 状態変更時は適切なイベントを発行
- イベントは冪等性を持つように設計
- リトライ処理を考慮した実装

## セキュリティルール

### 1. 認証・認可
- 認証はAWS Cognitoで統一
- 認可はゲートウェイレベルで実施
- 内部サービス間は信頼された通信とする

### 2. データ保護
- 個人情報は適切に暗号化
- ログに機密情報を出力しない
- 環境変数で機密情報を管理
