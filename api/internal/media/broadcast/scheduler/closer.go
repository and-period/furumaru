package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"path/filepath"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/mediaconvert"
	"github.com/and-period/furumaru/api/pkg/medialive"
	"github.com/and-period/furumaru/api/pkg/sfn"
	"github.com/and-period/furumaru/api/pkg/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediaconvert/types"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type closer struct {
	now         func() time.Time
	waitGroup   *sync.WaitGroup
	semaphore   *semaphore.Weighted
	db          *database.Database
	storage     storage.Bucket
	store       store.Service
	sfn         sfn.StepFunction
	media       medialive.MediaLive
	convert     mediaconvert.MediaConvert
	storageURL  func() *url.URL
	bucketName  string
	jobTemplate string
}

func NewCloser(params *Params, opts ...Option) Scheduler {
	defaultURL, _ := params.Storage.GetHost()
	dopts := &options{
		concurrency: 2,
		storageURL:  defaultURL,
	}
	for i := range opts {
		opts[i](dopts)
	}
	storageURL := func() *url.URL {
		url := *dopts.storageURL // copy
		return &url
	}
	return &closer{
		now:         jst.Now,
		waitGroup:   params.WaitGroup,
		semaphore:   semaphore.NewWeighted(dopts.concurrency),
		db:          params.Database,
		storage:     params.Storage,
		store:       params.Store,
		sfn:         params.StepFunction,
		media:       params.MediaLive,
		convert:     params.MediaConvert,
		storageURL:  storageURL,
		bucketName:  params.ArchiveBucketName,
		jobTemplate: params.ConvertJobTemplate,
	}
}

func (c *closer) Lambda(ctx context.Context) (err error) {
	slog.Debug("Started Lambda function", slog.Time("now", c.now()))
	defer func() {
		slog.Debug("Finished Lambda function", slog.Time("now", c.now()), log.Error(err))
	}()

	return c.run(ctx, c.now())
}

func (c *closer) Run(ctx context.Context, target time.Time) error {
	if err := c.run(ctx, target); err != nil {
		slog.Error("Failed to run", slog.Time("target", target), log.Error(err))
		return err
	}
	return nil
}

// run - ライブ配信のリソース削除と停止処理
func (c *closer) run(ctx context.Context, target time.Time) error {
	if err := c.removeChannel(ctx, target); err != nil {
		slog.Error("Failed to remove channel", slog.Time("target", target), log.Error(err))
		return err
	}
	if err := c.stopChannel(ctx, target); err != nil {
		slog.Error("Failed to stop channel", slog.Time("target", target), log.Error(err))
	}
	return nil
}

// stopChannel - ライブ配信を停止 (1時間後)
func (c *closer) stopChannel(ctx context.Context, target time.Time) error {
	slog.Debug("Stopping channel...", slog.Time("target", target))
	in := &store.ListSchedulesInput{
		EndAtGte: target.AddDate(0, 0, -7),   // 〜マルシェ開催終了7日経過
		EndAtLt:  target.Add(-1 * time.Hour), // マルシェ開催終了1時間経過〜
		NoLimit:  true,
	}
	schedules, total, err := c.store.ListSchedules(ctx, in)
	if err != nil || total == 0 {
		return err
	}
	slog.Debug("Got schedules to stop", slog.Int64("total", total))

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
				slog.Debug("Channels excluded from stop",
					slog.String("scheduleId", schedule.ID), slog.Int("status", int(broadcast.Status)))
				return nil // 起動中の場合のみ、停止処理を進める
			}

			if broadcast.MediaLiveChannelID == "" {
				slog.Error("Empty media live channel id", slog.String("scheduleId", schedule.ID))
				return fmt.Errorf("unexpected media live channel arn format. arn=%s", broadcast.MediaLiveChannelArn)
			}

			slog.Info("Calling to stop media live", slog.String("scheduleId", schedule.ID))
			if err := c.media.StopChannel(ctx, broadcast.MediaLiveChannelID); err != nil {
				slog.Error("Failed to stop media live", slog.String("scheduleId", schedule.ID), log.Error(err))
				return err
			}
			slog.Info("Succeeded to stop media live", slog.String("scheduleId", schedule.ID))

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
		EndAtGte: target.AddDate(0, 0, -7),   // 〜マルシェ開催終了7日経過
		EndAtLt:  target.Add(-1 * time.Hour), // マルシェ開催終了1時間経過〜
		NoLimit:  true,
	}
	schedules, total, err := c.store.ListSchedules(ctx, in)
	if err != nil || total == 0 {
		return err
	}
	slog.Debug("Got schedules to remove", slog.Int64("total", total))

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

			slog.Info("Calling to create convert job", slog.String("scheduleId", schedule.ID))
			if err := c.convert.CreateJob(ectx, c.jobTemplate, c.newMediaConvertJobSettings(broadcast)); err != nil {
				slog.Error("Failed to create convert job", slog.String("scheduleId", schedule.ID), log.Error(err))
				return err
			}
			slog.Info("Succeeded to create convert job", slog.String("scheduleId", schedule.ID))

			payload := &RemovePayload{
				CloudFrontDistributionARN: broadcast.CloudFrontDistributionArn,
				MediaLiveChannelID:        broadcast.MediaLiveChannelID,
				MediaStoreContainerARN:    broadcast.MediaStoreContainerArn,
			}
			slog.Info("Calling step function", slog.String("scheduleId", schedule.ID))
			if err := c.sfn.StartExecution(ectx, payload); err != nil {
				slog.Error("Failed step function", slog.String("scheduleId", schedule.ID), log.Error(err))
				return err
			}
			slog.Info("Succeeded step function", slog.String("scheduleId", schedule.ID))

			archiveURL := c.storageURL()
			archiveURL.Path = filepath.Join(newArchiveMP4Path(broadcast.ScheduleID), archiveFilename)

			params := &database.UpdateBroadcastParams{
				Status: entity.BroadcastStatusDisabled,
				UploadBroadcastArchiveParams: &database.UploadBroadcastArchiveParams{
					ArchiveURL: archiveURL.String(),
				},
			}
			return c.db.Broadcast.Update(ectx, broadcast.ID, params)
		})
	}
	return eg.Wait()
}

func (c *closer) newMediaConvertJobSettings(broadcast *entity.Broadcast) *types.JobSettings {
	src := c.storage.GenerateS3URI(filepath.Join(newArchiveHLSPath(broadcast.ScheduleID), playlistFilename))
	dst := c.storage.GenerateS3URI(newArchiveMP4Path(broadcast.ScheduleID))

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
