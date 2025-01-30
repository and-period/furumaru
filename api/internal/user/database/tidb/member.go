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

const memberTable = "members"

type member struct {
	db  *mysql.Client
	now func() time.Time
}

func NewMember(db *mysql.Client) database.Member {
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

	if len(fields) == 0 {
		fields = []string{"*"}
	}
	for i, field := range fields {
		fields[i] = fmt.Sprintf("members.%s", field)
	}

	stmt := m.db.Statement(ctx, m.db.DB, memberTable, fields...).
		Joins("INNER JOIN users ON members.user_id = users.id").
		Where("members.email = ?", email).
		Where("members.provider_type = ?", entity.ProviderTypeEmail).
		Where("users.deleted_at IS NULL")

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

func (m *member) UpdateUsername(ctx context.Context, userID, username string) error {
	params := map[string]interface{}{
		"username":   username,
		"updated_at": m.now(),
	}
	err := m.db.DB.WithContext(ctx).
		Table(memberTable).
		Where("user_id = ?", userID).
		Updates(params).Error
	return dbError(err)
}

func (m *member) UpdateAccountID(ctx context.Context, userID, accountID string) error {
	params := map[string]interface{}{
		"account_id": accountID,
		"updated_at": m.now(),
	}
	err := m.db.DB.WithContext(ctx).
		Table(memberTable).
		Where("user_id = ?", userID).
		Updates(params).Error
	return dbError(err)
}

func (m *member) UpdateThumbnailURL(ctx context.Context, userID, thumbnailURL string) error {
	params := map[string]interface{}{
		"thumbnail_url": thumbnailURL,
		"updated_at":    m.now(),
	}
	err := m.db.DB.WithContext(ctx).
		Table(memberTable).
		Where("user_id = ?", userID).
		Updates(params).Error
	return dbError(err)
}

func (m *member) Delete(ctx context.Context, userID string, auth func(ctx context.Context) error) error {
	err := m.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := m.now()
		memberParams := map[string]interface{}{
			"exists":     nil,
			"updated_at": now,
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

	stmt := m.db.Statement(ctx, tx, memberTable, fields...).Unscoped().Where("user_id = ?", userID)

	if err := stmt.First(&member).Error; err != nil {
		return nil, err
	}
	return member, nil
}
