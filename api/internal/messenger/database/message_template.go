package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
)

const messageTemplateTable = "message_templates"

var messageTemplateFields = []string{
	"id", "title_template", "body_template", "created_at", "updated_at",
}

type messageTemplate struct {
	db  *database.Client
	now func() time.Time
}

func NewMessageTemplate(db *database.Client) MessageTemplate {
	return &messageTemplate{
		db:  db,
		now: jst.Now,
	}
}

func (t *messageTemplate) Get(
	ctx context.Context, messageID string, fields ...string,
) (*entity.MessageTemplate, error) {
	var template *entity.MessageTemplate
	if len(fields) == 0 {
		fields = messageTemplateFields
	}

	err := t.db.DB.WithContext(ctx).
		Table(messageTemplateTable).Select(fields).
		Where("id = ?", messageID).
		First(&template).Error
	return template, exception.InternalError(err)
}
