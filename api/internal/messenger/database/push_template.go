package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
)

const pushTemplateTable = "push_templates"

var pushTemplateFields = []string{
	"id", "title_template", "body_template", "image_url", "created_at", "updated_at",
}

type pushTemplate struct {
	db  *database.Client
	now func() time.Time
}

func NewPushTemplate(db *database.Client) PushTemplate {
	return &pushTemplate{
		db:  db,
		now: jst.Now,
	}
}

func (t *pushTemplate) Get(ctx context.Context, pushID string, fields ...string) (*entity.PushTemplate, error) {
	var template *entity.PushTemplate
	if len(fields) == 0 {
		fields = pushTemplateFields
	}

	err := t.db.DB.WithContext(ctx).
		Table(pushTemplateTable).Select(fields).
		Where("id = ?", pushID).
		First(&template).Error
	return template, exception.InternalError(err)
}
