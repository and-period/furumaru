package user

import (
	"time"

	"github.com/and-period/furumaru/api/internal/common"
)

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

type ForgotAdminPasswordInput struct {
	Email string `validate:"required,max=256,email"`
}

type VerifyAdminPasswordInput struct {
	Email                string `validate:"required,max=256,email"`
	VerifyCode           string `validate:"required"`
	NewPassword          string `validate:"min=8,max=32,password"`
	PasswordConfirmation string `validate:"required,eqfield=NewPassword"`
}

type MultiGetAdminsInput struct {
	AdminIDs []string `validate:"dive,required"`
}

type MultiGetAdminDevicesInput struct {
	AdminIDs []string `validate:"dive,required"`
}

type GetAdminInput struct {
	AdminID string `validate:"required"`
}

type ListAdministratorsInput struct {
	Limit  int64 `validate:"required,max=200"`
	Offset int64 `validate:"min=0"`
}

type MultiGetAdministratorsInput struct {
	AdministratorIDs []string `validate:"dive,required"`
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
	PhoneNumber   string `validate:"required,e164"`
}

type UpdateAdministratorInput struct {
	AdministratorID string `validate:"required"`
	Lastname        string `validate:"required,max=16"`
	Firstname       string `validate:"required,max=16"`
	LastnameKana    string `validate:"required,max=32,hiragana"`
	FirstnameKana   string `validate:"required,max=32,hiragana"`
	PhoneNumber     string `validate:"required,e164"`
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
	Name   string `validate:"max=64"`
	Limit  int64  `validate:"required,max=200"`
	Offset int64  `validate:"min=0"`
}

type MultiGetCoordinatorsInput struct {
	CoordinatorIDs []string `validate:"dive,required"`
	WithDeleted    bool     `validate:""`
}

type GetCoordinatorInput struct {
	CoordinatorID string `validate:"required"`
	WithDeleted   bool   `validate:""`
}

type CreateCoordinatorInput struct {
	Lastname          string         `validate:"required,max=16"`
	Firstname         string         `validate:"required,max=16"`
	LastnameKana      string         `validate:"required,max=32,hiragana"`
	FirstnameKana     string         `validate:"required,max=32,hiragana"`
	MarcheName        string         `validate:"required,max=64"`
	Username          string         `validate:"required,max=64"`
	Profile           string         `validate:"max=2000"`
	ProductTypeIDs    []string       `validate:"dive,required"`
	ThumbnailURL      string         `validate:"omitempty,url"`
	HeaderURL         string         `validate:"omitempty,url"`
	PromotionVideoURL string         `validate:"omitempty,url"`
	BonusVideoURL     string         `validate:"omitempty,url"`
	InstagramID       string         `validate:"max=30"`
	FacebookID        string         `validate:"max=50"`
	Email             string         `validate:"required,max=256,email"`
	PhoneNumber       string         `validate:"required,e164"`
	PostalCode        string         `validate:"max=16,numeric"`
	PrefectureCode    int32          `validate:"required"`
	City              string         `validate:"max=32"`
	AddressLine1      string         `validate:"max=64"`
	AddressLine2      string         `validate:"max=64"`
	BusinessDays      []time.Weekday `validate:"max=7,unique"`
}

type UpdateCoordinatorInput struct {
	CoordinatorID     string         `validate:"required"`
	Lastname          string         `validate:"required,max=16"`
	Firstname         string         `validate:"required,max=16"`
	LastnameKana      string         `validate:"required,max=32,hiragana"`
	FirstnameKana     string         `validate:"required,max=32,hiragana"`
	MarcheName        string         `validate:"required,max=64"`
	Username          string         `validate:"required,max=64"`
	Profile           string         `validate:"max=2000"`
	ProductTypeIDs    []string       `validate:"dive,required"`
	ThumbnailURL      string         `validate:"omitempty,url"`
	HeaderURL         string         `validate:"omitempty,url"`
	PromotionVideoURL string         `validate:"omitempty,url"`
	BonusVideoURL     string         `validate:"omitempty,url"`
	InstagramID       string         `validate:"max=30"`
	FacebookID        string         `validate:"max=50"`
	PhoneNumber       string         `validate:"required,e164"`
	PostalCode        string         `validate:"max=16,numeric"`
	PrefectureCode    int32          `validate:"required"`
	City              string         `validate:"max=32"`
	AddressLine1      string         `validate:"max=64"`
	AddressLine2      string         `validate:"max=64"`
	BusinessDays      []time.Weekday `validate:"max=7,unique"`
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

type RemoveCoordinatorProductTypeInput struct {
	ProductTypeID string `validate:"required"`
}

type DeleteCoordinatorInput struct {
	CoordinatorID string `validate:"required"`
}

type AggregateRealatedProducersInput struct {
	CoordinatorIDs []string `validate:"dive,required"`
}

type ListProducersInput struct {
	CoordinatorID string `validate:""`
	Name          string `validate:"max=64"`
	Limit         int64  `validate:"required_without=CoordinatorID,max=200"`
	Offset        int64  `validate:"min=0"`
}

type MultiGetProducersInput struct {
	ProducerIDs []string `validate:"dive,required"`
	WithDeleted bool     `validate:""`
}

type GetProducerInput struct {
	ProducerID  string `validate:"required"`
	WithDeleted bool   `validate:""`
}

type CreateProducerInput struct {
	CoordinatorID     string `validate:"required"`
	Lastname          string `validate:"required,max=16"`
	Firstname         string `validate:"required,max=16"`
	LastnameKana      string `validate:"required,max=32,hiragana"`
	FirstnameKana     string `validate:"required,max=32,hiragana"`
	Username          string `validate:"required,max=64"`
	Profile           string `validate:"max=2000"`
	ThumbnailURL      string `validate:"omitempty,url"`
	HeaderURL         string `validate:"omitempty,url"`
	PromotionVideoURL string `validate:"omitempty,url"`
	BonusVideoURL     string `validate:"omitempty,url"`
	InstagramID       string `validate:"max=30"`
	FacebookID        string `validate:"max=50"`
	Email             string `validate:"omitempty,max=256,email"`
	PhoneNumber       string `validate:"omitempty,e164"`
	PostalCode        string `validate:"omitempty,max=16,numeric"`
	PrefectureCode    int32  `validate:"min=0"`
	City              string `validate:"max=32"`
	AddressLine1      string `validate:"max=64"`
	AddressLine2      string `validate:"max=64"`
}

type UpdateProducerInput struct {
	ProducerID        string `validate:"required"`
	Lastname          string `validate:"required,max=16"`
	Firstname         string `validate:"required,max=16"`
	LastnameKana      string `validate:"required,max=32,hiragana"`
	FirstnameKana     string `validate:"required,max=32,hiragana"`
	Username          string `validate:"required,max=64"`
	Profile           string `validate:"max=2000"`
	ThumbnailURL      string `validate:"omitempty,url"`
	HeaderURL         string `validate:"omitempty,url"`
	PromotionVideoURL string `validate:"omitempty,url"`
	BonusVideoURL     string `validate:"omitempty,url"`
	InstagramID       string `validate:"max=30"`
	FacebookID        string `validate:"max=50"`
	Email             string `validate:"omitempty,max=256,email"`
	PhoneNumber       string `validate:"omitempty,e164"`
	PostalCode        string `validate:"omitempty,max=16,numeric"`
	PrefectureCode    int32  `validate:"min=0"`
	City              string `validate:"max=32"`
	AddressLine1      string `validate:"max=64"`
	AddressLine2      string `validate:"max=64"`
}

type UpdateProducerThumbnailsInput struct {
	ProducerID string        `validate:"required"`
	Thumbnails common.Images `validate:""`
}

type UpdateProducerHeadersInput struct {
	ProducerID string        `validate:"required"`
	Headers    common.Images `validate:""`
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
	Limit          int64 `validate:"required,max=200"`
	Offset         int64 `validate:"min=0"`
	OnlyRegistered bool  `validate:""`
	OnlyVerified   bool  `validate:""`
}

type MultiGetUsersInput struct {
	UserIDs []string `validate:"dive,required"`
}

type MultiGetUserDevicesInput struct {
	UserIDs []string `validate:"dive,required"`
}

type GetUserInput struct {
	UserID string `validate:"required"`
}

type CreateUserInput struct {
	Username             string `validate:"required,max=32"`
	AccountID            string `validate:"required,max=32"`
	Lastname             string `validate:"required,max=16"`
	Firstname            string `validate:"required,max=16"`
	LastnameKana         string `validate:"required,max=32,hiragana"`
	FirstnameKana        string `validate:"required,max=32,hiragana"`
	Email                string `validate:"required,max=256,email"`
	PhoneNumber          string `validate:"required,e164"`
	Password             string `validate:"min=8,max=32,password"`
	PasswordConfirmation string `validate:"required,eqfield=Password"`
}

type VerifyUserInput struct {
	UserID     string `validate:"required"`
	VerifyCode string `validate:"required"`
}

type CreateUserWithOAuthInput struct {
	AccessToken   string `validate:"required"`
	Username      string `validate:"required,max=32"`
	AccountID     string `validate:"required,max=32"`
	Lastname      string `validate:"required,max=16"`
	Firstname     string `validate:"required,max=16"`
	LastnameKana  string `validate:"required,max=32,hiragana"`
	FirstnameKana string `validate:"required,max=32,hiragana"`
	PhoneNumber   string `validate:"required,e164"`
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

type UpdateUserThumbnailsInput struct {
	UserID     string        `validate:"required"`
	Thumbnails common.Images `validate:""`
}

type DeleteUserInput struct {
	UserID string `validate:"required"`
}

type ListAddressesInput struct {
	UserID string `validate:"required"`
	Limit  int64  `validate:"min=1,max=200"`
	Offset int64  `validate:"min=0"`
}

type ListDefaultAddressesInput struct {
	UserIDs []string `validate:"dive,required"`
}

type MultiGetAddressesInput struct {
	AddressIDs []string `validate:"dive,required"`
}

type MultiGetAddressesByRevisionInput struct {
	AddressRevisionIDs []int64 `validate:"dive,required"`
}

type GetAddressInput struct {
	AddressID string `validate:"required"`
	UserID    string `validate:"required"`
}

type GetDefaultAddressInput struct {
	UserID string `validate:"required"`
}

type CreateAddressInput struct {
	UserID         string `validate:"required"`
	Lastname       string `validate:"required,max=16"`
	Firstname      string `validate:"required,max=16"`
	LastnameKana   string `validate:"required,max=32,hiragana"`
	FirstnameKana  string `validate:"required,max=32,hiragana"`
	PostalCode     string `validate:"required,max=16,numeric"`
	PrefectureCode int32  `validate:"required"`
	City           string `validate:"required,max=32"`
	AddressLine1   string `validate:"required,max=64"`
	AddressLine2   string `validate:"max=64"`
	PhoneNumber    string `validate:"required,e164"`
	IsDefault      bool   `validate:""`
}

type UpdateAddressInput struct {
	AddressID      string `validate:"required"`
	UserID         string `validate:"required"`
	Lastname       string `validate:"required,max=16"`
	Firstname      string `validate:"required,max=16"`
	LastnameKana   string `validate:"required,max=32,hiragana"`
	FirstnameKana  string `validate:"required,max=32,hiragana"`
	PostalCode     string `validate:"required,max=16,numeric"`
	PrefectureCode int32  `validate:"required"`
	City           string `validate:"required,max=32"`
	AddressLine1   string `validate:"required,max=64"`
	AddressLine2   string `validate:"max=64"`
	PhoneNumber    string `validate:"required,e164"`
	IsDefault      bool   `validate:""`
}

type DeleteAddressInput struct {
	AddressID string `validate:"required"`
	UserID    string `validate:"required"`
}
