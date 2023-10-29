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
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		_, err = s.db.Schedule.Get(ectx, in.ScheduleID)
		return
	})
	eg.Go(func() (err error) {
		in := &user.GetProducerInput{
			ProducerID: in.ProducerID,
		}
		_, err = s.user.GetProducer(ectx, in)
		return
	})
	eg.Go(func() error {
		products, err := s.db.Product.MultiGet(ectx, in.ProductIDs)
		if err != nil {
			return err
		}
		if len(products) != len(in.ProductIDs) {
			return errUnmatchProducts
		}
		return nil
	})
	err := eg.Wait()
	if errors.Is(err, database.ErrNotFound) || errors.Is(err, exception.ErrNotFound) || errors.Is(err, errUnmatchProducts) {
		return nil, fmt.Errorf("api: invalid request: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
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
	if err := s.db.Live.Create(ctx, live); err != nil {
		return nil, internalError(err)
	}
	return live, nil
}

func (s *service) UpdateLive(ctx context.Context, in *store.UpdateLiveInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	if _, err := s.db.Live.Get(ctx, in.LiveID); err != nil {
		return internalError(err)
	}
	products, err := s.db.Product.MultiGet(ctx, in.ProductIDs)
	if err != nil {
		return internalError(err)
	}
	if len(products) != len(in.ProductIDs) {
		return fmt.Errorf("api: umatch product: %w", exception.ErrInvalidArgument)
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
