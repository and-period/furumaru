package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
	"github.com/and-period/furumaru/api/internal/user/database"
	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"golang.org/x/sync/errgroup"
)

func (s *service) ListUsers(ctx context.Context, in *user.ListUsersInput) (entity.Users, int64, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, 0, internalError(err)
	}
	params := &database.ListUsersParams{
		Limit:          int(in.Limit),
		Offset:         int(in.Offset),
		OnlyRegistered: in.OnlyRegistered,
		OnlyVerified:   in.OnlyVerified,
		WithDeleted:    in.WithDeleted,
	}
	var (
		users entity.Users
		total int64
	)
	eg, ectx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		users, err = s.db.User.List(ectx, params)
		return
	})
	eg.Go(func() (err error) {
		total, err = s.db.User.Count(ectx, params)
		return
	})
	if err := eg.Wait(); err != nil {
		return nil, 0, internalError(err)
	}
	return users, total, nil
}

func (s *service) MultiGetUsers(ctx context.Context, in *user.MultiGetUsersInput) (entity.Users, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	users, err := s.db.User.MultiGet(ctx, in.UserIDs)
	return users, internalError(err)
}

func (s *service) MultiGetUserDevices(_ context.Context, in *user.MultiGetUserDevicesInput) ([]string, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	// TODO: 詳細の実装
	return []string{}, nil
}

func (s *service) GetUser(ctx context.Context, in *user.GetUserInput) (*entity.User, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	u, err := s.db.User.Get(ctx, in.UserID)
	return u, internalError(err)
}

func (s *service) CreateUser(ctx context.Context, in *user.CreateUserInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	member, err := s.db.Member.GetByEmail(ctx, in.Email)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return "", internalError(err)
	}
	// 確認コードの検証ができていないユーザーがいる場合、コードの再送を実施
	if member != nil && member.VerifiedAt.IsZero() {
		if err := s.userAuth.ResendSignUpCode(ctx, member.CognitoID); err != nil {
			return "", internalError(err)
		}
		return member.UserID, nil
	}
	cognitoID := uuid.Base58Encode(uuid.New())
	params := &entity.NewUserParams{
		Registered:    true,
		CognitoID:     cognitoID,
		Username:      in.Username,
		AccountID:     in.AccountID,
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
		ProviderType:  entity.ProviderTypeEmail,
		Email:         in.Email,
		PhoneNumber:   in.PhoneNumber,
	}
	u := entity.NewUser(params)
	auth := func(ctx context.Context) error {
		params := &cognito.SignUpParams{
			Username:    cognitoID,
			Email:       in.Email,
			PhoneNumber: in.PhoneNumber,
			Password:    in.Password,
		}
		return s.userAuth.SignUp(ctx, params)
	}
	if err := s.db.Member.Create(ctx, u, auth); err != nil {
		return "", internalError(err)
	}
	return u.ID, nil
}

func (s *service) VerifyUser(ctx context.Context, in *user.VerifyUserInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	u, err := s.db.User.Get(ctx, in.UserID)
	if err != nil {
		return internalError(err)
	}
	err = s.userAuth.ConfirmSignUp(ctx, u.Member.CognitoID, in.VerifyCode)
	if errors.Is(err, cognito.ErrCodeExpired) {
		err = s.userAuth.ResendSignUpCode(ctx, u.Member.CognitoID)
		return internalError(err)
	}
	if err != nil {
		return internalError(err)
	}
	err = s.db.Member.UpdateVerified(ctx, in.UserID)
	return internalError(err)
}

func (s *service) CreateUserWithOAuth(
	ctx context.Context, in *user.CreateUserWithOAuthInput,
) (*entity.User, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	cuser, err := s.userAuth.GetUser(ctx, in.AccessToken)
	if err != nil {
		return nil, internalError(err)
	}
	params := &entity.NewUserParams{
		Registered:    true,
		CognitoID:     cuser.Username,
		Username:      in.Username,
		AccountID:     in.AccountID,
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
		ProviderType:  entity.ProviderTypeOAuth,
		Email:         cuser.Email,
		PhoneNumber:   in.PhoneNumber,
	}
	u := entity.NewUser(params)
	auth := func(ctx context.Context) error {
		return nil // Cognitoへはすでに登録済みのため何もしない
	}
	if err := s.db.Member.Create(ctx, u, auth); err != nil {
		return nil, internalError(err)
	}
	return u, nil
}

func (s *service) UpdateUserEmail(ctx context.Context, in *user.UpdateUserEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	username, err := s.userAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return internalError(err)
	}
	m, err := s.db.Member.GetByCognitoID(ctx, username, "provider_type", "email")
	if err != nil {
		return internalError(err)
	}
	if m.ProviderType != entity.ProviderTypeEmail {
		return fmt.Errorf("%w: %s", exception.ErrFailedPrecondition, "api: not allow provider type to change email")
	}
	params := &cognito.ChangeEmailParams{
		AccessToken: in.AccessToken,
		Username:    username,
		OldEmail:    m.Email,
		NewEmail:    in.Email,
	}
	err = s.userAuth.ChangeEmail(ctx, params)
	return internalError(err)
}

func (s *service) VerifyUserEmail(ctx context.Context, in *user.VerifyUserEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	username, err := s.userAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return internalError(err)
	}
	m, err := s.db.Member.GetByCognitoID(ctx, username, "user_id")
	if err != nil {
		return internalError(err)
	}
	params := &cognito.ConfirmChangeEmailParams{
		AccessToken: in.AccessToken,
		Username:    username,
		VerifyCode:  in.VerifyCode,
	}
	email, err := s.userAuth.ConfirmChangeEmail(ctx, params)
	if err != nil {
		return internalError(err)
	}
	err = s.db.Member.UpdateEmail(ctx, m.UserID, email)
	return internalError(err)
}

func (s *service) UpdateUserPassword(ctx context.Context, in *user.UpdateUserPasswordInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &cognito.ChangePasswordParams{
		AccessToken: in.AccessToken,
		OldPassword: in.OldPassword,
		NewPassword: in.NewPassword,
	}
	err := s.userAuth.ChangePassword(ctx, params)
	return internalError(err)
}

func (s *service) ForgotUserPassword(ctx context.Context, in *user.ForgotUserPasswordInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	m, err := s.db.Member.GetByEmail(ctx, in.Email, "cognito_id")
	if err != nil {
		return internalError(err)
	}
	if err := s.userAuth.ForgotPassword(ctx, m.CognitoID); err != nil {
		return fmt.Errorf("%w: %s", exception.ErrNotFound, err.Error())
	}
	return nil
}

func (s *service) VerifyUserPassword(ctx context.Context, in *user.VerifyUserPasswordInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	m, err := s.db.Member.GetByEmail(ctx, in.Email, "cognito_id")
	if err != nil {
		return internalError(err)
	}
	params := &cognito.ConfirmForgotPasswordParams{
		Username:    m.CognitoID,
		VerifyCode:  in.VerifyCode,
		NewPassword: in.NewPassword,
	}
	err = s.userAuth.ConfirmForgotPassword(ctx, params)
	return internalError(err)
}

func (s *service) UpdateUserThumbnails(ctx context.Context, in *user.UpdateUserThumbnailsInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Member.UpdateThumbnails(ctx, in.UserID, in.Thumbnails)
	return internalError(err)
}

func (s *service) DeleteUser(ctx context.Context, in *user.DeleteUserInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	u, err := s.db.User.Get(ctx, in.UserID)
	if err != nil {
		return internalError(err)
	}
	if u.Registered {
		auth := func(ctx context.Context) error {
			err := s.userAuth.DeleteUser(ctx, u.Member.CognitoID)
			if errors.Is(err, cognito.ErrNotFound) {
				return nil // すでに削除済み
			}
			return err
		}
		err = s.db.Member.Delete(ctx, u.ID, auth)
	} else {
		err = s.db.Guest.Delete(ctx, u.ID)
	}
	return internalError(err)
}
