//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/user/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"errors"
	"time"

	"github.com/and-period/furumaru/api/internal/user/entity"
)

var (
	ErrInvalidArgument    = &Error{err: errors.New("database: invalid argument")}
	ErrNotFound           = &Error{err: errors.New("database: not found")}
	ErrAlreadyExists      = &Error{err: errors.New("database: already exists")}
	ErrFailedPrecondition = &Error{err: errors.New("database: failed precondition")}
	ErrCanceled           = &Error{err: errors.New("database: canceled")}
	ErrDeadlineExceeded   = &Error{err: errors.New("database: deadline exceeded")}
	ErrInternal           = &Error{err: errors.New("database: internal error")}
	ErrUnknown            = &Error{err: errors.New("database: unknown")}
)

type Database struct {
	Address          Address
	Admin            Admin
	AdminGroup       AdminGroup
	AdminGroupRole   AdminGroupRole
	AdminGroupUser   AdminGroupUser
	AdminPolicy      AdminPolicy
	AdminRole        AdminRole
	AdminRolePolicy  AdminRolePolicy
	Administrator    Administrator
	Coordinator      Coordinator
	Guest            Guest
	Member           Member
	Producer         Producer
	User             User
	UserNotification UserNotification
}

/**
 * interface
 */
type Address interface {
	List(ctx context.Context, params *ListAddressesParams, fields ...string) (entity.Addresses, error)
	ListDefault(ctx context.Context, userIDs []string, fields ...string) (entity.Addresses, error)
	Count(ctx context.Context, params *ListAddressesParams) (int64, error)
	MultiGet(ctx context.Context, addressIDs []string, fields ...string) (entity.Addresses, error)
	MultiGetByRevision(ctx context.Context, revisionIDs []int64, fields ...string) (entity.Addresses, error)
	Get(ctx context.Context, addressID string, fields ...string) (*entity.Address, error)
	GetDefault(ctx context.Context, userID string, fields ...string) (*entity.Address, error)
	Create(ctx context.Context, address *entity.Address) error
	Update(ctx context.Context, addressID, userID string, params *UpdateAddressParams) error
	Delete(ctx context.Context, addressID, userID string) error
}

type ListAddressesParams struct {
	UserID string
	Limit  int
	Offset int
}

type UpdateAddressParams struct {
	Lastname       string
	Firstname      string
	LastnameKana   string
	FirstnameKana  string
	PostalCode     string
	PrefectureCode int32
	City           string
	AddressLine1   string
	AddressLine2   string
	PhoneNumber    string
	IsDefault      bool
}

type Admin interface {
	MultiGet(ctx context.Context, adminIDs []string, fields ...string) (entity.Admins, error)
	Get(ctx context.Context, adminID string, fields ...string) (*entity.Admin, error)
	GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.Admin, error)
	GetByEmail(ctx context.Context, email string, fields ...string) (*entity.Admin, error)
	UpdateEmail(ctx context.Context, adminID, email string) error
	UpdateDevice(ctx context.Context, adminID, device string) error
	UpdateSignInAt(ctx context.Context, adminID string) error
}

type AdminGroup interface {
	List(ctx context.Context, params *ListAdminGroupsParams, fields ...string) (entity.AdminGroups, error)
	Count(ctx context.Context, params *ListAdminGroupsParams) (int64, error)
	MultiGet(ctx context.Context, groupIDs []string, fields ...string) (entity.AdminGroups, error)
	Get(ctx context.Context, groupID string, fields ...string) (*entity.AdminGroup, error)
	Upsert(ctx context.Context, group *entity.AdminGroup) error
	Delete(ctx context.Context, groupID string) error
}

type ListAdminGroupsParams struct {
	Limit  int
	Offset int
}

type AdminGroupRole interface {
	List(ctx context.Context, params *ListAdminGroupRolesParams, fields ...string) (entity.AdminGroupRoles, error)
	Count(ctx context.Context, params *ListAdminGroupRolesParams) (int64, error)
	Get(ctx context.Context, groupID, roleID string, fields ...string) (*entity.AdminGroupRole, error)
	Upsert(ctx context.Context, role *entity.AdminGroupRole) error
	Delete(ctx context.Context, groupID, roleID string) error
}

type ListAdminGroupRolesParams struct {
	Limit  int
	Offset int
}

type AdminGroupUser interface {
	List(ctx context.Context, params *ListAdminGroupUsersParams, fields ...string) (entity.AdminGroupUsers, error)
	Count(ctx context.Context, params *ListAdminGroupUsersParams) (int64, error)
	Get(ctx context.Context, groupID, adminID string, fields ...string) (*entity.AdminGroupUser, error)
	Upsert(ctx context.Context, user *entity.AdminGroupUser) error
	Delete(ctx context.Context, groupID, adminID string) error
}

type ListAdminGroupUsersParams struct {
	Limit  int
	Offset int
}

type AdminPolicy interface {
	List(ctx context.Context, params *ListAdminPoliciesParams, fields ...string) (entity.AdminPolicies, error)
	Count(ctx context.Context, params *ListAdminPoliciesParams) (int64, error)
	MultiGet(ctx context.Context, policyIDs []string, fields ...string) (entity.AdminPolicies, error)
	Get(ctx context.Context, policyID string, fields ...string) (*entity.AdminPolicy, error)
}

type ListAdminPoliciesParams struct {
	Limit  int
	Offset int
}

type AdminRole interface {
	List(ctx context.Context, params *ListAdminRolesParams, fields ...string) (entity.AdminRoles, error)
	Count(ctx context.Context, params *ListAdminRolesParams) (int64, error)
	MultiGet(ctx context.Context, roleIDs []string, fields ...string) (entity.AdminRoles, error)
	Get(ctx context.Context, roleID string, fields ...string) (*entity.AdminRole, error)
}

type ListAdminRolesParams struct {
	Limit  int
	Offset int
}

type AdminRolePolicy interface {
	List(ctx context.Context, params *ListAdminRolePoliciesParams, fields ...string) (entity.AdminRolePolicies, error)
	Count(ctx context.Context, params *ListAdminRolePoliciesParams) (int64, error)
	Get(ctx context.Context, roleID, policyID string, fields ...string) (*entity.AdminRolePolicy, error)
}

type ListAdminRolePoliciesParams struct {
	Limit  int
	Offset int
}

type Administrator interface {
	List(ctx context.Context, params *ListAdministratorsParams, fields ...string) (entity.Administrators, error)
	Count(ctx context.Context, params *ListAdministratorsParams) (int64, error)
	MultiGet(ctx context.Context, administratorIDs []string, fields ...string) (entity.Administrators, error)
	Get(ctx context.Context, administratorID string, fields ...string) (*entity.Administrator, error)
	Create(ctx context.Context, administrator *entity.Administrator, auth func(ctx context.Context) error) error
	Update(ctx context.Context, administratorID string, params *UpdateAdministratorParams) error
	Delete(ctx context.Context, administratorID string, auth func(ctx context.Context) error) error
}

type ListAdministratorsParams struct {
	Limit  int
	Offset int
}

type UpdateAdministratorParams struct {
	Lastname      string
	Firstname     string
	LastnameKana  string
	FirstnameKana string
	PhoneNumber   string
}

type Coordinator interface {
	List(ctx context.Context, params *ListCoordinatorsParams, fields ...string) (entity.Coordinators, error)
	Count(ctx context.Context, params *ListCoordinatorsParams) (int64, error)
	MultiGet(ctx context.Context, coordinatorIDs []string, fields ...string) (entity.Coordinators, error)
	MultiGetWithDeleted(ctx context.Context, coordinatorIDs []string, fields ...string) (entity.Coordinators, error)
	Get(ctx context.Context, coordinatorID string, fields ...string) (*entity.Coordinator, error)
	GetWithDeleted(ctx context.Context, coordinatorID string, fields ...string) (*entity.Coordinator, error)
	Create(ctx context.Context, coordinator *entity.Coordinator, auth func(ctx context.Context) error) error
	Update(ctx context.Context, coordinatorID string, params *UpdateCoordinatorParams) error
	Delete(ctx context.Context, coordinatorID string, auth func(ctx context.Context) error) error
	RemoveProductTypeID(ctx context.Context, productTypeID string) error
}

type ListCoordinatorsParams struct {
	Name   string
	Limit  int
	Offset int
}

type UpdateCoordinatorParams struct {
	Lastname          string
	Firstname         string
	LastnameKana      string
	FirstnameKana     string
	MarcheName        string
	Username          string
	Profile           string
	ProductTypeIDs    []string
	ThumbnailURL      string
	HeaderURL         string
	PromotionVideoURL string
	BonusVideoURL     string
	InstagramID       string
	FacebookID        string
	PhoneNumber       string
	PostalCode        string
	PrefectureCode    int32
	City              string
	AddressLine1      string
	AddressLine2      string
	BusinessDays      []time.Weekday
}

type Guest interface {
	GetByEmail(ctx context.Context, email string, fields ...string) (*entity.Guest, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, userID string, params *UpdateGuestParams) error
	Delete(ctx context.Context, userID string) error
}

type UpdateGuestParams struct {
	Lastname      string
	Firstname     string
	LastnameKana  string
	FirstnameKana string
}

type Member interface {
	Get(ctx context.Context, userID string, fields ...string) (*entity.Member, error)
	GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.Member, error)
	GetByEmail(ctx context.Context, email string, fields ...string) (*entity.Member, error)
	Create(ctx context.Context, user *entity.User, auth func(ctx context.Context) error) error
	UpdateVerified(ctx context.Context, userID string) error
	UpdateEmail(ctx context.Context, userID, email string) error
	UpdateUsername(ctx context.Context, userID, username string) error
	UpdateAccountID(ctx context.Context, userID, accountID string) error
	UpdateThumbnailURL(ctx context.Context, userID, thumbnailURL string) error
	Delete(ctx context.Context, userID string, auth func(ctx context.Context) error) error
}

type Producer interface {
	List(ctx context.Context, params *ListProducersParams, fields ...string) (entity.Producers, error)
	Count(ctx context.Context, params *ListProducersParams) (int64, error)
	MultiGet(ctx context.Context, producerIDs []string, fields ...string) (entity.Producers, error)
	MultiGetWithDeleted(ctx context.Context, producerIDs []string, fields ...string) (entity.Producers, error)
	Get(ctx context.Context, producerID string, fields ...string) (*entity.Producer, error)
	GetWithDeleted(ctx context.Context, producerID string, fields ...string) (*entity.Producer, error)
	Create(ctx context.Context, producer *entity.Producer, auth func(ctx context.Context) error) error
	Update(ctx context.Context, producerID string, params *UpdateProducerParams) error
	Delete(ctx context.Context, producerID string, auth func(ctx context.Context) error) error
	AggregateByCoordinatorID(ctx context.Context, coordinatorIDs []string) (map[string]int64, error)
}

type ListProducersParams struct {
	CoordinatorID string
	Name          string
	Limit         int
	Offset        int
}

type UpdateProducerParams struct {
	Lastname          string
	Firstname         string
	LastnameKana      string
	FirstnameKana     string
	Username          string
	Profile           string
	ThumbnailURL      string
	HeaderURL         string
	PromotionVideoURL string
	BonusVideoURL     string
	InstagramID       string
	FacebookID        string
	Email             string
	PhoneNumber       string
	PostalCode        string
	PrefectureCode    int32
	City              string
	AddressLine1      string
	AddressLine2      string
}

type User interface {
	List(ctx context.Context, params *ListUsersParams, fields ...string) (entity.Users, error)
	Count(ctx context.Context, params *ListUsersParams) (int64, error)
	MultiGet(ctx context.Context, userIDs []string, fields ...string) (entity.Users, error)
	Get(ctx context.Context, userID string, fields ...string) (*entity.User, error)
}

type ListUsersParams struct {
	Limit          int
	Offset         int
	OnlyRegistered bool
	OnlyVerified   bool
	WithDeleted    bool
}

type UserNotification interface {
	MultiGet(ctx context.Context, userIDs []string, fields ...string) (entity.UserNotifications, error)
	Get(ctx context.Context, userID string, fields ...string) (*entity.UserNotification, error)
	Upsert(ctx context.Context, notification *entity.UserNotification) error
}

type Error struct {
	err error
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}
