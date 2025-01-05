package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const reportTemplateTable = "report_templates"

type reportTemplate struct {
	db  *mysql.Client
	now func() time.Time
}

func NewReportTemplate(db *mysql.Client) database.ReportTemplate {
	return &reportTemplate{
		db:  db,
		now: jst.Now,
	}
}

func (t *reportTemplate) Get(ctx context.Context, reportID entity.ReportTemplateID, fields ...string) (*entity.ReportTemplate, error) {
	var template *entity.ReportTemplate

	stmt := t.db.Statement(ctx, t.db.DB, reportTemplateTable, fields...).
		Where("id = ?", reportID)

	if err := stmt.First(&template).Error; err != nil {
		return nil, dbError(err)
	}
	return template, nil
}
