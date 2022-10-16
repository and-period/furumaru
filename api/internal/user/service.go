//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package user

import (
	"context"

	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Service interface {
	// 管理者サインイン
	SignInAdmin(ctx context.Context, in *SignInAdminInput) (*entity.AdminAuth, error)
	// 管理者サインアウト
	SignOutAdmin(ctx context.Context, in *SignOutAdminInput) error
	// 管理者認証情報取得
	GetAdminAuth(ctx context.Context, in *GetAdminAuthInput) (*entity.AdminAuth, error)
	// 管理者アクセストークンの更新
	RefreshAdminToken(ctx context.Context, in *RefreshAdminTokenInput) (*entity.AdminAuth, error)
	// 管理者デバイストークンの更新
	RegisterAdminDevice(ctx context.Context, in *RegisterAdminDeviceInput) error
	// 管理者メールアドレス更新
	UpdateAdminEmail(ctx context.Context, in *UpdateAdminEmailInput) error
	// 管理者メールアドレス更新後の確認
	VerifyAdminEmail(ctx context.Context, in *VerifyAdminEmailInput) error
	// 管理者パスワード更新
	UpdateAdminPassword(ctx context.Context, in *UpdateAdminPasswordInput) error
	// 管理者一覧取得(ID指定)
	MultiGetAdmins(ctx context.Context, in *MultiGetAdminsInput) (entity.Admins, error)
	// 管理者デバイストークン一覧取得
	MultiGetAdminDevices(ctx context.Context, in *MultiGetAdminDevicesInput) ([]string, error)
	// 管理者取得
	GetAdmin(ctx context.Context, in *GetAdminInput) (*entity.Admin, error)
	// システム管理者一覧取得
	ListAdministrators(ctx context.Context, in *ListAdministratorsInput) (entity.Administrators, int64, error)
	// システム管理者一覧取得(ID指定)
	MultiGetAdministrators(ctx context.Context, in *MultiGetAdministratorsInput) (entity.Administrators, error)
	// システム管理者取得
	GetAdministrator(ctx context.Context, in *GetAdministratorInput) (*entity.Administrator, error)
	// システム管理者登録
	CreateAdministrator(ctx context.Context, in *CreateAdministratorInput) (*entity.Administrator, error)
	// システム管理者更新
	UpdateAdministrator(ctx context.Context, in *UpdateAdministratorInput) error
	// システム管理者メールアドレス更新
	UpdateAdministratorEmail(ctx context.Context, in *UpdateAdministratorEmailInput) error
	// システム管理者パスワードリセット
	ResetAdministratorPassword(ctx context.Context, in *ResetAdministratorPasswordInput) error
	// 仲介者一覧取得
	ListCoordinators(ctx context.Context, in *ListCoordinatorsInput) (entity.Coordinators, int64, error)
	// 仲介者一覧取得(ID指定)
	MultiGetCoordinators(ctx context.Context, in *MultiGetCoordinatorsInput) (entity.Coordinators, error)
	// 仲介者取得
	GetCoordinator(ctx context.Context, in *GetCoordinatorInput) (*entity.Coordinator, error)
	// 仲介者登録
	CreateCoordinator(ctx context.Context, in *CreateCoordinatorInput) (*entity.Coordinator, error)
	// 仲介者更新
	UpdateCoordinator(ctx context.Context, in *UpdateCoordinatorInput) error
	// 仲介者メールアドレス更新
	UpdateCoordinatorEmail(ctx context.Context, in *UpdateCoordinatorEmailInput) error
	// 仲介者パスワードリセット
	ResetCoordinatorPassword(ctx context.Context, in *ResetCoordinatorPasswordInput) error
	// 生産者一覧取得
	ListProducers(ctx context.Context, in *ListProducersInput) (entity.Producers, int64, error)
	// 生産者一覧取得(ID指定)
	MultiGetProducers(ctx context.Context, in *MultiGetProducersInput) (entity.Producers, error)
	// 生産者取得
	GetProducer(ctx context.Context, in *GetProducerInput) (*entity.Producer, error)
	// 生産者登録
	CreateProducer(ctx context.Context, in *CreateProducerInput) (*entity.Producer, error)
	// 生産者更新
	UpdateProducer(ctx context.Context, in *UpdateProducerInput) error
	// 生産者メールアドレス更新
	UpdateProducerEmail(ctx context.Context, in *UpdateProducerEmailInput) error
	// 生産者パスワードリセット
	ResetProducerPassword(ctx context.Context, in *ResetProducerPasswordInput) error
	// 生産者関連付け
	RelatedProducer(ctx context.Context, in *RelatedProducerInput) error
	// 生産者関連付け解除
	UnrelatedProducer(ctx context.Context, in *UnrelatedProducerInput) error
	// 購入者サインイン
	SignInUser(ctx context.Context, in *SignInUserInput) (*entity.UserAuth, error)
	// 購入者サインアウト
	SignOutUser(ctx context.Context, in *SignOutUserInput) error
	// 購入者認証情報取得
	GetUserAuth(ctx context.Context, in *GetUserAuthInput) (*entity.UserAuth, error)
	// 購入者アクセストークン更新
	RefreshUserToken(ctx context.Context, in *RefreshUserTokenInput) (*entity.UserAuth, error)
	// 購入者一覧取得(ID指定)
	MultiGetUsers(ctx context.Context, in *MultiGetUsersInput) (entity.Users, error)
	// 購入者デバイストークン一覧取得
	MultiGetUserDevices(ctx context.Context, in *MultiGetUserDevicesInput) ([]string, error)
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
