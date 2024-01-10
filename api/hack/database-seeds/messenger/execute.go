package messenger

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/hack/database-seeds/common"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	database           = "messengers"
	srcMessageTemplate = "message-template.csv"
	srcPushTemplate    = "push-template.csv"
	srcReportTemplate  = "report-template.csv"
)

var errInvalidCSVFormat = errors.New("messenger: invalid csv format")

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
	a.logger.Info("Executing messengers database seeds...")
	if err := a.executeMessageTemplates(ctx); err != nil {
		return err
	}
	a.logger.Info("Completed message_templates table")
	if err := a.executePushTemplates(ctx); err != nil {
		return err
	}
	a.logger.Info("Completed push_templates table")
	if err := a.executeReportTemplates(ctx); err != nil {
		return err
	}
	a.logger.Info("Completed report_templates table")
	a.logger.Info("Completed messengers database seeds")
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

func (a *app) executeMessageTemplates(ctx context.Context) error {
	reader, file, err := a.newReader(srcMessageTemplate)
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
			// 形式）"テンプレートID","タイトル","本文"
			if len(records) < 3 {
				return errInvalidCSVFormat
			}

			now := a.now()
			template := &entity.MessageTemplate{
				TemplateID:    entity.MessageTemplateID(records[0]),
				TitleTemplate: records[1],
				BodyTemplate:  records[2],
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			updates := map[string]interface{}{
				"title_template": template.TitleTemplate,
				"body_template":  template.BodyTemplate,
				"updated_at":     now,
			}

			stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.Assignments(updates),
			})
			if err := stmt.Create(&template).Error; err != nil {
				return err
			}
		}
	})
}

func (a *app) executePushTemplates(ctx context.Context) error {
	reader, file, err := a.newReader(srcPushTemplate)
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
			// 形式）"テンプレートID","タイトル","本文","サムネイルURL"
			if len(records) < 4 {
				return errInvalidCSVFormat
			}

			now := a.now()
			template := &entity.PushTemplate{
				TemplateID:    entity.PushTemplateID(records[0]),
				TitleTemplate: records[1],
				BodyTemplate:  records[2],
				ImageURL:      records[3],
				CreatedAt:     now,
				UpdatedAt:     now,
			}
			updates := map[string]interface{}{
				"title_template": template.TitleTemplate,
				"body_template":  template.BodyTemplate,
				"image_url":      template.ImageURL,
				"updated_at":     now,
			}

			stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.Assignments(updates),
			})
			if err := stmt.Create(&template).Error; err != nil {
				return err
			}
		}
	})
}

func (a *app) executeReportTemplates(ctx context.Context) error {
	reader, file, err := a.newReader(srcReportTemplate)
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
			// 形式）"テンプレートID","テンプレート"
			if len(records) < 2 {
				a.logger.Debug("debug", zap.Strings("records", records))
				return errInvalidCSVFormat
			}

			now := a.now()
			template := &entity.ReportTemplate{
				TemplateID: entity.ReportTemplateID(records[0]),
				Template:   records[1],
				CreatedAt:  now,
				UpdatedAt:  now,
			}
			updates := map[string]interface{}{
				"template":   template.Template,
				"updated_at": now,
			}

			stmt := tx.WithContext(ctx).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.Assignments(updates),
			})
			if err := stmt.Create(&template).Error; err != nil {
				return err
			}
		}
	})
}
