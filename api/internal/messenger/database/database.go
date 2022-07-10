//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/messenger/$GOPACKAGE/$GOFILE
package database

import (
	"context"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/database"
)

type Params struct {
	Database *database.Client
}

type Database struct {
	Contact      Contact
	Notification Notification
}

func NewDatabase(params *Params) *Database {
	return &Database{
		Notification: NewNotification(params.Database),
	}
}

/**
 * interface
 */
type Contact interface {
	List(ctx context.Context, params *ListContactsParams, fields ...string) (entity.Contacts, error)
	Get(ctx context.Context, contactID string, fields ...string) (*entity.Contact, error)
	Create(ctx context.Context, contact *entity.Contact) error
	Update(ctx context.Context, contactID string, params *UpdateContactParams) error
	Delete(ctx context.Context, contactID string) error
}

type Notification interface {
	Create(ctx context.Context, notification *entity.Notification) error
}

/**
 * params
 */
type ListContactsParams struct {
	Limit  int
	Offset int
}

type UpdateContactParams struct {
	Status   entity.ContactStatus
	Priority entity.ContactPriority
	Note     string
}
