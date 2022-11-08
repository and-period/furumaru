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

	err := t.db.Statement(ctx, t.db.DB, pushTemplateTable, fields...).
		Where("id = ?", pushID).
		First(&template).Error
	return template, exception.InternalError(err)
}
