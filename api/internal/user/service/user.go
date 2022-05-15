package service

import (
	"context"
	"fmt"

	"github.com/and-period/marche/api/internal/user/entity"
	"github.com/and-period/marche/api/pkg/cognito"
	"github.com/and-period/marche/api/pkg/uuid"
)

func (s *userService) GetUser(ctx context.Context, in *GetUserInput) (*entity.User, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	u, err := s.db.User.Get(ctx, in.UserID)
	return u, userError(err)
}

func (s *userService) CreateUser(ctx context.Context, in *CreateUserInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", userError(err)
	}
	userID := uuid.Base58Encode(uuid.New())
	u := entity.NewUser(userID, userID, entity.ProviderTypeEmail, in.Email, in.PhoneNumber)
	if err := s.db.User.Create(ctx, u); err != nil {
		return "", userError(err)
	}
	params := &cognito.SignUpParams{
		Username:    u.CognitoID,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Password:    in.Password,
	}
	if err := s.userAuth.SignUp(ctx, params); err != nil {
		return "", userError(err)
	}
	return userID, nil
}

func (s *userService) VerifyUser(ctx context.Context, in *VerifyUserInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}
	if err := s.userAuth.ConfirmSignUp(ctx, in.UserID, in.VerifyCode); err != nil {
		return userError(err)
	}
	err := s.db.User.UpdateVerified(ctx, in.UserID)
	return userError(err)
}

func (s *userService) CreateUserWithOAuth(ctx context.Context, in *CreateUserWithOAuthInput) (*entity.User, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, userError(err)
	}
	auth, err := s.userAuth.GetUser(ctx, in.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnauthenticated, err.Error())
	}
	userID := uuid.Base58Encode(uuid.New())
	u := entity.NewUser(userID, auth.Username, entity.ProviderTypeOAuth, auth.Email, auth.PhoneNumber)
	if err := s.db.User.Create(ctx, u); err != nil {
		return nil, userError(err)
	}
	return u, nil
}

func (s *userService) InitializeUser(ctx context.Context, in *InitializeUserInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}

	err := s.db.User.UpdateAccount(ctx, in.UserID, in.AccountID, in.Username)
	return userError(err)
}

func (s *userService) UpdateUserEmail(ctx context.Context, in *UpdateUserEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}
	username, err := s.userAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnauthenticated, err.Error())
	}
	u, err := s.db.User.GetByCognitoID(ctx, username, "id", "provider_type", "email")
	if err != nil {
		return userError(err)
	}
	if u.ProviderType != entity.ProviderTypeEmail {
		return fmt.Errorf("%w: %s", ErrFailedPrecondition, "api: not allow provider type to change email")
	}
	params := &cognito.ChangeEmailParams{
		AccessToken: in.AccessToken,
		Username:    username,
		OldEmail:    u.Email,
		NewEmail:    in.Email,
	}
	err = s.userAuth.ChangeEmail(ctx, params)
	return userError(err)
}

func (s *userService) VerifyUserEmail(ctx context.Context, in *VerifyUserEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}
	username, err := s.userAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUnauthenticated, err.Error())
	}
	u, err := s.db.User.GetByCognitoID(ctx, username, "id")
	if err != nil {
		return userError(err)
	}
	params := &cognito.ConfirmChangeEmailParams{
		AccessToken: in.AccessToken,
		Username:    username,
		VerifyCode:  in.VerifyCode,
	}
	email, err := s.userAuth.ConfirmChangeEmail(ctx, params)
	if err != nil {
		return userError(err)
	}
	err = s.db.User.UpdateEmail(ctx, u.ID, email)
	return userError(err)
}

func (s *userService) UpdateUserPassword(ctx context.Context, in *UpdateUserPasswordInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}
	params := &cognito.ChangePasswordParams{
		AccessToken: in.AccessToken,
		OldPassword: in.OldPassword,
		NewPassword: in.NewPassword,
	}
	err := s.userAuth.ChangePassword(ctx, params)
	return userError(err)
}

func (s *userService) ForgotUserPassword(ctx context.Context, in *ForgotUserPasswordInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}
	u, err := s.db.User.GetByEmail(ctx, in.Email, "cognito_id")
	if err != nil {
		return userError(err)
	}
	if err := s.userAuth.ForgotPassword(ctx, u.CognitoID); err != nil {
		return fmt.Errorf("%w: %s", ErrNotFound, err.Error())
	}
	return nil
}

func (s *userService) VerifyUserPassword(ctx context.Context, in *VerifyUserPasswordInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}
	u, err := s.db.User.GetByEmail(ctx, in.Email, "cognito_id")
	if err != nil {
		return userError(err)
	}
	params := &cognito.ConfirmForgotPasswordParams{
		Username:    u.CognitoID,
		VerifyCode:  in.VerifyCode,
		NewPassword: in.NewPassword,
	}
	err = s.userAuth.ConfirmForgotPassword(ctx, params)
	return userError(err)
}

func (s *userService) DeleteUser(ctx context.Context, in *DeleteUserInput) error {
	if err := s.validator.Struct(in); err != nil {
		return userError(err)
	}
	u, err := s.db.User.Get(ctx, in.UserID)
	if err != nil {
		return userError(err)
	}
	if err := s.userAuth.DeleteUser(ctx, u.CognitoID); err != nil {
		return userError(err)
	}
	err = s.db.User.Delete(ctx, u.ID)
	return userError(err)
}
