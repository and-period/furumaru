package user

import "github.com/and-period/furumaru/api/internal/common"

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

type RegisterAdminDeviceInput struct {
	AdminID string `validate:"required"`
	Device  string `validate:"required"`
}

type RefreshAdminTokenInput struct {
	RefreshToken string `validate:"required"`
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

type MultiGetAdminsInput struct {
	AdminIDs []string `validate:"omitempty,dive,required"`
}

type MultiGetAdminDevicesInput struct {
	AdminIDs []string `validate:"omitempty,dive,required"`
}

type GetAdminInput struct {
	AdminID string `validate:"required"`
}

type ListAdministratorsInput struct {
	Limit  int64 `validate:"required,max=200"`
	Offset int64 `validate:"min=0"`
}

type MultiGetAdministratorsInput struct {
	AdministratorIDs []string `validate:"omitempty,dive,required"`
}

type GetAdministratorInput struct {
	AdministratorID string `validate:"required"`
}

type CreateAdministratorInput struct {
	Lastname      string `validate:"required,max=16"`
	Firstname     string `validate:"required,max=16"`
	LastnameKana  string `validate:"required,max=32,hiragana"`
	FirstnameKana string `validate:"required,max=32,hiragana"`
	Email         string `validate:"required,max=256,email"`
	PhoneNumber   string `validate:"min=12,max=18,phone_number"`
}

type UpdateAdministratorInput struct {
	AdministratorID string `validate:"required"`
	Lastname        string `validate:"required,max=16"`
	Firstname       string `validate:"required,max=16"`
	LastnameKana    string `validate:"required,max=32,hiragana"`
	FirstnameKana   string `validate:"required,max=32,hiragana"`
	PhoneNumber     string `validate:"min=12,max=18,phone_number"`
}

type UpdateAdministratorEmailInput struct {
	AdministratorID string `validate:"required"`
	Email           string `validate:"required,max=256,email"`
}

type ResetAdministratorPasswordInput struct {
	AdministratorID string `validate:"required"`
}

type DeleteAdministratorInput struct {
	AdministratorID string `validate:"required"`
}

type ListCoordinatorsInput struct {
	Limit  int64 `validate:"required,max=200"`
	Offset int64 `validate:"min=0"`
}

type MultiGetCoordinatorsInput struct {
	CoordinatorIDs []string `validate:"omitempty,dive,required"`
}

type GetCoordinatorInput struct {
	CoordinatorID string `validate:"required"`
}

type CreateCoordinatorInput struct {
	Lastname         string `validate:"required,max=16"`
	Firstname        string `validate:"required,max=16"`
	LastnameKana     string `validate:"required,max=32,hiragana"`
	FirstnameKana    string `validate:"required,max=32,hiragana"`
	CompanyName      string `validate:"required,max=64"`
	StoreName        string `validate:"required,max=64"`
	ThumbnailURL     string `validate:"omitempty,url"`
	HeaderURL        string `validate:"omitempty,url"`
	TwitterAccount   string `validate:"omitempty,max=15"`
	InstagramAccount string `validate:"omitempty,max=30"`
	FacebookAccount  string `validate:"omitempty,max=50"`
	Email            string `validate:"required,max=256,email"`
	PhoneNumber      string `validate:"min=12,max=18,phone_number"`
	PostalCode       string `validate:"omitempty,max=16,numeric"`
	Prefecture       string `validate:"omitempty,max=32"`
	City             string `validate:"omitempty,max=32"`
	AddressLine1     string `validate:"omitempty,max=64"`
	AddressLine2     string `validate:"omitempty,max=64"`
}

type UpdateCoordinatorInput struct {
	CoordinatorID    string `validate:"required"`
	Lastname         string `validate:"required,max=16"`
	Firstname        string `validate:"required,max=16"`
	LastnameKana     string `validate:"required,max=32,hiragana"`
	FirstnameKana    string `validate:"required,max=32,hiragana"`
	CompanyName      string `validate:"required,max=64"`
	StoreName        string `validate:"required,max=64"`
	ThumbnailURL     string `validate:"omitempty,url"`
	HeaderURL        string `validate:"omitempty,url"`
	TwitterAccount   string `validate:"omitempty,max=15"`
	InstagramAccount string `validate:"omitempty,max=30"`
	FacebookAccount  string `validate:"omitempty,max=50"`
	PhoneNumber      string `validate:"min=12,max=18,phone_number"`
	PostalCode       string `validate:"omitempty,max=16,numeric"`
	Prefecture       string `validate:"omitempty,max=32"`
	City             string `validate:"omitempty,max=32"`
	AddressLine1     string `validate:"omitempty,max=64"`
	AddressLine2     string `validate:"omitempty,max=64"`
}

type UpdateCoordinatorEmailInput struct {
	CoordinatorID string `validate:"required"`
	Email         string `validate:"required,max=256,email"`
}

type UpdateCoordinatorThumbnailsInput struct {
	CoordinatorID string        `validate:"required"`
	Thumbnails    common.Images `validate:""`
}

type UpdateCoordinatorHeadersInput struct {
	CoordinatorID string        `validate:"required"`
	Headers       common.Images `validate:""`
}

type ResetCoordinatorPasswordInput struct {
	CoordinatorID string `validate:"required"`
}

type DeleteCoordinatorInput struct {
	CoordinatorID string `validate:"required"`
}

type ListProducersInput struct {
	CoordinatorID string `validate:"omitempty"`
	Limit         int64  `validate:"required_without=CoordinatorID,max=200"`
	Offset        int64  `validate:"min=0"`
	OnlyUnrelated bool   `validate:""`
}

type MultiGetProducersInput struct {
	ProducerIDs []string `validate:"omitempty,dive,required"`
}

type GetProducerInput struct {
	ProducerID string `validate:"required"`
}

type CreateProducerInput struct {
	Lastname      string `validate:"required,max=16"`
	Firstname     string `validate:"required,max=16"`
	LastnameKana  string `validate:"required,max=32,hiragana"`
	FirstnameKana string `validate:"required,max=32,hiragana"`
	StoreName     string `validate:"required,max=64"`
	ThumbnailURL  string `validate:"omitempty,url"`
	HeaderURL     string `validate:"omitempty,url"`
	Email         string `validate:"required,max=256,email"`
	PhoneNumber   string `validate:"min=12,max=18,phone_number"`
	PostalCode    string `validate:"omitempty,max=16,numeric"`
	Prefecture    string `validate:"omitempty,max=32"`
	City          string `validate:"omitempty,max=32"`
	AddressLine1  string `validate:"omitempty,max=64"`
	AddressLine2  string `validate:"omitempty,max=64"`
}

type UpdateProducerInput struct {
	ProducerID    string `validate:"required"`
	Lastname      string `validate:"required,max=16"`
	Firstname     string `validate:"required,max=16"`
	LastnameKana  string `validate:"required,max=32,hiragana"`
	FirstnameKana string `validate:"required,max=32,hiragana"`
	StoreName     string `validate:"required,max=64"`
	ThumbnailURL  string `validate:"omitempty,url"`
	HeaderURL     string `validate:"omitempty,url"`
	PhoneNumber   string `validate:"min=12,max=18,phone_number"`
	PostalCode    string `validate:"omitempty,max=16,numeric"`
	Prefecture    string `validate:"omitempty,max=32"`
	City          string `validate:"omitempty,max=32"`
	AddressLine1  string `validate:"omitempty,max=64"`
	AddressLine2  string `validate:"omitempty,max=64"`
}

type UpdateProducerEmailInput struct {
	ProducerID string `validate:"required"`
	Email      string `validate:"required,max=256,email"`
}

type UpdateProducerThumbnailsInput struct {
	ProducerID string        `validate:"required"`
	Thumbnails common.Images `validate:""`
}

type UpdateProducerHeadersInput struct {
	ProducerID string        `validate:"required"`
	Headers    common.Images `validate:""`
}

type ResetProducerPasswordInput struct {
	ProducerID string `validate:"required"`
}

type RelateProducersInput struct {
	CoordinatorID string   `validate:"required"`
	ProducerIDs   []string `validate:"min=1,max=20,dive,required"`
}

type UnrelateProducerInput struct {
	ProducerID string `validate:"required"`
}

type DeleteProducerInput struct {
	ProducerID string `validate:"required"`
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

type ListUsersInput struct {
	Limit  int64 `validate:"required,max=200"`
	Offset int64 `validate:"min=0"`
}

type MultiGetUsersInput struct {
	UserIDs []string `validate:"omitempty,dive,required"`
}

type MultiGetUserDevicesInput struct {
	UserIDs []string `validate:"omitempty,dive,required"`
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
