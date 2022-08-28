package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
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

func (s *service) CreateNotification(
	ctx context.Context, in *messenger.CreateNotificationInput,
) (*entity.Notification, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	adminID := &user.GetAdminInput{
		AdminID: in.CreatedBy,
	}
	admin, err := s.user.GetAdmin(ctx, adminID)
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid admin id format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}

	if err != nil {
		return nil, exception.InternalError(err)
	}

	params := &entity.NewNotificationParams{
		CreatedBy:   admin.ID,
		CreatorName: admin.Name(),
		UpdatedBy:   admin.ID,
		Title:       in.Title,
		Body:        in.Body,
		Targets:     in.Targets,
		Public:      in.Public,
		PublishedAt: in.PublishedAt,
	}
	notification := entity.NewNotification(params)
	if err := s.db.Notification.Create(ctx, notification); err != nil {
		return nil, exception.InternalError(err)
	}
	return notification, nil
}
