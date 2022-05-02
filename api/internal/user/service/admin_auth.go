package service

import (
	"context"

	"github.com/and-period/marche/api/internal/user/entity"
)

func (s *userService) SignInAdmin(ctx context.Context, in *SignInAdminInput) (*entity.AdminAuth, error) {
	return nil, ErrNotImplemented
}

func (s *userService) SignOutAdmin(ctx context.Context, in *SignOutAdminInput) error {
	return ErrNotImplemented
}

func (s *userService) GetAdminAuth(ctx context.Context, in *GetAdminAuthInput) (*entity.AdminAuth, error) {
	return nil, ErrNotImplemented
}

func (s *userService) RefreshAdminToken(ctx context.Context, in *RefreshAdminTokenInput) (*entity.AdminAuth, error) {
	return nil, ErrNotImplemented
}

// func (s *userService) getAdminAuth(ctx context.Context, rs *cognito.AuthResult) (*entity.AdminAuth, error) {
// 	username, err := s.adminAuth.GetUsername(ctx, rs.AccessToken)
// 	if err != nil {
// 		return nil, err
// 	}
// 	out, err := s.db.Admin.GetByCognitoID(ctx, username, "id", "role")
// 	if err != nil {
// 		return nil, err
// 	}
// 	auth := entity.NewAdminAuth(out.ID, out.Role, rs)
// 	return auth, nil
// }
