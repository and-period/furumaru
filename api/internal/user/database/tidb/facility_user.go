package tidb

import (
	"context"
	"time"

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
