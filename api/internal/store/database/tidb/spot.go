package tidb

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const spotTable = "spots"

type spot struct {
	db  *mysql.Client
	now func() time.Time
}

func NewSpot(db *mysql.Client) database.Spot {
	return &spot{
		db:  db,
		now: jst.Now,
	}
}

type listSpotsParams database.ListSpotsParams

func (p listSpotsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.Name != "" {
		stmt = stmt.Where("`name` LIKE ?", fmt.Sprintf("%%%s%%", p.Name)).
			Or("`description` LIKE ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	if p.UserID != "" {
		stmt = stmt.Where("user_id = ?", p.UserID)
	}
	if p.ExcludeApproved {
		stmt = stmt.Where("approved = ?", false)
	}
	if p.ExcludeDisabled {
		stmt = stmt.Where("approved = ?", true)
	}
	return stmt
}

func (p listSpotsParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (s *spot) List(
	ctx context.Context,
	params *database.ListSpotsParams,
	fields ...string,
) (entity.Spots, error) {
	var spots entity.Spots

	prm := listSpotsParams(*params)

	stmt := s.db.Statement(ctx, s.db.DB, spotTable, fields...)
	stmt = prm.stmt(stmt)
	stmt = prm.pagination(stmt)

	if err := stmt.Find(&spots).Error; err != nil {
		return nil, dbError(err)
	}
	if err := spots.Fill(); err != nil {
		return nil, dbError(err)
	}
	return spots, nil
}

func (s *spot) ListByGeolocation(
	ctx context.Context, params *database.ListSpotsByGeolocationParams, fields ...string,
) (entity.Spots, error) {
	var spots entity.Spots

	// Haversine式を用いて2点間の距離を計算する
	// - 第1引数: latitude
	// - 第2引数: latitude
	// - 第3引数: longitude
	// - 第4引数: radius
	const distance = `ACOS(
    SIN(RADIANS(latitude)) * SIN(RADIANS(?)) +
    COS(RADIANS(latitude)) * COS(RADIANS(?)) *
    COS(RADIANS(longitude) - RADIANS(?))
  ) * 6371 <= ?`

	stmt := s.db.Statement(ctx, s.db.DB, spotTable, fields...).
		Where(distance, params.Latitude, params.Latitude, params.Longitude, params.Radius)
	if len(params.SpotTypeIDs) > 0 {
		stmt = stmt.Where("spot_type_id IN (?)", params.SpotTypeIDs)
	}
	if params.ExcludeDisabled {
		stmt = stmt.Where("approved = ?", true)
	}

	if err := stmt.Find(&spots).Error; err != nil {
		return nil, dbError(err)
	}
	if err := spots.Fill(); err != nil {
		return nil, dbError(err)
	}
	return spots, nil
}

func (s *spot) Count(ctx context.Context, params *database.ListSpotsParams) (int64, error) {
	prm := listSpotsParams(*params)

	total, err := s.db.Count(ctx, s.db.DB, &entity.Spot{}, prm.stmt)
	return total, dbError(err)
}

func (s *spot) Get(ctx context.Context, spotID string, fields ...string) (*entity.Spot, error) {
	var spot *entity.Spot

	stmt := s.db.Statement(ctx, s.db.DB, spotTable, fields...).Where("id = ?", spotID)

	if err := stmt.First(&spot).Error; err != nil {
		return nil, dbError(err)
	}
	if err := spot.Fill(); err != nil {
		return nil, dbError(err)
	}
	return spot, nil
}

func (s *spot) Create(ctx context.Context, spot *entity.Spot) error {
	now := s.now()
	spot.CreatedAt, spot.UpdatedAt = now, now

	err := s.db.DB.WithContext(ctx).Table(spotTable).Create(&spot).Error
	return dbError(err)
}

func (s *spot) Update(ctx context.Context, spotID string, params *database.UpdateSpotParams) error {
	updates := map[string]interface{}{
		"spot_type_id":  params.SpotTypeID,
		"name":          params.Name,
		"description":   params.Description,
		"thumbnail_url": params.ThumbnailURL,
		"longitude":     params.Longitude,
		"latitude":      params.Latitude,
		"postal_code":   params.PostalCode,
		"prefecture":    params.PrefectureCode,
		"city":          params.City,
		"address_line1": params.AddressLine1,
		"address_line2": params.AddressLine2,
		"updated_at":    s.now(),
	}
	stmt := s.db.DB.WithContext(ctx).Table(spotTable).Where("id = ?", spotID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (s *spot) Delete(ctx context.Context, spotID string) error {
	stmt := s.db.DB.WithContext(ctx).Table(spotTable).Where("id = ?", spotID)

	err := stmt.Delete(&entity.Spot{}).Error
	return dbError(err)
}

func (s *spot) Approve(
	ctx context.Context,
	spotID string,
	params *database.ApproveSpotParams,
) error {
	updates := map[string]interface{}{
		"approved":          true,
		"approved_admin_id": params.ApprovedAdminID,
		"updated_at":        s.now(),
	}
	stmt := s.db.DB.WithContext(ctx).Table(spotTable).Where("id = ?", spotID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}
