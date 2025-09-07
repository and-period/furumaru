# /project:load-rules

プロジェクトルールファイルを読み込んで理解します。

## 実行内容

`.claude/rules/` ディレクトリ内の全てのルールファイルを読み込み、Furumaruプロジェクトの制約、技術スタック、コーディングルール、開発ワークフローを理解します。

## 使用方法

```
/project:load-rules
```

このコマンドを実行すると、以下の処理が行われます：

1. `.claude/rules/` ディレクトリ内の全ファイルを検索
2. 各ルールファイルの内容を読み込み
3. プロジェクトの重要な制約やルールを理解
4. 以降の作業でこれらのルールを適用

## 対象ファイル

- **architecture.md**: アーキテクチャ設計ルール
- **tech-stack.md**: 技術スタックルール
- **coding.md**: コーディングルール
- **workflow.md**: 開発ワークフロー

## 重要なルール例

### アーキテクチャルール
- **マイクロサービス**: Gateway → Internal Services → Database
- **Frontend-Backend分離**: Nuxt 3 SPA + Go REST API/gRPC
- **非同期処理**: AWS SQS を使用した注文・通知処理

### 技術スタック
- **Frontend**: Nuxt 3 + Vue 3 + Vuetify/Tailwind CSS + TypeScript
- **Backend**: Go 1.25.1 + Gin + GORM + gRPC
- **Database**: MySQL/TiDB（ドメインごとに分離）
- **Infrastructure**: AWS (S3, Cognito, SQS, MediaLive) + Docker

### コーディングルール
- **Go**: gofmt, golangci-lint準拠、テーブル駆動テスト
- **Vue**: Composition API、TypeScript、ESLint/Prettier
- **命名規則**: Go (CamelCase/camelCase), TypeScript (camelCase)
- **Git**: 英語コミットメッセージ、プレフィックス使用

### ワークフロー
- **Branch Strategy**: main + feature/fix/hotfix ブランチ
- **開発環境**: Docker Compose + Make コマンド
- **PR**: 包括的レビュー + 自動CI/CD
- **リリース**: セマンティックバージョニング

新しいセッションを開始する際は、このコマンドを実行してプロジェクトルールを把握することを推奨します。
