//go:generate go tool mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package user

import (
	"context"

	"github.com/and-period/furumaru/api/internal/user/entity"
)

type Service interface {
	// Addresss - アドレス帳
	ListAddresses(ctx context.Context, in *ListAddressesInput) (entity.Addresses, int64, error)                      // 一覧取得
	ListDefaultAddresses(ctx context.Context, in *ListDefaultAddressesInput) (entity.Addresses, error)               // 一覧取得(デフォルト設定)
	MultiGetAddresses(ctx context.Context, in *MultiGetAddressesInput) (entity.Addresses, error)                     // 一覧取得(ID指定)
	MultiGetAddressesByRevision(ctx context.Context, in *MultiGetAddressesByRevisionInput) (entity.Addresses, error) // 一覧取得(変更履歴ID指定)
	GetAddress(ctx context.Context, in *GetAddressInput) (*entity.Address, error)                                    // １件取得
	GetDefaultAddress(ctx context.Context, in *GetDefaultAddressInput) (*entity.Address, error)                      // １件取得(デフォルト設定)
	CreateAddress(ctx context.Context, in *CreateAddressInput) (*entity.Address, error)                              // 登録
	UpdateAddress(ctx context.Context, in *UpdateAddressInput) error                                                 // 更新
	DeleteAddress(ctx context.Context, in *DeleteAddressInput) error                                                 // 削除
	// Admin - 管理者
	SignInAdmin(ctx context.Context, in *SignInAdminInput) (*entity.AdminAuth, error)                               // サインイン
	SignOutAdmin(ctx context.Context, in *SignOutAdminInput) error                                                  // サインアウト
	GetAdminAuth(ctx context.Context, in *GetAdminAuthInput) (*entity.AdminAuth, error)                             // 認証情報取得
	RefreshAdminToken(ctx context.Context, in *RefreshAdminTokenInput) (*entity.AdminAuth, error)                   // アクセストークンの更新
	RegisterAdminDevice(ctx context.Context, in *RegisterAdminDeviceInput) error                                    // デバイストークンの更新
	UpdateAdminEmail(ctx context.Context, in *UpdateAdminEmailInput) error                                          // メールアドレス更新
	VerifyAdminEmail(ctx context.Context, in *VerifyAdminEmailInput) error                                          // メールアドレス更新後の確認
	UpdateAdminPassword(ctx context.Context, in *UpdateAdminPasswordInput) error                                    // パスワード更新
	ForgotAdminPassword(ctx context.Context, in *ForgotAdminPasswordInput) error                                    // パスワードリセット (メール送信)
	VerifyAdminPassword(ctx context.Context, in *VerifyAdminPasswordInput) error                                    // パスワードリセット (パスワード更新)
	MultiGetAdmins(ctx context.Context, in *MultiGetAdminsInput) (entity.Admins, error)                             // 一覧取得(ID指定)
	MultiGetAdminDevices(ctx context.Context, in *MultiGetAdminDevicesInput) ([]string, error)                      // デバイストークン一覧取得
	GetAdmin(ctx context.Context, in *GetAdminInput) (*entity.Admin, error)                                         // １件取得
	ListAdminAuthProviders(ctx context.Context, in *ListAdminAuthProvidersInput) (entity.AdminAuthProviders, error) // 認証プロバイダ一覧取得
	InitialGoogleAdminAuth(ctx context.Context, in *InitialGoogleAdminAuthInput) (string, error)                    // Google認証開始
	ConnectGoogleAdminAuth(ctx context.Context, in *ConnectGoogleAdminAuthInput) error                              // Google認証連携
	InitialLINEAdminAuth(ctx context.Context, in *InitialLINEAdminAuthInput) (string, error)                        // LINE認証開始
	ConnectLINEAdminAuth(ctx context.Context, in *ConnectLINEAdminAuthInput) error                                  // LINE認証連携
	// AdminRole - 管理者ロール
	GenerateAdminRole(ctx context.Context, in *GenerateAdminRoleInput) (string, string, error) // ロール生成
	// Administrator - システム管理者
	ListAdministrators(ctx context.Context, in *ListAdministratorsInput) (entity.Administrators, int64, error)  // 一覧取得
	MultiGetAdministrators(ctx context.Context, in *MultiGetAdministratorsInput) (entity.Administrators, error) // 一覧取得(ID指定)
	GetAdministrator(ctx context.Context, in *GetAdministratorInput) (*entity.Administrator, error)             // １件取得
	CreateAdministrator(ctx context.Context, in *CreateAdministratorInput) (*entity.Administrator, error)       // 登録
	UpdateAdministrator(ctx context.Context, in *UpdateAdministratorInput) error                                // 更新
	UpdateAdministratorEmail(ctx context.Context, in *UpdateAdministratorEmailInput) error                      // メールアドレス更新
	ResetAdministratorPassword(ctx context.Context, in *ResetAdministratorPasswordInput) error                  // パスワードリセット
	DeleteAdministrator(ctx context.Context, in *DeleteAdministratorInput) error                                // 退会
	// Coordinator - コーディネータ
	ListCoordinators(ctx context.Context, in *ListCoordinatorsInput) (entity.Coordinators, int64, error)           // 一覧取得
	MultiGetCoordinators(ctx context.Context, in *MultiGetCoordinatorsInput) (entity.Coordinators, error)          // 一覧取得(ID指定)
	GetCoordinator(ctx context.Context, in *GetCoordinatorInput) (*entity.Coordinator, error)                      // １件取得
	CreateCoordinator(ctx context.Context, in *CreateCoordinatorInput) (*entity.Coordinator, string, error)        // 登録
	UpdateCoordinator(ctx context.Context, in *UpdateCoordinatorInput) error                                       // 更新
	UpdateCoordinatorEmail(ctx context.Context, in *UpdateCoordinatorEmailInput) error                             // メールアドレス更新
	ResetCoordinatorPassword(ctx context.Context, in *ResetCoordinatorPasswordInput) error                         // パスワードリセット
	AggregateRealatedProducers(ctx context.Context, in *AggregateRealatedProducersInput) (map[string]int64, error) // 担当生産者数の取得
	DeleteCoordinator(ctx context.Context, in *DeleteCoordinatorInput) error                                       // 退会
	// Guest - ゲスト
	UpsertGuest(ctx context.Context, in *UpsertGuestInput) (string, error) // ゲスト登録・更新
	// Member - 会員
	CreateMember(ctx context.Context, in *CreateMemberInput) (string, error)                           // 登録 (メールアドレス/SMS認証)
	VerifyMember(ctx context.Context, in *VerifyMemberInput) error                                     // 登録後の確認 (メールアドレス/SMS認証)
	UpdateMemberEmail(ctx context.Context, in *UpdateMemberEmailInput) error                           // メールアドレス更新
	VerifyMemberEmail(ctx context.Context, in *VerifyMemberEmailInput) error                           // メールアドレス更新後の確認
	UpdateMemberPassword(ctx context.Context, in *UpdateMemberPasswordInput) error                     // パスワード更新
	ForgotMemberPassword(ctx context.Context, in *ForgotMemberPasswordInput) error                     // パスワードリセット (メール送信)
	VerifyMemberPassword(ctx context.Context, in *VerifyMemberPasswordInput) error                     // パスワードリセット (パスワード更新)
	UpdateMemberUsername(ctx context.Context, in *UpdateMemberUsernameInput) error                     // 表示名更新
	UpdateMemberAccountID(ctx context.Context, in *UpdateMemberAccountIDInput) error                   // 検索名更新
	UpdateMemberThumbnailURL(ctx context.Context, in *UpdateMemberThumbnailURLInput) error             // サムネイルURL更新
	AuthMemberWithGoogle(ctx context.Context, in *AuthMemberWithGoogleInput) (string, error)           // 認証開始（Google認証）
	AuthMemberWithLINE(ctx context.Context, in *AuthMemberWithLINEInput) (string, error)               // 認証開始（LINE認証）
	CreateMemberWithGoogle(ctx context.Context, in *CreateMemberWithGoogleInput) (*entity.User, error) // 登録（Google認証）
	CreateMemberWithLINE(ctx context.Context, in *CreateMemberWithLINEInput) (*entity.User, error)     // 登録（LINE認証）
	// Producer - 生産者
	ListProducers(ctx context.Context, in *ListProducersInput) (entity.Producers, int64, error)  // 一覧取得
	MultiGetProducers(ctx context.Context, in *MultiGetProducersInput) (entity.Producers, error) // 一覧取得(ID指定)
	GetProducer(ctx context.Context, in *GetProducerInput) (*entity.Producer, error)             // １件取得
	CreateProducer(ctx context.Context, in *CreateProducerInput) (*entity.Producer, error)       // 登録
	UpdateProducer(ctx context.Context, in *UpdateProducerInput) error                           // 更新
	DeleteProducer(ctx context.Context, in *DeleteProducerInput) error                           // 退会
	// User - 購入者
	SignInUser(ctx context.Context, in *SignInUserInput) (*entity.UserAuth, error)             // サインイン
	SignOutUser(ctx context.Context, in *SignOutUserInput) error                               // サインアウト
	GetUserAuth(ctx context.Context, in *GetUserAuthInput) (*entity.UserAuth, error)           // 認証情報取得
	RefreshUserToken(ctx context.Context, in *RefreshUserTokenInput) (*entity.UserAuth, error) // アクセストークン更新
	ListUsers(ctx context.Context, in *ListUsersInput) (entity.Users, int64, error)            // 一覧取得
	MultiGetUsers(ctx context.Context, in *MultiGetUsersInput) (entity.Users, error)           // 一覧取得(ID指定)
	MultiGetUserDevices(ctx context.Context, in *MultiGetUserDevicesInput) ([]string, error)   // デバイストークン一覧取得
	GetUser(ctx context.Context, in *GetUserInput) (*entity.User, error)                       // １件取得
	DeleteUser(ctx context.Context, in *DeleteUserInput) error                                 // 退会
	// UserNotification - 購入者通知設定
	MultiGetUserNotifications(ctx context.Context, in *MultiGetUserNotificationsInput) (entity.UserNotifications, error) // 一覧取得(ID指定)
	GetUserNotification(ctx context.Context, in *GetUserNotificationInput) (*entity.UserNotification, error)             // １件取得
	UpdateUserNotification(ctx context.Context, in *UpdateUserNotificationInput) error                                   // 更新
}
