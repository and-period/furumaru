//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/user/$GOPACKAGE/$GOFILE
package database

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"gorm.io/gorm"
)

type Params struct {
	Database *database.Client
}

type Database struct {
	AdminAuth     AdminAuth
	Administrator Administrator
	Coordinator   Coordinator
	Producer      Producer
	User          User
}

func NewDatabase(params *Params) *Database {
	return &Database{
		AdminAuth:     NewAdminAuth(params.Database),
		Administrator: NewAdministrator(params.Database),
		Coordinator:   NewCoordinator(params.Database),
		Producer:      NewProducer(params.Database),
		User:          NewUser(params.Database),
	}
}

/**
 * interface
 */
type AdminAuth interface {
	MultiGet(ctx context.Context, adminIDs []string, fields ...string) (entity.AdminAuths, error)
	GetByAdminID(ctx context.Context, adminID string, fields ...string) (*entity.AdminAuth, error)
	GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.AdminAuth, error)
	UpdateDevice(ctx context.Context, adminID, device string) error
}

type Administrator interface {
	List(ctx context.Context, params *ListAdministratorsParams, fields ...string) (entity.Administrators, error)
	Count(ctx context.Context, params *ListAdministratorsParams) (int64, error)
	MultiGet(ctx context.Context, administratorIDs []string, fields ...string) (entity.Administrators, error)
	Get(ctx context.Context, administratorID string, fields ...string) (*entity.Administrator, error)
	Create(ctx context.Context, auth *entity.AdminAuth, administrator *entity.Administrator) error
	Update(ctx context.Context, administratorID string, params *UpdateAdministratorParams) error
	UpdateEmail(ctx context.Context, administratorID, email string) error
}

type Coordinator interface {
	List(ctx context.Context, params *ListCoordinatorsParams, fields ...string) (entity.Coordinators, error)
	Count(ctx context.Context, params *ListCoordinatorsParams) (int64, error)
	MultiGet(ctx context.Context, coordinatorIDs []string, fields ...string) (entity.Coordinators, error)
	Get(ctx context.Context, coordinatorID string, fields ...string) (*entity.Coordinator, error)
	Create(ctx context.Context, auth *entity.AdminAuth, coordinator *entity.Coordinator) error
	Update(ctx context.Context, coordinatorID string, params *UpdateCoordinatorParams) error
	UpdateEmail(ctx context.Context, coordinatorID, email string) error
}

type Producer interface {
	List(ctx context.Context, params *ListProducersParams, fields ...string) (entity.Producers, error)
	Count(ctx context.Context, params *ListProducersParams) (int64, error)
	MultiGet(ctx context.Context, producerIDs []string, fields ...string) (entity.Producers, error)
	Get(ctx context.Context, producerID string, fields ...string) (*entity.Producer, error)
	Create(ctx context.Context, auth *entity.AdminAuth, producer *entity.Producer) error
	Update(ctx context.Context, producerID string, params *UpdateProducerParams) error
	UpdateEmail(ctx context.Context, producerID, email string) error
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
	Orders []*ListAdministratorsOrder
}
type ListAdministratorsOrder struct {
	Key        entity.AdministratorOrderBy
	OrderByASC bool
}

func (p *ListAdministratorsParams) stmt(stmt *gorm.DB) *gorm.DB {
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("%s ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("%s DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	return stmt
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
	Orders []*ListCoordinatorsOrder
}

type ListCoordinatorsOrder struct {
	Key        entity.CoordinatorOrderBy
	OrderByASC bool
}

func (p *ListCoordinatorsParams) stmt(stmt *gorm.DB) *gorm.DB {
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("%s ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("%s DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	return stmt
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
	Orders []*ListProducersOrder
}

type ListProducersOrder struct {
	Key        entity.ProducerOrderBy
	OrderByASC bool
}

func (p *ListProducersParams) stmt(stmt *gorm.DB) *gorm.DB {
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("%s ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("%s DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
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
