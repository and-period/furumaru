package api

import (
	"context"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/pkg/uuid"
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
	u, err := s.db.User.Get(ctx, req.UserId)
	if err != nil {
		return nil, gRPCError(err)
	}
	res := &user.GetUserResponse{
		User: u.Proto(),
	}
	return res, nil
}

func (s *userService) CreateUser(
	ctx context.Context, req *user.CreateUserRequest,
) (*user.CreateUserResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if req.Password != req.PasswordConfirmation {
		return nil, status.Error(codes.InvalidArgument, "password is unmatch")
	}
	userID := uuid.Base58Encode(uuid.New())
	u := entity.NewUser(userID, userID, entity.ProviderTypeEmail, req.Email, req.PhoneNumber)
	if err := s.db.User.Create(ctx, u); err != nil {
		return nil, gRPCError(err)
	}
	params := &cognito.SignUpParams{
		Username:    u.CognitoID,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Password:    req.Password,
	}
	if err := s.userAuth.SignUp(ctx, params); err != nil {
		return nil, gRPCError(err)
	}
	res := &user.CreateUserResponse{
		UserId: userID,
	}
	return res, nil
}

func (s *userService) VerifyUser(
	ctx context.Context, req *user.VerifyUserRequest,
) (*user.VerifyUserResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.userAuth.ConfirmSignUp(ctx, req.UserId, req.VerifyCode); err != nil {
		return nil, gRPCError(err)
	}
	if err := s.db.User.UpdateVerified(ctx, req.UserId); err != nil {
		return nil, gRPCError(err)
	}
	return &user.VerifyUserResponse{}, nil
}

func (s *userService) CreateUserWithOAuth(
	ctx context.Context, req *user.CreateUserWithOAuthRequest,
) (*user.CreateUserWithOAuthResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	auth, err := s.userAuth.GetUser(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	userID := uuid.Base58Encode(uuid.New())
	u := entity.NewUser(userID, auth.Username, entity.ProviderTypeOAuth, auth.Email, auth.PhoneNumber)
	if err := s.db.User.Create(ctx, u); err != nil {
		return nil, gRPCError(err)
	}
	res := &user.CreateUserWithOAuthResponse{
		User: u.Proto(),
	}
	return res, nil
}

func (s *userService) UpdateUserEmail(
	ctx context.Context, req *user.UpdateUserEmailRequest,
) (*user.UpdateUserEmailResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	username, err := s.userAuth.GetUsername(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	u, err := s.db.User.GetByCognitoID(ctx, username, "id", "provider_type", "email")
	if err != nil {
		return nil, gRPCError(err)
	}
	if u.ProviderType != entity.ProviderTypeEmail {
		return nil, status.Error(codes.FailedPrecondition, "api: not allow provider type to change email")
	}
	params := &cognito.ChangeEmailParams{
		AccessToken: req.AccessToken,
		Username:    username,
		OldEmail:    u.Email,
		NewEmail:    req.Email,
	}
	if err := s.userAuth.ChangeEmail(ctx, params); err != nil {
		return nil, gRPCError(err)
	}
	return &user.UpdateUserEmailResponse{}, nil
}

func (s *userService) VerifyUserEmail(
	ctx context.Context, req *user.VerifyUserEmailRequest,
) (*user.VerifyUserEmailResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	username, err := s.userAuth.GetUsername(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	u, err := s.db.User.GetByCognitoID(ctx, username, "id")
	if err != nil {
		return nil, gRPCError(err)
	}
	params := &cognito.ConfirmChangeEmailParams{
		AccessToken: req.AccessToken,
		Username:    username,
		VerifyCode:  req.VerifyCode,
	}
	email, err := s.userAuth.ConfirmChangeEmail(ctx, params)
	if err != nil {
		return nil, gRPCError(err)
	}
	if err := s.db.User.UpdateEmail(ctx, u.ID, email); err != nil {
		return nil, gRPCError(err)
	}
	return &user.VerifyUserEmailResponse{}, nil
}

func (s *userService) UpdateUserPassword(
	ctx context.Context, req *user.UpdateUserPasswordRequest,
) (*user.UpdateUserPasswordResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if req.NewPassword != req.PasswordConfirmation {
		return nil, status.Error(codes.InvalidArgument, "password is unmatch")
	}
	params := &cognito.ChangePasswordParams{
		AccessToken: req.AccessToken,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
	if err := s.userAuth.ChangePassword(ctx, params); err != nil {
		return nil, gRPCError(err)
	}
	return &user.UpdateUserPasswordResponse{}, nil
}

func (s *userService) ForgotUserPassword(
	ctx context.Context, req *user.ForgotUserPasswordRequest,
) (*user.ForgotUserPasswordResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	u, err := s.db.User.GetByEmail(ctx, req.Email, "cognito_id")
	if err != nil {
		return nil, gRPCError(err)
	}
	if err := s.userAuth.ForgotPassword(ctx, u.CognitoID); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &user.ForgotUserPasswordResponse{}, nil
}

func (s *userService) VerifyUserPassword(
	ctx context.Context, req *user.VerifyUserPasswordRequest,
) (*user.VerifyUserPasswordResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if req.NewPassword != req.PasswordConfirmation {
		return nil, status.Error(codes.InvalidArgument, "password is unmatch")
	}
	u, err := s.db.User.GetByEmail(ctx, req.Email, "cognito_id")
	if err != nil {
		return nil, gRPCError(err)
	}
	params := &cognito.ConfirmForgotPasswordParams{
		Username:    u.CognitoID,
		VerifyCode:  req.VerifyCode,
		NewPassword: req.NewPassword,
	}
	if err := s.userAuth.ConfirmForgotPassword(ctx, params); err != nil {
		return nil, gRPCError(err)
	}
	return &user.VerifyUserPasswordResponse{}, nil
}

func (s *userService) DeleteUser(
	ctx context.Context, req *user.DeleteUserRequest,
) (*user.DeleteUserResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	u, err := s.db.User.Get(ctx, req.UserId)
	if err != nil {
		return nil, gRPCError(err)
	}
	if err := s.userAuth.DeleteUser(ctx, u.CognitoID); err != nil {
		return nil, gRPCError(err)
	}
	if err := s.db.User.Delete(ctx, u.ID); err != nil {
		return nil, gRPCError(err)
	}
	return &user.DeleteUserResponse{}, nil
}
