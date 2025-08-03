package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/log"
	"github.com/and-period/furumaru/api/pkg/medialive"
	"github.com/and-period/furumaru/api/pkg/sfn"
	"github.com/and-period/furumaru/api/pkg/storage"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type starter struct {
	now        func() time.Time
	waitGroup  *sync.WaitGroup
	semaphore  *semaphore.Weighted
	db         *database.Database
	storage    storage.Bucket
	store      store.Service
	sfn        sfn.StepFunction
	media      medialive.MediaLive
	env        string
	bucketName string
}

func NewStarter(params *Params, opts ...Option) Scheduler {
	dopts := &options{
		concurrency: 2,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &starter{
		now:        jst.Now,
		waitGroup:  params.WaitGroup,
		semaphore:  semaphore.NewWeighted(dopts.concurrency),
		db:         params.Database,
		storage:    params.Storage,
		store:      params.Store,
		sfn:        params.StepFunction,
		media:      params.MediaLive,
		env:        params.Environment,
		bucketName: params.ArchiveBucketName,
	}
}

func (s *starter) Lambda(ctx context.Context) (err error) {
	slog.Debug("Started Lambda function", slog.Time("now", s.now()))
	defer func() {
		slog.Debug("Finished Lambda function", slog.Time("now", s.now()), log.Error(err))
	}()

	return s.run(ctx, s.now())
}

func (s *starter) Run(ctx context.Context, target time.Time) error {
	if err := s.run(ctx, target); err != nil {
		slog.Error("Failed to run", slog.Time("target", target), log.Error(err))
		return err
	}
	return nil
}

// run - ライブ配信のリソース作成と開始処理
func (s *starter) run(ctx context.Context, target time.Time) error {
	if err := s.startChannel(ctx, target); err != nil {
		slog.Error("Failed to start channel", slog.Time("target", target), log.Error(err))
		return err
	}
	if err := s.createChannel(ctx, target); err != nil {
		slog.Error("Failed to create channel", slog.Time("target", target), log.Error(err))
		return err
	}
	return nil
}

// startChannel - ライブ配信を開始 (5分前)
func (s *starter) startChannel(ctx context.Context, target time.Time) error {
	slog.Debug("Starting channel...", slog.Time("target", target))
	in := &store.ListSchedulesInput{
		StartAtLt: target.Add(15 * time.Minute), // マルシェ開催開始15分前〜
		EndAtGte:  target,                       // 〜マルシェ開催終了
		NoLimit:   true,
	}
	schedules, total, err := s.store.ListSchedules(ctx, in)
	if err != nil || total == 0 {
		return err
	}
	slog.Debug("Got schedules to start", slog.Int64("total", total))

	eg, ectx := errgroup.WithContext(ctx)
	for i := range schedules {
		if err := s.semaphore.Acquire(ctx, 1); err != nil {
			return err
		}

		schedule := schedules[i]
		eg.Go(func() error {
			defer s.semaphore.Release(1)
			broadcast, err := s.db.Broadcast.GetByScheduleID(ectx, schedule.ID)
			if err != nil {
				return err
			}
			if broadcast.Status != entity.BroadcastStatusIdle {
				slog.Debug("Channels excluded from start",
					slog.String("scheduleId", schedule.ID), slog.Int("status", int(broadcast.Status)))
				return nil // 停止中の場合のみ、起動処理を進める
			}
			if broadcast.MediaLiveChannelID == "" {
				slog.Error("Empty media live channel id", slog.String("scheduleId", schedule.ID))
				return fmt.Errorf("unexpected media live channel arn format. arn=%s", broadcast.MediaLiveChannelArn)
			}

			settings := &medialive.CreateScheduleParams{
				ChannelID: broadcast.MediaLiveChannelID,
				Settings:  s.newStartScheduleSettings(schedule, broadcast),
			}
			slog.Info("Calling to create media live schedule", slog.String("scheduleId", schedule.ID))
			if err := s.media.CreateSchedule(ctx, settings); err != nil {
				slog.Error("Failed to create media live schedule",
					slog.String("scheduleId", schedule.ID),
					slog.String("channelId", broadcast.MediaLiveChannelID),
					slog.Any("settings", settings),
					log.Error(err))
				return err
			}
			slog.Info("Succeeded to create media live schedule", slog.String("scheduleId", schedule.ID))

			slog.Info("Calling to start media live", slog.String("scheduleId", schedule.ID))
			if err := s.media.StartChannel(ctx, broadcast.MediaLiveChannelID); err != nil {
				slog.Error("Failed to start media live", slog.String("scheduleId", schedule.ID), log.Error(err))
				return err
			}
			slog.Info("Succeeded to start media live", slog.String("scheduleId", schedule.ID))

			params := &database.UpdateBroadcastParams{
				Status: entity.BroadcastStatusActive,
			}
			return s.db.Broadcast.Update(ctx, broadcast.ID, params)
		})
	}
	return eg.Wait()
}

func (s *starter) newStartScheduleSettings(schedule *sentity.Schedule, broadcast *entity.Broadcast) []*medialive.ScheduleSetting {
	sourceURL, _ := s.storage.ReplaceURLToS3URI(schedule.OpeningVideoURL)
	// ライブ配信開始時は一律、オープニング動画再生から始めるようにする
	return []*medialive.ScheduleSetting{
		{
			Name:       fmt.Sprintf("%s immediate-input-mp4", jst.Format(s.now(), time.DateTime)),
			ActionType: medialive.ScheduleActionTypeInputSwitch,
			StartType:  medialive.ScheduleStartTypeImmediate,
			Reference:  broadcast.MediaLiveMP4InputName,
			Source:     sourceURL,
		},
	}
}

// createChannel - ライブ配信リソースの作成を開始 (30分前)
func (s *starter) createChannel(ctx context.Context, target time.Time) error {
	slog.Debug("Creating channel...", slog.Time("target", target))
	in := &store.ListSchedulesInput{
		StartAtLt: target.Add(30 * time.Minute), // マルシェ開催開始30分前〜
		EndAtGte:  target,                       // 〜マルシェ開催終了
		NoLimit:   true,
	}
	schedules, total, err := s.store.ListSchedules(ctx, in)
	if err != nil || total == 0 {
		return err
	}
	slog.Debug("Got schedules to create", slog.Int64("total", total))

	eg, ectx := errgroup.WithContext(ctx)
	for i := range schedules {
		if err := s.semaphore.Acquire(ctx, 1); err != nil {
			return err
		}

		schedule := schedules[i]
		eg.Go(func() error {
			defer s.semaphore.Release(1)
			broadcast, err := s.db.Broadcast.GetByScheduleID(ectx, schedule.ID)
			if err != nil {
				return err
			}
			if broadcast.Status != entity.BroadcastStatusDisabled {
				slog.Debug("Channels excluded from creation",
					slog.String("scheduleId", schedule.ID), slog.Int("status", int(broadcast.Status)))
				return nil // リソース未作成の場合のみ、作成処理を進める
			}
			payload := &CreatePayload{
				ScheduleID: broadcast.ScheduleID,
				Channel: &CreateChannelPayload{
					Name:                   fmt.Sprintf("%s-%s", s.env, schedule.ID),
					StartTime:              schedule.StartAt.Format(time.RFC3339),
					InputLossImageSlateURI: schedule.ImageURL,
				},
				MP4Input: &CreateMp4InputPayload{
					InputURL: dynamicMP4InputURL,
				},
				RtmpInput: &CreateRtmpInputPayload{
					StreamName: streamName,
				},
				RtmpOutputs: s.createRtmpOuputPayload(broadcast),
				Archive: &CreateArchivePayload{
					BucketName: s.bucketName,
					Path:       newArchiveHLSPath(schedule.ID),
				},
			}
			slog.Info("Calling step function", slog.String("scheduleId", schedule.ID))
			if err := s.sfn.StartExecution(ectx, payload); err != nil {
				slog.Error("Failed step function", slog.String("scheduleId", schedule.ID), log.Error(err))
				return err
			}
			slog.Info("Succeeded step function", slog.String("scheduleId", schedule.ID))
			params := &database.UpdateBroadcastParams{
				Status: entity.BroadcastStatusWaiting,
			}
			return s.db.Broadcast.Update(ctx, broadcast.ID, params)
		})
	}
	return eg.Wait()
}

// createRtmpOuputPayload - 配信リソース(MediaLive RTMP Pushアウトプット)
func (s *starter) createRtmpOuputPayload(broadcast *entity.Broadcast) []*CreateRtmpOutputPayload {
	outputs := make([]*CreateRtmpOutputPayload, 0, 2)
	// Youtube配信設定
	if broadcast.YoutubeStreamKey != "" {
		if broadcast.YoutubeStreamURL != "" {
			payload := &CreateRtmpOutputPayload{
				Name:      "youtube",
				StreamURL: broadcast.YoutubeStreamURL,
				StreamKey: broadcast.YoutubeStreamKey,
			}
			outputs = append(outputs, payload)
		}
		// MediaLive側をシングルパイプライン設定で作成しているため、バックアップ配信は不要
		// if broadcast.YoutubeBackupURL != "" {
		// 	payload := &CreateRtmpOutputPayload{
		// 		Name:      "youtube-backup",
		// 		StreamURL: broadcast.YoutubeBackupURL,
		// 		StreamKey: broadcast.YoutubeStreamKey,
		// 	}
		// 	outputs = append(outputs, payload)
		// }
	}
	return outputs
}
