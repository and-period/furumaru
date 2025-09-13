# クイックスタート（ローカル開発環境）

Codex・Claude 共通の開発環境セットアップ・基本操作ガイドです。

## 初期セットアップ

```bash
# 初回セットアップ：コンテナビルド、依存関係インストール、Swagger生成
make setup

# 環境ファイルを作成・編集
cp .env.temp .env
# .env を編集してAWS認証情報を設定
```

## サービス起動

```bash
make start          # 全サービスを起動
make start-user     # 購入者向けサービスのみ起動
make start-admin    # 管理者向けサービスのみ起動
make migrate        # データベースマイグレーションを実行
```

## API 開発

```bash
cd api
make test           # カバレッジ付きで全テストを実行
make mockgen        # インターフェース変更後にモックを生成
make lint-fix       # Lintエラーを修正
make fmt-fix        # フォーマットを修正
make start-dev SERVICE=gateway/admin  # ホットリロード付きでサービスを起動
```

## Web 開発

### 管理者ポータル
```bash
cd web/admin
yarn dev            # 開発サーバーを起動
yarn typecheck      # 型チェックを実行
yarn lint           # Lintを実行
yarn format         # Lintエラーを自動修正
```

### 購入者ポータル
```bash
cd web/user
yarn dev            # HTTPS付き開発サーバーを起動
yarn test           # Vitestテストを実行
yarn coverage       # カバレッジ付きでテストを実行
yarn typecheck      # 型チェックを実行
yarn lint           # Lintを実行
yarn format         # Lintエラーを自動修正
```

## Git と PR の作法

- `main` から `feature/*`、`fix/*`、`hotfix/*` を切って作業。
- コミットは焦点を絞り、メッセージは規約に従う。
- PR は範囲とテスト観点を明記し、ワークフローのチェックリストに従う。

## 詳細リンク

- [Git ワークフロー詳細](../contributing/git-workflow.md)
- [よく使うコマンド](../knowledge/commands.md)
- [トラブルシューティング](../knowledge/troubleshooting.md)
