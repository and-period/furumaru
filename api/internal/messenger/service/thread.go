package service

import (
	"context"
	"errors"
	"fmt"

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

func (s *service) CreateThread(ctx context.Context, in *messenger.CreateThreadInput) (*entity.Thread, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	contactIn := &messenger.GetContactInput{
		ContactID: in.ContactID,
	}
	_, err := s.GetContact(ctx, contactIn)
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid contact id: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, exception.InternalError(err)
	}
	params := &entity.NewThreadParams{
		ContactID: in.ContactID,
		UserType:  in.UserType,
		Content:   in.Content,
	}
	thread := entity.NewThread(params)
	thread.Fill(in.UserID)
	if err := s.db.Thread.Create(ctx, thread); err != nil {
		return nil, exception.InternalError(err)
	}
	return thread, nil
}
