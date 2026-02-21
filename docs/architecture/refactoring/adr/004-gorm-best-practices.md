# ADR-004: GORMベストプラクティスの適用

## ステータス

承認済み

## 日付

2026-02-21

## コンテキスト

現行の GORM 利用パターンは全体として 7.5/10 の品質評価であり、TiDB リトライ機構やトランザクション管理は堅牢である。しかし、以下の改善点が特定された。

### 1. Preload 未使用 (fill() パターン)
- 関連エンティティの取得に手動 fill() パターンを使用している
- N+1 問題の潜在的リスクがある
- Preload / Joins を活用することで推定 -30% のコード削減が可能

### 2. JSON ラッパー型の重複
- 7箇所以上で internalXXX 構造体による JSON ラッパー型が重複実装されている
- 各モジュールで同様のシリアライゼーション/デシリアライゼーションロジックが繰り返されている

### 3. Statement() カスタム関数
- GORM 標準の `WithContext()` で代替可能なカスタム `Statement()` 関数が使用されている
- GORM のバージョンアップ時に互換性リスクとなる

### 4. MySQL コネクション設定不足
- コネクションプールの設定 (`MaxOpenConns`, `MaxIdleConns`, `ConnMaxLifetime` 等) が未設定または不十分

## 決定

GORM の利用パターンを以下の方針で改善する。

### 1. GORM Preload の導入

手動 fill() パターンを GORM Preload / Joins に段階的に置換する。

**Before**:
```go
func (d *database) ListProducts(ctx context.Context) (entity.Products, error) {
    var products entity.Products
    err := d.db.Find(&products).Error
    if err != nil {
        return nil, err
    }
    d.fill(ctx, products) // 手動で関連エンティティを取得
    return products, nil
}
```

**After**:
```go
func (d *database) ListProducts(ctx context.Context) (entity.Products, error) {
    var products entity.Products
    err := d.db.WithContext(ctx).
        Preload("Category").
        Preload("Producer").
        Find(&products).Error
    if err != nil {
        return nil, err
    }
    return products, nil
}
```

### 2. カスタム GORM Type による JSON カラム統一

ジェネリクスを活用したカスタム GORM Type を定義し、JSON ラッパー型を統一する。

```go
// pkg/gorm/types/json.go
type JSON[T any] struct {
    Data T
}

func (j JSON[T]) Value() (driver.Value, error) {
    return json.Marshal(j.Data)
}

func (j *JSON[T]) Scan(value interface{}) error {
    bytes, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("failed to scan JSON: %v", value)
    }
    return json.Unmarshal(bytes, &j.Data)
}
```

### 3. Statement() の廃止

カスタム Statement() 関数を GORM 標準の WithContext() に置換する。

### 4. コネクション設定の追加

MySQL/TiDB 向けの適切なコネクションプール設定を追加する。

## 検討した選択肢

### Preload について

#### 選択肢1: fill() パターン継続
- 利点: 変更コストゼロ
- 欠点: N+1 リスク、冗長コードの維持
- **却下**: 技術的負債の蓄積

#### 選択肢2: GORM Preload / Joins 導入 (採用)
- 利点: コード -30% 削減、N+1 解消、GORM 標準パターンに準拠
- 欠点: クエリパターンの変化による性能変動の可能性
- **採用**: ベンチマークで性能を確認しつつ段階的に移行

#### 選択肢3: sqlc への全面移行
- 利点: SQL 直接制御、最高のパフォーマンス
- 欠点: 大規模な書き換え、Phase 5 の検討範囲
- **保留**: 将来的に検討 (ADR を別途作成)

### ORM 自体について

2026年の ORM 評価:
- **GORM**: プロトタイプ・中小規模向き。エコシステムが最大
- **sqlc**: SQL 中心、パフォーマンス重視で推奨上昇中
- **ent**: 複雑なドメインモデル向き

現行の GORM を最適化した上で継続使用し、将来的に sqlc への段階的移行を検討する。

## 結果

### 期待される効果
- fill() パターン排除によるコード -30% 削減
- JSON ラッパー型統一による重複コード排除 (7箇所以上)
- GORM 標準 API 準拠によるバージョンアップ追従性向上
- コネクション設定追加による DB 接続の安定性向上

### 受け入れるトレードオフ
- Preload 導入によるクエリパターン変化 (ベンチマークで検証必要)
- カスタム GORM Type 導入の学習コスト
- 段階的移行に伴う一時的なパターン混在

### 実施上の制約
- モジュールごとに段階的に移行し、一括変更は避ける
- 各移行後にベンチマークテストを実施し、性能劣化がないことを確認する
- DB スキーマの変更は伴わない (アプリケーション層の変更のみ)

## 関連

- [全体設計書](../overview.md)
- [詳細設計書 Phase 4](../detailed-design.md#phase-4-gormdb-改善-1-2ヶ月)
- [既存 DB 設計](../../api/database-design.md)
