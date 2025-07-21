package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/store"
	"github.com/and-period/furumaru/api/internal/store/entity"
)

func (s *service) SearchPostalCode(
	ctx context.Context,
	in *store.SearchPostalCodeInput,
) (*entity.PostalCode, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	code, err := s.postalCode.Search(ctx, in.PostlCode)
	if err != nil {
		return nil, internalError(err)
	}
	postalCode, err := entity.NewPostalCode(code)
	if err != nil {
		return nil, internalError(err)
	}
	return postalCode, nil
}
