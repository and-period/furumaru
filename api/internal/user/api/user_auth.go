package api

import (
	"context"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *userService) SignInUser(
	ctx context.Context, req *user.SignInUserRequest,
) (*user.SignInUserResponse, error) {
	rs, err := s.userAuth.SignIn(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	username, err := s.userAuth.GetUsername(ctx, rs.AccessToken)
	if err != nil {
		return nil, gRPCError(err)
	}
	out, err := s.db.User.GetByCognitoID(ctx, username, "id")
	if err != nil {
		return nil, gRPCError(err)
	}
	auth := entity.NewAuth(out.ID, rs)
	res := &user.SignInUserResponse{
		Auth: auth.Proto(),
	}
	return res, nil
}

func (s *userService) SignOutUser(
	ctx context.Context, req *user.SignOutUserRequest,
) (*user.SignOutUserResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.userAuth.SignOut(ctx, req.AccessToken); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &user.SignOutUserResponse{}, nil
}

func (s *userService) RefreshUserToken(
	ctx context.Context, req *user.RefreshUserTokenRequest,
) (*user.RefreshUserTokenResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	rs, err := s.userAuth.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	username, err := s.userAuth.GetUsername(ctx, rs.AccessToken)
	if err != nil {
		return nil, gRPCError(err)
	}
	out, err := s.db.User.GetByCognitoID(ctx, username, "id")
	if err != nil {
		return nil, gRPCError(err)
	}
	auth := entity.NewAuth(out.ID, rs)
	res := &user.RefreshUserTokenResponse{
		Auth: auth.Proto(),
	}
	return res, nil
}
