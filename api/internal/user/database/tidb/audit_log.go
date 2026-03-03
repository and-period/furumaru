package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const auditLogTable = "audit_logs"

type auditLog struct {
	db  *mysql.Client
	now func() time.Time
}

func NewAuditLog(db *mysql.Client) database.AuditLog {
	return &auditLog{
		db:  db,
		now: jst.Now,
	}
}

func (a *auditLog) List(
	ctx context.Context, params *database.ListAuditLogsParams, fields ...string,
) (entity.AuditLogs, error) {
	var logs entity.AuditLogs

	stmt := a.db.Statement(ctx, a.db.DB, auditLogTable, fields...)
	stmt = a.filterStatement(stmt, params)
	stmt = stmt.Order("created_at DESC")

	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	if err := stmt.Find(&logs).Error; err != nil {
		return nil, dbError(err)
	}
	return logs, nil
}

func (a *auditLog) Count(ctx context.Context, params *database.ListAuditLogsParams) (int64, error) {
	var count int64

	stmt := a.db.DB.WithContext(ctx).Table(auditLogTable)
	stmt = a.filterStatement(stmt, params)

	if err := stmt.Count(&count).Error; err != nil {
		return 0, dbError(err)
	}
	return count, nil
}

func (a *auditLog) Create(ctx context.Context, log *entity.AuditLog) error {
	now := a.now()
	log.CreatedAt = now
	log.UpdatedAt = now
	err := a.db.DB.WithContext(ctx).Table(auditLogTable).Create(log).Error
	return dbError(err)
}

func (a *auditLog) BatchCreate(ctx context.Context, logs entity.AuditLogs) error {
	if len(logs) == 0 {
		return nil
	}
	now := a.now()
	for _, log := range logs {
		log.CreatedAt = now
		log.UpdatedAt = now
	}
	err := a.db.DB.WithContext(ctx).Table(auditLogTable).Create(&logs).Error
	return dbError(err)
}

func (a *auditLog) filterStatement(stmt *gorm.DB, params *database.ListAuditLogsParams) *gorm.DB {
	if params.AdminID != "" {
		stmt = stmt.Where("admin_id = ?", params.AdminID)
	}
	if params.ResourceType != "" {
		stmt = stmt.Where("resource_type = ?", params.ResourceType)
	}
	if params.Action > 0 {
		stmt = stmt.Where("action = ?", params.Action)
	}
	if !params.StartAt.IsZero() {
		stmt = stmt.Where("created_at >= ?", params.StartAt)
	}
	if !params.EndAt.IsZero() {
		stmt = stmt.Where("created_at <= ?", params.EndAt)
	}
	return stmt
}
