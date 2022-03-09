package database

import (
	"context"
	"time"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/database"
	"github.com/and-period/marche/api/pkg/jst"
	"gorm.io/gorm"
)

const userTable = "users"

var userFields = []string{
	"id", "cognito_id", "provider_type", "email", "phone_number",
	"created_at", "updated_at", "verified_at", "deleted_at",
}

type user struct {
	db  *database.Client
	now func() time.Time
}

func NewUser(db *database.Client) User {
	return &user{
		db:  db,
		now: jst.Now,
	}
}

func (u *user) Get(ctx context.Context, userID string, fields ...string) (*entity.User, error) {
	var user *entity.User
	if len(fields) == 0 {
		fields = userFields
	}

	stmt := u.db.DB.
		Table(userTable).Select(fields).
		Where("id = ?", userID)

	if err := stmt.First(&user).Error; err != nil {
		return nil, dbError(err)
	}
	return user, nil
}

func (u *user) GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.User, error) {
	var user *entity.User
	if len(fields) == 0 {
		fields = userFields
	}

	stmt := u.db.DB.
		Table(userTable).Select(fields).
		Where("cognito_id = ?", cognitoID)

	if err := stmt.First(&user).Error; err != nil {
		return nil, dbError(err)
	}
	return user, nil
}

func (u *user) GetByEmail(ctx context.Context, email string, fields ...string) (*entity.User, error) {
	var user *entity.User
	if len(fields) == 0 {
		fields = userFields
	}

	stmt := u.db.DB.
		Table(userTable).Select(fields).
		Where("email = ?", email).
		Where("provider_type = ?", entity.ProviderTypeEmail)

	if err := stmt.First(&user).Error; err != nil {
		return nil, dbError(err)
	}
	return user, nil
}

func (u *user) Create(ctx context.Context, user *entity.User) error {
	_, err := u.db.Transaction(func(tx *gorm.DB) (interface{}, error) {
		now := u.now()
		user.CreatedAt, user.UpdatedAt = now, now

		err := tx.Table(userTable).Create(&user).Error
		return nil, err
	})
	return dbError(err)
}

func (u *user) UpdateVerified(ctx context.Context, userID string) error {
	_, err := u.db.Transaction(func(tx *gorm.DB) (interface{}, error) {
		var current *entity.User
		err := tx.Table(userTable).Select("id", "verified_at").
			Where("id = ?", userID).
			First(&current).Error
		if err != nil || current.ID == "" {
			return nil, err
		}
		if !current.VerifiedAt.IsZero() {
			return nil, ErrFailedPrecondition
		}

		now := u.now()
		params := map[string]interface{}{
			"verified_at": now,
			"updated_at":  now,
		}
		err = tx.Table(userTable).
			Where("id = ?", userID).
			Updates(params).Error
		return nil, err
	})
	return dbError(err)
}
