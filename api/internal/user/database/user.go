package database

import (
	"context"
	"errors"
	"time"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/database"
	"github.com/and-period/marche/api/pkg/jst"
	"gorm.io/gorm"
)

const userTable = "users"

var userFields = []string{
	"id", "account_id", "cognito_id", "provider_type",
	"username", "email", "phone_number", "thumbnail_url",
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

	stmt := u.db.DB.WithContext(ctx).
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

	stmt := u.db.DB.WithContext(ctx).
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

	stmt := u.db.DB.WithContext(ctx).
		Table(userTable).Select(fields).
		Where("email = ?", email).
		Where("provider_type = ?", entity.ProviderTypeEmail)

	if err := stmt.First(&user).Error; err != nil {
		return nil, dbError(err)
	}
	return user, nil
}

func (u *user) Create(ctx context.Context, user *entity.User) error {
	_, err := u.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := u.now()
		user.CreatedAt, user.UpdatedAt = now, now

		err := tx.WithContext(ctx).Table(userTable).Create(&user).Error
		return nil, err
	})
	return dbError(err)
}

func (u *user) UpdateVerified(ctx context.Context, userID string) error {
	_, err := u.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		var current *entity.User
		err := tx.WithContext(ctx).
			Table(userTable).Select("id", "verified_at").
			Where("id = ?", userID).
			First(&current).Error
		if err != nil {
			return nil, err
		}
		if !current.VerifiedAt.IsZero() {
			return nil, ErrFailedPrecondition
		}

		now := u.now()
		params := map[string]interface{}{
			"id":          current.ID,
			"verified_at": now,
			"updated_at":  now,
		}
		err = tx.WithContext(ctx).
			Table(userTable).
			Where("id = ?", current.ID).
			Updates(params).Error
		return nil, err
	})
	return dbError(err)
}

func (u *user) UpdateAccountInfo(ctx context.Context, userID, accountID, userName string) error {
	_, err := u.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		var current *entity.User
		err := tx.WithContext(ctx).
			Table(userTable).Select("id").
			Where("id != ?", userID).
			Where("account_id = ?", accountID).
			First(&current).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if current.ID != "" {
			return nil, ErrAlreadyExists
		}

		err = tx.WithContext(ctx).
			Table(userTable).Select("id").
			Where("id = ?", userID).
			First(&current).Error
		if err != nil {
			return nil, ErrNotFound
		}

		params := map[string]interface{}{
			"account_id": accountID,
			"username":   userName,
			"updated_at": u.now(),
		}
		err = tx.WithContext(ctx).
			Table(userTable).
			Where("id = ?", userID).
			Updates(params).Error
		return nil, err
	})
	return dbError(err)
}

func (u *user) UpdateEmail(ctx context.Context, userID, email string) error {
	_, err := u.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		var current *entity.User
		err := tx.WithContext(ctx).
			Table(userTable).Select("id", "provider_type").
			Where("id = ?", userID).
			First(&current).Error
		if err != nil {
			return nil, err
		}
		if current.ProviderType != entity.ProviderTypeEmail {
			return nil, ErrFailedPrecondition
		}

		params := map[string]interface{}{
			"id":         current.ID,
			"email":      email,
			"updated_at": u.now(),
		}
		err = tx.WithContext(ctx).
			Table(userTable).
			Where("id = ?", userID).
			Updates(params).Error
		return nil, err
	})
	return dbError(err)
}

func (u *user) Delete(ctx context.Context, userID string) error {
	_, err := u.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		var current *entity.User
		err := tx.WithContext(ctx).
			Table(userTable).Select("id").
			Where("id = ?", userID).
			First(&current).Error
		if err != nil {
			return nil, err
		}

		now := u.now()
		params := map[string]interface{}{
			"id":         current.ID,
			"updated_at": now,
			"deleted_at": now,
		}
		err = tx.WithContext(ctx).
			Table(userTable).
			Where("id = ?", current.ID).
			Updates(params).Error
		return nil, err
	})
	return dbError(err)
}
