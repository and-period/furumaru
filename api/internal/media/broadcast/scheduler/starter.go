package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"github.com/and-period/furumaru/api/internal/store"
	sentity "github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/and-period/furumaru/api/pkg/medialive"
	"github.com/and-period/furumaru/api/pkg/sfn"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type starter struct {
	now        func() time.Time
	logger     *zap.Logger
	waitGroup  *sync.WaitGroup
	semaphore  *semaphore.Weighted
	db         *database.Database
	store      store.Service
	sfn        sfn.StepFunction
	media      medialive.MediaLive
	env        string
	bucketName string
}

func NewStarter(params *Params, opts ...Option) Scheduler {
	dopts := &options{
		logger:      zap.NewNop(),
		concurrency: 2,
	}
	for i := range opts {
		opts[i](dopts)
	}
	return &starter{
		now:        jst.Now,
		logger:     dopts.logger,
		waitGroup:  params.WaitGroup,
		semaphore:  semaphore.NewWeighted(dopts.concurrency),
		db:         params.Database,
		store:      params.Store,
		sfn:        params.StepFunction,
		media:      params.MediaLive,
		env:        params.Environment,
		bucketName: params.ArchiveBucketName,
	}
}

func (s *starter) Lambda(ctx context.Context) (err error) {
	s.logger.Debug("Started Lambda function", zap.Time("now", s.now()))
	defer func() {
		s.logger.Debug("Finished Lambda function", zap.Time("now", s.now()), zap.Error(err))
	}()

	return s.run(ctx, s.now())
}

func (s *starter) Run(ctx context.Context, target time.Time) error {
	if err := s.run(ctx, target); err != nil {
		s.logger.Error("Failed to run", zap.Time("target", target), zap.Error(err))
		return err
	}
	return nil
}

// run - ライブ配信のリソース作成と開始処理
func (s *starter) run(ctx context.Context, target time.Time) error {
	if err := s.startChannel(ctx, target); err != nil {
		s.logger.Error("Failed to start channel", zap.Time("target", target), zap.Error(err))
		return err
	}
	if err := s.createChannel(ctx, target); err != nil {
		s.logger.Error("Failed to create channel", zap.Time("target", target), zap.Error(err))
		return err
	}
	return nil
}

// startChannel - ライブ配信を開始 (5分前)
func (s *starter) startChannel(ctx context.Context, target time.Time) error {
	s.logger.Debug("Starting channel...", zap.Time("target", target))
	in := &store.ListSchedulesInput{
		StartAtLt: target.Add(5 * time.Minute), // マルシェ開催開始5分前〜
		EndAtGte:  target,                      // 〜マルシェ開催終了
		NoLimit:   true,
	}
	schedules, total, err := s.store.ListSchedules(ctx, in)
	if err != nil || total == 0 {
		return err
	}
	s.logger.Debug("Got schedules to start", zap.Int64("total", total))

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
				s.logger.Debug("Channels excluded from start",
					zap.String("scheduleId", schedule.ID), zap.Int("status", int(broadcast.Status)))
				return nil // 停止中の場合のみ、起動処理を進める
			}
			if broadcast.MediaLiveChannelID == "" {
				s.logger.Error("Empty media live channel id", zap.String("scheduleId", schedule.ID))
				return fmt.Errorf("unexpected media live channel arn format. arn=%s", broadcast.MediaLiveChannelArn)
			}

			settings := &medialive.CreateScheduleParams{
				ChannelID: broadcast.MediaLiveChannelID,
				Settings:  s.newStartScheduleSettings(schedule, broadcast),
			}
			s.logger.Info("Calling to create media live schedule", zap.String("scheduleId", schedule.ID))
			if err := s.media.CreateSchedule(ctx, settings); err != nil {
				s.logger.Error("Failed to create media live schedule",
					zap.String("scheduleId", schedule.ID),
					zap.String("channelId", broadcast.MediaLiveChannelID),
					zap.Any("settings", settings),
					zap.Error(err))
				return err
			}
			s.logger.Info("Succeeded to create media live schedule", zap.String("scheduleId", schedule.ID))

			s.logger.Info("Calling to start media live", zap.String("scheduleId", schedule.ID))
			if err := s.media.StartChannel(ctx, broadcast.MediaLiveChannelID); err != nil {
				s.logger.Error("Failed to start media live", zap.String("scheduleId", schedule.ID), zap.Error(err))
				return err
			}
			s.logger.Info("Succeeded to start media live", zap.String("scheduleId", schedule.ID))

			params := &database.UpdateBroadcastParams{
				Status: entity.BroadcastStatusActive,
			}
			return s.db.Broadcast.Update(ctx, broadcast.ID, params)
		})
	}
	return eg.Wait()
}

func (s *starter) newStartScheduleSettings(schedule *sentity.Schedule, broadcast *entity.Broadcast) []*medialive.ScheduleSetting {
	switch {
	case s.now().After(schedule.StartAt.Add(-1 * time.Minute)): // ライブ配信開始まで1分を切っている
		return []*medialive.ScheduleSetting{{
			Name:       fmt.Sprintf("immediate-input-rtmp %s", jst.Format(s.now(), time.DateTime)),
			ActionType: medialive.ScheduleActionTypeInputSwitch,
			StartType:  medialive.ScheduleStartTypeImmediate,
			Reference:  broadcast.MediaLiveRTMPInputName,
		}}
	case s.now().After(schedule.StartAt.Add(-6 * time.Minute)): // ライブ配信開始まで6分を切っている
		return []*medialive.ScheduleSetting{
			{
				Name:       fmt.Sprintf("immediate-input-mp4 %s", jst.Format(s.now(), time.DateTime)),
				ActionType: medialive.ScheduleActionTypeInputSwitch,
				StartType:  medialive.ScheduleStartTypeImmediate,
				Reference:  broadcast.MediaLiveMP4InputName,
			},
			{
				Name:       fmt.Sprintf("fixed-input-rtmp %s", jst.Format(schedule.StartAt, time.DateTime)),
				ActionType: medialive.ScheduleActionTypeInputSwitch,
				StartType:  medialive.ScheduleStartTypeFixed,
				ExecutedAt: schedule.StartAt,
				Reference:  broadcast.MediaLiveRTMPInputName,
			},
		}
	default: // ライブ配信開始まで6分以上時間あり
		return []*medialive.ScheduleSetting{
			{
				Name:       fmt.Sprintf("immediate-input-rtmp %s", jst.Format(s.now(), time.DateTime)),
				ActionType: medialive.ScheduleActionTypeInputSwitch,
				StartType:  medialive.ScheduleStartTypeImmediate,
				Reference:  broadcast.MediaLiveRTMPInputName,
			},
			{
				Name:       fmt.Sprintf("fixed-input-mp4 %s", jst.Format(schedule.StartAt.Add(-5*time.Minute), time.DateTime)),
				ActionType: medialive.ScheduleActionTypeInputSwitch,
				StartType:  medialive.ScheduleStartTypeFixed,
				ExecutedAt: schedule.StartAt.Add(-5 * time.Minute),
				Reference:  broadcast.MediaLiveMP4InputName,
			},
			{
				Name:       fmt.Sprintf("fixed-input-rtmp %s", jst.Format(schedule.StartAt, time.DateTime)),
				ActionType: medialive.ScheduleActionTypeInputSwitch,
				StartType:  medialive.ScheduleStartTypeFixed,
				ExecutedAt: schedule.StartAt,
				Reference:  broadcast.MediaLiveRTMPInputName,
			},
		}
	}
}

// createChannel - ライブ配信リソースの作成を開始 (30分前)
func (s *starter) createChannel(ctx context.Context, target time.Time) error {
	s.logger.Debug("Creating channel...", zap.Time("target", target))
	in := &store.ListSchedulesInput{
		StartAtLt: target.Add(30 * time.Minute), // マルシェ開催開始30分前〜
		EndAtGte:  target,                       // 〜マルシェ開催終了
		NoLimit:   true,
	}
	schedules, total, err := s.store.ListSchedules(ctx, in)
	if err != nil || total == 0 {
		return err
	}
	s.logger.Debug("Got schedules to create", zap.Int64("total", total))

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
				s.logger.Debug("Channels excluded from creation",
					zap.String("scheduleId", schedule.ID), zap.Int("status", int(broadcast.Status)))
				return nil // リソース未作成の場合のみ、作成処理を進める
			}
			payload := &CreatePayload{
				ScheduleID: broadcast.ScheduleID,
				ChannelInput: &CreateChannelPayload{
					Name:                   fmt.Sprintf("%s-%s", s.env, schedule.ID),
					StartTime:              schedule.StartAt.Format(time.RFC3339),
					InputLossImageSlateURI: schedule.ImageURL,
				},
				MP4Input: &CreateMp4Payload{
					OpeningVideoURL: schedule.OpeningVideoURL,
				},
				RtmpInput: &CreateRtmpPayload{
					StreamName: streamName,
				},
				ArchiveInput: &CreateArchivePayload{
					BucketName: s.bucketName,
					Path:       newArchiveHLSPath(schedule.ID),
				},
			}
			s.logger.Info("Calling step function", zap.String("scheduleId", schedule.ID))
			if err := s.sfn.StartExecution(ectx, payload); err != nil {
				s.logger.Error("Failed step function", zap.String("scheduleId", schedule.ID), zap.Error(err))
				return err
			}
			s.logger.Info("Succeeded step function", zap.String("scheduleId", schedule.ID))
			params := &database.UpdateBroadcastParams{
				Status: entity.BroadcastStatusWaiting,
			}
			return s.db.Broadcast.Update(ctx, broadcast.ID, params)
		})
	}
	return eg.Wait()
}
