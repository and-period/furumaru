package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/database"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

func (s *service) ListContactCategories(ctx context.Context, in *messenger.ListContactCategoriesInput) (entity.ContactCategories, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	params := &database.ListContactCategoriesParams{
		Limit:  int(in.Limit),
		Offset: int(in.Offset),
	}
	categories, err := s.db.ContactCategory.List(ctx, params)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	return categories, nil
}

func (s *service) GetContactCategory(ctx context.Context, in *messenger.GetContactCategoryInput) (*entity.ContactCategory, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	category, err := s.db.ContactCategory.Get(ctx, in.CategoryID)
	return category, exception.InternalError(err)
}
