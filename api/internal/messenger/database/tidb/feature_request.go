package tidb

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const featureRequestTable = "feature_requests"

type featureRequest struct {
	db  *mysql.Client
	now func() time.Time
}

func NewFeatureRequest(db *mysql.Client) database.FeatureRequest {
	return &featureRequest{
		db:  db,
		now: jst.Now,
	}
}

func (f *featureRequest) List(ctx context.Context, params *database.ListFeatureRequestsParams, fields ...string) (entity.FeatureRequests, error) {
	if params == nil {
		params = &database.ListFeatureRequestsParams{}
	}
	var featureRequests entity.FeatureRequests

	stmt := f.db.Statement(ctx, f.db.DB, featureRequestTable, fields...)
	if params.SubmittedBy != "" {
		stmt = stmt.Where("submitted_by = ?", params.SubmittedBy)
	}
	if params.Limit > 0 {
		stmt = stmt.Limit(params.Limit)
	}
	if params.Offset > 0 {
		stmt = stmt.Offset(params.Offset)
	}

	err := stmt.Find(&featureRequests).Error
	return featureRequests, dbError(err)
}

func (f *featureRequest) Count(ctx context.Context, params *database.ListFeatureRequestsParams) (int64, error) {
	if params == nil {
		params = &database.ListFeatureRequestsParams{}
	}
	stmt := f.db.DB.WithContext(ctx).Table(featureRequestTable)
	if params.SubmittedBy != "" {
		stmt = stmt.Where("submitted_by = ?", params.SubmittedBy)
	}
	var total int64
	err := stmt.Count(&total).Error
	return total, dbError(err)
}

func (f *featureRequest) Get(ctx context.Context, featureRequestID string, fields ...string) (*entity.FeatureRequest, error) {
	featureRequest, err := f.get(ctx, f.db.DB, featureRequestID, fields...)
	return featureRequest, dbError(err)
}

func (f *featureRequest) Create(ctx context.Context, featureRequest *entity.FeatureRequest) error {
	now := f.now()
	featureRequest.CreatedAt, featureRequest.UpdatedAt = now, now

	err := f.db.DB.WithContext(ctx).Table(featureRequestTable).Create(&featureRequest).Error
	return dbError(err)
}

func (f *featureRequest) Update(ctx context.Context, featureRequestID string, params *database.UpdateFeatureRequestParams) error {
	updates := map[string]interface{}{
		"status":     params.Status,
		"note":       params.Note,
		"updated_at": f.now(),
	}
	stmt := f.db.DB.WithContext(ctx).
		Table(featureRequestTable).
		Where("id = ?", featureRequestID)

	result := stmt.Updates(updates)
	if err := result.Error; err != nil {
		return dbError(err)
	}
	if result.RowsAffected == 0 {
		return dbError(fmt.Errorf("%w: feature request not found (id=%s)", gorm.ErrRecordNotFound, featureRequestID))
	}
	return nil
}

func (f *featureRequest) Delete(ctx context.Context, featureRequestID string) error {
	params := map[string]interface{}{
		"deleted_at": f.now(),
	}
	stmt := f.db.DB.WithContext(ctx).
		Table(featureRequestTable).
		Where("id = ?", featureRequestID)

	result := stmt.Updates(params)
	if err := result.Error; err != nil {
		return dbError(err)
	}
	if result.RowsAffected == 0 {
		return dbError(fmt.Errorf("%w: feature request not found (id=%s)", gorm.ErrRecordNotFound, featureRequestID))
	}
	return nil
}

func (f *featureRequest) get(ctx context.Context, tx *gorm.DB, featureRequestID string, fields ...string) (*entity.FeatureRequest, error) {
	var featureRequest *entity.FeatureRequest

	stmt := f.db.Statement(ctx, tx, featureRequestTable, fields...).
		Where("id = ?", featureRequestID)

	if err := stmt.First(&featureRequest).Error; err != nil {
		return nil, err
	}
	return featureRequest, nil
}
