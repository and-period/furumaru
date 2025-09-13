# Furumaru アーキテクチャ概要

## システム構成

### バックエンドアーキテクチャ

**モジュラーモノリス + レイヤードアーキテクチャ**
```
api/ (単一アプリケーション)
├── cmd/                    # エントリポイント
├── internal/              # ビジネスモジュール
│   ├── gateway/           # APIゲートウェイ
│   ├── user/              # ユーザー管理（レイヤード）
│   ├── store/             # EC機能（レイヤード）
│   ├── media/             # メディア管理（レイヤード）
│   └── messenger/         # 通知機能（レイヤード）
└── pkg/                   # 共通パッケージ
```

**各モジュール構成（レイヤード）**:
- **api/**: プレゼンテーション層（HTTPハンドラー）
- **service/**: ビジネスロジック層
- **database/**: データアクセス層
- **entity/**: エンティティ・ドメインモデル

### フロントエンドアーキテクチャ

**マルチアプリケーション構成**
- **web/admin**: 管理者ポータル（Nuxt 4 + Vuetify）
- **web/user**: 購入者ポータル（Nuxt 3 + Tailwind）
- **web/liff**: LINEアプリ（Nuxt 3 + LINE LIFF）
- **web/shared**: 共通コンポーネント（Vue 3 + Vite）

## 認証・認可アーキテクチャ

### 認証プロバイダー
- **AWS Cognito**: JWTトークン発行
- **Google OAuth**: ソーシャルログイン
- **LINE OAuth**: ソーシャルログイン

### 認証フロー

#### Bearer Token Authentication
```
Client -> Gateway -> Internal Module
  |         |           |
  JWT       JWT         Function Call with UserID
```

#### Cookie Session Authentication  
```
Client -> Gateway -> Internal Module
  |         |           |
  Cookie    SessionID   Function Call with SessionID
```

### 認証スキーム適用パターン

| エンドポイントタイプ | 認証方式 | 理由 |
|---|---|---|
| ゲスト機能 | cookieauth | セッションベースで一時的 |
| カート機能 | cookieauth | ログイン前でも利用可能 |
| OAuth認証 | cookieauth | 認証プロセス中 |
| ユーザー管理 | bearerauth | 確実な本人確認が必要 |
| 決済処理 | 両方対応 | ゲスト・会員両方の利用 |

## データフロー

### リクエストフロー
```
Web Frontend -> API Gateway -> Internal Module -> Database
     |              |               |               |
   (HTTP/JSON)   (HTTP/JSON)    (Function Call)   (MySQL/TiDB)
```

**モノリス内通信**: 各モジュール間は関数呼び出しで連携（直接的なGoの関数呼び出し）

### 横断的機能
- **バリデーション**: Gateway層で実施
- **認証・認可**: Gateway層で検証
- **エラーハンドリング**: 統一されたErrorResponseフォーマット
- **ログ**: 各層でのリクエスト追跡

## ビジネスドメイン

### EC機能
- **Product**: 商品管理、在庫管理
- **Cart**: カート機能、セッション管理
- **Order**: 注文処理、決済連携  
- **Shipping**: 配送管理

### コンテンツ機能
- **Experience**: 体験商品
- **Live**: ライブ配信、ライブコマース
- **Video**: オンデマンド動画
- **Review**: レビュー・評価システム

### ユーザー管理
- **Member**: 購入者会員
- **Guest**: ゲストユーザー
- **Coordinator**: コーディネーター（販売者）
- **Producer**: 生産者
- **Admin**: 管理者

## 決済アーキテクチャ

### 決済プロバイダー
- **Komoju**: 主要決済サービス
- **Stripe**: クレジットカード決済
- 各種電子マネー・QR決済

### 決済フロー
```
Checkout Request -> Payment System Check -> Provider API -> Redirect URL
                        |                       |              |
                    メンテナンス確認           決済情報作成      ユーザーリダイレクト
```

### エラーハンドリング
- **403**: 決済システムメンテナンス
- **412**: 在庫不足、無効プロモーション等の前提条件エラー

## データベース設計

### マイクロサービス別DB
- 各サービスが独自のデータベース／スキーマを保持
- **user_db**: ユーザー、認証情報
- **store_db**: 商品、注文、決済
- **media_db**: メディアコンテンツ、コメント
- **messenger_db**: 通知、メッセージ

### 主要エンティティ関係

#### User系
```
User (基底) -> Member, Guest, Admin, Coordinator, Producer
              |
              Address (複数住所管理)
```

#### Store系  
```
Product -> ProductReview -> ProductReviewReaction
    |
    Cart -> Order -> OrderPayment
    |          |
    Experience -> ExperienceReview -> ExperienceReviewReaction
```

#### Media系
```
Video -> VideoComment
Live -> LiveComment  
Broadcast -> BroadcastComment
```

## 開発・運用パターン

### テストアプローチ
- **単体テスト**: テーブル駆動テスト、モック活用
- **統合テスト**: Docker Compose環境
- **カバレッジ**: サービス層重視

### エラーハンドリング統一
```go
// 共通エラーレスポンス
type ErrorResponse struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

### ログ・モニタリング
- **分散トレーシング**: リクエストID追跡
- **メトリクス収集**: Prometheus/Grafana想定
- **アラート**: 業務エラー（在庫不足等）とシステムエラーを区別

## セキュリティ考慮事項

### 機密情報管理
- AWS Secrets Manager活用
- 環境変数での設定分離
- コミット時の機密情報検査

### API セキュリティ
- CORS設定
- Rate Limiting
- Input Validation
- SQL Injection対策（ORMマッパー使用）

### 認証セキュリティ  
- JWT有効期限管理
- リフレッシュトークン
- セッション固定攻撃対策
