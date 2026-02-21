# ADR-005: コンテナベーステストへの移行

## ステータス

承認済み

## 日付

2026-02-21

## コンテキスト

### 現行のテスト環境

現行のテスト環境には以下の問題がある。

1. **Docker Compose 依存**: テスト実行前に外部 DB の起動が必要であり、環境依存性が高い
2. **外部キー制約の無効化**: テストデータ投入時に `SET foreign_key_checks = 0` を使用しており、参照整合性のバグを見逃すリスクがある
3. **テスト間の干渉**: 共有 DB を使用するため、テスト間でデータが干渉する可能性がある

### 2026年のテスト動向

- **testcontainers-go** がコンテナベーステストのデファクトスタンダードとなっている
- テスト内でコンテナを起動・破棄するセルフコンテインドなアプローチが主流

### 参考: mytec プロジェクトのパターン

mytec プロジェクトでは以下のパターンでコンテナベーステストを実現している。

```go
func TestMain(m *testing.M) {
    ctx := context.Background()
    if os.Getenv("DISABLE_CONTAINER_DB") == "true" {
        testDB, cleanup, _ = newTestDB(ctx)
    } else {
        testDB, cleanup, _ = newTestDBWithContainer(ctx) // MySQL 8.0 コンテナ
    }
    defer cleanup()

    // マイグレーション実行
    migrate.Create(ctx, testDB.Ent().Schema, tables, migrate.WithForeignKeys(false))

    code := m.Run()
    os.Exit(code)
}
```

このパターンにより:
- ローカル開発では `DISABLE_CONTAINER_DB=true` で既存 DB を使用可能
- CI 環境ではテストごとにクリーンな DB コンテナが自動起動
- テスト終了後にコンテナが自動破棄

## 決定

**testcontainers-go を導入し、コンテナベーステストに段階的に移行する。** 同時に、テスト時の外部キー制約無効化を段階的に廃止する。

### 実装方針

#### 1. 共通テストヘルパーの作成

`pkg/testutil/` 配下にコンテナベーステスト用のヘルパーを作成する。

```go
// pkg/testutil/database.go
package testutil

func NewTestDBWithContainer(ctx context.Context) (*gorm.DB, func(), error) {
    req := testcontainers.ContainerRequest{
        Image:        "mysql:8.0",
        ExposedPorts: []string{"3306/tcp"},
        Env: map[string]string{
            "MYSQL_ROOT_PASSWORD": "test",
            "MYSQL_DATABASE":     "test",
        },
        WaitingFor: wait.ForListeningPort("3306/tcp"),
    }
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    // ... DB 接続とクリーンアップ関数を返す
}
```

#### 2. TestMain パターンの標準化

mytec のパターンを参考に、環境変数による切り替えを実現する。

#### 3. 外部キー制約の正常化

テストデータ投入順序を外部キー制約に準拠するよう修正し、`foreign_key_checks = 0` を段階的に廃止する。

### 移行順序

1. `pkg/testutil/` にコンテナベーステストヘルパーを作成
2. User モジュールへの試験導入
3. 効果検証後、Store / Media / Messenger モジュールへ展開
4. CI パイプラインの更新
5. テスト外部キー制約の正常化

## 検討した選択肢

### 選択肢1: 現状維持 (Docker Compose ベース)
- 利点: 変更コストゼロ
- 欠点: 環境依存性、外部キー制約問題の継続、テスト間干渉リスク
- **却下**: テスト品質の改善機会を逃す

### 選択肢2: testcontainers-go 導入 (採用)
- 利点: セルフコンテインド、テスト間分離、CI 環境の簡素化、業界標準
- 欠点: テスト実行時間の増加 (コンテナ起動時間)、Docker-in-Docker 対応が必要
- **採用**: テスト品質と CI 信頼性の向上が期待できる

### 選択肢3: SQLite インメモリ DB
- 利点: 最速のテスト実行、外部依存なし
- 欠点: MySQL/TiDB との SQL 互換性の差異、本番環境との乖離
- **却下**: TiDB 固有機能のテストが不可能

### 選択肢4: テスト専用 TiDB コンテナ
- 利点: 本番環境と完全一致
- 欠点: TiDB コンテナの起動時間が長い (30秒以上)
- **保留**: MySQL 8.0 コンテナでの運用を先行し、必要に応じて検討

## 結果

### 期待される効果
- テスト環境のセットアップ手順の簡素化
- テスト間の完全な分離 (データ干渉の排除)
- CI パイプラインの信頼性向上
- 外部キー制約有効化による参照整合性バグの早期発見

### 受け入れるトレードオフ
- テスト実行時間の増加 (コンテナ起動に5-10秒)
- Docker-in-Docker 対応の CI 設定が必要
- 段階的移行期間中の二つのテストパターン並存

### テスト実行時間の緩和策
- TestMain でモジュール単位のコンテナ共有 (テストごとではなくパッケージごとに1コンテナ)
- 並列テスト実行の活用
- `DISABLE_CONTAINER_DB` 環境変数によるローカル開発時のスキップ

## 関連

- [全体設計書](../architecture/api/backend-refactoring-overview.md)
- [詳細設計書 Phase 4](../architecture/api/backend-refactoring-detailed-design.md#phase-4-gormdb-改善-1-2ヶ月)
- [ADR-004: GORMベストプラクティスの適用](./004-gorm-best-practices.md)
