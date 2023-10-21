package mysql

import (
	"context"
	"errors"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
)

const memberTable = "members"

type member struct {
	db  *mysql.Client
	now func() time.Time
}

func newMember(db *mysql.Client) database.Member {
	return &member{
		db:  db,
		now: jst.Now,
	}
}

func (m *member) Get(ctx context.Context, userID string, fields ...string) (*entity.Member, error) {
	member, err := m.get(ctx, m.db.DB, userID, fields...)
	return member, dbError(err)
}

func (m *member) GetByCognitoID(ctx context.Context, cognitoID string, fields ...string) (*entity.Member, error) {
	var member *entity.Member

	stmt := m.db.Statement(ctx, m.db.DB, memberTable, fields...).
		Where("cognito_id = ?", cognitoID)

	if err := stmt.First(&member).Error; err != nil {
		return nil, dbError(err)
	}
	return member, nil
}

func (m *member) GetByEmail(ctx context.Context, email string, fields ...string) (*entity.Member, error) {
	var member *entity.Member

	stmt := m.db.Statement(ctx, m.db.DB, memberTable, fields...).
		Where("email = ?", email).
		Where("provider_type = ?", entity.ProviderTypeEmail)

	if err := stmt.First(&member).Error; err != nil {
		return nil, dbError(err)
	}
	return member, nil
}

func (m *member) Create(ctx context.Context, user *entity.User, auth func(ctx context.Context) error) error {
	err := m.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := m.now()
		user.CreatedAt, user.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(userTable).Create(&user).Error; err != nil {
			return err
		}
		user.Member.CreatedAt, user.Member.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(memberTable).Create(&user.Member).Error; err != nil {
			return err
		}
		user.Customer.CreatedAt, user.Customer.UpdatedAt = now, now
		if err := tx.WithContext(ctx).Table(customerTable).Create(&user.Customer).Error; err != nil {
			return err
		}
		return auth(ctx)
	})
	return dbError(err)
}

func (m *member) UpdateVerified(ctx context.Context, userID string) error {
	err := m.db.Transaction(ctx, func(tx *gorm.DB) error {
		current, err := m.get(ctx, tx, userID, "verified_at")
		if err != nil {
			return err
		}
		if !current.VerifiedAt.IsZero() {
			return database.ErrFailedPrecondition
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
		return err
	})
	return dbError(err)
}

func (m *member) UpdateAccount(ctx context.Context, userID, accountID, username string) error {
	err := m.db.Transaction(ctx, func(tx *gorm.DB) error {
		var current *entity.Member
		err := m.db.Statement(ctx, tx, memberTable, "user_id").
			Where("user_id != ?", userID).
			Where("account_id = ?", accountID).
			First(&current).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if current.UserID != "" {
			return database.ErrAlreadyExists
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
		return err
	})
	return dbError(err)
}

func (m *member) UpdateEmail(ctx context.Context, userID, email string) error {
	err := m.db.Transaction(ctx, func(tx *gorm.DB) error {
		current, err := m.get(ctx, tx, userID, "provider_type")
		if err != nil {
			return err
		}
		if current.ProviderType != entity.ProviderTypeEmail {
			return database.ErrFailedPrecondition
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
		return err
	})
	return dbError(err)
}

func (m *member) Delete(ctx context.Context, userID string, auth func(ctx context.Context) error) error {
	err := m.db.Transaction(ctx, func(tx *gorm.DB) error {
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
		if err != nil {
			return err
		}
		return auth(ctx)
	})
	return dbError(err)
}

func (m *member) get(ctx context.Context, tx *gorm.DB, userID string, fields ...string) (*entity.Member, error) {
	var member *entity.Member

	stmt := m.db.Statement(ctx, tx, memberTable, fields...).
		Where("user_id = ?", userID)

	if err := stmt.First(&member).Error; err != nil {
		return nil, err
	}
	return member, nil
}