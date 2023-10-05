package mysql

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
)

const broadcastTable = "broadcasts"

type broadcast struct {
	db  *mysql.Client
	now func() time.Time
}

func newBroadcast(db *mysql.Client) database.Broadcast {
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
		return nil, dbError(err)
	}
	return broadcast, nil
}

func (b *broadcast) Create(ctx context.Context, broadcast *entity.Broadcast) error {
	now := b.now()
	broadcast.CreatedAt, broadcast.UpdatedAt = now, now

	err := b.db.DB.WithContext(ctx).Table(broadcastTable).Create(&broadcast).Error
	return dbError(err)
}

func (b *broadcast) Update(ctx context.Context, broadcastID string, params *database.UpdateBroadcastParams) error {
	updates := map[string]interface{}{
		"status":     params.Status,
		"updated_at": b.now(),
	}
	if params.InitializeBroadcastParams != nil {
		updates["input_url"] = params.InputURL
		updates["output_url"] = params.OutputURL
		updates["cloud_front_distribution_arn"] = params.CloudFrontDistributionArn
		updates["media_live_channel_arn"] = params.MediaLiveChannelArn
		updates["media_live_channel_id"] = params.MediaLiveChannelID
		updates["media_live_rtmp_input_arn"] = params.MediaLiveRTMPInputArn
		updates["media_live_rtmp_input_name"] = params.MediaLiveRTMPInputName
		updates["media_live_mp4_input_arn"] = params.MediaLiveMP4InputArn
		updates["media_live_mp4_input_name"] = params.MediaLiveMP4InputName
		updates["media_store_container_arn"] = params.MediaStoreContainerArn
	}
	if params.Status == entity.BroadcastStatusDisabled {
		updates["input_url"] = ""
		updates["output_url"] = ""
	}
	stmt := b.db.DB.WithContext(ctx).
		Table(broadcastTable).
		Where("id = ?", broadcastID)

	err := stmt.Updates(updates).Error
	return dbError(err)
}
