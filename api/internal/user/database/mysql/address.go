package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const (
	addressTable         = "addresses"
	addressRevisionTable = "address_revisions"
)

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

	stmt := a.db.Statement(ctx, a.db.DB, addressTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&addresses).Error; err != nil {
		return nil, dbError(err)
	}
	if err := a.fill(ctx, a.db.DB, addresses...); err != nil {
		return nil, dbError(err)
	}
	return addresses, nil
}

func (a *address) ListDefault(ctx context.Context, userIDs []string, fields ...string) (entity.Addresses, error) {
	var addresses entity.Addresses

	stmt := a.db.Statement(ctx, a.db.DB, addressTable, fields...).
		Where("user_id IN (?)", userIDs).
		Where("is_default = ?", true)

	if err := stmt.Find(&addresses).Error; err != nil {
		return nil, dbError(err)
	}
	if err := a.fill(ctx, a.db.DB, addresses...); err != nil {
		return nil, dbError(err)
	}
	return addresses, nil
}

func (a *address) Count(ctx context.Context, params *database.ListAddressesParams) (int64, error) {
	p := listAddressesParams(*params)

	total, err := a.db.Count(ctx, a.db.DB, &entity.Address{}, p.stmt)
	return total, dbError(err)
}

func (a *address) MultiGet(ctx context.Context, addressIDs []string, fields ...string) (entity.Addresses, error) {
	addresses, err := a.multiGet(ctx, a.db.DB, addressIDs, fields...)
	return addresses, dbError(err)
}

func (a *address) MultiGetByRevision(ctx context.Context, revisionIDs []int64, fields ...string) (entity.Addresses, error) {
	var revisions entity.AddressRevisions

	stmt := a.db.Statement(ctx, a.db.DB, addressRevisionTable).
		Where("id IN (?)", revisionIDs)

	if err := stmt.Find(&revisions).Error; err != nil {
		return nil, dbError(err)
	}
	if len(revisions) == 0 {
		return entity.Addresses{}, nil
	}
	revisions.Fill()

	addresses, err := a.multiGet(ctx, a.db.DB, revisions.AddressIDs(), fields...)
	if err != nil {
		return nil, err
	}
	if len(addresses) == 0 {
		return entity.Addresses{}, nil
	}

	res, err := revisions.Merge(addresses.Map())
	if err != nil {
		return nil, dbError(err)
	}
	return res, nil
}

func (a *address) Get(ctx context.Context, addressID string, fields ...string) (*entity.Address, error) {
	address, err := a.get(ctx, a.db.DB, addressID, fields...)
	return address, dbError(err)
}

func (a *address) GetDefault(ctx context.Context, userID string, fields ...string) (*entity.Address, error) {
	var address *entity.Address

	stmt := a.db.Statement(ctx, a.db.DB, addressTable, fields...).
		Where("user_id = ?", userID).
		Where("is_default = ?", true)

	if err := stmt.First(&address).Error; err != nil {
		return nil, dbError(err)
	}
	if err := a.fill(ctx, a.db.DB, address); err != nil {
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
				Table(addressTable).
				Where("user_id = ?", address.UserID).
				Where("is_default = ?", true)
			if err := stmt.Updates(updates).Error; err != nil {
				return err
			}
		}

		address.CreatedAt, address.UpdatedAt = now, now
		address.AddressRevision.CreatedAt, address.AddressRevision.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(addressTable).Create(&address).Error; err != nil {
			return err
		}
		return tx.WithContext(ctx).Table(addressRevisionTable).Create(&address.AddressRevision).Error
	})
	return dbError(err)
}

func (a *address) Update(ctx context.Context, addressID, userID string, params *database.UpdateAddressParams) error {
	now := a.now()
	rparams := &entity.NewAddressRevisionParams{
		AddressID:      addressID,
		Lastname:       params.Lastname,
		Firstname:      params.Firstname,
		LastnameKana:   params.LastnameKana,
		FirstnameKana:  params.FirstnameKana,
		PostalCode:     params.PostalCode,
		PrefectureCode: params.PrefectureCode,
		City:           params.City,
		AddressLine1:   params.AddressLine1,
		AddressLine2:   params.AddressLine2,
		PhoneNumber:    params.PhoneNumber,
	}
	revision, err := entity.NewAddressRevision(rparams)
	if err != nil {
		return fmt.Errorf("mysql: failed to new address revision: %w: %s", database.ErrInvalidArgument, err.Error())
	}

	err = a.db.Transaction(ctx, func(tx *gorm.DB) error {
		if params.IsDefault {
			// 現在デフォルト設定になっているものを解除する
			current := map[string]interface{}{
				"is_default": false,
				"updated_at": now,
			}
			stmt := tx.WithContext(ctx).
				Table(addressTable).
				Where("user_id = ?", userID).
				Where("is_default = ?", true)
			if err := stmt.Updates(current).Error; err != nil {
				return err
			}

			// デフォルト設定に変更
			updates := map[string]interface{}{
				"is_default": true,
				"updated_at": now,
			}
			stmt = tx.WithContext(ctx).
				Table(addressTable).
				Where("id = ?", addressID).
				Where("user_id = ?", userID)
			if err := stmt.Updates(updates).Error; err != nil {
				return err
			}
		}

		revision.CreatedAt, revision.UpdatedAt = now, now
		return tx.WithContext(ctx).Table(addressRevisionTable).Create(&revision).Error
	})
	return dbError(err)
}

func (a *address) Delete(ctx context.Context, addressID, userID string) error {
	now := a.now()

	// 現在デフォルト設定になっている場合に解除も同時に行う
	updates := map[string]interface{}{
		"is_default": false,
		"updated_at": now,
		"deleted_at": now,
	}
	stmt := a.db.DB.WithContext(ctx).
		Table(addressTable).
		Where("id = ?", addressID).
		Where("user_id = ?", userID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}

func (a *address) multiGet(ctx context.Context, tx *gorm.DB, addressIDs []string, fields ...string) (entity.Addresses, error) {
	var addresses entity.Addresses

	stmt := a.db.Statement(ctx, tx, addressTable, fields...).Unscoped().Where("id IN (?)", addressIDs)

	if err := stmt.Find(&addresses).Error; err != nil {
		return nil, dbError(err)
	}
	if err := a.fill(ctx, tx, addresses...); err != nil {
		return nil, dbError(err)
	}
	return addresses, nil
}

func (a *address) get(ctx context.Context, tx *gorm.DB, addressID string, fields ...string) (*entity.Address, error) {
	var address *entity.Address

	stmt := a.db.Statement(ctx, tx, addressTable, fields...).
		Where("id = ?", addressID)

	if err := stmt.First(&address).Error; err != nil {
		return nil, err
	}
	if err := a.fill(ctx, tx, address); err != nil {
		return nil, err
	}
	return address, nil
}

func (a *address) fill(ctx context.Context, tx *gorm.DB, addresses ...*entity.Address) error {
	var revisions entity.AddressRevisions

	ids := entity.Addresses(addresses).IDs()
	if len(ids) == 0 {
		return nil
	}

	sub := tx.Table(addressRevisionTable).
		Select("MAX(id)").
		Where("address_id IN (?)", ids).
		Group("address_id")
	stmt := a.db.Statement(ctx, tx, addressRevisionTable).
		Where("id IN (?)", sub)

	if err := stmt.Find(&revisions).Error; err != nil {
		return err
	}
	revisions.Fill()
	entity.Addresses(addresses).Fill(revisions.MapByAddressID())
	return nil
}
