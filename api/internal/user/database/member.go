package database

import (
	"context"
	"errors"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
)

const memberTable = "members"

type member struct {
	db  *database.Client
	now func() time.Time
}

func NewMember(db *database.Client) Member {
	return &member{
		db:  db,
		now: jst.Now,
	}
}

func (m *member) Get(ctx context.Context, userID string, fields ...string) (*entity.Member, error) {
	member, err := m.get(ctx, m.db.DB, userID, fields...)
	return member, exception.InternalError(err)
}

func (m *member) GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.Member, error) {
	var member *entity.Member

	stmt := m.db.Statement(ctx, m.db.DB, memberTable, fields...).
		Where("cognito_id = ?", cognitoID)

	err := stmt.First(&member).Error
	return member, exception.InternalError(err)
}

func (m *member) GetByEmail(ctx context.Context, email string, fields ...string) (*entity.Member, error) {
	var member *entity.Member

	stmt := m.db.Statement(ctx, m.db.DB, memberTable, fields...).
		Where("email = ?", email).
		Where("provider_type = ?", entity.ProviderTypeEmail)

	err := stmt.First(&member).Error
	return member, exception.InternalError(err)
}

func (m *member) Create(ctx context.Context, user *entity.User, auth func(ctx context.Context) error) error {
	_, err := m.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		now := m.now()
		user.CreatedAt, user.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(userTable).Create(&user).Error; err != nil {
			return nil, err
		}
		user.Member.CreatedAt, user.Member.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(memberTable).Create(&user.Member).Error; err != nil {
			return nil, err
		}
		user.Customer.CreatedAt, user.Customer.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(customerTable).Create(&user.Customer).Error; err != nil {
			return nil, err
		}
		return nil, auth(ctx)
	})
	return exception.InternalError(err)
}

func (m *member) UpdateVerified(ctx context.Context, userID string) error {
	_, err := m.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		current, err := m.get(ctx, tx, userID, "verified_at")
		if err != nil {
			return nil, err
		}
		if !current.VerifiedAt.IsZero() {
			return nil, exception.ErrFailedPrecondition
		}

		now := m.now()
		params := map[string]interface{}{
			"updated_at":  now,
			"verified_at": now,
		}
		err = tx.WithContext(ctx).
			Table(memberTable).
			Where("user_id = ?", userID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (m *member) UpdateAccount(ctx context.Context, userID, accountID, username string) error {
	_, err := m.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := m.get(ctx, tx, userID); err != nil {
			return nil, err
		}

		var current *entity.Member
		err := m.db.Statement(ctx, tx, memberTable, "user_id").
			Where("user_id != ?", userID).
			Where("account_id = ?", accountID).
			First(&current).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if current.UserID != "" {
			return nil, exception.ErrAlreadyExists
		}

		now := m.now()
		params := map[string]interface{}{
			"account_id": accountID,
			"username":   username,
			"updated_at": now,
		}
		err = tx.WithContext(ctx).
			Table(memberTable).
			Where("user_id = ?", userID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (m *member) UpdateEmail(ctx context.Context, userID, email string) error {
	_, err := m.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		current, err := m.get(ctx, tx, userID, "provider_type")
		if err != nil {
			return nil, err
		}
		if current.ProviderType != entity.ProviderTypeEmail {
			return nil, exception.ErrFailedPrecondition
		}

		now := m.now()
		params := map[string]interface{}{
			"email":      email,
			"updated_at": now,
		}
		err = tx.WithContext(ctx).
			Table(memberTable).
			Where("user_id = ?", userID).
			Updates(params).Error
		return nil, err
	})
	return exception.InternalError(err)
}

func (m *member) Delete(ctx context.Context, userID string, auth func(ctx context.Context) error) error {
	_, err := m.db.Transaction(ctx, func(tx *gorm.DB) (interface{}, error) {
		if _, err := m.get(ctx, tx, userID, "user_id"); err != nil {
			return nil, err
		}

		now := m.now()
		memberParams := map[string]interface{}{
			"exists":     nil,
			"updated_at": now,
			"deleted_at": now,
		}
		err := tx.WithContext(ctx).
			Table(memberTable).
			Where("user_id = ?", userID).
			Updates(memberParams).Error
		if err != nil {
			return nil, err
		}
		userParams := map[string]interface{}{
			"updated_at": now,
			"deleted_at": now,
		}
		err = tx.WithContext(ctx).
			Table(userTable).
			Where("id = ?", userID).
			Updates(userParams).Error
		if err != nil {
			return nil, err
		}
		return nil, auth(ctx)
	})
	return exception.InternalError(err)
}

func (m *member) get(ctx context.Context, tx *gorm.DB, userID string, fields ...string) (*entity.Member, error) {
	var member *entity.Member

	err := m.db.Statement(ctx, tx, memberTable, fields...).
		Where("user_id = ?", userID).
		First(&member).Error
	return member, err
}
