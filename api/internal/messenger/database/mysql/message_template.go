package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const messageTemplateTable = "message_templates"

type messageTemplate struct {
	db  *mysql.Client
	now func() time.Time
}

func NewMessageTemplate(db *mysql.Client) database.MessageTemplate {
	return &messageTemplate{
		db:  db,
		now: jst.Now,
	}
}

func (t *messageTemplate) Get(ctx context.Context, messageID entity.MessageTemplateID, fields ...string) (*entity.MessageTemplate, error) {
	var template *entity.MessageTemplate

	stmt := t.db.Statement(ctx, t.db.DB, messageTemplateTable, fields...).
		Where("id = ?", messageID)

	if err := stmt.First(&template).Error; err != nil {
		return nil, dbError(err)
	}
	return template, nil
}
