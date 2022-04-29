package service

import (
	"regexp"

	validator "github.com/go-playground/validator/v10"
)

const (
	passwordString    = "^[a-zA-Z0-9_!@#$_%^&*.?()-=+]*$"
	phoneNumberString = "^\\+[0-9]{11,17}$"
)

var (
	passwordRegex    = regexp.MustCompile(passwordString)
	phoneNumberRegex = regexp.MustCompile(phoneNumberString)
)

func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return passwordRegex.MatchString(fl.Field().String())
	})
	v.RegisterValidation("phone_number", func(fl validator.FieldLevel) bool {
		return phoneNumberRegex.MatchString(fl.Field().String())
	})
	return v
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
