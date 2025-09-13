# .claude ディレクトリ

Claude Code の作業用ディレクトリです。エージェント共通の正式ドキュメントは `docs/` 配下に移設されています。

## ディレクトリ構造

```
.claude/
├── README.md              # このファイル（このディレクトリの説明）
├── commands/              # Claude Code カスタムコマンド定義
├── debug/                 # 作業セッション記録・デバッグ情報（重要：残す）
├── settings.local.json    # Claude Code ローカル設定
├── session.json          # セッション管理情報
└── session-manager.sh    # セッション管理スクリプト
```

## 正式ドキュメントの参照先

エージェント向けの正式ドキュメントは以下に移設されています：

### メインインデックス
- **[AGENTS.md](../AGENTS.md)** - エージェント共通の運用ガイド（メインインデックス）

### 詳細ドキュメント
- **[docs/agents/](../docs/agents/)** - エージェント向け詳細ガイド
- **[docs/rules/](../docs/rules/)** - 守るべき規約・原則
- **[docs/architecture/](../docs/architecture/)** - システム設計意図・構造
- **[docs/knowledge/](../docs/knowledge/)** - 実装知見・手順・Tips

## .claude ディレクトリの役割

### commands/ - Claude Code カスタムコマンド
Claude Code の `/load-rules`、`/save` などのカスタムコマンド定義。

### debug/ - 作業セッション記録（重要：保持）
- 作業セッションの詳細ログ
- エラー発生時の調査・解決過程
- 試行錯誤の記録
- ファイル命名規則: `YYYY-MM-DD_[作業内容].md`

### 設定ファイル
- **settings.local.json**: Claude Code のローカル設定
- **session.json**: セッション管理情報
- **session-manager.sh**: セッション管理スクリプト

## 注意点

- **作業記録は必ず残してください**: `.claude/debug/` の内容は Claude Code の作業に必要です。
- **恒常的なドキュメント**: 一般的に参照される内容は `docs/` 配下に移設済みです。
- **一時的な調査メモ**: `.claude/debug/` に記録し、汎用化できる知見は `docs/knowledge/` に整理してください。
