# AI エージェント共通運用ガイド

このリポジトリは、Codex・Claude Code 共通で利用できるよう整備されています。エージェントは以下の方針に従って安全かつ効率的に作業してください。

## プロジェクト概要

**Furumaru** は地域特産品を扱う日本の EC マーケットプレイスプラットフォーム（ライブコマース機能付き）です。Go マイクロサービス + Nuxt フロントエンドのモノレポ構成。

## クイックスタート

```bash
make setup                    # 初回セットアップ
cp .env.temp .env            # 環境ファイル作成・編集
make start                   # 全サービス起動
cd api && make test          # API テスト
cd web/user && yarn dev      # Web 開発
```

## 詳細ガイド

### 基本方針・運用規約
- **[docs/agents/charter.md](docs/agents/charter.md)** - エージェントの基本方針と運用規約

### ナレッジ・作業記録
- **[docs/agents/knowledge-management.md](docs/agents/knowledge-management.md)** - ナレッジ管理・作業記録の運用
- **[docs/agents/quick-start.md](docs/agents/quick-start.md)** - 開発環境セットアップ・基本操作

### 参照ドキュメント（推奨順）

#### 規範的ドキュメント（変更に強い）
- [docs/rules/architecture.md](docs/rules/architecture.md) - アーキテクチャ設計ルール
- [docs/rules/coding.md](docs/rules/coding.md) - コーディング規約
- [docs/rules/tech-stack.md](docs/rules/tech-stack.md) - 技術スタック選定規約
- [docs/rules/workflow.md](docs/rules/workflow.md) - Git・PR ワークフロー

#### 設計意図・構造
- [docs/architecture/overview.md](docs/architecture/overview.md) - システム全体概要
- [docs/architecture/api/documentation-patterns.md](docs/architecture/api/documentation-patterns.md) - API ドキュメント記述パターン

#### 実務知見・手順
- [docs/knowledge/README.md](docs/knowledge/README.md) - 実装パターンと技術的知見
- [docs/knowledge/commands.md](docs/knowledge/commands.md) - よく使うコマンド
- [docs/knowledge/troubleshooting.md](docs/knowledge/troubleshooting.md) - トラブルシューティング

#### 作業用（Claude Code のみ）
- [.claude/debug/](/.claude/debug/) - 作業セッション記録・デバッグ情報

## 判断に迷ったら

1. **アーキテクチャ・コマンド**: 上記の参照ドキュメントを優先参照
2. **セッション記録**: `.claude/debug/` に一時記録、汎用知見は `docs/knowledge/` へ
3. **破壊的操作**: 事前に確認を取る

## 詳細

エージェント向けの詳細ガイドは **[docs/agents/README.md](docs/agents/README.md)** を参照してください。
