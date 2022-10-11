package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"github.com/and-period/furumaru/api/internal/user"
	"golang.org/x/sync/errgroup"
)

func (s *service) CreateSchedule(ctx context.Context, in *store.CreateScheduleInput) (*entity.Schedule, entity.Lives, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, nil, exception.InternalError(err)
	}
	producerIDs := make([]string, len(in.Lives))
	for i := range in.Lives {
		producerIDs[i] = in.Lives[i].ProducerID
	}
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		in := &user.MultiGetProducersInput{
			ProducerIDs: producerIDs,
		}
		producers, err := s.user.MultiGetProducers(ectx, in)
		if len(producers) != len(producerIDs) {
			return exception.ErrInvalidArgument
		}
		return err
	})
	// eg.Go(func() (err error) {
	// 	for i := range in.Lives {
	// 		in := &store.MultiGetProductsInput{
	// 			ProductIDs: in.Lives[i].Recommends,
	// 		}
	// 		_, err = s.db.Product.MultiGet(ectx, in.ProductIDs) // TODO
	// 	}
	// 	return err
	// })
	err := eg.Wait()
	if err != nil {
		return nil, nil, exception.InternalError(err)
	}
	sparams := &entity.NewScheduleParams{
		Title:        in.Title,
		Description:  in.Description,
		ThumbnailURL: in.ThumbnailURL,
		StartAt:      in.StartAt,
		EndAt:        in.EndAt,
	}
	schedule := entity.NewSchedule(sparams)
	lives := make(entity.Lives, len(in.Lives))
	for i := range in.Lives {
		l := in.Lives[i]
		lparams := &entity.NewLiveParams{
			ScheduleID:  schedule.ID,
			Title:       l.Title,
			Description: l.Description,
			ProducerID:  l.ProducerID,
			StartAt:     l.StartAt,
			EndAt:       l.EndAt,
			Recommends:  l.Recommends,
		}
		lives[i] = entity.NewLive(lparams)
	}
	if err := s.db.Schedule.Create(ctx, schedule, lives); err != nil {
		return nil, nil, exception.InternalError(err)
	}
	return schedule, lives, nil
}
