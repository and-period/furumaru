package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const paymentSystemTable = "payment_systems"

type paymentSystem struct {
	db  *mysql.Client
	now func() time.Time
}

func NewPaymentSystem(db *mysql.Client) database.PaymentSystem {
	return &paymentSystem{
		db:  db,
		now: jst.Now,
	}
}

func (s *paymentSystem) MultiGet(
	ctx context.Context, methodTypes []entity.PaymentMethodType, fields ...string,
) (entity.PaymentSystems, error) {
	var systems entity.PaymentSystems

	err := s.db.Statement(ctx, s.db.DB, paymentSystemTable, fields...).
		Where("method_type IN (?)", methodTypes).
		Find(&systems).Error
	return systems, dbError(err)
}

func (s *paymentSystem) Get(
	ctx context.Context, methodType entity.PaymentMethodType, fields ...string,
) (*entity.PaymentSystem, error) {
	var system *entity.PaymentSystem

	stmt := s.db.Statement(ctx, s.db.DB, paymentSystemTable, fields...).
		Where("method_type = ?", methodType)

	if err := stmt.First(&system).Error; err != nil {
		return nil, dbError(err)
	}
	return system, nil
}

func (s *paymentSystem) Update(
	ctx context.Context, methodType entity.PaymentMethodType, status entity.PaymentSystemStatus,
) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": s.now(),
	}
	stmt := s.db.DB.WithContext(ctx).
		Table(paymentSystemTable).
		Where("method_type = ?", methodType)

	err := stmt.Updates(updates).Error
	return dbError(err)
}
