//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/user/$GOPACKAGE/$GOFILE
package database

import (
	"context"

	"github.com/and-period/furumaru/api/internal/common"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"gorm.io/gorm"
)

type Params struct {
	Database *database.Client
}

type Database struct {
	Admin         Admin
	Administrator Administrator
	Coordinator   Coordinator
	Member        Member
	Producer      Producer
	User          User
}

func NewDatabase(params *Params) *Database {
	return &Database{
		Admin:         NewAdmin(params.Database),
		Administrator: NewAdministrator(params.Database),
		Coordinator:   NewCoordinator(params.Database),
		Member:        NewMember(params.Database),
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
	Create(ctx context.Context, administrator *entity.Administrator, auth func(ctx context.Context) error) error
	Update(ctx context.Context, administratorID string, params *UpdateAdministratorParams) error
	Delete(ctx context.Context, administratorID string, auth func(ctx context.Context) error) error
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
}

type User interface {
	List(ctx context.Context, params *ListUsersParams, fields ...string) (entity.Users, error)
	Count(ctx context.Context, params *ListUsersParams) (int64, error)
	MultiGet(ctx context.Context, userIDs []string, fields ...string) (entity.Users, error)
	Get(ctx context.Context, userID string, fields ...string) (*entity.User, error)
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
	CoordinatorID string
	Limit         int
	Offset        int
	OnlyUnrelated bool
}

func (p *ListProducersParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.CoordinatorID != "" {
		stmt = stmt.Where("coordinator_id = ?", p.CoordinatorID)
	}
	if p.OnlyUnrelated {
		stmt = stmt.Where("coordinator_id IS NULL")
	}
	return stmt
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

type ListUsersParams struct {
	Limit  int
	Offset int
}
