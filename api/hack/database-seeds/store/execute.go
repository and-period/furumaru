package store

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/hack/database-seeds/common"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	database       = "stores"
	srcCategory    = "category.csv"
	srcProductType = "product-type.csv"
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
	a.logger.Info("Executing store database seeds...")
	if err := a.executeCategories(ctx); err != nil {
		return err
	}
	if err := a.executeProductTypes(ctx); err != nil {
		return err
	}
	a.logger.Info("Completed store database seeds")
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
