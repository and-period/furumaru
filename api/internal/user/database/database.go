//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/user/$GOPACKAGE/$GOFILE
package database

import (
	"context"

	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/database"
)

type Params struct {
	Database *database.Client
}

type Database struct {
	Admin         Admin
	Administrator Administrator
	Coordinator   Coordinator
	Producer      Producer
	User          User
}

func NewDatabase(params *Params) *Database {
	return &Database{
		Admin:         NewAdmin(params.Database),
		Administrator: NewAdministrator(params.Database),
		Coordinator:   NewCoordinator(params.Database),
		Producer:      NewProducer(params.Database),
		User:          NewUser(params.Database),
	}
}

/**
 * interface
 */
type Admin interface {
	MultiGet(ctx context.Context, adminIDs []string, fields ...string) (entity.Admins, error)
	Get(ctx context.Context, adminID string, fields ...string) (*entity.Admin, error)
	GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.Admin, error)
	UpdateEmail(ctx context.Context, adminID, email string) error
	UpdateDevice(ctx context.Context, adminID, device string) error
}

type Administrator interface {
	List(ctx context.Context, params *ListAdministratorsParams, fields ...string) (entity.Administrators, error)
	Count(ctx context.Context, params *ListAdministratorsParams) (int64, error)
	MultiGet(ctx context.Context, administratorIDs []string, fields ...string) (entity.Administrators, error)
	Get(ctx context.Context, administratorID string, fields ...string) (*entity.Administrator, error)
	Create(ctx context.Context, admin *entity.Admin, administrator *entity.Administrator) error
	Update(ctx context.Context, administratorID string, params *UpdateAdministratorParams) error
}

type Coordinator interface {
	List(ctx context.Context, params *ListCoordinatorsParams, fields ...string) (entity.Coordinators, error)
	Count(ctx context.Context, params *ListCoordinatorsParams) (int64, error)
	MultiGet(ctx context.Context, coordinatorIDs []string, fields ...string) (entity.Coordinators, error)
	Get(ctx context.Context, coordinatorID string, fields ...string) (*entity.Coordinator, error)
	Create(ctx context.Context, admin *entity.Admin, coordinator *entity.Coordinator) error
	Update(ctx context.Context, coordinatorID string, params *UpdateCoordinatorParams) error
}

type Producer interface {
	List(ctx context.Context, params *ListProducersParams, fields ...string) (entity.Producers, error)
	Count(ctx context.Context, params *ListProducersParams) (int64, error)
	MultiGet(ctx context.Context, producerIDs []string, fields ...string) (entity.Producers, error)
	Get(ctx context.Context, producerID string, fields ...string) (*entity.Producer, error)
	Create(ctx context.Context, admin *entity.Admin, producer *entity.Producer) error
	Update(ctx context.Context, producerID string, params *UpdateProducerParams) error
}

type User interface {
	MultiGet(ctx context.Context, userIDs []string, fields ...string) (entity.Users, error)
	Get(ctx context.Context, userID string, fields ...string) (*entity.User, error)
	GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string, fields ...string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	UpdateVerified(ctx context.Context, userID string) error
	UpdateAccount(ctx context.Context, userID, accountID, username string) error
	UpdateEmail(ctx context.Context, userID, email string) error
	Delete(ctx context.Context, userID string) error
}

/**
 * params
 */
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

type ListCoordinatorsParams struct {
	Limit  int
	Offset int
}

type UpdateCoordinatorParams struct {
	Lastname         string
	Firstname        string
	LastnameKana     string
	FirstnameKana    string
	CompanyName      string
	StoreName        string
	ThumbnailURL     string
	HeaderURL        string
	TwitterAccount   string
	InstagramAccount string
	FacebookAccount  string
	PhoneNumber      string
	PostalCode       string
	Prefecture       string
	City             string
	AddressLine1     string
	AddressLine2     string
}

type ListProducersParams struct {
	Limit  int
	Offset int
}

type UpdateProducerParams struct {
	Lastname      string
	Firstname     string
	LastnameKana  string
	FirstnameKana string
	StoreName     string
	ThumbnailURL  string
	HeaderURL     string
	PhoneNumber   string
	PostalCode    string
	Prefecture    string
	City          string
	AddressLine1  string
	AddressLine2  string
}
