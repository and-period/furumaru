package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/pkg/log"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListNotifications(ctx context.Context, in *messenger.ListNotificationsInput) (entity.Notifications, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	orders, err := s.newListNotificationsOrders(in.Orders)
	if err != nil {
		return nil, 0, fmt.Errorf("service: invalid list notifications orders: err=%s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	params := &database.ListNotificationsParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
		Since:  in.Since,
		Until:  in.Until,
		Orders: orders,
	}
	var (
		notifications entity.Notifications
		total         int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		notifications, err = s.db.Notification.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Notification.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return notifications, total, nil
}

func (s *service) newListNotificationsOrders(in []*messenger.ListNotificationsOrder) ([]*database.ListNotificationsOrder, error) {
	res := make([]*database.ListNotificationsOrder, len(in))
	for i := range in {
		var key database.ListNotificationsOrderKey
		switch in[i].Key {
		case messenger.ListNotificationsOrderByTitle:
			key = database.ListNotificationsOrderByTitle
		case messenger.ListNotificationsOrderByPublishedAt:
			key = database.ListNotificationsOrderByPublishedAt
		default:
			return nil, errors.New("service: invalid list notifications order key")
		}
		res[i] = &database.ListNotificationsOrder{
			Key:        key,
			OrderByASC: in[i].OrderByASC,
		}
	}
	return res, nil
}

func (s *service) GetNotification(ctx context.Context, in *messenger.GetNotificationInput) (*entity.Notification, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	notification, err := s.db.Notification.Get(ctx, in.NotificationID)
	return notification, internalError(err)
}

func (s *service) CreateNotification(
	ctx context.Context, in *messenger.CreateNotificationInput,
) (*entity.Notification, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		req := &user.GetAdminInput{
			AdminID: in.CreatedBy,
		}
		_, err := s.user.GetAdmin(ectx, req)
		return err
	})
	eg.Go(func() error {
		if in.Type != entity.NotificationTypePromotion {
			return nil
		}
		req := &store.GetPromotionInput{
			PromotionID: in.PromotionID,
		}
		_, err := s.store.GetPromotion(ectx, req)
		return err
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("service: not found reference: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewNotificationParams{
		Type:        in.Type,
		Targets:     in.Targets,
		Title:       in.Title,
		Body:        in.Body,
		Note:        in.Note,
		PublishedAt: in.PublishedAt,
		PromotionID: in.PromotionID,
		CreatedBy:   in.CreatedBy,
	}
	notification := entity.NewNotification(params)
	if err := notification.Validate(s.now()); err != nil {
		return nil, fmt.Errorf("service: invalid notification: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err := s.db.Notification.Create(ctx, notification); err != nil {
		return nil, internalError(err)
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		in := &messenger.ReserveNotificationInput{
			NotificationID: notification.ID,
		}
		if err := s.ReserveNotification(context.Background(), in); err != nil {
			slog.Error("Failed to reserve notification", slog.String("notificationId", notification.ID), log.Error(err))
		}
	}()
	return notification, nil
}

func (s *service) UpdateNotification(ctx context.Context, in *messenger.UpdateNotificationInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	adminIn := &user.GetAdminInput{
		AdminID: in.UpdatedBy,
	}
	_, err := s.user.GetAdmin(ctx, adminIn)
	if errors.Is(err, exception.ErrNotFound) {
		return fmt.Errorf("api: invalid admin id format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return internalError(err)
	}
	notification, err := s.db.Notification.Get(ctx, in.NotificationID)
	if err != nil {
		return internalError(err)
	}
	if s.now().After(notification.PublishedAt) {
		// すでに投稿済みの場合は更新できない
		return fmt.Errorf("api: already notified: %w", exception.ErrFailedPrecondition)
	}
	notification.Targets = in.Targets
	notification.Title = in.Title
	notification.Body = in.Body
	notification.Note = in.Note
	notification.PublishedAt = in.PublishedAt
	if err := notification.Validate(s.now()); err != nil {
		return fmt.Errorf("api: invalid notification: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	params := &database.UpdateNotificationParams{
		Targets:     in.Targets,
		Title:       in.Title,
		Body:        in.Body,
		Note:        in.Note,
		PublishedAt: in.PublishedAt,
		UpdatedBy:   in.UpdatedBy,
	}
	if err := s.db.Notification.Update(ctx, in.NotificationID, params); err != nil {
		return internalError(err)
	}
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		in := &messenger.ReserveNotificationInput{
			NotificationID: notification.ID,
		}
		if err := s.ReserveNotification(context.Background(), in); err != nil {
			slog.Error("Failed to reserve notification", slog.String("notificationId", notification.ID), log.Error(err))
		}
	}()
	return nil
}

func (s *service) DeleteNotification(ctx context.Context, in *messenger.DeleteNotificationInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Notification.Delete(ctx, in.NotificationID)
	return internalError(err)
}
