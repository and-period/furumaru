package scheduler

import (
	"context"
	"time"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"golang.org/x/sync/errgroup"
)

func (s *scheduler) dispatchNotification(ctx context.Context, target time.Time) error {
	// 通知対象スケジュール,お知らせ一覧の取得
	var (
		schedules     entity.Schedules
		notifications entity.Notifications
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		params := &database.ListSchedulesParams{
			Types: []entity.ScheduleType{entity.ScheduleTypeNotification},
			Since: jst.BeginningOfDay(target),
			Until: target,
		}
		schedules, err = s.db.Schedule.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		params := &database.ListNotificationsParams{
			Since: jst.BeginningOfDay(target),
			Until: target,
		}
		notifications, err = s.db.Notification.List(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return err
	}

	// 通知スケジュールが取得できないものは生成
	scheduleMap := schedules.Map()
	for _, n := range notifications {
		if _, ok := scheduleMap[n.ID]; ok {
			continue
		}
		params := &entity.NewScheduleParams{
			MessageType: entity.ScheduleTypeNotification,
			MessageID:   n.ID,
			SentAt:      n.PublishedAt,
		}
		schedules = append(schedules, entity.NewSchedule(params))
	}

	// 通知スケジュールにしたがってそれぞれ実行
	eg, ectx = errgroup.WithContext(ctx)
	for i := range schedules {
		if err := s.semaphore.Acquire(ctx, 1); err != nil {
			return err
		}

		schedule := schedules[i]
		eg.Go(func() (err error) {
			defer s.semaphore.Release(1)
			err = s.executeNotification(ectx, schedule)
			return
		})
	}
	return eg.Wait()
}

func (s *scheduler) executeNotification(ctx context.Context, schedule *entity.Schedule) error {
	now := s.now()
	// 通知前処理
	if schedule.ShouldCancel(now) {
		return s.db.Schedule.UpdateCancel(ctx, schedule.MessageType, schedule.MessageID)
	}
	if !schedule.Executable(now) {
		return nil
	}
	if err := s.db.Schedule.UpsertProcessing(ctx, schedule); err != nil {
		return err
	}
	// 通知処理
	in := &messenger.NotifyNotificationInput{
		NotificationID: schedule.MessageID,
	}
	if err := s.messenger.NotifyNotification(ctx, in); err != nil {
		return err
	}
	// 通知後処理
	return s.db.Schedule.UpdateDone(ctx, schedule.MessageType, schedule.MessageID)
}
