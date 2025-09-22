# 管理者外部認証プロバイダー連携仕様

## 概要

管理者アカウントに外部認証プロバイダー（Google、LINE）を連携する機能の仕様変更について記載します。

## 変更日時

2025-01-22

## 変更背景

Cognito の仕様により、すでにサインイン済みのユーザーアカウントは外部プロバイダーと連携できない制約がありました。
この問題に対処するため、実装方針を変更しました。

## 変更内容

### 従来の実装（削除）

```go
// 旧実装：外部アカウントを一度削除してから連携
if err := s.adminAuth.DeleteUser(ctx, user.Username); err != nil {
    return internalError(err)
}
```

以下の理由により削除：
1. **破壊的操作**: 外部アカウントの削除は予期しない副作用を引き起こす可能性がある
2. **データ整合性**: 削除と再作成の間でデータの不整合が発生するリスク
3. **Cognito の仕様変更**: 現在のCognito仕様では、この削除操作なしでも連携が可能

### 新しい実装

```go
// 新実装：削除操作なしで直接連携
linkParams := &cognito.LinkProviderParams{
    Username:     admin.CognitoID,
    ProviderType: provider.ProviderType.ToCognito(),
    AccountID:    provider.AccountID,
}
if err := s.adminAuth.LinkProvider(ctx, linkParams); err != nil {
    return internalError(err)
}
```

## 影響範囲

### 修正ファイル

1. **api/internal/user/service/admin_auth.go**
   - `connectAdminAuth` 関数から `DeleteUser` の呼び出しを削除

2. **api/internal/user/service/admin_auth_test.go**
   - モックの期待値から `DeleteUser` の呼び出しを削除
   - エラーケースのテストを更新

## エラーハンドリングの変更

### 変更前
- `ErrInvalidAdminAuthUsername` → `exception.ErrAlreadyExists`

### 変更後
- `ErrInvalidAdminAuthProviderType` → `exception.ErrForbidden`

プロバイダータイプの不一致エラーのみを適切に処理するよう変更しました。

## テスト

以下のテストケースが正常に動作することを確認：
- `TestConnectGoogleAdminAuth`
- `TestConnectLINEAdminAuth`
- `TestConnectAdminAuth`

## 注意事項

1. この変更により、Cognitoの外部アカウント管理がより安全になります
2. 既存の連携済みアカウントには影響ありません
3. 新規連携時のフローが簡潔になり、エラーリスクが減少します

## 関連リソース

AWS Cognito のユーザープール連携に関する公式ドキュメント：
- [Linking federated users to an existing user profile](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-identity-federation.html)
- [AdminLinkProviderForUser API](https://docs.aws.amazon.com/cognito-user-identity-pools/latest/APIReference/API_AdminLinkProviderForUser.html)

※ Cognitoの具体的な仕様変更内容については、AWS の変更履歴やリリースノートを参照してください。