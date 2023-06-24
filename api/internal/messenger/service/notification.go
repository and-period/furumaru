package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/user"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListNotifications(ctx context.Context, in *messenger.ListNotificationsInput) (entity.Notifications, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	orders := make([]*database.ListNotificationsOrder, len(in.Orders))
	for i := range in.Orders {
		orders[i] = &database.ListNotificationsOrder{
			Key:        in.Orders[i].Key,
			OrderByASC: in.Orders[i].OrderByASC,
		}
	}
	params := &database.ListNotificationsParams{
		Limit:         int(in.Limit),
		Offset:        int(in.Offset),
		Since:         in.Since,
		Until:         in.Until,
		OnlyPublished: in.OnlyPublished,
		Orders:        orders,
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
		return nil, 0, exception.InternalError(err)
	}
	return notifications, total, nil
}

func (s *service) GetNotification(ctx context.Context, in *messenger.GetNotificationInput) (*entity.Notification, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	notification, err := s.db.Notification.Get(ctx, in.NotificationID)
	return notification, exception.InternalError(err)
}

func (s *service) CreateNotification(
	ctx context.Context, in *messenger.CreateNotificationInput,
) (*entity.Notification, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
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
		return nil, fmt.Errorf("api: not found reference: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, exception.InternalError(err)
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
		return nil, fmt.Errorf("api: invalid notification: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err := s.db.Notification.Create(ctx, notification); err != nil {
		return nil, exception.InternalError(err)
	}
	return notification, nil
}

func (s *service) UpdateNotification(ctx context.Context, in *messenger.UpdateNotificationInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	adminIn := &user.GetAdminInput{
		AdminID: in.UpdatedBy,
	}
	_, err := s.user.GetAdmin(ctx, adminIn)
	if errors.Is(err, exception.ErrNotFound) {
		return fmt.Errorf("api: invalid admin id format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return exception.InternalError(err)
	}
	notification, err := s.db.Notification.Get(ctx, in.NotificationID)
	if err != nil {
		return exception.InternalError(err)
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
	err = s.db.Notification.Update(ctx, in.NotificationID, params)
	return exception.InternalError(err)
}

func (s *service) DeleteNotification(ctx context.Context, in *messenger.DeleteNotificationInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Notification.Delete(ctx, in.NotificationID)
	return exception.InternalError(err)
}
