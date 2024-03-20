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
)

func (s *service) CreateMember(ctx context.Context, in *user.CreateMemberInput) (string, error) {
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

func (s *service) VerifyMember(ctx context.Context, in *user.VerifyMemberInput) error {
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

func (s *service) CreateMemberWithOAuth(
	ctx context.Context, in *user.CreateMemberWithOAuthInput,
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
	auth := func(_ context.Context) error {
		return nil // Cognitoへはすでに登録済みのため何もしない
	}
	if err := s.db.Member.Create(ctx, u, auth); err != nil {
		return nil, internalError(err)
	}
	return u, nil
}

func (s *service) UpdateMemberEmail(ctx context.Context, in *user.UpdateMemberEmailInput) error {
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

func (s *service) VerifyMemberEmail(ctx context.Context, in *user.VerifyMemberEmailInput) error {
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

func (s *service) UpdateMemberPassword(ctx context.Context, in *user.UpdateMemberPasswordInput) error {
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

func (s *service) ForgotMemberPassword(ctx context.Context, in *user.ForgotMemberPasswordInput) error {
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

func (s *service) VerifyMemberPassword(ctx context.Context, in *user.VerifyMemberPasswordInput) error {
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

func (s *service) UpdateMemberThumbnailURL(ctx context.Context, in *user.UpdateMemberThumbnailURLInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Member.UpdateThumbnailURL(ctx, in.UserID, in.ThumbnailURL)
	return internalError(err)
}

func (s *service) UpdateMemberThumbnails(ctx context.Context, in *user.UpdateMemberThumbnailsInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Member.UpdateThumbnails(ctx, in.UserID, in.Thumbnails)
	return internalError(err)
}
