# 仕様書・API定義

このディレクトリには Furumaru プロジェクトの仕様書と API 定義を格納します。

## ディレクトリ構成

```
spec/
├── external/              # 外部API仕様
│   └── ...               # 外部サービスとの連携仕様
├── internal/             # 内部設計仕様
│   └── ...               # ふるマル内部設計・仕様書
└── swagger/              # OpenAPI/Swagger定義
    ├── admin/v1/         # 管理者ポータルAPI
    └── user/             # 購入者ポータルAPI
        ├── facility/     # 施設関連API
        └── v1/           # ユーザーAPI v1
```

## 各ディレクトリの役割

### 📋 external/
外部サービスとの連携仕様書を格納
- 決済サービス連携仕様
- 配送業者API仕様
- SNS連携仕様
- その他外部API仕様

### 🏗️ internal/
ふるマル機能実装の仕様書を格納
- 機能仕様書
- 実装仕様書
- ビジネスルール定義
- 実装ガイドライン

### 📊 swagger/
OpenAPI 3.0 形式のAPI定義を格納
- **admin/v1/**: 管理者ポータル用REST API定義
- **user/facility/**: 施設関連API定義
- **user/v1/**: 購入者ポータル用REST API定義

## ファイル管理ルール

### 命名規則
- ファイル名: `kebab-case.yaml` または `kebab-case.md`
- API仕様: OpenAPI 3.0 形式（YAML推奨）
- 設計書: Markdown形式

### 更新基準
- API変更時は対応するSwagger定義を同時更新
- 外部連携変更時は該当仕様書を更新
- 内部設計変更時は関連文書を更新

### バージョン管理
- API仕様はセマンティックバージョニングに従う
- 破壊的変更は新バージョンディレクトリで管理
- 旧バージョンは一定期間保持

## 参照順序

仕様確認時は以下の順序で参照することを推奨：

1. **swagger/**: 現在のAPI仕様
2. **internal/**: 機能実装仕様
3. **external/**: 外部連携仕様

## 関連ドキュメント

- [アーキテクチャ概要](../architecture/overview.md)
- [API サービス概要](../architecture/api/api-services-overview.md)
- [API ドキュメントパターン](../architecture/api/documentation-patterns.md)
