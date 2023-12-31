package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const pushTemplateTable = "push_templates"

type pushTemplate struct {
	db  *mysql.Client
	now func() time.Time
}

func newPushTemplate(db *mysql.Client) database.PushTemplate {
	return &pushTemplate{
		db:  db,
		now: jst.Now,
	}
}

func (t *pushTemplate) Get(ctx context.Context, pushID entity.PushTemplateID, fields ...string) (*entity.PushTemplate, error) {
	var template *entity.PushTemplate

	stmt := t.db.Statement(ctx, t.db.DB, pushTemplateTable, fields...).
		Where("id = ?", pushID)

	if err := stmt.First(&template).Error; err != nil {
		return nil, dbError(err)
	}
	return template, nil
}
