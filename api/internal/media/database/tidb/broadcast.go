package tidb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mysql"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const broadcastTable = "broadcasts"

type broadcast struct {
	db  *mysql.Client
	now func() time.Time
}

func NewBroadcast(db *mysql.Client) database.Broadcast {
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
	var internal internalBroadcasts

	p := listBroadcastsParams(*params)

	stmt := b.db.Statement(ctx, b.db.DB, broadcastTable, fields...)
	stmt = p.stmt(stmt)
	stmt = p.pagination(stmt)

	if err := stmt.Find(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	broadcasts, err := internal.entities()
	if err != nil {
		return nil, dbError(err)
	}
	return broadcasts, nil
}

func (b *broadcast) Count(ctx context.Context, params *database.ListBroadcastsParams) (int64, error) {
	p := listBroadcastsParams(*params)

	total, err := b.db.Count(ctx, b.db.DB, &entity.Broadcast{}, p.stmt)
	return total, dbError(err)
}

func (b *broadcast) Get(ctx context.Context, broadcastID string, fields ...string) (*entity.Broadcast, error) {
	var internal *internalBroadcast

	stmt := b.db.Statement(ctx, b.db.DB, broadcastTable, fields...).
		Where("id = ?", broadcastID)

	if err := stmt.First(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	broadcast, err := internal.entity()
	if err != nil {
		return nil, dbError(err)
	}
	return broadcast, nil
}

func (b *broadcast) GetByScheduleID(
	ctx context.Context, scheduleID string, fields ...string,
) (*entity.Broadcast, error) {
	var internal *internalBroadcast

	stmt := b.db.Statement(ctx, b.db.DB, broadcastTable, fields...).
		Where("schedule_id = ?", scheduleID)

	if err := stmt.First(&internal).Error; err != nil {
		return nil, dbError(err)
	}
	broadcast, err := internal.entity()
	if err != nil {
		return nil, dbError(err)
	}
	return broadcast, nil
}

func (b *broadcast) Create(ctx context.Context, broadcast *entity.Broadcast) error {
	now := b.now()
	broadcast.CreatedAt, broadcast.UpdatedAt = now, now

	internal, err := newInternalBroadcast(broadcast)
	if err != nil {
		return dbError(err)
	}

	err = b.db.DB.WithContext(ctx).Table(broadcastTable).Create(&internal).Error
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
		updates["archive_url"] = params.UploadBroadcastArchiveParams.ArchiveURL
		updates["archive_fixed"] = params.ArchiveFixed
	}
	if params.UpdateBroadcastArchiveParams != nil {
		metadata, err := json.Marshal(params.ArchiveMetadata)
		if err != nil {
			return dbError(err)
		}
		updates["archive_url"] = params.UpdateBroadcastArchiveParams.ArchiveURL
		updates["archive_metadata"] = metadata
	}
	if params.UpsertYoutubeBroadcastParams != nil {
		updates["youtube_account"] = params.YoutubeAccount
		updates["youtube_broadcast_id"] = params.YoutubeBroadcastID
		updates["youtube_stream_id"] = params.YoutubeStreamID
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

type internalBroadcast struct {
	entity.Broadcast    `gorm:"embedded"`
	ArchiveMetadataJSON datatypes.JSON `gorm:"default:null;column:archive_metadata"` // アーカイブメタデータ(JSON)
}

type internalBroadcasts []*internalBroadcast

func newInternalBroadcast(broadcast *entity.Broadcast) (*internalBroadcast, error) {
	metadata, err := json.Marshal(broadcast.ArchiveMetadata)
	if err != nil {
		return nil, fmt.Errorf("database: failed to marshal archive metadata: %w", err)
	}
	internal := &internalBroadcast{
		Broadcast:           *broadcast,
		ArchiveMetadataJSON: metadata,
	}
	return internal, nil
}

func (b *internalBroadcast) entity() (*entity.Broadcast, error) {
	if err := b.unmarshalArchiveMetadata(); err != nil {
		return nil, err
	}
	return &b.Broadcast, nil
}

func (b *internalBroadcast) unmarshalArchiveMetadata() error {
	if b == nil || b.ArchiveMetadataJSON == nil {
		return nil
	}
	metadata := &entity.BroadcastArchiveMetadata{}
	if err := json.Unmarshal(b.ArchiveMetadataJSON, &metadata); err != nil {
		return fmt.Errorf("database: failed to unmarshal archive metadata: %w", err)
	}
	b.Broadcast.ArchiveMetadata = metadata
	return nil
}

func (es internalBroadcasts) entities() (entity.Broadcasts, error) {
	res := make(entity.Broadcasts, len(es))
	for i := range es {
		b, err := es[i].entity()
		if err != nil {
			return nil, err
		}
		res[i] = b
	}
	return res, nil
}
