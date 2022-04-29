package service

import (
	"context"
	"fmt"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/cognito"
)

func (s *userService) SignInUser(ctx context.Context, in *SignInUserInput) (*entity.UserAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	rs, err := s.userAuth.SignIn(ctx, in.Key, in.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnauthenticated, err.Error())
	}
	auth, err := s.getUserAuth(ctx, rs)
	return auth, userError(err)
}

func (s *userService) SignOutUser(ctx context.Context, in *SignOutUserInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}
	if err := s.userAuth.SignOut(ctx, in.AccessToken); err != nil {
		return fmt.Errorf("%w: %s", ErrUnauthenticated, err.Error())
	}
	return nil
}

func (s *userService) GetUserAuth(ctx context.Context, in *GetUserAuthInput) (*entity.UserAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	rs := &cognito.AuthResult{AccessToken: in.AccessToken}
	auth, err := s.getUserAuth(ctx, rs)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnauthenticated, err.Error())
	}
	return auth, nil
}

func (s *userService) RefreshUserToken(ctx context.Context, in *RefreshUserTokenInput) (*entity.UserAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	rs, err := s.userAuth.RefreshToken(ctx, in.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnauthenticated, err.Error())
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
