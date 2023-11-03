package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const adddressTable = "addresses"

type address struct {
	db  *mysql.Client
	now func() time.Time
}

func newAddress(db *mysql.Client) database.Address {
	return &address{
		db:  db,
		now: jst.Now,
	}
}

type listAddressesParams database.ListAddressesParams

func (p listAddressesParams) stmt(stmt *gorm.DB) *gorm.DB {
	if p.UserID != "" {
		stmt = stmt.Where("user_id = ?", p.UserID)
	}
	// デフォルト設定を優先して表示
	stmt = stmt.Order("is_default DESC")
	return stmt
}

func (p listAddressesParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (a *address) List(ctx context.Context, params *database.ListAddressesParams, fields ...string) (entity.Addresses, error) {
	var addresses entity.Addresses

	p := listAddressesParams(*params)

	stmt := a.db.Statement(ctx, a.db.DB, adddressTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	err := stmt.Find(&addresses).Error
	return addresses, dbError(err)
}

func (a *address) Count(ctx context.Context, params *database.ListAddressesParams) (int64, error) {
	p := listAddressesParams(*params)

	total, err := a.db.Count(ctx, a.db.DB, &entity.Address{}, p.stmt)
	return total, dbError(err)
}

func (a *address) MultiGet(ctx context.Context, addressIDs []string, fields ...string) (entity.Addresses, error) {
	var addresses entity.Addresses

	err := a.db.Statement(ctx, a.db.DB, adddressTable, fields...).
		Where("id IN (?)", addressIDs).
		Find(&addresses).Error
	return addresses, dbError(err)
}

func (a *address) Get(ctx context.Context, addressID string, fields ...string) (*entity.Address, error) {
	var address *entity.Address

	stmt := a.db.Statement(ctx, a.db.DB, adddressTable, fields...).
		Where("id = ?", addressID)

	if err := stmt.First(&address).Error; err != nil {
		return nil, dbError(err)
	}
	return address, nil
}

func (a *address) Create(ctx context.Context, address *entity.Address) error {
	err := a.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := a.now()

		if address.IsDefault {
			// 現在デフォルト設定になっているものを解除する
			updates := map[string]interface{}{
				"is_default": false,
				"updated_at": now,
			}
			stmt := tx.WithContext(ctx).
				Table(adddressTable).
				Where("user_id = ?", address.UserID).
				Where("is_default = ?", true)
			if err := stmt.Updates(updates).Error; err != nil {
				return err
			}
		}

		address.CreatedAt, address.UpdatedAt = now, now
		return tx.WithContext(ctx).Table(adddressTable).Create(&address).Error
	})
	return dbError(err)
}

func (a *address) Update(ctx context.Context, addressID, userID string, params *database.UpdateAddressParams) error {
	err := a.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := a.now()

		if params.IsDefault {
			// 現在デフォルト設定になっているものを解除する
			updates := map[string]interface{}{
				"is_default": false,
				"updated_at": now,
			}
			stmt := tx.WithContext(ctx).
				Table(adddressTable).
				Where("user_id = ?", userID).
				Where("is_default = ?", true)
			if err := stmt.Updates(updates).Error; err != nil {
				return err
			}
		}

		updates := map[string]interface{}{
			"lastname":      params.Lastname,
			"firstname":     params.Firstname,
			"postal_code":   params.PostalCode,
			"prefecture":    params.Prefecture,
			"city":          params.City,
			"address_line1": params.AddressLine1,
			"address_line2": params.AddressLine2,
			"phone_number":  params.PhoneNumber,
			"is_default":    params.IsDefault,
			"updated_at":    a.now(),
		}
		stmt := tx.WithContext(ctx).
			Table(adddressTable).
			Where("id = ?", addressID)
		return stmt.Updates(updates).Error
	})
	return dbError(err)
}

func (a *address) Delete(ctx context.Context, addressID string) error {
	stmt := a.db.DB.WithContext(ctx).
		Table(adddressTable).
		Where("id = ?", addressID)

	err := stmt.Delete(&entity.Address{}).Error
	return dbError(err)
}
