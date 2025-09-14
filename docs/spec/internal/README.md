# Internal Specifications

ふるマルの機能実装仕様書を格納するディレクトリです。

## 目的

このディレクトリは以下の目的で使用します：

- **機能仕様書**: 実装する機能の詳細仕様
- **ビジネスルール**: 業務ロジックの定義
- **実装ガイドライン**: 開発時の実装方針
- **データフロー仕様**: 機能間のデータ連携仕様

## ドキュメント作成ルール

### ファイル命名規則
- フォーマット: `${作成日時}_${タイトル}.md`
- 例: `20250914_admin-dashboard-improvement.md`
- 作成日時は `YYYYMMDD` 形式

### テンプレート
- 必ず `template.md` のフォーマットに従って作成すること
- 各セクションは省略せず、該当なしの場合は「特になし」と記載

## ファイル構成例

```
internal/
├── template.md                              # ドキュメントテンプレート
├── 20250914_admin-dashboard-improvement.md  # 管理画面ダッシュボード改善
├── user-registration-spec.md                # ユーザー登録機能仕様
├── order-processing-spec.md                 # 注文処理機能仕様
├── payment-flow-spec.md                     # 決済フロー仕様
└── live-commerce-spec.md                    # ライブコマース機能仕様
```

## 注意事項

- **アーキテクチャ設計**: `docs/architecture/` を参照
- **API仕様**: `docs/spec/swagger/` を参照
- **外部連携**: `docs/spec/external/` を参照

このディレクトリは機能実装の具体的な仕様のみを扱います。
