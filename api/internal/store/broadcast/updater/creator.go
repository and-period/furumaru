package updater

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"go.uber.org/zap"
)

type Creator interface {
	Lambda(ctx context.Context, event CreatePayload) error
}

type creator struct {
	now        func() time.Time
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
	db         *database.Database
	maxRetries int64
}

func NewCreator(params *Params, opts ...Option) Creator {
	dopts := &options{
		logger:     zap.NewNop(),
		maxRetries: 3,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &creator{
		now:        jst.Now,
		logger:     dopts.logger,
		waitGroup:  params.WaitGroup,
		db:         params.Database,
		maxRetries: dopts.maxRetries,
	}
}

func (c *creator) Lambda(ctx context.Context, payload CreatePayload) error {
	c.logger.Debug("Received event", zap.Any("payload", payload))
	u, err := url.Parse(payload.CloudFrontURL)
	if err != nil {
		c.logger.Error("Failed to parse cloudfront url", zap.Error(err), zap.String("scheduleId", payload.ScheduleID))
		return err
	}
	u.Path = path.Join(u.Path, payload.MediaLiveRtmpStreamName)
	params := &database.UpdateBroadcastParams{
		Status: entity.BroadcastStatusIdle,
		InitializeBroadcastParams: &database.InitializeBroadcastParams{
			InputURL:                  payload.MediaLiveRtmpInputURL,
			OutputURL:                 fmt.Sprintf("%s.m3u8", u.String()),
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
	if err := c.db.Broadcast.Update(ctx, payload.ScheduleID, params); err != nil {
		c.logger.Error("Failed to update broadcast", zap.Error(err), zap.String("scheduleId", payload.ScheduleID))
		return err
	}
	c.logger.Info("Succeeded to update broadcast", zap.String("scheduleId", payload.ScheduleID))
	return nil
}
