package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/database"
	"github.com/and-period/furumaru/api/internal/store/entity"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListExperienceTypes(ctx context.Context, in *store.ListExperienceTypesInput) (entity.ExperienceTypes, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListExperienceTypesParams{
		Name:   in.Name,
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	var (
		types entity.ExperienceTypes
		total int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		types, err = s.db.ExperienceType.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.ExperienceType.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return types, total, nil
}

func (s *service) MultiGetExperienceTypes(ctx context.Context, in *store.MultiGetExperienceTypesInput) (entity.ExperienceTypes, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return entity.ExperienceTypes{}, nil
}

func (s *service) GetExperienceType(ctx context.Context, in *store.GetExperienceTypeInput) (*entity.ExperienceType, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return &entity.ExperienceType{}, nil
}

func (s *service) CreateExperienceType(ctx context.Context, in *store.CreateExperienceTypeInput) (*entity.ExperienceType, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return &entity.ExperienceType{}, nil
}

func (s *service) UpdateExperienceType(ctx context.Context, in *store.UpdateExperienceTypeInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	// TODO: 詳細の実装
	return nil
}

func (s *service) DeleteExperienceType(ctx context.Context, in *store.DeleteExperienceTypeInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	// TODO: 詳細の実装
	return nil
}
