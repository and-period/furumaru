package updater

import (
	"context"
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
	params := &database.UpdateBroadcastParams{
		Status: entity.BroadcastStatusIdle,
		InitializeBroadcastParams: database.InitializeBroadcastParams{
			InputURL:                  payload.MediaLiveRtmpInputURL,
			OutputURL:                 payload.CloudFrontURL,
			CloudFrontDistributionArn: payload.CloudFrontDistributionARN,
			MediaLiveChannelArn:       payload.MediaLiveChannelARN,
			MediaLiveRTMPInputArn:     payload.MediaLiveRtmpInputARN,
			MediaLiveMP4InputArn:      payload.MediaLiveMp4InputARN,
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
