//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package user

import (
	"context"

	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Service interface {
	// 管理者
	SignInAdmin(ctx context.Context, in *SignInAdminInput) (*entity.AdminAuth, error)             // サインイン
	SignOutAdmin(ctx context.Context, in *SignOutAdminInput) error                                // サインアウト
	GetAdminAuth(ctx context.Context, in *GetAdminAuthInput) (*entity.AdminAuth, error)           // 認証情報取得
	RefreshAdminToken(ctx context.Context, in *RefreshAdminTokenInput) (*entity.AdminAuth, error) // アクセストークンの更新
	RegisterAdminDevice(ctx context.Context, in *RegisterAdminDeviceInput) error                  // デバイストークンの更新
	UpdateAdminEmail(ctx context.Context, in *UpdateAdminEmailInput) error                        // メールアドレス更新
	VerifyAdminEmail(ctx context.Context, in *VerifyAdminEmailInput) error                        // メールアドレス更新後の確認
	UpdateAdminPassword(ctx context.Context, in *UpdateAdminPasswordInput) error                  // パスワード更新
	ForgotAdminPassword(ctx context.Context, in *ForgotAdminPasswordInput) error                  // パスワードリセット (メール送信)
	VerifyAdminPassword(ctx context.Context, in *VerifyAdminPasswordInput) error                  // パスワードリセット (パスワード更新)
	MultiGetAdmins(ctx context.Context, in *MultiGetAdminsInput) (entity.Admins, error)           // 一覧取得(ID指定)
	MultiGetAdminDevices(ctx context.Context, in *MultiGetAdminDevicesInput) ([]string, error)    // デバイストークン一覧取得
	GetAdmin(ctx context.Context, in *GetAdminInput) (*entity.Admin, error)                       // １件取得
	// システム管理者
	ListAdministrators(ctx context.Context, in *ListAdministratorsInput) (entity.Administrators, int64, error)  // 一覧取得
	MultiGetAdministrators(ctx context.Context, in *MultiGetAdministratorsInput) (entity.Administrators, error) // 一覧取得(ID指定)
	GetAdministrator(ctx context.Context, in *GetAdministratorInput) (*entity.Administrator, error)             // １件取得
	CreateAdministrator(ctx context.Context, in *CreateAdministratorInput) (*entity.Administrator, error)       // 登録
	UpdateAdministrator(ctx context.Context, in *UpdateAdministratorInput) error                                // 更新
	UpdateAdministratorEmail(ctx context.Context, in *UpdateAdministratorEmailInput) error                      // メールアドレス更新
	ResetAdministratorPassword(ctx context.Context, in *ResetAdministratorPasswordInput) error                  // パスワードリセット
	DeleteAdministrator(ctx context.Context, in *DeleteAdministratorInput) error                                // 退会
	// コーディネータ
	ListCoordinators(ctx context.Context, in *ListCoordinatorsInput) (entity.Coordinators, int64, error)           // 一覧取得
	MultiGetCoordinators(ctx context.Context, in *MultiGetCoordinatorsInput) (entity.Coordinators, error)          // 一覧取得(ID指定)
	GetCoordinator(ctx context.Context, in *GetCoordinatorInput) (*entity.Coordinator, error)                      // １件取得
	CreateCoordinator(ctx context.Context, in *CreateCoordinatorInput) (*entity.Coordinator, error)                // 登録
	UpdateCoordinator(ctx context.Context, in *UpdateCoordinatorInput) error                                       // 更新
	UpdateCoordinatorEmail(ctx context.Context, in *UpdateCoordinatorEmailInput) error                             // メールアドレス更新
	UpdateCoordinatorThumbnails(ctx context.Context, in *UpdateCoordinatorThumbnailsInput) error                   // サムネイル(リサイズ済み)更新
	UpdateCoordinatorHeaders(ctx context.Context, in *UpdateCoordinatorHeadersInput) error                         // ヘッダー画像(リサイズ済み)更新
	ResetCoordinatorPassword(ctx context.Context, in *ResetCoordinatorPasswordInput) error                         // パスワードリセット
	RemoveCoordinatorProductType(ctx context.Context, in *RemoveCoordinatorProductTypeInput) error                 // 品目削除によるデータ変更
	AggregateRealatedProducers(ctx context.Context, in *AggregateRealatedProducersInput) (map[string]int64, error) // 担当生産者数の取得
	DeleteCoordinator(ctx context.Context, in *DeleteCoordinatorInput) error                                       // 退会
	// 生産者
	ListProducers(ctx context.Context, in *ListProducersInput) (entity.Producers, int64, error)  // 一覧取得
	MultiGetProducers(ctx context.Context, in *MultiGetProducersInput) (entity.Producers, error) // 一覧取得(ID指定)
	GetProducer(ctx context.Context, in *GetProducerInput) (*entity.Producer, error)             // １件取得
	CreateProducer(ctx context.Context, in *CreateProducerInput) (*entity.Producer, error)       // 登録
	UpdateProducer(ctx context.Context, in *UpdateProducerInput) error                           // 更新
	UpdateProducerEmail(ctx context.Context, in *UpdateProducerEmailInput) error                 // メールアドレス更新
	UpdateProducerThumbnails(ctx context.Context, in *UpdateProducerThumbnailsInput) error       // サムネイル(リサイズ済み)更新
	UpdateProducerHeaders(ctx context.Context, in *UpdateProducerHeadersInput) error             // ヘッダー画像(リサイズ済み)更新
	ResetProducerPassword(ctx context.Context, in *ResetProducerPasswordInput) error             // パスワードリセット
	RelateProducers(ctx context.Context, in *RelateProducersInput) error                         // コーディネータとの関連付け
	UnrelateProducer(ctx context.Context, in *UnrelateProducerInput) error                       // コーディネータとの関連付け解除
	DeleteProducer(ctx context.Context, in *DeleteProducerInput) error                           // 退会
	// 購入者
	SignInUser(ctx context.Context, in *SignInUserInput) (*entity.UserAuth, error)               // サインイン
	SignOutUser(ctx context.Context, in *SignOutUserInput) error                                 // サインアウト
	GetUserAuth(ctx context.Context, in *GetUserAuthInput) (*entity.UserAuth, error)             // 認証情報取得
	RefreshUserToken(ctx context.Context, in *RefreshUserTokenInput) (*entity.UserAuth, error)   // アクセストークン更新
	ListUsers(ctx context.Context, in *ListUsersInput) (entity.Users, int64, error)              // 一覧取得
	MultiGetUsers(ctx context.Context, in *MultiGetUsersInput) (entity.Users, error)             // 一覧取得(ID指定)
	MultiGetUserDevices(ctx context.Context, in *MultiGetUserDevicesInput) ([]string, error)     // デバイストークン一覧取得
	GetUser(ctx context.Context, in *GetUserInput) (*entity.User, error)                         // １件取得
	CreateUser(ctx context.Context, in *CreateUserInput) (string, error)                         // 登録 (メールアドレス/SMS認証)
	VerifyUser(ctx context.Context, in *VerifyUserInput) error                                   // 登録後の確認 (メールアドレス/SMS認証)
	CreateUserWithOAuth(ctx context.Context, in *CreateUserWithOAuthInput) (*entity.User, error) // 登録 (OAuth認証)
	InitializeUser(ctx context.Context, in *InitializeUserInput) error                           // Deprecated: 登録後の初回更新
	UpdateUserEmail(ctx context.Context, in *UpdateUserEmailInput) error                         // メールアドレス更新
	VerifyUserEmail(ctx context.Context, in *VerifyUserEmailInput) error                         // メールアドレス更新後の確認
	UpdateUserPassword(ctx context.Context, in *UpdateUserPasswordInput) error                   // パスワード更新
	ForgotUserPassword(ctx context.Context, in *ForgotUserPasswordInput) error                   // パスワードリセット (メール送信)
	VerifyUserPassword(ctx context.Context, in *VerifyUserPasswordInput) error                   // パスワードリセット (パスワード更新)
	DeleteUser(ctx context.Context, in *DeleteUserInput) error                                   // 退会
	// アドレス帳
	ListAddresses(ctx context.Context, in *ListAddressesInput) (entity.Addresses, int64, error)                      // 一覧取得
	ListDefaultAddresses(ctx context.Context, in *ListDefaultAddressesInput) (entity.Addresses, error)               // 一覧取得(デフォルト設定)
	MultiGetAddresses(ctx context.Context, in *MultiGetAddressesInput) (entity.Addresses, error)                     // 一覧取得(ID指定)
	MultiGetAddressesByRevision(ctx context.Context, in *MultiGetAddressesByRevisionInput) (entity.Addresses, error) // 一覧取得(変更履歴ID指定)
	GetAddress(ctx context.Context, in *GetAddressInput) (*entity.Address, error)                                    // １件取得
	GetDefaultAddress(ctx context.Context, in *GetDefaultAddressInput) (*entity.Address, error)                      // １件取得(デフォルト設定)
	CreateAddress(ctx context.Context, in *CreateAddressInput) (*entity.Address, error)                              // 登録
	UpdateAddress(ctx context.Context, in *UpdateAddressInput) error                                                 // 更新
	DeleteAddress(ctx context.Context, in *DeleteAddressInput) error                                                 // 削除
}
