//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package user

import (
	"context"

	"github.com/and-period/furumaru/api/internal/user/entity"
)

//nolint:revive
type UserService interface {
	// 管理者サインイン
	SignInAdmin(ctx context.Context, in *SignInAdminInput) (*entity.AdminAuth, error)
	// 管理者サインアウト
	SignOutAdmin(ctx context.Context, in *SignOutAdminInput) error
	// 管理者認証情報取得
	GetAdminAuth(ctx context.Context, in *GetAdminAuthInput) (*entity.AdminAuth, error)
	// 管理者アクセストークンの更新
	RefreshAdminToken(ctx context.Context, in *RefreshAdminTokenInput) (*entity.AdminAuth, error)
	// 管理者一覧取得
	ListAdmins(ctx context.Context, in *ListAdminsInput) (entity.Admins, error)
	// 管理者一覧取得 (ID指定)
	MultiGetAdmins(ctx context.Context, in *MultiGetAdminsInput) (entity.Admins, error)
	// 管理者取得
	GetAdmin(ctx context.Context, in *GetAdminInput) (*entity.Admin, error)
	// 管理者登録
	CreateAdmin(ctx context.Context, in *CreateAdminInput) (*entity.Admin, error)
	// 管理者メールアドレス更新
	UpdateAdminEmail(ctx context.Context, in *UpdateAdminEmailInput) error
	// 管理者メールアドレス更新後の確認
	VerifyAdminEmail(ctx context.Context, in *VerifyAdminEmailInput) error
	// 管理者パスワード更新
	UpdateAdminPassword(ctx context.Context, in *UpdateAdminPasswordInput) error
	// 購入者サインイン
	SignInUser(ctx context.Context, in *SignInUserInput) (*entity.UserAuth, error)
	// 購入者サインアウト
	SignOutUser(ctx context.Context, in *SignOutUserInput) error
	// 購入者認証情報取得
	GetUserAuth(ctx context.Context, in *GetUserAuthInput) (*entity.UserAuth, error)
	// 購入者アクセストークン更新
	RefreshUserToken(ctx context.Context, in *RefreshUserTokenInput) (*entity.UserAuth, error)
	// 購入者取得
	GetUser(ctx context.Context, in *GetUserInput) (*entity.User, error)
	// 購入者登録 (メールアドレス/SMS認証)
	CreateUser(ctx context.Context, in *CreateUserInput) (string, error)
	// 購入者登録後の確認 (メールアドレス/SMS認証)
	VerifyUser(ctx context.Context, in *VerifyUserInput) error
	// 購入者登録 (OAuth認証)
	CreateUserWithOAuth(ctx context.Context, in *CreateUserWithOAuthInput) (*entity.User, error)
	// 購入者登録後の初回更新
	InitializeUser(ctx context.Context, in *InitializeUserInput) error
	// 購入者メールアドレス更新
	UpdateUserEmail(ctx context.Context, in *UpdateUserEmailInput) error
	// 購入者メールアドレス更新後の確認
	VerifyUserEmail(ctx context.Context, in *VerifyUserEmailInput) error
	// 購入者パスワード更新
	UpdateUserPassword(ctx context.Context, in *UpdateUserPasswordInput) error
	// 購入者パスワードリセット (メール送信)
	ForgotUserPassword(ctx context.Context, in *ForgotUserPasswordInput) error
	// 購入者パスワードリセット (パスワード更新)
	VerifyUserPassword(ctx context.Context, in *VerifyUserPasswordInput) error
	// 購入者退会
	DeleteUser(ctx context.Context, in *DeleteUserInput) error
}
