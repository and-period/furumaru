package service

import "github.com/and-period/marche/api/internal/user/entity"

type SignInAdminInput struct {
	Key      string `validate:"required"`
	Password string `validate:"required"`
}

type SignOutAdminInput struct {
	AccessToken string `validate:"required"`
}

type GetAdminAuthInput struct {
	AccessToken string `validate:"required"`
}

type RefreshAdminTokenInput struct {
	RefreshToken string `validate:"required"`
}

type ListAdminsInput struct {
	Roles  []entity.AdminRole `validate:""`
	Limit  int64              `validate:"required,max=200"`
	Offset int64              `validate:"min=0"`
}

type MultiGetAdminsInput struct {
	AdminIDs []string `validate:"min=1,dive,required"`
}

type GetAdminInput struct {
	AdminID string `validate:"required"`
}

type CreateAdminInput struct {
	Lastname      string           `validate:"required,max=16"`
	Firstname     string           `validate:"required,max=16"`
	LastnameKana  string           `validate:"required,max=32,hiragana"`
	FirstnameKana string           `validate:"required,max=32,hiragana"`
	Email         string           `validate:"required,max=256,email"`
	Role          entity.AdminRole `validate:""`
}

type UpdateAdminEmailInput struct {
	AccessToken string `validate:"required"`
	Email       string `validate:"required,max=256,email"`
}

type VerifyAdminEmailInput struct {
	AccessToken string `validate:"required"`
	VerifyCode  string `validate:"required"`
}

type UpdateAdminPasswordInput struct {
	AccessToken          string `validate:"required"`
	OldPassword          string `validate:"required"`
	NewPassword          string `validate:"min=8,max=32,password"`
	PasswordConfirmation string `validate:"required,eqfield=NewPassword"`
}

type SignInUserInput struct {
	Key      string `validate:"required"`
	Password string `validate:"required"`
}

type SignOutUserInput struct {
	AccessToken string `validate:"required"`
}

type GetUserAuthInput struct {
	AccessToken string `validate:"required"`
}

type RefreshUserTokenInput struct {
	RefreshToken string `validate:"required"`
}

type GetUserInput struct {
	UserID string `validate:"required"`
}

type CreateUserInput struct {
	Email                string `validate:"required,max=256,email"`
	PhoneNumber          string `validate:"min=12,max=18,phone_number"`
	Password             string `validate:"min=8,max=32,password"`
	PasswordConfirmation string `validate:"required,eqfield=Password"`
}

type VerifyUserInput struct {
	UserID     string `validate:"required"`
	VerifyCode string `validate:"required"`
}

type InitializeUserInput struct {
	UserID    string `validate:"required"`
	Username  string `validate:"required,max=32"`
	AccountID string `validate:"required,max=32"`
}

type CreateUserWithOAuthInput struct {
	AccessToken string `validate:"required"`
}

type UpdateUserEmailInput struct {
	AccessToken string `validate:"required"`
	Email       string `validate:"required,max=256,email"`
}

type VerifyUserEmailInput struct {
	AccessToken string `validate:"required"`
	VerifyCode  string `validate:"required"`
}

type UpdateUserPasswordInput struct {
	AccessToken          string `validate:"required"`
	OldPassword          string `validate:"required"`
	NewPassword          string `validate:"min=8,max=32,password"`
	PasswordConfirmation string `validate:"required,eqfield=NewPassword"`
}

type ForgotUserPasswordInput struct {
	Email string `validate:"required,max=256,email"`
}

type VerifyUserPasswordInput struct {
	Email                string `validate:"required,max=256,email"`
	VerifyCode           string `validate:"required"`
	NewPassword          string `validate:"min=8,max=32,password"`
	PasswordConfirmation string `validate:"required,eqfield=NewPassword"`
}

type DeleteUserInput struct {
	UserID string `validate:"required"`
}
