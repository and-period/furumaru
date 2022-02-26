package api

import (
	"context"

	"github.com/and-period/marche/api/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *userService) SignInUser(
	ctx context.Context, req *user.SignInUserRequest,
) (*user.SignInUserResponse, error) {
	// TODO: 詳細の実装
	return &user.SignInUserResponse{}, nil
}

func (s *userService) SignOutUser(
	ctx context.Context, req *user.SignOutUserRequest,
) (*user.SignOutUserResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// TODO: 詳細の実装
	return &user.SignOutUserResponse{}, nil
}

func (s *userService) RefreshUserToken(
	ctx context.Context, req *user.RefreshUserTokenRequest,
) (*user.RefreshUserTokenResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// TODO: 詳細の実装
	return &user.RefreshUserTokenResponse{}, nil
}
