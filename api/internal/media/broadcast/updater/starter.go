package updater

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"go.uber.org/zap"
)

type starter struct {
	now        func() time.Time
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
	db         *database.Database
	maxRetries int64
}

func NewStarter(params *Params, opts ...Option) Updater {
	dopts := &options{
		logger:     zap.NewNop(),
		maxRetries: 3,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &starter{
		now:        jst.Now,
		logger:     dopts.logger,
		waitGroup:  params.WaitGroup,
		db:         params.Database,
		maxRetries: dopts.maxRetries,
	}
}

func (s *starter) Lambda(ctx context.Context, payload CreatePayload) error {
	s.logger.Debug("Received event", zap.Any("payload", payload))
	u := &url.URL{
		Scheme: "https",
		Host:   payload.CloudFrontURL,
		Path:   fmt.Sprintf("%s.m3u8", payload.MediaLiveRtmpStreamName),
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, payload.ScheduleID)
	if err != nil {
		s.logger.Error("Not found broadcast", zap.Error(err), zap.String("scheduleId", payload.ScheduleID))
		return nil
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
			MediaStoreContainerArn:    payload.MediaStoreContainerARN,
		},
	}
	if err := s.db.Broadcast.Update(ctx, broadcast.ID, params); err != nil {
		s.logger.Error("Failed to update broadcast", zap.Error(err), zap.String("scheduleId", payload.ScheduleID))
		return err
	}
	s.logger.Info("Succeeded to update broadcast", zap.String("scheduleId", payload.ScheduleID))
	return nil
}
