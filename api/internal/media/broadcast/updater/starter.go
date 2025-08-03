package updater

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
)

type starter struct {
	now        func() time.Time
	waitGroup  *sync.WaitGroup
	db         *database.Database
	maxRetries int64
	storageURL func() *url.URL
}

func NewStarter(params *Params, opts ...Option) Updater {
	dopts := &options{
		maxRetries: 3,
	}
	for i := range opts {
		opts[i](dopts)
	}
	storageURL := func() *url.URL {
		url := *params.StorageURL // copy
		return &url
	}
	return &starter{
		now:        jst.Now,
		waitGroup:  params.WaitGroup,
		db:         params.Database,
		maxRetries: dopts.maxRetries,
		storageURL: storageURL,
	}
}

func (s *starter) Lambda(ctx context.Context, payload CreatePayload) error {
	slog.Debug("Received event", slog.Any("payload", payload))
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, payload.ScheduleID)
	if err != nil {
		slog.Error("Not found broadcast", log.Error(err), slog.String("scheduleId", payload.ScheduleID))
		return nil
	}
	dir := fmt.Sprintf(entity.BroadcastArchiveHLSPath, payload.ScheduleID)
	u := &url.URL{
		Scheme: s.storageURL().Scheme,
		Host:   s.storageURL().Host,
		Path:   fmt.Sprintf("%s/live.m3u8", dir),
	}
	params := &database.UpdateBroadcastParams{
		Status: entity.BroadcastStatusIdle,
		InitializeBroadcastParams: &database.InitializeBroadcastParams{
			InputURL:                  payload.MediaLiveRtmpInputURL,
			OutputURL:                 u.String(),
			CloudFrontDistributionArn: payload.CloudFrontDistributionARN,
			MediaLiveChannelArn:       payload.MediaLiveChannelARN,
			MediaLiveChannelID:        payload.MediaLiveChannelID,
			MediaLiveRTMPInputArn:     payload.MediaLiveRtmpInputARN,
			MediaLiveRTMPInputName:    payload.MediaLiveRtmpInputName,
			MediaLiveMP4InputArn:      payload.MediaLiveMp4InputARN,
			MediaLiveMP4InputName:     payload.MediaLiveMp4InputName,
			MediaStoreContainerArn:    payload.MediaStoreContainerARN, // TODO: MediaStoreが廃止されるため、他実装削除後に削除
		},
	}
	if err := s.db.Broadcast.Update(ctx, broadcast.ID, params); err != nil {
		slog.Error("Failed to update broadcast", log.Error(err), slog.String("scheduleId", payload.ScheduleID))
		return err
	}
	slog.Info("Succeeded to update broadcast", slog.String("scheduleId", payload.ScheduleID))
	return nil
}
