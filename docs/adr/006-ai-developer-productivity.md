# ADR-006: AI活用を前提とした開発生産性向上

## ステータス

承認済み

## 日付

2026-02-21

## コンテキスト

Furumaru プロジェクトでは Claude Code / Codex を開発ワークフローに統合しており、CLAUDE.md やエージェント向けドキュメント (`docs/agents/`) を整備している。

AI エージェントを活用した開発をさらに効率化するためには、コードベース自体が AI にとって理解しやすく、安全に変更可能な構造であることが重要となる。

### 現状の課題

1. **インターフェース定義のばらつき**: サービス層のインターフェースが一部で不明瞭であり、AI がコード生成する際の制約が不十分
2. **テストパターンの非統一**: テスト記述のパターンが統一されておらず、AI によるテスト自動生成の品質にばらつきが生じる
3. **ドキュメントの断片化**: アーキテクチャ知見がコードコメント、docs/、CLAUDE.md に分散している

### AI 開発における原則

AI エージェントが効果的にコード生成・修正を行うための前提条件:
- **明確なインターフェース**: 入出力の型が明示されていること
- **一貫したパターン**: 同じ種類の処理が同じパターンで実装されていること
- **テストの自動実行**: 変更の正しさを自動的に検証できること
- **十分なコンテキスト**: AI が判断に必要な情報にアクセスできること

## 決定

AI 活用を前提とした開発生産性向上のため、以下の施策を継続的に実施する。

### 1. インターフェース駆動設計の徹底

全てのサービス層に明確なインターフェースを定義し、実装との対応を明示する。

```go
// service/interface.go
type ProductService interface {
    ListProducts(ctx context.Context, params *ListProductsParams) (entity.Products, int64, error)
    GetProduct(ctx context.Context, productID string) (*entity.Product, error)
    CreateProduct(ctx context.Context, params *CreateProductParams) (*entity.Product, error)
    UpdateProduct(ctx context.Context, params *UpdateProductParams) error
    DeleteProduct(ctx context.Context, productID string) error
}
```

これにより AI エージェントは:
- インターフェースを見て呼び出し側のコードを正確に生成できる
- モック生成による単体テストを自動作成できる
- 新規メソッド追加時に既存パターンを模倣できる

### 2. テストパターンの標準化

テーブル駆動テストの統一パターンを定義し、AI がテスト生成時に参照できるようにする。

```go
func TestProductService_GetProduct(t *testing.T) {
    tests := []struct {
        name      string
        setup     func(t *testing.T, mocks *mocks)
        productID string
        expect    *entity.Product
        expectErr error
    }{
        {
            name: "正常系: 商品取得",
            setup: func(t *testing.T, mocks *mocks) {
                mocks.db.EXPECT().GetProduct(gomock.Any(), "product-id").
                    Return(&entity.Product{ID: "product-id"}, nil)
            },
            productID: "product-id",
            expect:    &entity.Product{ID: "product-id"},
            expectErr: nil,
        },
        // ... 異常系テストケース
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ...
        })
    }
}
```

### 3. CLAUDE.md と docs/ の継続的改善

リファクタリングの進捗に合わせて、以下のドキュメントを継続的に更新する。

- **CLAUDE.md**: プロジェクト概要とクイックスタート
- **docs/agents/**: エージェント向けのガイドラインと作業記録
- **docs/rules/**: コーディング規約とアーキテクチャルール
- **docs/knowledge/**: 実装パターンと技術的知見

### 4. AI フレンドリーなコード構成

- 1ファイル1責務の原則を徹底する (Gateway handler 分割と整合)
- 関数のサイズを小さく保つ (50行以内を目安)
- コメントで「なぜ」を記述する (「何を」はコードで表現)

## 検討した選択肢

### 選択肢1: 現状維持
- 利点: 変更コストゼロ
- 欠点: AI 活用の効果が限定的
- **却下**: AI 活用の潜在力を引き出せない

### 選択肢2: インターフェース駆動 + テスト標準化 + ドキュメント改善 (採用)
- 利点: 人間の開発者にも有益、AI 活用効率の向上、段階的に実施可能
- 欠点: ドキュメント維持の継続コスト
- **採用**: 人間・AI 双方に利益があり、投資対効果が高い

### 選択肢3: AI 専用の設定ファイル・メタデータ導入
- 利点: AI 向けに最適化された情報提供
- 欠点: 人間の開発者にとって冗長、維持コストが高い
- **却下**: 人間にも有益な方法で AI 支援を実現すべき

## 結果

### 期待される効果
- AI によるコード生成の精度向上 (インターフェースに基づく正確な型推論)
- AI によるテスト自動生成の品質向上 (統一パターンの模倣)
- 新規参加者 (人間・AI 問わず) のオンボーディング時間短縮
- コードレビューの効率化 (パターン準拠の自動チェック)

### 受け入れるトレードオフ
- ドキュメント維持の継続的なコスト
- インターフェース定義の追加作業
- テストパターン統一のための既存テスト修正

### 測定指標

以下の指標で効果を測定する。
- AI エージェントによるPR作成の成功率 (テスト PASS 率)
- コードレビューでの指摘事項数の推移
- 新規機能実装にかかる時間の推移

## 関連

- [全体設計書](../architecture/api/backend-refactoring-overview.md)
- [詳細設計書 Phase 5](../architecture/api/backend-refactoring-detailed-design.md#phase-5-将来的検討事項)
- [CLAUDE.md](../../CLAUDE.md)
- [エージェント運用ガイド](../agents/README.md)
- [コーディング規約](../rules/coding.md)
