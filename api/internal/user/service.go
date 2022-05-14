//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package user

import (
	"context"

	"github.com/and-period/marche/api/internal/user/entity"
)

//nolint:revive
type UserService interface {
	SignInAdmin(ctx context.Context, in *SignInAdminInput) (*entity.AdminAuth, error)
	SignOutAdmin(ctx context.Context, in *SignOutAdminInput) error
	GetAdminAuth(ctx context.Context, in *GetAdminAuthInput) (*entity.AdminAuth, error)
	RefreshAdminToken(ctx context.Context, in *RefreshAdminTokenInput) (*entity.AdminAuth, error)
	ListAdmins(ctx context.Context, in *ListAdminsInput) (entity.Admins, error)
	MultiGetAdmins(ctx context.Context, in *MultiGetAdminsInput) (entity.Admins, error)
	GetAdmin(ctx context.Context, in *GetAdminInput) (*entity.Admin, error)
	CreateAdmin(ctx context.Context, in *CreateAdminInput) (*entity.Admin, error)
	UpdateAdminEmail(ctx context.Context, in *UpdateAdminEmailInput) error
	VerifyAdminEmail(ctx context.Context, in *VerifyAdminEmailInput) error
	UpdateAdminPassword(ctx context.Context, in *UpdateAdminPasswordInput) error
	SignInUser(ctx context.Context, in *SignInUserInput) (*entity.UserAuth, error)
	SignOutUser(ctx context.Context, in *SignOutUserInput) error
	GetUserAuth(ctx context.Context, in *GetUserAuthInput) (*entity.UserAuth, error)
	RefreshUserToken(ctx context.Context, in *RefreshUserTokenInput) (*entity.UserAuth, error)
	GetUser(ctx context.Context, in *GetUserInput) (*entity.User, error)
	CreateUser(ctx context.Context, in *CreateUserInput) (string, error)
	VerifyUser(ctx context.Context, in *VerifyUserInput) error
	CreateUserWithOAuth(ctx context.Context, in *CreateUserWithOAuthInput) (*entity.User, error)
	InitializeUser(ctx context.Context, in *InitializeUserInput) error
	UpdateUserEmail(ctx context.Context, in *UpdateUserEmailInput) error
	VerifyUserEmail(ctx context.Context, in *VerifyUserEmailInput) error
	UpdateUserPassword(ctx context.Context, in *UpdateUserPasswordInput) error
	ForgotUserPassword(ctx context.Context, in *ForgotUserPasswordInput) error
	VerifyUserPassword(ctx context.Context, in *VerifyUserPasswordInput) error
	DeleteUser(ctx context.Context, in *DeleteUserInput) error
}
