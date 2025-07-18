# /project:push - 差分コミット＆PR作成/更新

## 概要

作業中の全差分を適切にコミットし、PRを作成または更新するコマンドです。
PRガイドラインに従って、レビュワーを意識したPR作成を行います。

## 使用方法

```
/project:push
```

## 処理フロー

1. **差分確認**
   - `git status` で現在の変更を確認
   - 変更されたディレクトリを特定（api/web）

2. **ブランチ確認・作成**
   - 現在のブランチを確認
   - `main`ブランチの場合は自動的に新しいフィーチャーブランチを作成
   - ブランチ名: `feature/YYYY-MM-DD-description` または `fix/YYYY-MM-DD-description`

3. **コミット作成**
   - 変更内容に応じた適切なコミットメッセージを生成
   - コミットフォーマット: `{type}: {subject}`
   - type: feat/fix/docs/style/refactor/test/chore

4. **PR作成/更新**
   - 既存PRがある場合: pushのみ実行
   - 新規PRの場合: 
     - テンプレートに従った説明文を生成
     - `main`ブランチをベースブランチに指定
     - draft PRとして作成

## コミットメッセージ規則

### typeの選び方
- **feat**: 新機能の追加
- **fix**: バグ修正
- **docs**: ドキュメントのみの変更
- **style**: フォーマット、セミコロン、空白などの修正
- **refactor**: リファクタリング
- **test**: テストの追加・修正
- **chore**: ビルドプロセスやツールの変更

### 対象の指定
- API変更: `feat(api): `
- フロントエンド変更: `fix(web): `
- 両方の変更: `refactor: `
- 設定関連: `chore: `

## PR説明文テンプレート

```markdown
## 概要
[変更の概要を簡潔に記述]

## 変更内容
- [主な変更点1]
- [主な変更点2]
- [主な変更点3]

## 背景・目的
[なぜこの変更が必要なのか]

## テスト
- [ ] API: make test (Go)
- [ ] API: make lint-fix
- [ ] フロントエンド: yarn typecheck
- [ ] フロントエンド: yarn lint
- [ ] 手動テスト完了

## スクリーンショット（UIの変更がある場合）
[必要に応じてスクリーンショットを添付]

## レビューポイント
[特に見てもらいたい箇所]

## 関連情報
- Issue: #[issue-number]
- 関連PR: #[pr-number]

🤖 Generated with [Claude Code](https://claude.ai/code)
```

## 実行時の注意事項

1. **ブランチ運用**
   - **mainブランチには直接pushしない**
   - mainブランチでの実行時は自動的にフィーチャーブランチを作成
   - 全てのコードはPRを経由してmainブランチにマージ

2. **変更範囲**
   - 可能な限り機能単位でPRを作成
   - APIとフロントエンドの変更は分けることを推奨

3. **コード量**
   - 500行以下を目安に
   - 大きくなる場合は機能を分割

4. **CI/CD**
   - draftで作成し、CI通過後にready for reviewへ
   - mainブランチへのマージが基本

5. **品質チェック**
   - lint/formatチェックを事前実行
   - TypeScript型チェックの通過確認
   - Go テストとlintの通過確認

## 実装詳細

このコマンドは以下の処理を自動化します：

1. 現在のブランチ確認と必要に応じてフィーチャーブランチ作成
2. 差分の分析（api/web）
3. 適切なコミットメッセージの生成
4. `git add -A && git commit`
5. 既存PRの確認（`gh pr list`）
6. PR作成（`gh pr create`）またはpush（`git push`）
7. PR URLの表示

## Furumaru固有の考慮事項

- **マイクロサービス**: 複数サービスの変更は影響範囲を明記
- **データベース変更**: マイグレーションファイルの確認
- **環境変数**: .envファイルの変更は注意深く確認

レビュワーを意識し、分かりやすいPRを作成することで、
スムーズなレビューとマージを実現します。
