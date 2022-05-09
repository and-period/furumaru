package service

import (
	"context"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/cognito"
)

func (s *userService) SignInUser(ctx context.Context, in *SignInUserInput) (*entity.UserAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	rs, err := s.userAuth.SignIn(ctx, in.Key, in.Password)
	if err != nil {
		return nil, userError(err)
	}
	auth, err := s.getUserAuth(ctx, rs)
	return auth, userError(err)
}

func (s *userService) SignOutUser(ctx context.Context, in *SignOutUserInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}
	err := s.userAuth.SignOut(ctx, in.AccessToken)
	return userError(err)
}

func (s *userService) GetUserAuth(ctx context.Context, in *GetUserAuthInput) (*entity.UserAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	rs := &cognito.AuthResult{AccessToken: in.AccessToken}
	auth, err := s.getUserAuth(ctx, rs)
	return auth, userError(err)
}

func (s *userService) RefreshUserToken(ctx context.Context, in *RefreshUserTokenInput) (*entity.UserAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	rs, err := s.userAuth.RefreshToken(ctx, in.RefreshToken)
	if err != nil {
		return nil, userError(err)
	}
	auth, err := s.getUserAuth(ctx, rs)
	return auth, userError(err)
}

func (s *userService) getUserAuth(ctx context.Context, rs *cognito.AuthResult) (*entity.UserAuth, error) {
	username, err := s.userAuth.GetUsername(ctx, rs.AccessToken)
	if err != nil {
		return nil, err
	}
	out, err := s.db.User.GetByCognitoID(ctx, username, "id")
	if err != nil {
		return nil, err
	}
	auth := entity.NewUserAuth(out.ID, rs)
	return auth, nil
}
