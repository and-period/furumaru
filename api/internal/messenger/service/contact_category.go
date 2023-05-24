package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/messenger"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

func (s *service) GetContactCategory(ctx context.Context, in *messenger.GetContactCategoryInput) (*entity.ContactCategory, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	category, err := s.db.ContactCategory.Get(ctx, in.CategoryID)
	return category, exception.InternalError(err)
}
