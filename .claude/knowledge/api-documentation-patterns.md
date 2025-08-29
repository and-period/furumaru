# API Documentation Patterns - Furumaru

## Swagger/OpenAPI ドキュメント生成アーキテクチャ

### 1. ドキュメント生成の仕組み
- **自動生成**: [swaggo/swag](https://github.com/swaggo/swag)を使用してGoコメントからOpenAPI仕様を生成
- **手動作成**: `docs/swagger/user/`に手動で作成されたOpenAPI仕様が存在
- **統合版**: `/tmp/data/swagger/user/docs/openapi/openapi.yaml`に統合されたドキュメント
- **生成コマンド**: `make swagger`でコメントから自動生成

### 2. コメント記述パターン

#### 基本構造
```go
// @tag.name        TagName
// @tag.description タグの説明
func (h *handler) routeGroup(rg *gin.RouterGroup) {
    // ルート定義
}

// @Summary     エンドポイントの概要
// @Description 詳細な説明文
// @Tags        TagName
// @Router      /path/{param} [method]
// @Security    auth_type
// @Param       param path/query/body type required/optional "説明"
// @Accept      json
// @Produce     json
// @Success     200 {object} ResponseType "成功時の説明"
// @Failure     400 {object} util.ErrorResponse "エラーの説明"
func (h *handler) HandlerFunction(ctx *gin.Context) {
    // 実装
}
```

#### 日本語パターン
- **Summary**: 簡潔な機能説明（例：「商品レビュー作成」）
- **Description**: 詳細説明（例：「商品のレビューを作成します。」）
- **Tags**: 機能グループ（例：「ProductReview」）
- **失敗レスポンス**: 日本語での説明（例：「バリデーションエラー」）

### 3. 認証スキーム

#### 利用パターン
- **cookieauth**: セッションベース認証（ゲスト、カート機能）
- **bearerauth**: JWTトークン認証（ユーザー専用機能）
- **複数認証**: 両方対応の場合は複数の`@Security`行を記述

```go
// 単一認証
// @Security    cookieauth

// 複数認証対応
// @Security    bearerauth
// @Security    cookieauth
```

#### 重要な命名規則
- 手動docs: `cookieAuth`, `bearerAuth` (キャメルケース)  
- swaggo: `cookieauth`, `bearerauth` (小文字) ← **これが正しい**

### 4. HTTPステータスコードの体系

#### 成功レスポンス
- **200**: データ取得、情報返却
- **204**: 作成、更新、削除（レスポンスボディなし）

#### エラーレスポンス
- **400**: バリデーションエラー（必須）
- **401**: 認証エラー（認証必須エンドポイント）
- **403**: アクセス権限エラー、サービス停止
  - 決済システムメンテナンス
  - 店舗利用停止
  - 商品非公開状態
- **404**: リソースが存在しない
- **409**: 重複・競合エラー
  - 既存メールアドレス
  - 重複データ作成
- **412**: 前提条件エラー（ビジネスルール違反）
  - 在庫不足
  - 無効なプロモーション  
  - 同一データ更新
  - 配信未公開状態

### 5. パラメータ記述パターン

```go
// パスパラメータ
// @Param       userId path string true "ユーザーID"

// クエリパラメータ
// @Param       limit query int64 false "取得件数" default(20)

// リクエストボディ  
// @Param       request body request.CreateUserRequest true "ユーザー作成"
```

### 6. ドキュメントファイル構造

```
docs/swagger/user/
├── openapi.yaml              # メインファイル
├── components/
│   ├── schemas/             # データ型定義
│   └── securitySchemes/     # 認証スキーム
└── paths/                   # エンドポイント定義
    └── v1/
        ├── auth/
        ├── users/
        ├── products/
        └── ...
```

### 7. エラーハンドリングの対応関係

| ビジネスケース | HTTPコード | 説明テンプレート |
|---|---|---|
| バリデーション失敗 | 400 | "バリデーションエラー" |
| 未認証 | 401 | "認証エラー" |  
| システムメンテナンス | 403 | "決済システムがメンテナンス中" |
| サービス停止 | 403 | "店舗が利用停止中" |
| リソース不存在 | 404 | "{リソース名}が存在しない" |
| データ重複 | 409 | "すでに存在する{項目名}" |
| 在庫不足 | 412 | "商品在庫が不足している" |
| 無効状態 | 412 | "{リソース名}が公開されていない" |

### 8. 開発フロー

1. **ハンドラ実装**: 機能実装とエラーハンドリング
2. **Swaggerコメント追加**: 上記パターンに従って記述
3. **手動docsとの整合性確認**: 特にエラーコード
4. **swagger生成**: `make swagger`で自動生成
5. **統合確認**: 生成された仕様の検証

### 9. 品質保証のチェックポイント

- [ ] 全エンドポイントにSwaggerコメントが存在
- [ ] 認証スキームが正しい（小文字）
- [ ] エラーコードが網羅的（特に403, 409, 412）
- [ ] レスポンスコードが実装と一致（200 vs 204）
- [ ] 日本語説明が統一的
- [ ] パラメータの必須/任意が正確

### 10. 参考ファイル

- **パターン参考**: `api/internal/gateway/user/facility/handler/auth.go`
- **手動docs**: `docs/swagger/user/paths/`
- **統合版**: `/tmp/data/swagger/user/docs/openapi/openapi.yaml`
- **既存実装**: `api/internal/gateway/user/v1/handler/`の全ファイル