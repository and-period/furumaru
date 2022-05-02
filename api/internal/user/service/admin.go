package service

import (
	"context"

	"github.com/and-period/marche/api/internal/user/entity"
)

func (s *userService) GetAdmin(ctx context.Context, in *GetAdminInput) (*entity.Admin, error) {
	return nil, ErrNotImplemented
}

func (s *userService) UpdateAdminEmail(ctx context.Context, in *UpdateAdminEmailInput) error {
	return ErrNotImplemented
}

func (s *userService) VerifyAdminEmail(ctx context.Context, in *VerifyAdminEmailInput) error {
	return ErrNotImplemented
}

func (s *userService) UpdateAdminPassword(ctx context.Context, in *UpdateAdminPasswordInput) error {
	return ErrNotImplemented
}
