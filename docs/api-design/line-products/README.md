# LINE向け商品一覧API設計書

## 概要

本ドキュメントは、ふるマル（Furumaru）のLINEアプリ向け商品一覧取得APIの設計仕様を定義します。

### 目的

- LINEアプリ内でふるマルの商品一覧を表示できるようにする
- 既存のWebサイト向けAPIとは別に、LINE特有の要件に対応したAPIを提供する
- モバイル環境に最適化されたレスポンスを返す

### 対象システム

- **クライアント**: LINE公式アカウント / LINEミニアプリ
- **サーバー**: ふるマルAPIゲートウェイ（gateway/user）

## API仕様

### エンドポイント

```
GET /v1/line/products
```

### 認証

- 認証不要（公開API）
- 誰でもアクセス可能

### リクエストパラメータ

| パラメータ | 型 | 必須 | デフォルト | 説明 |
|-----------|-----|------|------------|------|
| limit | integer | No | 10 | 取得件数（最大20件） |
| offset | integer | No | 0 | 取得開始位置 |
| category_id | string | No | - | カテゴリーID |
| producer_id | string | No | - | 生産者ID |
| prefecture_code | integer | No | - | 都道府県コード |
| sort | string | No | recommended | ソート順（recommended, new, price_asc, price_desc） |

### レスポンス

#### 成功時（200 OK）

```json
{
  "products": [
    {
      "id": "product_123456",
      "name": "新潟県産コシヒカリ",
      "description": "新潟県の豊かな土壌で育った...",
      "thumbnail_url": "https://...",
      "price": 3000,
      "price_text": "¥3,000",
      "producer": {
        "id": "producer_123",
        "name": "田中農園"
      },
      "prefecture": "新潟県",
      "tags": ["お米", "新潟県産"],
      "is_limited": true,
      "is_out_of_stock": false,
      "line_url": "https://liff.line.me/..."
    }
  ],
  "pagination": {
    "total": 150,
    "limit": 10,
    "offset": 0,
    "has_next": true
  },
  "display_settings": {
    "show_producer": true,
    "show_prefecture": true,
    "show_tags": true
  }
}
```

#### エラー時

```json
{
  "error": {
    "code": "INVALID_PARAMETER",
    "message": "limitは1以上20以下である必要があります",
    "details": {
      "field": "limit",
      "value": 50
    }
  }
}
```

### レスポンスフィールド説明

#### product オブジェクト

| フィールド | 型 | 説明 |
|-----------|-----|------|
| id | string | 商品ID |
| name | string | 商品名（最大40文字） |
| description | string | 商品説明（最大100文字） |
| thumbnail_url | string | サムネイル画像URL（正方形推奨） |
| price | integer | 価格（税込） |
| price_text | string | 表示用価格テキスト |
| producer | object | 生産者情報 |
| prefecture | string | 都道府県名 |
| tags | array[string] | 商品タグ（最大3つ） |
| is_limited | boolean | 限定商品フラグ |
| is_out_of_stock | boolean | 在庫切れフラグ |
| line_url | string | LINE内商品詳細URL |

### LINE特有の考慮事項

1. **レスポンスサイズの最適化**
   - 商品説明文は100文字以内に要約
   - 画像はサムネイルURLのみ返却
   - 不要なフィールドは除外

2. **表示最適化**
   - Flex Messageでの表示を考慮したデータ構造
   - カルーセル表示用に最大10件を推奨
   - 商品名は全角20文字以内を推奨

3. **パフォーマンス**
   - レスポンスタイム目標: 1秒以内
   - キャッシュ: 5分間（CDN層）

4. **LINE連携機能**
   - 商品詳細はLIFFアプリで表示
   - 購入はLINE Pay決済に対応予定
   - お気に入り機能はLINEアカウントと連携

## 実装方針

### アーキテクチャ

```
LINE App → LINE Platform → API Gateway → Store Service → Database
                               ↓
                          Cache Layer
```

### 実装場所

- **ハンドラー**: `/api/internal/gateway/user/v1/handler/line_product.go`
- **サービス**: `/api/internal/gateway/user/v1/service/line_product.go`
- **レスポンス**: `/api/internal/gateway/user/v1/response/line_product.go`

### 既存APIとの差分

| 項目 | 既存API | LINE API |
|------|---------|----------|
| エンドポイント | /v1/products | /v1/line/products |
| デフォルト取得件数 | 20件 | 10件 |
| 最大取得件数 | 200件 | 20件 |
| 商品説明文 | 全文 | 100文字要約 |
| 画像 | 複数URL | サムネイルのみ |
| 関連データ | 全て含む | 最小限 |
| レスポンスサイズ | 大 | 小（モバイル最適化） |

### セキュリティ

1. **アクセス制御**
   - CORS設定で適切なオリジンを許可
   - 必要に応じてRefererチェック

2. **データ保護**
   - 個人情報は含まない
   - HTTPS必須

## 今後の拡張

1. **パーソナライゼーション**
   - ユーザーの購入履歴に基づく推奨
   - 地域に基づく商品表示

2. **リッチコンテンツ**
   - 動画サムネイル対応
   - ライブコマース連携

3. **通知連携**
   - 在庫復活通知
   - セール情報のプッシュ通知

## 関連ドキュメント

- [API仕様書（OpenAPI）](./openapi.yaml)
- [実装ガイドライン](./implementation-guide.md)
- [テスト仕様書](./test-specification.md)