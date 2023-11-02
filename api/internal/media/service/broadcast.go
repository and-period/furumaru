package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/media"
	"github.com/and-period/furumaru/api/internal/media/database"
	"github.com/and-period/furumaru/api/internal/media/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListBroadcasts(ctx context.Context, in *media.ListBroadcastsInput) (entity.Broadcasts, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	orders := make([]*database.ListBroadcastsOrder, len(in.Orders))
	for i := range in.Orders {
		orders[i] = &database.ListBroadcastsOrder{
			Key:        in.Orders[i].Key,
			OrderByASC: in.Orders[i].OrderByASC,
		}
	}
	params := &database.ListBroadcastsParams{
		ScheduleIDs:  in.ScheduleIDs,
		OnlyArchived: in.OnlyArchived,
		Limit:        int(in.Limit),
		Offset:       int(in.Offset),
		Orders:       orders,
	}
	var (
		broadcasts entity.Broadcasts
		total      int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		broadcasts, err = s.db.Broadcast.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Broadcast.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return broadcasts, total, nil
}

func (s *service) GetBroadcastByScheduleID(ctx context.Context, in *media.GetBroadcastByScheduleIDInput) (*entity.Broadcast, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	broadcast, err := s.db.Broadcast.GetByScheduleID(ctx, in.ScheduleID)
	return broadcast, internalError(err)
}

func (s *service) CreateBroadcast(ctx context.Context, in *media.CreateBroadcastInput) (*entity.Broadcast, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewBroadcastParams{
		ScheduleID: in.ScheduleID,
	}
	broadcast := entity.NewBroadcast(params)
	if err := s.db.Broadcast.Create(ctx, broadcast); err != nil {
		return nil, internalError(err)
	}
	return broadcast, nil
}
