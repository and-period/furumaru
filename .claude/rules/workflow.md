# 開発ワークフロー

## Git運用

### ブランチ戦略
- **main**: 本番環境用の安定ブランチ
- **feature/**: 機能開発用ブランチ
- **fix/**: バグ修正用ブランチ
- **hotfix/**: 緊急対応用ブランチ

### ブランチ命名規則
```
feature/{issue-number}-{brief-description}
fix/{issue-number}-{brief-description}
hotfix/{brief-description}
```

例：
- `feature/123-user-authentication`
- `fix/456-order-calculation-bug`
- `hotfix/critical-security-patch`

### コミットメッセージ規則
```
<type>(<scope>): <description>

<body>

<footer>
```

**Types:**
- `feat`: 新機能
- `fix`: バグ修正
- `refactor`: リファクタリング
- `docs`: ドキュメント更新
- `style`: コードスタイル修正
- `test`: テスト追加・修正
- `chore`: ビルド・設定変更

**例：**
```
feat(store): add product review functionality

- Add review submission API endpoint
- Implement review display component
- Add rating calculation logic

Closes #123
```

## プルリクエスト（PR）

### PRの作成
1. 機能開発・バグ修正が完了したらPRを作成
2. 適切なタイトルと説明を記述
3. 関連するissueをリンク
4. レビュアーを指定

### PRテンプレート
```markdown
## 概要
このPRの目的と実装内容を簡潔に説明

## 変更内容
- [ ] 機能A
- [ ] 機能B
- [ ] テストの追加

## テスト
- [ ] 単体テスト実行済み
- [ ] 結合テスト実行済み
- [ ] 手動テスト実行済み

## チェックリスト
- [ ] コードレビュー観点を満たしている
- [ ] ドキュメントが更新されている
- [ ] 破壊的変更がある場合は適切に記述されている

## スクリーンショット（該当する場合）

## 関連Issue
Closes #xxx
```

### レビュー基準
1. **機能性**: 要求仕様を満たしているか
2. **コード品質**: コーディング規約に準拠しているか
3. **テスト**: 適切なテストが実装されているか
4. **セキュリティ**: セキュリティ要件を満たしているか
5. **パフォーマンス**: パフォーマンスに悪影響がないか

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

### ローカル開発フロー

#### 初回セットアップ
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

#### 日常の開発フロー
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

#### コードレビュー後
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

### バージョニング
- セマンティックバージョニングを使用
- `v{major}.{minor}.{patch}` 形式
- 例：`v1.2.3`

### リリース作業
1. **リリースブランチ作成**
2. **CHANGELOG.md更新**
3. **バージョン番号更新**
4. **最終テスト実行**
5. **PRレビュー・マージ**
6. **本番デプロイ**
7. **動作確認**
8. **リリースノート作成**

## トラブル時の対応

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
