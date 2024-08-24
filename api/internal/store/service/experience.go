package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListExperiences(ctx context.Context, in *store.ListExperiencesInput) (entity.Experiences, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListExperiencesParams{
		Name:          in.Name,
		CoordinatorID: in.CoordinatorID,
		ProducerID:    in.ProducerID,
		Limit:         int(in.Limit),
		Offset:        int(in.Offset),
	}
	var (
		experiences entity.Experiences
		total       int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		experiences, err = s.db.Experience.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.Experience.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return experiences, total, nil
}

func (s *service) MultiGetExperiences(ctx context.Context, in *store.MultiGetExperiencesInput) (entity.Experiences, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	experiences, err := s.db.Experience.MultiGet(ctx, in.ExperienceIDs)
	return experiences, internalError(err)
}

func (s *service) MultiGetExperiencesByRevision(
	ctx context.Context, in *store.MultiGetExperiencesByRevisionInput,
) (entity.Experiences, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	experiences, err := s.db.Experience.MultiGetByRevision(ctx, in.ExperienceRevisionIDs)
	return experiences, internalError(err)
}

func (s *service) GetExperience(ctx context.Context, in *store.GetExperienceInput) (*entity.Experience, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	experience, err := s.db.Experience.Get(ctx, in.ExperienceID)
	return experience, internalError(err)
}

func (s *service) CreateExperience(ctx context.Context, in *store.CreateExperienceInput) (*entity.Experience, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return &entity.Experience{}, nil
}

func (s *service) UpdateExperience(ctx context.Context, in *store.UpdateExperienceInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	// TODO: 詳細の実装
	return nil
}

func (s *service) DeleteExperience(ctx context.Context, in *store.DeleteExperienceInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	// TODO: 詳細の実装
	return nil
}
