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

func (s *service) ListLivesByScheduleID(ctx context.Context, in *store.ListLivesByScheduleIDInput) (entity.Lives, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	lives, err := s.db.Live.ListByScheduleID(ctx, in.ScheduleID)
	return lives, exception.InternalError(err)
}

func (s *service) GetLive(ctx context.Context, in *store.GetLiveInput) (*entity.Live, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	live, err := s.db.Live.Get(ctx, in.LiveID)
	return live, exception.InternalError(err)
}

func (s *service) CreateLive(ctx context.Context, in *store.CreateLiveInput) (*entity.Live, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
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
			return fmt.Errorf("api: unmatch product: %w", exception.ErrInvalidArgument)
		}
		return nil
	})
	err := eg.Wait()
	if errors.Is(err, exception.ErrNotFound) {
		return nil, fmt.Errorf("api: invalid request: %s: %w", err.Error(), exception.ErrInvalidArgument)
	}
	if err != nil {
		return nil, exception.InternalError(err)
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
		return nil, exception.InternalError(err)
	}
	return live, nil
}

func (s *service) UpdateLive(ctx context.Context, in *store.UpdateLiveInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	if _, err := s.db.Live.Get(ctx, in.LiveID); err != nil {
		return exception.InternalError(err)
	}
	products, err := s.db.Product.MultiGet(ctx, in.ProductIDs)
	if err != nil {
		return exception.InternalError(err)
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
	return exception.InternalError(err)
}

func (s *service) DeleteLive(ctx context.Context, in *store.DeleteLiveInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.Live.Delete(ctx, in.LiveID)
	return exception.InternalError(err)
}
