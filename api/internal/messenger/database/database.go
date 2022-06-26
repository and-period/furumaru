//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../../mock/store/$GOPACKAGE/$GOFILE
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

type Notification interface {
	Create(ctx context.Context, notification *entity.Notification) error
}

/**
* params
 */
type ListNotificationsParams struct {
	Limit  int
	Offset int
}
