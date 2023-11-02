package scheduler

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/mediaconvert"
	"github.com/and-period/furumaru/api/pkg/medialive"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert/types"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type closer struct {
	now         func() time.Time
	logger      *zap.Logger
	waitGroup   *sync.WaitGroup
	semaphore   *semaphore.Weighted
	db          *database.Database
	store       store.Service
	media       medialive.MediaLive
	convert     mediaconvert.MediaConvert
	bucketName  string
	jobTemplate string
}

func NewCloser(params *Params, opts ...Option) Scheduler {
	dopts := &options{
		logger:      zap.NewNop(),
		concurrency: 2,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &closer{
		now:         jst.Now,
		logger:      dopts.logger,
		waitGroup:   params.WaitGroup,
		semaphore:   semaphore.NewWeighted(dopts.concurrency),
		db:          params.Database,
		store:       params.Store,
		media:       params.MediaLive,
		convert:     params.MediaConvert,
		bucketName:  params.ArchiveBucketName,
		jobTemplate: params.ConvertJobTemplate,
	}
}

func (c *closer) Lambda(ctx context.Context) (err error) {
	c.logger.Debug("Started Lambda function", zap.Time("now", c.now()))
	defer func() {
		c.logger.Debug("Finished Lambda function", zap.Time("now", c.now()), zap.Error(err))
	}()

	return c.run(ctx, c.now())
}

func (c *closer) Run(ctx context.Context, target time.Time) error {
	if err := c.run(ctx, target); err != nil {
		c.logger.Error("Failed to run", zap.Time("target", target), zap.Error(err))
		return err
	}
	return nil
}

// run - ライブ配信のリソース削除と停止処理
func (c *closer) run(ctx context.Context, target time.Time) error {
	if err := c.removeChannel(ctx, target); err != nil {
		c.logger.Error("Failed to remove channel", zap.Time("target", target), zap.Error(err))
		return err
	}
	if err := c.stopChannel(ctx, target); err != nil {
		c.logger.Error("Failed to stop channel", zap.Time("target", target), zap.Error(err))
	}
	return nil
}

// stopChannel - ライブ配信を停止 (1時間後)
func (c *closer) stopChannel(ctx context.Context, target time.Time) error {
	c.logger.Debug("Stopping channel...", zap.Time("target", target))
	in := &store.ListSchedulesInput{
		EndAtGte: target.AddDate(0, 0, -1),   // マルシェ開催終了1日経過〜
		EndAtLt:  target.Add(-1 * time.Hour), // 〜マルシェ開催終了1時間経過
		NoLimit:  true,
	}
	schedules, total, err := c.store.ListSchedules(ctx, in)
	if err != nil || total == 0 {
		return err
	}
	c.logger.Debug("Got schedules to stop", zap.Int64("total", total))

	eg, ectx := errgroup.WithContext(ctx)
	for i := range schedules {
		if err := c.semaphore.Acquire(ctx, 1); err != nil {
			return err
		}

		schedule := schedules[i]
		eg.Go(func() error {
			defer c.semaphore.Release(1)
			broadcast, err := c.db.Broadcast.GetByScheduleID(ectx, schedule.ID)
			if err != nil {
				return err
			}
			if broadcast.Status != entity.BroadcastStatusActive {
				c.logger.Debug("Channels excluded from stop",
					zap.String("scheduleId", schedule.ID), zap.Int("status", int(broadcast.Status)))
				return nil // 起動中の場合のみ、停止処理を進める
			}

			if broadcast.MediaLiveChannelID == "" {
				c.logger.Error("Empty media live channel id", zap.String("scheduleId", schedule.ID))
				return fmt.Errorf("unexpected media live channel arn format. arn=%s", broadcast.MediaLiveChannelArn)
			}

			c.logger.Info("Calling to stop media live", zap.String("scheduleId", schedule.ID))
			if err := c.media.StopChannel(ctx, broadcast.MediaLiveChannelID); err != nil {
				c.logger.Error("Failed to stop media live", zap.String("scheduleId", schedule.ID), zap.Error(err))
				return err
			}
			c.logger.Info("Succeeded to stop media live", zap.String("scheduleId", schedule.ID))

			params := &database.UpdateBroadcastParams{
				Status: entity.BroadcastStatusIdle,
			}
			return c.db.Broadcast.Update(ctx, broadcast.ID, params)
		})
	}
	return eg.Wait()
}

// removeChannel - ライブ配信を削除 (1時間後&&停止中)
func (c *closer) removeChannel(ctx context.Context, target time.Time) error {
	in := &store.ListSchedulesInput{
		EndAtGte: target.AddDate(0, 0, -1),   // マルシェ開催終了1日後〜
		EndAtLt:  target.Add(-1 * time.Hour), // 〜マルシェ開催終了1時間後
		NoLimit:  true,
	}
	schedules, total, err := c.store.ListSchedules(ctx, in)
	if err != nil || total == 0 {
		return err
	}
	c.logger.Debug("Got schedules to remove", zap.Int64("total", total))

	eg, ectx := errgroup.WithContext(ctx)
	for i := range schedules {
		if err := c.semaphore.Acquire(ctx, 1); err != nil {
			return err
		}

		schedule := schedules[i]
		eg.Go(func() error {
			defer c.semaphore.Release(1)
			broadcast, err := c.db.Broadcast.GetByScheduleID(ectx, schedule.ID)
			if err != nil {
				return err
			}
			if broadcast.Status != entity.BroadcastStatusIdle {
				return nil // 停止中の場合のみ、削除処理を進める
			}

			c.logger.Warn("Not implemented to remove channel", zap.String("schedule", schedule.ID))

			c.logger.Info("Calling to create convert job", zap.String("scheduleId", schedule.ID))
			if err := c.convert.CreateJob(ctx, c.jobTemplate, c.newMediaConvertJobSettings(broadcast)); err != nil {
				c.logger.Error("Failed to create convert job", zap.String("scheduleId", schedule.ID), zap.Error(err))
				return err
			}
			c.logger.Info("Succeeded to create convert job", zap.String("scheduleId", schedule.ID))

			params := &database.UpdateBroadcastParams{
				Status: entity.BroadcastStatusDisabled,
			}
			return c.db.Broadcast.Update(ctx, broadcast.ID, params)
		})
	}
	return eg.Wait()
}

func (c *closer) newMediaConvertJobSettings(broadcast *entity.Broadcast) *types.JobSettings {
	src := fmt.Sprintf("s3://%s", filepath.Join(c.bucketName, newArchiveHLSPath(broadcast.ScheduleID), playlistFilename))
	dst := fmt.Sprintf("s3://%s", filepath.Join(c.bucketName, newArchiveMP4Path(broadcast.ScheduleID)))

	return &types.JobSettings{
		Inputs: []types.Input{{
			FileInput:      aws.String(src),
			TimecodeSource: types.InputTimecodeSourceZerobased,
		}},
		OutputGroups: []types.OutputGroup{{
			OutputGroupSettings: &types.OutputGroupSettings{
				Type: types.OutputGroupTypeFileGroupSettings,
				FileGroupSettings: &types.FileGroupSettings{
					Destination: aws.String(dst),
				},
			},
		}},
	}
}
