package service

import (
	"context"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/cognito"
)

func (s *userService) SignInAdmin(ctx context.Context, in *SignInAdminInput) (*entity.AdminAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	rs, err := s.adminAuth.SignIn(ctx, in.Key, in.Password)
	if err != nil {
		return nil, userError(err)
	}
	auth, err := s.getAdminAuth(ctx, rs)
	return auth, userError(err)
}

func (s *userService) SignOutAdmin(ctx context.Context, in *SignOutAdminInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}
	err := s.adminAuth.SignOut(ctx, in.AccessToken)
	return userError(err)
}

func (s *userService) GetAdminAuth(ctx context.Context, in *GetAdminAuthInput) (*entity.AdminAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	rs := &cognito.AuthResult{AccessToken: in.AccessToken}
	auth, err := s.getAdminAuth(ctx, rs)
	return auth, userError(err)
}

func (s *userService) RefreshAdminToken(ctx context.Context, in *RefreshAdminTokenInput) (*entity.AdminAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	rs, err := s.adminAuth.RefreshToken(ctx, in.RefreshToken)
	if err != nil {
		return nil, userError(err)
	}
	auth, err := s.getAdminAuth(ctx, rs)
	return auth, userError(err)
}

func (s *userService) getAdminAuth(ctx context.Context, rs *cognito.AuthResult) (*entity.AdminAuth, error) {
	username, err := s.adminAuth.GetUsername(ctx, rs.AccessToken)
	if err != nil {
		return nil, err
	}
	out, err := s.db.Admin.GetByCognitoID(ctx, username, "id", "role")
	if err != nil {
		return nil, err
	}
	auth := entity.NewAdminAuth(out.ID, out.Role, rs)
	return auth, nil
}
