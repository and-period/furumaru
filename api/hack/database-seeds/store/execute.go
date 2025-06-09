package store

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/hack/database-seeds/common"
	"github.com/and-period/furumaru/api/hack/database-seeds/master"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	database           = "stores"
	srcCategory        = "category.csv"
	srcProductType     = "product-type.csv"
	srcPaymentStatus   = "payment-status.csv"
	srcMessageTemplate = "message-template.csv"
	srcPushTemplate    = "push-template.csv"
	srcReportTemplate  = "report-template.csv"
)

var (
	errInvalidCSVFormat    = errors.New("store: invalid csv format")
	errInvalidCategoryName = errors.New("store: invalid category name")
)

type app struct {
	db     *mysql.Client
	logger *zap.Logger
	now    func() time.Time
	srcDir string
}

func NewClient(params *common.Params) (common.Client, error) {
	db, err := common.NewDBClient(params, database)
	if err != nil {
		return nil, err
	}
	return &app{
		db:     db,
		logger: params.Logger,
		now:    jst.Now,
		srcDir: params.SrcDir,
	}, nil
}

func (a *app) Execute(ctx context.Context) error {
	a.logger.Info("Executing stores database seeds...")
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := a.executeCategories(ectx); err != nil {
			return fmt.Errorf("failed to execute categories table: %w", err)
		}
		a.logger.Info("Completed categories table")
		return nil
	})
	eg.Go(func() error {
		if err := a.executeProductTypes(ectx); err != nil {
			return fmt.Errorf("failed to execute product_types table: %w", err)
		}
		a.logger.Info("Completed product_types table")
		return nil
	})
	eg.Go(func() error {
		if err := a.executeShipping(ectx); err != nil {
			return fmt.Errorf("failed to execute shippings table: %w", err)
		}
		a.logger.Info("Completed shippings table")
		return nil
	})
	eg.Go(func() error {
		if err := a.executePaymentSystems(ectx); err != nil {
			return fmt.Errorf("failed to execute payment_systems table: %w", err)
		}
		a.logger.Info("Completed payment_systems table")
		return nil
	})
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to execute stores database seeds: %w", err)
	}
	a.logger.Info("Completed stores database seeds")
	return nil
}

func (a *app) newReader(src string) (*csv.Reader, *os.File, error) {
	filename := strings.Join([]string{a.srcDir, src}, "/")
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	reader := csv.NewReader(file)
	// ヘッダー行は読み飛ばしたいため、１行分読み取りをする
	if _, err := reader.Read(); err != nil {
		return nil, nil, err
	}
	return reader, file, nil
}

func (a *app) executeCategories(ctx context.Context) error {
	reader, file, err := a.newReader(srcCategory)
	if err != nil {
		return err
	}
	defer file.Close()

	return a.db.Transaction(ctx, func(tx *gorm.DB) error {
		for {
			records, err := reader.Read()
			if errors.Is(err, io.EOF) {
				return nil
			}
			// 形式）"カテゴリ"
			if len(records) < 1 {
				return errInvalidCSVFormat
			}

			params := &entity.NewCategoryParams{
				Name: records[0],
			}
			now := a.now()
			category := entity.NewCategory(params)
			category.CreatedAt = now
			category.UpdatedAt = now

			updates := map[string]interface{}{
				"name":       category.Name,
				"updated_at": now,
			}
			stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "name"}},
				DoUpdates: clause.Assignments(updates),
				DoNothing: true, // 他のカラムが追加されたらfalseにする
			})
			if err := stmt.Create(&category).Error; err != nil {
				return err
			}
		}
	})
}

func (a *app) executeProductTypes(ctx context.Context) error {
	reader, file, err := a.newReader(srcProductType)
	if err != nil {
		return err
	}
	defer file.Close()

	return a.db.Transaction(ctx, func(tx *gorm.DB) error {
		var categories entity.Categories
		if err := a.db.Statement(ctx, tx, "categories").Find(&categories).Error; err != nil {
			return err
		}
		categoryMap := categories.MapByName()

		for {
			records, err := reader.Read()
			if errors.Is(err, io.EOF) {
				return nil
			}
			// 形式）"カテゴリ","品目"
			if len(records) < 2 {
				return errInvalidCSVFormat
			}

			category, ok := categoryMap[records[0]]
			if !ok {
				return fmt.Errorf("%w: name=%s", errInvalidCategoryName, records[0])
			}
			params := &entity.NewProductTypeParams{
				Name:       records[1],
				CategoryID: category.ID,
			}
			now := a.now()
			productType := entity.NewProductType(params)
			productType.CreatedAt = now
			productType.UpdatedAt = now
			updates := map[string]interface{}{
				"name":        productType.Name,
				"category_id": productType.CategoryID,
				"updated_at":  now,
			}

			stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "category_id"}, {Name: "name"}},
				DoUpdates: clause.Assignments(updates),
				DoNothing: true, // 他のカラムが追加されたらfalseにする
			})
			if err := stmt.Create(&productType).Error; err != nil {
				return err
			}
		}
	})
}

type internalShippingRevision struct {
	entity.ShippingRevision `gorm:"embedded"`
	Box60RatesJSON          datatypes.JSON `gorm:"default:null;column:box60_rates"`  // 箱サイズ60の通常便配送料一覧(JSON)
	Box80RatesJSON          datatypes.JSON `gorm:"default:null;column:box80_rates"`  // 箱サイズ80の通常便配送料一覧(JSON)
	Box100RatesJSON         datatypes.JSON `gorm:"default:null;column:box100_rates"` // 箱サイズ100の通常便配送料一覧(JSON)
}

func newInternalShippingRevision(revision *entity.ShippingRevision) (*internalShippingRevision, error) {
	box60Rates, err := revision.Box60Rates.Marshal()
	if err != nil {
		return nil, fmt.Errorf("tidb: failed to marshal box60 rates: %w", err)
	}
	box80Rates, err := revision.Box80Rates.Marshal()
	if err != nil {
		return nil, fmt.Errorf("tidb: failed to marshal box80 rates: %w", err)
	}
	box100Rates, err := revision.Box100Rates.Marshal()
	if err != nil {
		return nil, fmt.Errorf("tidb: failed to marshal box100 rates: %w", err)
	}
	internal := &internalShippingRevision{
		ShippingRevision: *revision,
		Box60RatesJSON:   box60Rates,
		Box80RatesJSON:   box80Rates,
		Box100RatesJSON:  box100Rates,
	}
	return internal, nil
}

func (a *app) executeShipping(ctx context.Context) error {
	return a.db.Transaction(ctx, func(tx *gorm.DB) error {
		now := a.now()

		shipping := master.DefaultShipping
		shipping.CreatedAt = now
		shipping.UpdatedAt = now

		stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true,
		})
		if err := stmt.Create(&shipping).Error; err != nil {
			return err
		}

		revision, err := newInternalShippingRevision(master.DefaultShippingRevision)
		if err != nil {
			return err
		}

		revision.CreatedAt = now
		revision.UpdatedAt = now

		stmt = tx.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoNothing: true, // 過去の配送履歴等に影響するため、変更があっても反映しない
		})
		return stmt.Table("shipping_revisions").Create(&revision).Error
	})
}

func (a *app) executePaymentSystems(ctx context.Context) error {
	reader, file, err := a.newReader(srcPaymentStatus)
	if err != nil {
		return err
	}
	defer file.Close()

	return a.db.Transaction(ctx, func(tx *gorm.DB) error {
		for {
			records, err := reader.Read()
			if errors.Is(err, io.EOF) {
				return nil
			}
			// 形式）"決済ID","状態種別"
			if len(records) < 2 {
				return errInvalidCSVFormat
			}

			methodType, err := strconv.ParseInt(records[0], 10, 32)
			if err != nil {
				return err
			}
			status, err := strconv.ParseInt(records[1], 10, 32)
			if err != nil {
				return err
			}
			now := a.now()
			system := &entity.PaymentSystem{
				MethodType: entity.PaymentMethodType(methodType),
				Status:     entity.PaymentSystemStatus(status),
				CreatedAt:  now,
				UpdatedAt:  now,
			}

			stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "method_type"}},
				DoNothing: true, // ステータス以外を変更する場合はfalseにする
			})
			if err := stmt.Create(&system).Error; err != nil {
				return err
			}
		}
	})
}
