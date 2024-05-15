package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/gorm"
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

type listBroadcastsParams database.ListBroadcastsParams

func (p listBroadcastsParams) stmt(stmt *gorm.DB) *gorm.DB {
	if len(p.ScheduleIDs) > 0 {
		stmt = stmt.Where("schedule_id IN (?)", p.ScheduleIDs)
	}
	if p.CoordinatorID != "" {
		stmt = stmt.Where("coordinator_id = ?", p.CoordinatorID)
	}
	if p.OnlyArchived {
		stmt = stmt.Where("archive_url != ''")
	}
	for i := range p.Orders {
		var value string
		if p.Orders[i].OrderByASC {
			value = fmt.Sprintf("`%s` ASC", p.Orders[i].Key)
		} else {
			value = fmt.Sprintf("`%s` DESC", p.Orders[i].Key)
		}
		stmt = stmt.Order(value)
	}
	return stmt
}

func (p listBroadcastsParams) pagination(stmt *gorm.DB) *gorm.DB {
	if p.Limit > 0 {
		stmt = stmt.Limit(p.Limit)
	}
	if p.Offset > 0 {
		stmt = stmt.Offset(p.Offset)
	}
	return stmt
}

func (b *broadcast) List(ctx context.Context, params *database.ListBroadcastsParams, fields ...string) (entity.Broadcasts, error) {
	var broadcasts entity.Broadcasts

	p := listBroadcastsParams(*params)

	stmt := b.db.Statement(ctx, b.db.DB, broadcastTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	err := stmt.Find(&broadcasts).Error
	return broadcasts, dbError(err)
}

func (b *broadcast) Count(ctx context.Context, params *database.ListBroadcastsParams) (int64, error) {
	p := listBroadcastsParams(*params)

	total, err := b.db.Count(ctx, b.db.DB, &entity.Broadcast{}, p.stmt)
	return total, dbError(err)
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
		"updated_at": b.now(),
	}
	if params.Status != entity.BroadcastStatusUnknown {
		updates["status"] = params.Status
	}
	if params.InitializeBroadcastParams != nil {
		updates["input_url"] = params.InputURL
		updates["output_url"] = params.OutputURL
		updates["archive_url"] = ""
		updates["cloud_front_distribution_arn"] = params.CloudFrontDistributionArn
		updates["media_live_channel_arn"] = params.MediaLiveChannelArn
		updates["media_live_channel_id"] = params.MediaLiveChannelID
		updates["media_live_rtmp_input_arn"] = params.MediaLiveRTMPInputArn
		updates["media_live_rtmp_input_name"] = params.MediaLiveRTMPInputName
		updates["media_live_mp4_input_arn"] = params.MediaLiveMP4InputArn
		updates["media_live_mp4_input_name"] = params.MediaLiveMP4InputName
		updates["media_store_container_arn"] = params.MediaStoreContainerArn
	}
	if params.UploadBroadcastArchiveParams != nil {
		updates["archive_url"] = params.ArchiveURL
		updates["archive_fixed"] = params.ArchiveFixed
	}
	if params.UpsertYoutubeBroadcastParams != nil {
		updates["youtube_account"] = params.YoutubeAccount
		updates["youtube_stream_key"] = params.YoutubeStreamKey
		updates["youtube_stream_url"] = params.YoutubeStreamURL
		updates["youtube_backup_url"] = params.YoutubeBackupURL
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
