package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListThreadsByContactID(ctx context.Context, in *messenger.ListThreadsByContactIDInput) (entity.Threads, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, exception.InternalError(err)
	}
	params := &database.ListThreadsByContactIDParams{
		ContactID: in.ContactID,
		Limit:     int(in.Limit),
		Offset:    int(in.Offset),
	}
	var (
		threads entity.Threads
		total   int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		threads, err = s.db.Thread.ListByContactID(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Thread.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, exception.InternalError(err)
	}

	return threads, total, nil
}

func (s *service) GetThread(ctx context.Context, in *messenger.GetThreadInput) (*entity.Thread, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}

	thread, err := s.db.Thread.Get(ctx, in.ThreadID)
	if err != nil {
		return nil, exception.InternalError(err)
	}

	return thread, exception.InternalError(err)
}
