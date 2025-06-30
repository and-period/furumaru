# .claude ディレクトリ

このディレクトリはClaude Codeのプロジェクト固有の情報を管理するためのディレクトリです。

## ディレクトリ構造

```
.claude/
├── README.md              # このファイル
├── rules/                 # プロジェクトルール
│   ├── architecture.md    # アーキテクチャ設計ルール
│   ├── tech-stack.md      # 技術スタックルール
│   ├── coding.md          # コーディングルール
│   └── workflow.md        # 開発ワークフロー
├── knowledge/             # 知見・学習記録
│   ├── knowledge.md       # 実装パターン
│   ├── improvements.md    # 教訓と改善
│   ├── commands.md        # コマンド一覧
│   └── troubleshooting.md # トラブルシューティング
└── debug/                 # デバッグ関連ファイル
    └── README.md          # デバッグ用README
```

## 使用方法

### rules/ - プロジェクトルール
プロジェクト固有のルールや規約を記録します。
- **architecture.md**: システム設計の原則、マイクロサービス間の通信ルールなど
- **tech-stack.md**: 使用する技術スタックの選定基準、バージョン管理など
- **coding.md**: コーディング規約、命名規則、コードレビュー基準など
- **workflow.md**: Git運用、CI/CD、リリースフローなど

### knowledge/ - 知見・学習記録
開発中に得られた知見を記録します。
- **knowledge.md**: 有用な実装パターン、ベストプラクティス
- **improvements.md**: 過去の失敗から学んだ教訓、改善提案
- **commands.md**: よく使うコマンドの一覧とその説明
- **troubleshooting.md**: よくあるエラーとその解決方法

### debug/ - デバッグ関連ファイル
デバッグに役立つ情報を保存します。
- **README.md**: デバッグ手順、ツールの使い方など
