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

func (s *service) ListThreads(ctx context.Context, in *messenger.ListThreadsInput) (entity.Threads, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListThreadsParams{
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
		threads, err = s.db.Thread.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Thread.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	updateReadParams := &database.UpdateContactReadParams{
		ContactID: in.ContactID,
		UserID:    in.UserID,
		Read:      true,
	}
	err := s.db.ContactRead.Update(ctx, updateReadParams)
	if err != nil {
		return nil, 0, internalError(err)
	}
	return threads, total, nil
}

func (s *service) GetThread(ctx context.Context, in *messenger.GetThreadInput) (*entity.Thread, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	thread, err := s.db.Thread.Get(ctx, in.ThreadID)
	return thread, internalError(err)
}

func (s *service) CreateThread(ctx context.Context, in *messenger.CreateThreadInput) (*entity.Thread, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	_, err := s.db.Contact.Get(ctx, in.ContactID)
	if errors.Is(err, database.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid contact id: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewThreadParams{
		ContactID: in.ContactID,
		UserType:  in.UserType,
		Content:   in.Content,
	}
	thread := entity.NewThread(params)
	thread.Fill(in.UserID)
	if err := s.db.Thread.Create(ctx, thread); err != nil {
		return nil, internalError(err)
	}
	_, err = s.db.ContactRead.GetByContactIDAndUserID(ctx, in.ContactID, in.UserID)
	if errors.Is(err, database.ErrNotFound) {
		params := &entity.NewContactReadParams{
			ContactID: in.ContactID,
			UserID:    in.UserID,
			UserType:  entity.ContactUserType(in.UserType), // FIXME: ThreadUserTypeと統合
			Read:      false,
		}
		contactRead, err := entity.NewContactRead(params)
		if err != nil {
			return nil, internalError(err)
		}
		if err := s.db.ContactRead.Create(ctx, contactRead); err != nil {
			return nil, internalError(err)
		}
	}
	return thread, nil
}

func (s *service) UpdateThread(ctx context.Context, in *messenger.UpdateThreadInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}

	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if in.UserType != 1 {
			return nil
		}
		adminID := &user.GetAdminInput{
			AdminID: in.UserID,
		}
		_, err := s.user.GetAdmin(ectx, adminID)
		return err
	})
	eg.Go(func() error {
		if in.UserType != 2 {
			return nil
		}
		userID := &user.GetUserInput{
			UserID: in.UserID,
		}
		_, err := s.user.GetUser(ectx, userID)
		return err
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		return fmt.Errorf("api: invalid user id format: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return internalError(err)
	}
	params := &database.UpdateThreadParams{
		Content:  in.Content,
		UserID:   in.UserID,
		UserType: in.UserType,
	}
	err = s.db.Thread.Update(ctx, in.ThreadID, params)
	return internalError(err)
}

func (s *service) DeleteThread(ctx context.Context, in *messenger.DeleteThreadInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Thread.Delete(ctx, in.ThreadID)
	return internalError(err)
}
