//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/user/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"errors"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

var (
	ErrInvalidArgument    = errors.New("database: invalid argument")
	ErrNotFound           = errors.New("database: not found")
	ErrAlreadyExists      = errors.New("database: already exists")
	ErrFailedPrecondition = errors.New("database: failed precondition")
	ErrCanceled           = errors.New("database: canceled")
	ErrDeadlineExceeded   = errors.New("database: deadline exceeded")
	ErrInternal           = errors.New("database: internal error")
	ErrUnknown            = errors.New("database: unknown")
)

type Database struct {
	Admin         Admin
	Administrator Administrator
	Coordinator   Coordinator
	Member        Member
	Producer      Producer
	User          User
}

/**
 * interface
 */
type Admin interface {
	MultiGet(ctx context.Context, adminIDs []string, fields ...string) (entity.Admins, error)
	Get(ctx context.Context, adminID string, fields ...string) (*entity.Admin, error)
	GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.Admin, error)
	GetByEmail(ctx context.Context, email string, fields ...string) (*entity.Admin, error)
	UpdateEmail(ctx context.Context, adminID, email string) error
	UpdateDevice(ctx context.Context, adminID, device string) error
	UpdateSignInAt(ctx context.Context, adminID string) error
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
	Get(ctx context.Context, coordinatorID string, fields ...string) (*entity.Coordinator, error)
	Create(ctx context.Context, coordinator *entity.Coordinator, auth func(ctx context.Context) error) error
	Update(ctx context.Context, coordinatorID string, params *UpdateCoordinatorParams) error
	UpdateThumbnails(ctx context.Context, coordinatorID string, thumbnails common.Images) error
	UpdateHeaders(ctx context.Context, coordinatorID string, headers common.Images) error
	Delete(ctx context.Context, coordinatorID string, auth func(ctx context.Context) error) error
}

type ListCoordinatorsParams struct {
	Username string
	Limit    int
	Offset   int
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
	Prefecture        int64
	City              string
	AddressLine1      string
	AddressLine2      string
}

type Member interface {
	Get(ctx context.Context, userID string, fields ...string) (*entity.Member, error)
	GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.Member, error)
	GetByEmail(ctx context.Context, email string, fields ...string) (*entity.Member, error)
	Create(ctx context.Context, user *entity.User, auth func(ctx context.Context) error) error
	UpdateVerified(ctx context.Context, userID string) error
	UpdateAccount(ctx context.Context, userID, accountID, username string) error
	UpdateEmail(ctx context.Context, userID, email string) error
	Delete(ctx context.Context, userID string, auth func(ctx context.Context) error) error
}

type Producer interface {
	List(ctx context.Context, params *ListProducersParams, fields ...string) (entity.Producers, error)
	Count(ctx context.Context, params *ListProducersParams) (int64, error)
	MultiGet(ctx context.Context, producerIDs []string, fields ...string) (entity.Producers, error)
	Get(ctx context.Context, producerID string, fields ...string) (*entity.Producer, error)
	Create(ctx context.Context, producer *entity.Producer, auth func(ctx context.Context) error) error
	Update(ctx context.Context, producerID string, params *UpdateProducerParams) error
	UpdateThumbnails(ctx context.Context, producerID string, thumbnails common.Images) error
	UpdateHeaders(ctx context.Context, producerID string, headers common.Images) error
	UpdateRelationship(ctx context.Context, coordinatorID string, producerIDs ...string) error
	Delete(ctx context.Context, producerID string, auth func(ctx context.Context) error) error
	AggregateByCoordinatorID(ctx context.Context, coordinatorIDs []string) (map[string]int64, error)
}

type ListProducersParams struct {
	CoordinatorID string
	Username      string
	Limit         int
	Offset        int
	OnlyUnrelated bool
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
	PhoneNumber       string
	PostalCode        string
	Prefecture        int64
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
	Limit  int
	Offset int
}
