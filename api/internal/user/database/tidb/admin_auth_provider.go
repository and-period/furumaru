package tidb

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm/clause"
)

const adminAuthProviderTable = "admin_auth_providers"

type adminAuthProvider struct {
	db  *mysql.Client
	now func() time.Time
}

func NewAdminAuthProvider(db *mysql.Client) database.AdminAuthProvider {
	return &adminAuthProvider{
		db:  db,
		now: time.Now,
	}
}

func (p *adminAuthProvider) List(
	ctx context.Context, params *database.ListAdminAuthProvidersParams, fields ...string,
) (entity.AdminAuthProviders, error) {
	var providers entity.AdminAuthProviders

	stmt := p.db.Statement(ctx, p.db.DB, adminAuthProviderTable, fields...).
		Where("admin_id = ?", params.AdminID)

	err := stmt.Find(&providers).Error
	return providers, dbError(err)
}

func (p *adminAuthProvider) Get(
	ctx context.Context, adminID string, providerType entity.AdminAuthProviderType, fields ...string,
) (*entity.AdminAuthProvider, error) {
	var provider *entity.AdminAuthProvider

	stmt := p.db.Statement(ctx, p.db.DB, adminAuthProviderTable, fields...).
		Where("admin_id = ?", adminID).
		Where("provider_type = ?", providerType)

	err := stmt.First(&provider).Error
	if err != nil {
		return nil, dbError(err)
	}
	return provider, nil
}

func (p *adminAuthProvider) Upsert(ctx context.Context, provider *entity.AdminAuthProvider) error {
	now := p.now()
	provider.CreatedAt, provider.UpdatedAt = now, now

	updates := map[string]interface{}{
		"account_id": provider.AccountID,
		"email":      provider.Email,
		"updated_at": now,
	}
	clauses := clause.OnConflict{
		Columns:   []clause.Column{{Name: "admin_id"}, {Name: "provider_type"}},
		DoUpdates: clause.Assignments(updates),
	}
	err := p.db.DB.WithContext(ctx).Clauses(clauses).Create(provider).Error
	return dbError(err)
}
