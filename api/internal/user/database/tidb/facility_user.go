package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const facilityUserTable = "facility_users"

type facilityUser struct {
	db  *mysql.Client
	now func() time.Time
}

func NewFacilityUser(db *mysql.Client) *facilityUser {
	return &facilityUser{
		db:  db,
		now: time.Now,
	}
}

func (f *facilityUser) GetByExternalID(
	ctx context.Context,
	providerType entity.UserAuthProviderType,
	externalID, producerID string,
	fields ...string,
) (*entity.FacilityUser, error) {
	var facilityUser *entity.FacilityUser

	stmt := f.db.Statement(ctx, f.db.DB, facilityUserTable, fields...).
		Where("provider_type = ?", providerType).
		Where("external_id = ?", externalID).
		Where("producer_id = ?", producerID)

	if err := stmt.First(&facilityUser).Error; err != nil {
		return nil, dbError(err)
	}
	return facilityUser, nil
}

func (f *facilityUser) Create(ctx context.Context, user *entity.User) error {
	err := f.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := f.now()
		user.CreatedAt, user.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(userTable).Create(&user).Error; err != nil {
			return err
		}
		user.FacilityUser.CreatedAt, user.FacilityUser.UpdatedAt = now, now
		return tx.WithContext(ctx).Table(facilityUserTable).Create(&user.FacilityUser).Error
	})
	return dbError(err)
}

func (f *facilityUser) Update(ctx context.Context, userID string, params *database.UpdateFacilityUserParams) error {
	updates := map[string]interface{}{
		"lastname":         params.Lastname,
		"firstname":        params.Firstname,
		"lastname_kana":    params.LastnameKana,
		"firstname_kana":   params.FirstnameKana,
		"phone_number":     params.PhoneNumber,
		"last_check_in_at": params.LastCheckInAt,
		"updated_at":       f.now(),
	}
	err := f.db.DB.WithContext(ctx).
		Table(facilityUserTable).
		Where("user_id = ?", userID).
		Updates(updates).Error
	return dbError(err)
}
