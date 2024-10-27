package tidb

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/database/mysql"
	apmysql "github.com/and-period/furumaru/api/pkg/mysql"
	gmysql "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase(db *apmysql.Client) *database.Database {
	client := mysql.NewDatabase(db)
	return &database.Database{
		Contact:         client.Contact,
		ContactCategory: client.ContactCategory,
		ContactRead:     client.ContactRead,
		Message:         client.Message,
		MessageTemplate: client.MessageTemplate,
		Notification:    client.Notification,
		PushTemplate:    client.PushTemplate,
		ReceivedQueue:   client.ReceivedQueue,
		ReportTemplate:  client.ReportTemplate,
		Schedule:        client.Schedule,
	}
}
