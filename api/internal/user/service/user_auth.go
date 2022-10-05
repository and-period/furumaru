package service

import (
	"context"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
)

func (s *service) SignInUser(ctx context.Context, in *user.SignInUserInput) (*entity.UserAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	rs, err := s.userAuth.SignIn(ctx, in.Key, in.Password)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	auth, err := s.getUserAuth(ctx, rs)
	return auth, exception.InternalError(err)
}

func (s *service) SignOutUser(ctx context.Context, in *user.SignOutUserInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	err := s.userAuth.SignOut(ctx, in.AccessToken)
	return exception.InternalError(err)
}

func (s *service) GetUserAuth(ctx context.Context, in *user.GetUserAuthInput) (*entity.UserAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	rs := &cognito.AuthResult{AccessToken: in.AccessToken}
	auth, err := s.getUserAuth(ctx, rs)
	return auth, exception.InternalError(err)
}

func (s *service) RefreshUserToken(ctx context.Context, in *user.RefreshUserTokenInput) (*entity.UserAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	rs, err := s.userAuth.RefreshToken(ctx, in.RefreshToken)
	if err != nil {
		return nil, exception.InternalError(err)
	}
	auth, err := s.getUserAuth(ctx, rs)
	return auth, exception.InternalError(err)
}

func (s *service) getUserAuth(ctx context.Context, rs *cognito.AuthResult) (*entity.UserAuth, error) {
	username, err := s.userAuth.GetUsername(ctx, rs.AccessToken)
	if err != nil {
		return nil, err
	}
	out, err := s.db.Member.GetByCognitoID(ctx, username, "user_id")
	if err != nil {
		return nil, err
	}
	auth := entity.NewUserAuth(out.UserID, rs)
	return auth, nil
}
