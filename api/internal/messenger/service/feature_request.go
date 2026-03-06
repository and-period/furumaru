package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListFeatureRequests(ctx context.Context, in *messenger.ListFeatureRequestsInput) (entity.FeatureRequests, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}

	params := &database.ListFeatureRequestsParams{
		SubmittedBy: in.SubmittedBy,
		Limit:       int(in.Limit),
		Offset:      int(in.Offset),
	}
	var (
		featureRequests entity.FeatureRequests
		total           int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		featureRequests, err = s.db.FeatureRequest.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.FeatureRequest.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return featureRequests, total, nil
}

func (s *service) GetFeatureRequest(ctx context.Context, in *messenger.GetFeatureRequestInput) (*entity.FeatureRequest, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	featureRequest, err := s.db.FeatureRequest.Get(ctx, in.FeatureRequestID)
	return featureRequest, internalError(err)
}

func (s *service) CreateFeatureRequest(ctx context.Context, in *messenger.CreateFeatureRequestInput) (*entity.FeatureRequest, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewFeatureRequestParams{
		Title:       in.Title,
		Description: in.Description,
		Category:    in.Category,
		Priority:    in.Priority,
		SubmittedBy: in.SubmittedBy,
	}
	featureRequest := entity.NewFeatureRequest(params)
	if err := s.db.FeatureRequest.Create(ctx, featureRequest); err != nil {
		return nil, internalError(err)
	}
	return featureRequest, nil
}

func (s *service) UpdateFeatureRequest(ctx context.Context, in *messenger.UpdateFeatureRequestInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	if _, err := s.db.FeatureRequest.Get(ctx, in.FeatureRequestID); err != nil {
		return internalError(err)
	}
	params := &database.UpdateFeatureRequestParams{
		Status: in.Status,
		Note:   in.Note,
	}
	err := s.db.FeatureRequest.Update(ctx, in.FeatureRequestID, params)
	return internalError(err)
}

func (s *service) DeleteFeatureRequest(ctx context.Context, in *messenger.DeleteFeatureRequestInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.FeatureRequest.Delete(ctx, in.FeatureRequestID)
	return internalError(err)
}
