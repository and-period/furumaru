package service

import (
	"context"

	"github.com/and-period/marche/api/internal/user/entity"
)

func (s *userService) MultiGetShops(ctx context.Context, in *MultiGetShopsInput) (entity.Shops, error) {
	return nil, ErrNotImplemented
}
