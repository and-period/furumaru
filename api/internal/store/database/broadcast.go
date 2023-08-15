package database

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/database"
	"github.com/and-period/furumaru/api/pkg/jst"
	"gorm.io/gorm"
)

const broadcastTable = "broadcasts"

type broadcast struct {
	db  *database.Client
	now func() time.Time
}

func NewBroadcast(db *database.Client) Broadcast {
	return &broadcast{
		db:  db,
		now: jst.Now,
	}
}

func (b *broadcast) GetByScheduleID(
	ctx context.Context, scheduleID string, fields ...string,
) (*entity.Broadcast, error) {
	var broadcast *entity.Broadcast

	stmt := b.db.Statement(ctx, b.db.DB, broadcastTable, fields...).
		Where("schedule_id = ?", scheduleID)

	if err := stmt.First(&broadcast).Error; err != nil {
		return nil, exception.InternalError(err)
	}
	return broadcast, nil
}

func (b *broadcast) Update(ctx context.Context, broadcastID string, params *UpdateBroadcastParams) error {
	err := b.db.Transaction(ctx, func(tx *gorm.DB) error {
		if _, err := b.get(ctx, tx, broadcastID); err != nil {
			return err
		}

		updates := map[string]interface{}{
			"status":     params.Status,
			"updated_at": b.now(),
		}
		if params.Status == entity.BroadcastStatusActive {
			updates["input_url"] = params.InputURL
			updates["output_url"] = params.OutputURL
			updates["cloud_front_distribution_arn"] = params.CloudFrontDistributionArn
			updates["media_live_channel_arn"] = params.MediaLiveChannelArn
			updates["media_live_rtmp_input_arn"] = params.MediaLiveRTMPInputArn
			updates["media_live_mp4_input_arn"] = params.MediaLiveMP4InputArn
			updates["media_store_container_arn"] = params.MediaStoreContainerArn
		}
		if params.Status == entity.BroadcastStatusDisabled {
			updates["input_url"] = ""
			updates["output_url"] = ""
		}

		err := tx.WithContext(ctx).
			Table(broadcastTable).
			Where("id = ?", broadcastID).
			Updates(updates).Error
		return err
	})
	return exception.InternalError(err)
}

func (b *broadcast) get(
	ctx context.Context, tx *gorm.DB, broadcastID string, fields ...string,
) (*entity.Broadcast, error) {
	var broadcast *entity.Broadcast

	stmt := b.db.Statement(ctx, tx, broadcastTable, fields...).Where("id = ?", broadcastID)

	if err := stmt.First(&broadcast).Error; err != nil {
		return nil, err
	}
	return broadcast, nil
}
