# Furumaru プロジェクトドキュメント

このディレクトリには Furumaru プロジェクトの正式なドキュメントを格納します。

## ディレクトリ構成（現行）

```
docs/
├── agents/                # AIエージェント向け詳細ガイド
├── architecture/          # アーキテクチャ設計書
│   ├── README.md
│   ├── overview.md
│   ├── api/
│   │   ├── README.md
│   │   └── documentation-patterns.md
│   ├── web/              # フロントエンドアーキテクチャ
│   │   ├── README.md
│   │   ├── components.md
│   │   ├── state-management.md
│   │   └── api-integration.md
│   └── shared/
│       └── README.md
├── knowledge/             # 実装パターンと技術的知見
├── rules/                 # 守るべき規約
│   ├── README.md
│   ├── coding-standards.md
│   ├── design-principles.md
│   ├── git-workflow.md
│   └── tech-stack.md
└── spec/                  # 仕様書・API定義
    ├── external/          # 外部API仕様
    ├── internal/          # 内部設計仕様
    └── swagger/           # OpenAPI/Swagger定義
        ├── admin/v1/
        └── user/
            ├── facility/
            └── v1/
```

## ドキュメント管理ルール

### 作成・更新基準
- エンジニアと AI エージェント双方が再利用できる内容
- 重要な意思決定・長期的に維持すべき情報
- 新メンバーのオンボーディングに有用

### 更新頻度
- アーキテクチャ変更時に必ず更新
- 新機能リリース時に関連ドキュメントを更新
- 四半期ごとに整合性レビュー

### 品質基準
- 技術的に正確で最新の情報
- 第三者が理解可能な記述
- 図表やサンプルコードを適切に含む
- バージョン管理された状態
