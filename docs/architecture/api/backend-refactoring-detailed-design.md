# バックエンドリファクタリング 詳細設計書

本ドキュメントでは、[全体設計書](./backend-refactoring-overview.md) で定義した各フェーズの具体的な実装方針と ToDo を記載する。

---

## Phase 1: セキュリティ・緊急対応 (1-2週間) ✅ 完了

セキュリティ脆弱性を持つライブラリおよびメンテナンス終了ライブラリの移行を最優先で実施する。

### 1.1 satori/go.uuid → google/uuid 移行 (PR #3297)

**背景**: `satori/go.uuid v1.2.0` は CVE-2021-3538 (予測可能な UUID 生成) の脆弱性を持ち、リポジトリも放棄済みである。

**方針**:
- `google/uuid` パッケージに完全移行する
- UUID v4 生成を `uuid.New()` に統一する
- 既存の UUID 文字列パース処理は `uuid.Parse()` に置換する
- `uuid.FromString()` → `uuid.Parse()` の対応に注意する

**影響範囲**: 全モジュール (user, store, media, messenger) のエンティティ生成処理

**ToDo**:
- [x] go.mod から `satori/go.uuid` を削除し `google/uuid` を追加
- [x] 全ソースコード中の `satori/go.uuid` インポートを `google/uuid` に置換
- [x] UUID 生成箇所を `uuid.New()` / `uuid.NewString()` に統一
- [x] UUID パース箇所を `uuid.Parse()` に統一
- [x] 全テストが PASS することを確認
- [x] `go mod tidy` で不要な依存を除去

### 1.2 golang/mock → go.uber.org/mock 完全移行 (PR #3298)

**背景**: `golang/mock v1.6.0` は Google によるメンテナンスが終了しており、`go.uber.org/mock` がコミュニティフォークとして活発にメンテナンスされている。

**方針**:
- `go.uber.org/mock` に完全移行する
- mockgen のバイナリも `go.uber.org/mock/mockgen` に切り替える
- `//go:generate` ディレクティブを全て更新する

**影響範囲**: 全モジュールのモック生成とテストコード

**ToDo**:
- [x] go.mod から `github.com/golang/mock` を削除し `go.uber.org/mock` を追加
- [x] 全ソースコード中の `github.com/golang/mock/gomock` インポートを `go.uber.org/mock/gomock` に置換
- [x] 全ソースコード中のモック生成 `github.com/golang/mock/mockgen` インポートを `go.uber.org/mock/mockgen` に置換
- [x] `//go:generate mockgen` ディレクティブのパッケージパスを更新
- [x] Makefile 内の mockgen 関連コマンドを更新
- [x] モック再生成を実行し差分を確認
- [x] 全テストが PASS することを確認

### 1.3 go-grpc-middleware v1 → v2 完全移行 (PR #3299)

**背景**: `grpc-ecosystem/go-grpc-middleware v1` は deprecated となっており、v2 への移行が推奨されている。

**方針**:
- v2 の新しいパッケージ構成に合わせてインポートパスを更新する
- v1 の `grpc_middleware.ChainUnaryServer()` → v2 のインターセプタチェーン構成に移行
- v2 ではミドルウェアが個別パッケージ (`logging`, `recovery`, `auth` 等) に分割されている点に注意

**主な変更点**:
```go
// v1
import grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
grpc_middleware.ChainUnaryServer(...)

// v2
import "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
import "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
```

**ToDo**:
- [x] go.mod の `go-grpc-middleware` を v2 に更新
- [x] v1 のチェーンミドルウェア構成を v2 のインターセプタ構成に移行
- [x] logging ミドルウェアの更新
- [x] recovery ミドルウェアの更新
- [x] auth ミドルウェアの更新 (使用している場合)
- [x] gRPC サーバー初期化コードの更新
- [x] 全テストが PASS することを確認
- [x] gRPC エンドポイントの動作確認

### 1.4 Go パッチバージョン適用 (PR #3304)

**方針**:
- Go 1.25 系の最新パッチバージョンを適用する
- CI/CD パイプラインの Go バージョン指定を更新する

**ToDo**:
- [x] go.mod の Go バージョン指定を更新
- [x] Dockerfile の Go バージョンを更新
- [x] GitHub Actions の Go バージョンを更新
- [x] ローカル開発環境のセットアップ手順を更新
- [x] CI で全テストが PASS することを確認

---

## Phase 2: ライブラリアップグレード (2-4週間) ✅ 完了

メジャーバージョン遅れの主要ライブラリを段階的にアップグレードする。

### 2.1 stripe-go v73 → v84 マイグレーション (PR #3300)

**背景**: 11メジャーバージョンの遅延は決済系ライブラリとしてセキュリティリスクが高い。Stripe API のバージョニングポリシーにより、各メジャーバージョンは特定の API バージョンに対応する。

**方針**:
- Stripe の公式マイグレーションガイドに従い、段階的に移行する
- 一度に複数メジャーバージョンをスキップする場合、Breaking Changes の累積影響を精査する
- ステージング環境での決済フロー全パステストを必須とする

**影響範囲**: `api/internal/store/` 配下の決済関連コード

**ToDo**:
- [x] Stripe 公式の Breaking Changes 一覧 (v74-v84) を確認・整理
- [x] 影響を受ける API 呼び出し箇所を洗い出し
- [x] go.mod の stripe-go バージョンを v84 に更新
- [x] 型名・メソッド名の変更に対応
- [x] 決済作成フロー (PaymentIntent / Checkout Session) の動作確認
- [x] Webhook イベント処理の互換性確認
- [x] 返金・キャンセル処理の動作確認
- [x] ステージング環境での決済フロー E2E テスト
- [x] 全テストが PASS することを確認

### 2.2 line-bot-sdk-go v7 → v8 マイグレーション (PR #3301)

**背景**: v7 は非推奨であり、v8 は OpenAPI ベースで再設計されている。API インターフェースが大幅に変更されている。

**方針**:
- v8 の OpenAPI ベース API に段階的に移行する
- LINE Messaging API, Login API の呼び出し箇所を全て更新する
- 認証連携 (LINE Login) のフローを重点的にテストする

**影響範囲**: `api/internal/messenger/` 配下の LINE 連携コード、`api/internal/user/` 配下の LINE 認証コード

**ToDo**:
- [x] v8 の公式マイグレーションガイドを確認
- [x] LINE Bot クライアント初期化コードの更新
- [x] メッセージ送信 API の更新 (Push / Reply / Multicast)
- [x] リッチメニュー API の更新 (使用している場合)
- [x] LINE Login / OAuth 認証フローの更新
- [x] Webhook イベント処理の更新
- [x] LINE 開発環境での動作確認
- [x] 全テストが PASS することを確認

### 2.3 その他ライブラリ更新 (PR #3302)

**対象**:
- sentry-go: v0.36.2 → v0.40.0
- Firebase SDK: v4.18.0 → v4.19.0
- gRPC: v1.77.0 → v1.79.x
- AWS SDK: 最新パッチ適用

**方針**:
- マイナーバージョンアップのため、Breaking Changes のリスクは低い
- 一括で更新し、テストで回帰を確認する

**ToDo**:
- [x] sentry-go v0.40.0 に更新
- [x] sentry の初期化オプション変更確認
- [x] Firebase SDK v4.19.0 に更新
- [x] gRPC v1.79.x に更新
- [x] AWS SDK の最新パッチ適用
- [x] `go mod tidy` で依存解決
- [x] 全テストが PASS することを確認

### 2.4 Go 1.26 移行検討 (PRs #3304, #3305)

**背景**: Go 1.26 は 2026-02-10 にリリースされ、以下の注目機能を含む。
- **Green Tea GC**: GC オーバーヘッドを 10-40% 削減
- **new() 初期値対応**: コンストラクタパターンの簡素化
- Go 1.24 の tool ディレクティブ、Swiss Table マップ (大規模マップ30%改善) も利用可能に

**方針**:
- Phase 1-2 完了後、Go 1.26 のリリースノートと互換性を精査する
- 安定性が確認された後 (リリースから1-2ヶ月後を目安) に移行を判断する

**ToDo**:
- [x] Go 1.26 リリースノートの互換性影響を確認
- [x] Green Tea GC のベンチマーク (現行ワークロード比較)
- [x] go.mod の Go バージョンを 1.26 に更新
- [x] Dockerfile, CI の Go バージョンを更新
- [x] 全テストが PASS することを確認
- [x] ステージングでの負荷テスト

---

## Phase 3: Gateway・構成改善 (1-2ヶ月)

Gateway 層の可読性と保守性を向上させる。

### 3.1 Gateway handler 責務分割 — ✅ 完了（既にドメイン別に分割済み）

**背景**: 現在の handler.go は300行超で30+エンドポイントを1ファイルで管理しており、変更時のコンフリクトリスクが高い。

**方針**:
- ドメインごとにハンドラーファイルを分割する
- 既存のルーティング構造を変更せず、ファイル分割のみを行う
- mytec プロジェクトの `handler/http/{admin,user}/` パターンを参考にする

**分割案**:
```
api/internal/gateway/
├── handler.go              # 共通処理・ルーター初期化
├── handler_auth.go         # 認証・認可関連エンドポイント
├── handler_user.go         # ユーザー管理エンドポイント
├── handler_store.go        # 商品・注文・決済エンドポイント
├── handler_media.go        # メディア・ライブ配信エンドポイント
├── handler_messenger.go    # 通知・メッセージエンドポイント
└── handler_webhook.go      # Webhook エンドポイント
```

**ToDo**:
- [x] 現行 handler.go のエンドポイント一覧を作成
- [x] ドメインごとのグルーピングを決定
- [x] handler_auth.go にauth関連ハンドラーを切り出し
- [x] handler_user.go にユーザー管理ハンドラーを切り出し
- [x] handler_store.go にEC関連ハンドラーを切り出し
- [x] handler_media.go にメディア関連ハンドラーを切り出し
- [x] handler_messenger.go に通知関連ハンドラーを切り出し
- [x] handler_webhook.go にWebhook処理を切り出し
- [x] handler.go にルーター初期化と共通処理のみ残す
- [x] 全エンドポイントのルーティングが変わっていないことを確認
- [x] 全テストが PASS することを確認

> **注記**: 既存コードベースにおいてハンドラーは既にドメイン別に分割されていたため、追加の分割作業は不要であった。

### 3.2 registry.go DI初期化の分割 (PR #3308)

**背景**: registry.go の inject() 関数が310行超の巨大な単一関数となっており、可読性と変更容易性が低い。

**方針**:
- 関心事ごとに初期化関数を分割する
- ファイル分割は必須ではなく、まず関数の分割から始める

**分割案**:
```go
// registry.go - メイン初期化
func (r *registry) inject() {
    r.injectAWS()
    r.injectDatabase()
    r.injectServices()
    r.injectMiddleware()
}

// registry_aws.go - AWS関連初期化
func (r *registry) injectAWS() { ... }

// registry_database.go - DB関連初期化
func (r *registry) injectDatabase() { ... }

// registry_services.go - サービス層初期化
func (r *registry) injectServices() { ... }
```

**ToDo**:
- [x] 現行 inject() 関数の依存グラフを分析
- [x] AWS 関連初期化を injectAWS() に抽出
- [x] DB 関連初期化を injectDatabase() に抽出
- [x] サービス層初期化を injectServices() に抽出
- [x] ミドルウェア初期化を injectMiddleware() に抽出
- [x] 各初期化関数を適切なファイルに配置
- [x] 初期化順序の依存関係が壊れないことを確認
- [x] 全テストが PASS することを確認

### 3.3 RBAC同期間隔の短縮 (PR #3307)

**背景**: 現行の RBAC 同期は5分間隔で行われており、権限変更の反映に最大5分の遅延が生じる。

**方針**:
- 同期間隔を5分から1分に短縮する
- 将来的にはイベント駆動での即時同期を検討する

**ToDo**:
- [x] RBAC 同期間隔の設定箇所を特定
- [x] 同期間隔を5分から1分に変更
- [x] 同期処理の負荷が問題ないことを確認
- [x] 権限変更から反映までの遅延が短縮されたことを確認

### 3.4 Presenter パターンの導入検討

**背景**: mytec プロジェクトでは `handler/http/staff/presenter/` でレスポンス変換を分離しており、ハンドラーの責務がシンプルになっている。

**方針**:
- まず新規エンドポイントから presenter パターンを試験導入する
- 効果を確認後、既存エンドポイントへの段階的適用を検討する
- handler からドメインオブジェクト → レスポンス DTO の変換ロジックを分離する

**ToDo**:
- [ ] 現行のレスポンス変換処理のパターンを調査
- [ ] presenter パッケージの設計 (インターフェース定義)
- [ ] 1-2 エンドポイントで試験導入
- [ ] コードの見通し改善効果を評価
- [ ] 展開範囲を決定

---

## Phase 4: GORM・DB 改善 (1-2ヶ月)

データアクセス層のベストプラクティス適用とテスト基盤の強化を行う。

### 4.1 GORM Preload 移行 — 現状維持を推奨

**背景**: 現行は手動 fill() パターンで関連エンティティを取得しており、N+1 問題のリスクと冗長なコードが存在する。GORM Preload を活用することで、推定 -30% のコード削減が可能。

**方針**:
- 段階的に fill() パターンを GORM Preload に置換する
- パフォーマンスが重要な箇所はベンチマークで比較する
- 一括置換ではなく、モジュールごとに移行する

**注意点**:
- Preload は JOIN ベースではなく別クエリ発行のため、N+1 は解消するがクエリ数は増加する場合がある
- `Joins()` による Eager Loading と使い分ける

> **評価結果**: 現行の fill() パターンは errgroup による並行取得で最適化済み。GORM Preload は逐次実行のため性能劣化リスクあり。現状維持を推奨。

**ToDo**:
- [ ] fill() パターンの使用箇所を全て洗い出し
- [ ] Preload / Joins への置換候補を優先度付けで整理
- [ ] User モジュールの fill() を Preload に移行
- [ ] Store モジュールの fill() を Preload に移行
- [ ] Media モジュールの fill() を Preload に移行
- [ ] Messenger モジュールの fill() を Preload に移行
- [ ] 各移行後のクエリ性能をベンチマーク
- [ ] 全テストが PASS することを確認

### 4.2 JSON カラム処理の統一 (PRs #3311, #3313, #3314) ✅ 完了

**背景**: 7箇所以上で internalXXX 構造体による JSON ラッパー型が重複実装されている。カスタム GORM Type で統一することで、重複を排除し保守性を向上させる。

**方針**:
- `database/sql/driver.Valuer` と `sql.Scanner` インターフェースを実装するカスタム型を定義する
- 共通パッケージ (`pkg/` 配下) にカスタム GORM Type を配置する
- 各モジュールの internalXXX 構造体をカスタム型に段階的に置換する

**実装例**:
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

**ToDo**:
- [x] 現行の internalXXX 構造体パターンを全て洗い出し
- [x] カスタム GORM Type のジェネリック型を設計
- [x] pkg/ 配下にカスタム GORM Type パッケージを作成
- [x] 各モジュールの JSON ラッパー型をカスタム型に置換
- [x] マイグレーションの互換性確認 (DB上のデータ形式は変更なし)
- [x] 全テストが PASS することを確認

### 4.3 Statement() → WithContext() 標準化 — 現状維持

**背景**: カスタム `Statement()` 関数が使用されているが、GORM 標準の `WithContext()` で代替可能である。標準 API に寄せることで、GORM のバージョンアップ追従が容易になる。

> **評価結果**: Statement() は179箇所で利用される WithContext+Table+Select ラッパーであり、有用な抽象化として機能している。置換コストに対してメリットが薄いため、現状維持とする。

**ToDo**:
- [ ] Statement() の使用箇所と用途を全て洗い出し
- [ ] WithContext() への置換が可能な箇所を特定
- [ ] 段階的に WithContext() に置換
- [ ] カスタム Statement() 関数を廃止
- [ ] 全テストが PASS することを確認

### 4.4 MySQL コネクション設定の追加 (PR #3310) ✅ 完了

**背景**: TiDB 向けのリトライ機構は充実しているが、MySQL 互換のコネクションプール設定が不足している。

**ToDo**:
- [x] 現行のコネクション設定を確認
- [x] `SetMaxOpenConns()` の適切な値を検討・設定
- [x] `SetMaxIdleConns()` の適切な値を検討・設定
- [x] `SetConnMaxLifetime()` の適切な値を検討・設定
- [x] `SetConnMaxIdleTime()` の適切な値を検討・設定
- [x] 負荷テストでコネクションプール挙動を確認

### 4.5 testcontainers-go 導入

**背景**: 現行のテストは Docker Compose ベースで外部 DB を前提としており、環境依存性が高い。`testcontainers-go` により、テスト内で DB コンテナを起動・破棄するセルフコンテインドなテスト環境を実現する。

**方針**:
- mytec プロジェクトの TestMain パターンを参考にする
- 既存テストとの並存期間を設ける (`DISABLE_CONTAINER_DB` 環境変数で切り替え)
- CI 環境では testcontainers-go を標準とする

**実装例** (mytec 参考):
```go
func TestMain(m *testing.M) {
    ctx := context.Background()
    if os.Getenv("DISABLE_CONTAINER_DB") == "true" {
        testDB, cleanup, _ = newTestDB(ctx)
    } else {
        testDB, cleanup, _ = newTestDBWithContainer(ctx)
    }
    defer cleanup()
    code := m.Run()
    os.Exit(code)
}
```

**ToDo**:
- [x] testcontainers-go の依存追加 (PR #3312)
- [x] TiDB/MySQL コンテナ起動ヘルパーの実装 (PR #3312)
- [x] TestMain パターンの共通化 (pkg/testutil 等) (PR #3312)
- [ ] User モジュールへの試験導入
- [ ] 他モジュールへの展開
- [ ] CI パイプラインの更新 (Docker-in-Docker 対応)
- [ ] テスト実行時間の計測・最適化

### 4.6 テスト外部キー制約の正常化

**背景**: テスト環境で `SET foreign_key_checks = 0` を使用しており、参照整合性のバグを見逃すリスクがある。

**方針**:
- テストデータ投入順序を外部キー制約に準拠するよう修正する
- テスト後のクリーンアップも依存順序を考慮する
- 既存テストを段階的に修正する

**ToDo**:
- [ ] `foreign_key_checks = 0` の使用箇所を洗い出し
- [ ] テストデータ投入順序の依存関係を整理
- [ ] テストヘルパーにデータ投入順序を組み込み
- [ ] 段階的に foreign_key_checks を有効化
- [ ] 全テストが外部キー制約有効状態で PASS することを確認

---

## Phase 5: 将来的検討事項

中長期的な技術戦略として、以下の項目を継続的に検討する。これらは Phase 1-4 完了後に具体的な計画を策定する。

### 5.1 Connect-Go (connectrpc) の検討

**背景**: Connect-Go はプロキシ不要で net/http に統合可能な gRPC 互換フレームワークであり、2026年に推奨度が上昇している。

**検討ポイント**:
- 現行の gRPC-Gateway + Gin 構成との比較
- 移行コスト vs 運用簡素化のトレードオフ
- 既存のフロントエンドとの互換性

**ToDo**:
- [ ] Connect-Go のプロトタイプ実装 (1-2エンドポイント)
- [ ] パフォーマンスベンチマーク
- [ ] 移行コストの見積もり
- [ ] 移行判断の ADR 作成

### 5.2 sqlc への段階的移行検討

**背景**: sqlc は SQL 中心のアプローチで型安全なコード生成を行い、パフォーマンス重視のプロジェクトで推奨されている。

**検討ポイント**:
- GORM との並存戦略
- 複雑なクエリが多い箇所から段階的に移行
- コード生成のワークフロー統合

**ToDo**:
- [ ] sqlc のプロトタイプ実装 (読み取り系クエリ)
- [ ] GORM との並存パターンの設計
- [ ] 移行判断の ADR 作成

### 5.3 Go 1.26 新機能活用

**背景**: Go 1.26 で導入される主要機能を活用し、コード品質とパフォーマンスを向上させる。

**検討ポイント**:
- **Green Tea GC**: GC オーバーヘッド 10-40% 削減の恩恵をベンチマークで確認
- **Range-over-function iterators** (Go 1.23): コレクション操作の簡素化
- **Swiss Table マップ** (Go 1.24): 大規模マップ操作の 30% 改善
- **tool ディレクティブ** (Go 1.24): 開発ツール管理の統一

**ToDo**:
- [ ] Green Tea GC のベンチマーク結果をまとめる
- [ ] iterators を活用できる箇所の洗い出し
- [ ] tool ディレクティブによるツール管理の統一
- [ ] 新機能活用のコーディングガイドライン追加

### 5.4 AI 開発生産性向上施策

**背景**: Claude Code / Codex を活用した開発ワークフローを更に強化し、開発効率を向上させる。

**検討ポイント**:
- インターフェース駆動設計の徹底 (AI がコード生成しやすい構造)
- テスト自動化率の向上 (AI によるテスト生成の前提条件整備)
- CLAUDE.md と docs/ の継続的改善

**ToDo**:
- [ ] インターフェース定義のカバレッジ向上
- [ ] テスト生成テンプレートの整備
- [ ] CLAUDE.md のリファクタリング対応更新
- [ ] AI エージェント向けコーディングガイドラインの拡充

---

## 補足: 横断的考慮事項

### テスト戦略

全フェーズを通じて以下のテスト戦略を適用する。

| レベル | 対象 | 方針 |
|--------|------|------|
| ユニットテスト | 各モジュールの service / database 層 | 既存テストの PASS を最低条件。カバレッジ低下を許容しない |
| 統合テスト | モジュール間連携 | Phase 4 で testcontainers-go 導入後に強化 |
| E2E テスト | 決済フロー、認証フロー | Phase 2 の stripe / LINE SDK 移行時に重点実施 |
| ベンチマーク | DB クエリ、GC 挙動 | Phase 4 の GORM 改善、Phase 5 の Go 1.26 移行時 |

### PR 分割方針

各 ToDo 項目は可能な限り独立した PR として提出する。以下の粒度を目安とする。

- **1 PR = 1 ライブラリ移行** (Phase 1-2)
- **1 PR = 1 ファイル分割** (Phase 3)
- **1 PR = 1 モジュール分の改善** (Phase 4)

### 関連ドキュメント

- [全体設計書](./backend-refactoring-overview.md)
- [ADR-001: 依存ライブラリのアップグレード戦略](../../adr/001-dependency-upgrade-strategy.md)
- [ADR-002: モジュラモノリスアーキテクチャの継続](../../adr/002-modular-monolith-continuation.md)
- [ADR-003: Gateway層のリストラクチャリング](../../adr/003-gateway-restructuring.md)
- [ADR-004: GORMベストプラクティスの適用](../../adr/004-gorm-best-practices.md)
- [ADR-005: コンテナベーステストへの移行](../../adr/005-container-based-testing.md)
- [ADR-006: AI活用を前提とした開発生産性向上](../../adr/006-ai-developer-productivity.md)
