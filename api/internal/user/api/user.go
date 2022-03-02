package api

import (
	"context"

	"github.com/and-period/marche/api/proto/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *userService) GetUser(
	ctx context.Context, req *user.GetUserRequest,
) (*user.GetUserResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// TODO: 詳細の実装
	return &user.GetUserResponse{}, nil
}

func (s *userService) CreateUser(
	ctx context.Context, req *user.CreateUserRequest,
) (*user.CreateUserResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// TODO: 詳細の実装
	return &user.CreateUserResponse{}, nil
}

func (s *userService) VerifyUser(
	ctx context.Context, req *user.VerifyUserRequest,
) (*user.VerifyUserResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// TODO: 詳細の実装
	return &user.VerifyUserResponse{}, nil
}

func (s *userService) CreateUserWithOAuth(
	ctx context.Context, req *user.CreateUserWithOAuthRequest,
) (*user.CreateUserWithOAuthResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// TODO: 詳細の実装
	return &user.CreateUserWithOAuthResponse{}, nil
}

func (s *userService) UpdateUserEmail(
	ctx context.Context, req *user.UpdateUserEmailRequest,
) (*user.UpdateUserEmailResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// TODO: 詳細の実装
	return &user.UpdateUserEmailResponse{}, nil
}

func (s *userService) VerifyUserEmail(
	ctx context.Context, req *user.VerifyUserEmailRequest,
) (*user.VerifyUserEmailResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// TODO: 詳細の実装
	return &user.VerifyUserEmailResponse{}, nil
}

func (s *userService) UpdateUserPassword(
	ctx context.Context, req *user.UpdateUserPasswordRequest,
) (*user.UpdateUserPasswordResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// TODO: 詳細の実装
	return &user.UpdateUserPasswordResponse{}, nil
}

func (s *userService) ForgotUserPassword(
	ctx context.Context, req *user.ForgotUserPasswordRequest,
) (*user.ForgotUserPasswordResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// TODO: 詳細の実装
	return &user.ForgotUserPasswordResponse{}, nil
}

func (s *userService) VerifyUserPassword(
	ctx context.Context, req *user.VerifyUserPasswordRequest,
) (*user.VerifyUserPasswordResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// TODO: 詳細の実装
	return &user.VerifyUserPasswordResponse{}, nil
}

func (s *userService) DeleteUser(
	ctx context.Context, req *user.DeleteUserRequest,
) (*user.DeleteuserResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// TODO: 詳細の実装
	return &user.DeleteuserResponse{}, nil
}
