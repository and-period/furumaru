# 開発手順

## CI/CD

### 自動化されるプロセス
1. **コード検証**:
   - Lint実行
   - 型チェック
   - 単体テスト実行
   - カバレッジ測定

2. **ビルド検証**:
   - Docker イメージビルド
   - 成果物の検証

3. **デプロイ**（本番）:
   - ステージング環境でのテスト
   - 本番環境へのデプロイ
   - ヘルスチェック

## ローカル開発フロー

### 初回セットアップ
```bash
# リポジトリクローン
git clone <repository-url>
cd furumaru

# 環境設定
cp .env.temp .env
# .envファイルを編集（AWS認証情報など）

# 初期セットアップ
make setup

# サービス起動
make start
```

### 日常の開発フロー
```bash
# 最新コードの取得
git checkout main
git pull origin main

# 新しいブランチの作成
git checkout -b feature/123-new-feature

# 開発
# ... コード編集 ...

# テスト実行
cd api && make test
cd web/user && yarn test

# コミット
git add .
git commit -m "feat: implement new feature"

# プッシュ
git push origin feature/123-new-feature

# PRの作成（GitHub Web UI）
```

### コードレビュー後の手順
```bash
# レビュー修正
# ... コード編集 ...
git add .
git commit -m "fix: address review comments"
git push origin feature/123-new-feature

# マージ後のクリーンアップ
git checkout main
git pull origin main
git branch -d feature/123-new-feature
```

## リリース手順

### リリース作業
1. **リリースブランチ作成**
2. **CHANGELOG.md更新**
3. **バージョン番号更新**
4. **最終テスト実行**
5. **PRレビュー・マージ**
6. **本番デプロイ**
7. **動作確認**
8. **リリースノート作成**

## トラブル時の対応手順

### ホットフィックス手順
1. `main`ブランチから`hotfix/`ブランチを作成
2. 修正を実装・テスト
3. 緊急PRを作成（レビューは最小限）
4. マージ後、即座にデプロイ
5. 事後検証・ドキュメント更新

### ロールバック手順
1. 問題の影響範囲を特定
2. 前バージョンへのロールバック実行
3. 問題の根本原因調査
4. 修正版の準備・テスト
5. 再デプロイ

## プルリクエスト作成手順

### PRの作成手順
1. 機能開発・バグ修正が完了したらPRを作成
2. 適切なタイトルと説明を記述
3. 関連するissueをリンク
4. レビュアーを指定
5. 自動テストの結果を確認
6. レビュー指摘への対応

### レビュー後の対応
1. レビューコメントを確認
2. 必要な修正を実施
3. 修正内容をコミット・プッシュ
4. レビュアーに再レビューを依頼
5. 承認後のマージ