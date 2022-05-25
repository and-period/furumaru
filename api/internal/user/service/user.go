package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/uuid"
)

func (s *userService) GetUser(ctx context.Context, in *user.GetUserInput) (*entity.User, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	u, err := s.db.User.Get(ctx, in.UserID)
	return u, exception.InternalError(err)
}

func (s *userService) CreateUser(ctx context.Context, in *user.CreateUserInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", exception.InternalError(err)
	}
	userID := uuid.Base58Encode(uuid.New())
	u := entity.NewUser(userID, userID, entity.ProviderTypeEmail, in.Email, in.PhoneNumber)
	if err := s.db.User.Create(ctx, u); err != nil {
		return "", exception.InternalError(err)
	}
	params := &cognito.SignUpParams{
		Username:    u.CognitoID,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Password:    in.Password,
	}
	if err := s.userAuth.SignUp(ctx, params); err != nil {
		return "", exception.InternalError(err)
	}
	return userID, nil
}

func (s *userService) VerifyUser(ctx context.Context, in *user.VerifyUserInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	if err := s.userAuth.ConfirmSignUp(ctx, in.UserID, in.VerifyCode); err != nil {
		return exception.InternalError(err)
	}
	err := s.db.User.UpdateVerified(ctx, in.UserID)
	return exception.InternalError(err)
}

func (s *userService) CreateUserWithOAuth(
	ctx context.Context, in *user.CreateUserWithOAuthInput,
) (*entity.User, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, exception.InternalError(err)
	}
	auth, err := s.userAuth.GetUser(ctx, in.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrUnauthenticated, err.Error())
	}
	userID := uuid.Base58Encode(uuid.New())
	u := entity.NewUser(userID, auth.Username, entity.ProviderTypeOAuth, auth.Email, auth.PhoneNumber)
	if err := s.db.User.Create(ctx, u); err != nil {
		return nil, exception.InternalError(err)
	}
	return u, nil
}

func (s *userService) InitializeUser(ctx context.Context, in *user.InitializeUserInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}

	err := s.db.User.UpdateAccount(ctx, in.UserID, in.AccountID, in.Username)
	return exception.InternalError(err)
}

func (s *userService) UpdateUserEmail(ctx context.Context, in *user.UpdateUserEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	username, err := s.userAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return fmt.Errorf("%w: %s", exception.ErrUnauthenticated, err.Error())
	}
	u, err := s.db.User.GetByCognitoID(ctx, username, "id", "provider_type", "email")
	if err != nil {
		return exception.InternalError(err)
	}
	if u.ProviderType != entity.ProviderTypeEmail {
		return fmt.Errorf("%w: %s", exception.ErrFailedPrecondition, "api: not allow provider type to change email")
	}
	params := &cognito.ChangeEmailParams{
		AccessToken: in.AccessToken,
		Username:    username,
		OldEmail:    u.Email,
		NewEmail:    in.Email,
	}
	err = s.userAuth.ChangeEmail(ctx, params)
	return exception.InternalError(err)
}

func (s *userService) VerifyUserEmail(ctx context.Context, in *user.VerifyUserEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	username, err := s.userAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return fmt.Errorf("%w: %s", exception.ErrUnauthenticated, err.Error())
	}
	u, err := s.db.User.GetByCognitoID(ctx, username, "id")
	if err != nil {
		return exception.InternalError(err)
	}
	params := &cognito.ConfirmChangeEmailParams{
		AccessToken: in.AccessToken,
		Username:    username,
		VerifyCode:  in.VerifyCode,
	}
	email, err := s.userAuth.ConfirmChangeEmail(ctx, params)
	if err != nil {
		return exception.InternalError(err)
	}
	err = s.db.User.UpdateEmail(ctx, u.ID, email)
	return exception.InternalError(err)
}

func (s *userService) UpdateUserPassword(ctx context.Context, in *user.UpdateUserPasswordInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	params := &cognito.ChangePasswordParams{
		AccessToken: in.AccessToken,
		OldPassword: in.OldPassword,
		NewPassword: in.NewPassword,
	}
	err := s.userAuth.ChangePassword(ctx, params)
	return exception.InternalError(err)
}

func (s *userService) ForgotUserPassword(ctx context.Context, in *user.ForgotUserPasswordInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	u, err := s.db.User.GetByEmail(ctx, in.Email, "cognito_id")
	if err != nil {
		return exception.InternalError(err)
	}
	if err := s.userAuth.ForgotPassword(ctx, u.CognitoID); err != nil {
		return fmt.Errorf("%w: %s", exception.ErrNotFound, err.Error())
	}
	return nil
}

func (s *userService) VerifyUserPassword(ctx context.Context, in *user.VerifyUserPasswordInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	u, err := s.db.User.GetByEmail(ctx, in.Email, "cognito_id")
	if err != nil {
		return exception.InternalError(err)
	}
	params := &cognito.ConfirmForgotPasswordParams{
		Username:    u.CognitoID,
		VerifyCode:  in.VerifyCode,
		NewPassword: in.NewPassword,
	}
	err = s.userAuth.ConfirmForgotPassword(ctx, params)
	return exception.InternalError(err)
}

func (s *userService) DeleteUser(ctx context.Context, in *user.DeleteUserInput) error {
	if err := s.validator.Struct(in); err != nil {
		return exception.InternalError(err)
	}
	u, err := s.db.User.Get(ctx, in.UserID)
	if err != nil {
		return exception.InternalError(err)
	}
	if err := s.userAuth.DeleteUser(ctx, u.CognitoID); err != nil {
		return exception.InternalError(err)
	}
	err = s.db.User.Delete(ctx, u.ID)
	return exception.InternalError(err)
}
