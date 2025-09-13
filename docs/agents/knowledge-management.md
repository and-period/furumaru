# ナレッジ管理・作業記録の運用

Codex・Claude 共通の知識管理とセッション記録のルールです。

## ドキュメント階層

### 1. 規範的ドキュメント（docs/rules/）
変更に強い、プロジェクトの基本的な守るべき規約：
- [アーキテクチャ設計ルール](../rules/architecture.md)
- [コーディング規約](../rules/coding.md)
- [技術スタック選定規約](../rules/tech-stack.md)
- [Git・PR ワークフロー](../rules/workflow.md)

### 2. 設計意図・構造（docs/architecture/）
システム全体の設計思想と意思決定の記録：
- [システム全体概要](../architecture/overview.md)
- [API ドキュメント記述パターン](../architecture/api/documentation-patterns.md)

### 3. 実務知見（docs/knowledge/）
実装で得た知見・手順・トラブルシューティング：
- [実装パターンと技術的知見](../knowledge/README.md)
- [よく使うコマンド](../knowledge/commands.md)
- [トラブルシューティング](../knowledge/troubleshooting.md)

### 4. セッション記録（.claude/debug/）
作業セッションの詳細・一時的な調査記録：
- 作業セッションの詳細ログ
- エラー発生時の調査・解決過程
- 試行錯誤の記録

## 記録のタイミング

### 作業開始時
- セッションログを `.claude/debug/YYYY-MM-DD_[作業内容].md` 形式で作成

### 新しい理解を得た時
- 即座に `docs/knowledge/` に記録（実装パターン、Tips など）

### 問題発生・解決時
- デバッグ過程を `.claude/debug/` に記録
- 解決策を `docs/knowledge/troubleshooting.md` に追加

### 作業完了時
- 重要な知識を整理し適切なディレクトリに保存
- 一般化できる知見を `.claude/debug/` から `docs/knowledge/` へ抽出

## Claude Code 専用コマンド

### /load-rules コマンド
実行時に以下から学習：
- `docs/rules/` - 守るべき規約と原則
- `docs/architecture/` - システム設計意図と構造
- `docs/knowledge/` - 実装知見と手順・Tips

### /save コマンド
実行時に学習内容を以下に保存：
1. 今回のセッション内容を `.claude/debug/` に記録
2. 学んだ知識を `docs/knowledge/` に整理保存
3. 必要に応じて `docs/rules/` のルールを更新

## Codex から Claude 資産を再利用する方法

### セッションログ
- 命名規則（例: `YYYY-MM-DD_session.md`）に従い、`.claude/debug/` に作業メモや調査過程を保存

### ナレッジ
- 再利用可能なパターンやドメイン知見は `docs/knowledge/` に追加・更新
- ディレクトリの重複は作らない。Codex でも同じツリーをそのまま利用

### 対応づけの目安
- Claude の `/load-rules` → Codex では必要時に `docs/` を読み、要点を計画に反映
- Claude の `/save` → Codex では `.claude/debug/` に発見事項を追記し、汎用化できる知見は `docs/knowledge/` に整理
