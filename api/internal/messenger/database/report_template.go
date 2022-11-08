package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
)

const reportTemplateTable = "report_templates"

type reportTemplate struct {
	db  *database.Client
	now func() time.Time
}

func NewReportTemplate(db *database.Client) ReportTemplate {
	return &reportTemplate{
		db:  db,
		now: jst.Now,
	}
}

func (t *reportTemplate) Get(ctx context.Context, reportID string, fields ...string) (*entity.ReportTemplate, error) {
	var template *entity.ReportTemplate

	err := t.db.Statement(ctx, t.db.DB, reportTemplateTable, fields...).
		Where("id = ?", reportID).
		First(&template).Error
	return template, exception.InternalError(err)
}
