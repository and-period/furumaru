# Git運用規約

## ブランチ戦略
- **main**: 本番環境用の安定ブランチ
- **feature/**: 機能開発用ブランチ
- **fix/**: バグ修正用ブランチ
- **hotfix/**: 緊急対応用ブランチ

## 命名規則

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

### コミットメッセージ規約
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

## プルリクエスト規約

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

## リリース規約

### バージョニング
- セマンティックバージョニングを使用
- `v{major}.{minor}.{patch}` 形式
- 例：`v1.2.3`