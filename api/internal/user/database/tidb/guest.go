package tidb

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

const guestTable = "guests"

type guest struct {
	db  *mysql.Client
	now func() time.Time
}

func NewGuest(db *mysql.Client) database.Guest {
	return &guest{
		db:  db,
		now: jst.Now,
	}
}

func (g *guest) GetDummy(ctx context.Context, fields ...string) (*entity.Guest, error) {
	var guest *entity.Guest

	stmt := g.db.Statement(ctx, g.db.DB, guestTable, fields...).
		Joins("INNER JOIN users ON guests.user_id = users.id").
		Where("email LIKE ?", "%@example.com").
		Where("users.deleted_at IS NULL").
		Order("RAND()")

	if err := stmt.First(&guest).Error; err != nil {
		return nil, dbError(err)
	}
	return guest, nil
}

func (g *guest) GetByEmail(ctx context.Context, email string, fields ...string) (*entity.Guest, error) {
	var guest *entity.Guest

	if len(fields) == 0 {
		fields = []string{"*"}
	}
	for i, field := range fields {
		fields[i] = fmt.Sprintf("guests.%s", field)
	}

	stmt := g.db.Statement(ctx, g.db.DB, guestTable, fields...).
		Joins("INNER JOIN users ON guests.user_id = users.id").
		Where("email = ?", email).
		Where("users.deleted_at IS NULL")

	if err := stmt.First(&guest).Error; err != nil {
		return nil, dbError(err)
	}
	return guest, nil
}

func (g *guest) Create(ctx context.Context, user *entity.User) error {
	err := g.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := g.now()
		user.CreatedAt, user.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(userTable).Create(&user).Error; err != nil {
			return err
		}
		user.Guest.CreatedAt, user.Guest.UpdatedAt = now, now
		return tx.WithContext(ctx).Table(guestTable).Create(&user.Guest).Error
	})
	return dbError(err)
}

func (g *guest) Update(ctx context.Context, userID string, params *database.UpdateGuestParams) error {
	updates := map[string]interface{}{
		"lastname":       params.Lastname,
		"firstname":      params.Firstname,
		"lastname_kana":  params.LastnameKana,
		"firstname_kana": params.FirstnameKana,
		"updated_at":     g.now(),
	}
	err := g.db.DB.WithContext(ctx).
		Table(guestTable).
		Where("user_id = ?", userID).
		Updates(updates).Error
	return dbError(err)
}

func (g *guest) Delete(ctx context.Context, userID string) error {
	err := g.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := g.now()
		guestParams := map[string]interface{}{
			"exists":     nil,
			"updated_at": now,
		}
		err := tx.WithContext(ctx).
			Table(guestTable).
			Where("user_id = ?", userID).
			Updates(guestParams).Error
		if err != nil {
			return err
		}
		userParams := map[string]interface{}{
			"updated_at": now,
			"deleted_at": now,
		}
		err = tx.WithContext(ctx).
			Table(userTable).
			Where("id = ?", userID).
			Updates(userParams).Error
		return err
	})
	return dbError(err)
}
