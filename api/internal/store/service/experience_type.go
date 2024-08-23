package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) ListExperienceTypes(ctx context.Context, in *store.ListExperienceTypesInput) (entity.ExperienceTypes, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	// TODO: 詳細の実装
	return entity.ExperienceTypes{}, 0, nil
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
