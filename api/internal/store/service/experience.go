package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) ListExperiences(ctx context.Context, in *store.ListExperiencesInput) (entity.Experiences, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	// TODO: 詳細の実装
	return entity.Experiences{}, 0, nil
}

func (s *service) MultiGetExperiences(ctx context.Context, in *store.MultiGetExperiencesInput) (entity.Experiences, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return entity.Experiences{}, nil
}

func (s *service) MultiGetExperiencesByRevision(
	ctx context.Context, in *store.MultiGetExperiencesByRevisionInput,
) (entity.Experiences, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return entity.Experiences{}, nil
}

func (s *service) GetExperience(ctx context.Context, in *store.GetExperienceInput) (*entity.Experience, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return &entity.Experience{}, nil
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
