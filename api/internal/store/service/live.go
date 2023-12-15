package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListLives(ctx context.Context, in *store.ListLivesInput) (entity.Lives, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListLivesParams{
		ScheduleIDs: in.ScheduleIDs,
		ProducerID:  in.ProducerID,
		Limit:       int(in.Limit),
		Offset:      int(in.Offset),
	}
	var (
		lives entity.Lives
		total int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		lives, err = s.db.Live.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Live.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	if len(lives) == 0 || !in.OnlyPublished {
		return lives, total, nil
	}
	products, err := s.db.Product.MultiGet(ctx, lives.ProductIDs())
	if err != nil {
		return nil, 0, internalError(err)
	}
	lives.ExcludeProductIDs(products.FilterByPublished().Map())
	return lives, total, nil
}

func (s *service) GetLive(ctx context.Context, in *store.GetLiveInput) (*entity.Live, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	live, err := s.db.Live.Get(ctx, in.LiveID)
	return live, internalError(err)
}

func (s *service) CreateLive(ctx context.Context, in *store.CreateLiveInput) (*entity.Live, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewLiveParams{
		ScheduleID: in.ScheduleID,
		ProducerID: in.ProducerID,
		ProductIDs: in.ProductIDs,
		Comment:    in.Comment,
		StartAt:    in.StartAt,
		EndAt:      in.EndAt,
	}
	live := entity.NewLive(params)
	if err := s.validateLive(ctx, live); err != nil {
		return nil, err
	}
	if err := s.db.Live.Create(ctx, live); err != nil {
		return nil, internalError(err)
	}
	return live, nil
}

func (s *service) UpdateLive(ctx context.Context, in *store.UpdateLiveInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	live, err := s.db.Live.Get(ctx, in.LiveID)
	if err != nil {
		return internalError(err)
	}
	live.StartAt = in.StartAt
	live.EndAt = in.EndAt
	if err := s.validateLive(ctx, live); err != nil {
		return err
	}
	params := &database.UpdateLiveParams{
		ProductIDs: in.ProductIDs,
		Comment:    in.Comment,
		StartAt:    in.StartAt,
		EndAt:      in.EndAt,
	}
	err = s.db.Live.Update(ctx, in.LiveID, params)
	return internalError(err)
}

func (s *service) DeleteLive(ctx context.Context, in *store.DeleteLiveInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Live.Delete(ctx, in.LiveID)
	return internalError(err)
}

func (s *service) validateLive(ctx context.Context, live *entity.Live) error {
	var (
		schedule *entity.Schedule
		lives    entity.Lives
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		schedule, err = s.db.Schedule.Get(ectx, live.ScheduleID)
		return
	})
	eg.Go(func() (err error) {
		params := &database.ListLivesParams{
			ScheduleIDs: []string{live.ScheduleID},
		}
		lives, err = s.db.Live.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		in := &user.GetProducerInput{
			ProducerID: live.ProducerID,
		}
		_, err = s.user.GetProducer(ectx, in)
		return
	})
	eg.Go(func() error {
		products, err := s.db.Product.MultiGet(ectx, live.ProductIDs)
		if err != nil {
			return err
		}
		if len(products) != len(live.ProductIDs) {
			return errUnmatchProducts
		}
		return nil
	})
	err := eg.Wait()
	if errors.Is(err, database.ErrNotFound) || errors.Is(err, exception.ErrNotFound) || errors.Is(err, errUnmatchProducts) {
		return fmt.Errorf("api: invalid request: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return internalError(err)
	}
	if err := live.Validate(schedule, lives); err != nil {
		return fmt.Errorf("api: invalid live schedule: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	return nil
}
