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
	"go.uber.org/zap"
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
		Email:         in.Email,
		ProviderType:  entity.UserAuthProviderTypeEmail,
		Username:      in.Username,
		AccountID:     in.AccountID,
		Lastname:      in.Lastname,
		Firstname:     in.Firstname,
		LastnameKana:  in.LastnameKana,
		FirstnameKana: in.FirstnameKana,
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
	err = s.userAuth.ConfirmSignUp(ctx, u.CognitoID, in.VerifyCode)
	if errors.Is(err, cognito.ErrCodeExpired) {
		err = s.userAuth.ResendSignUpCode(ctx, u.CognitoID)
		return internalError(err)
	}
	if err != nil {
		return internalError(err)
	}
	err = s.db.Member.UpdateVerified(ctx, in.UserID)
	return internalError(err)
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
	if m.ProviderType != entity.UserAuthProviderTypeEmail {
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

func (s *service) UpdateMemberUsername(ctx context.Context, in *user.UpdateMemberUsernameInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Member.UpdateUsername(ctx, in.UserID, in.Username)
	return internalError(err)
}

func (s *service) UpdateMemberAccountID(ctx context.Context, in *user.UpdateMemberAccountIDInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Member.UpdateAccountID(ctx, in.UserID, in.AccountID)
	return internalError(err)
}

func (s *service) UpdateMemberThumbnailURL(ctx context.Context, in *user.UpdateMemberThumbnailURLInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Member.UpdateThumbnailURL(ctx, in.UserID, in.ThumbnailURL)
	return internalError(err)
}

func (s *service) AuthMemberWithGoogle(ctx context.Context, in *user.AuthMemberWithGoogleInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	params := &authMemberWithOAuthParams{
		payload:      &in.AuthMemberDetailWithOAuth,
		providerType: entity.UserAuthProviderTypeGoogle,
		redirectURI:  s.userAuthGoogleRedirectURL,
	}
	if in.RedirectURI != "" {
		params.redirectURI = in.RedirectURI
	}
	return s.authMemberWithOAuth(ctx, params)
}

func (s *service) CreateMemberWithGoogle(ctx context.Context, in *user.CreateMemberWithGoogleInput) (*entity.User, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &createMemberWithOAuthParams{
		payload:      &in.CreateMemberDetailWithOAuth,
		providerType: entity.UserAuthProviderTypeGoogle,
		redirectURI:  s.userAuthGoogleRedirectURL,
	}
	if in.RedirectURI != "" {
		params.redirectURI = in.RedirectURI
	}
	return s.createMemberWithOAuth(ctx, params)
}

func (s *service) AuthMemberWithLINE(ctx context.Context, in *user.AuthMemberWithLINEInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	params := &authMemberWithOAuthParams{
		payload:      &in.AuthMemberDetailWithOAuth,
		providerType: entity.UserAuthProviderTypeLINE,
		redirectURI:  s.userAuthLINERedirectURL,
	}
	if in.RedirectURI != "" {
		params.redirectURI = in.RedirectURI
	}
	return s.authMemberWithOAuth(ctx, params)
}

func (s *service) CreateMemberWithLINE(ctx context.Context, in *user.CreateMemberWithLINEInput) (*entity.User, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &createMemberWithOAuthParams{
		payload:      &in.CreateMemberDetailWithOAuth,
		providerType: entity.UserAuthProviderTypeLINE,
		redirectURI:  s.userAuthLINERedirectURL,
	}
	if in.RedirectURI != "" {
		params.redirectURI = in.RedirectURI
	}
	return s.createMemberWithOAuth(ctx, params)
}

type authMemberWithOAuthParams struct {
	payload      *user.AuthMemberDetailWithOAuth
	providerType entity.UserAuthProviderType
	redirectURI  string
}

func (s *service) authMemberWithOAuth(ctx context.Context, params *authMemberWithOAuthParams) (string, error) {
	eventParams := &entity.UserAuthEventParams{
		SessionID:    params.payload.SessionID,
		ProviderType: params.providerType,
		Now:          s.now(),
		TTL:          s.userAuthTTL,
	}
	event := entity.NewUserAuthEvent(eventParams)
	if err := s.cache.Insert(ctx, event); err != nil {
		return "", internalError(err)
	}
	authParams := &cognito.GenerateAuthURLParams{
		State:        params.payload.State,
		Nonce:        event.Nonce,
		ProviderType: params.providerType.ToCognito(),
		RedirectURI:  params.redirectURI,
	}
	authURL, err := s.userAuth.GenerateAuthURL(ctx, authParams)
	return authURL, internalError(err)
}

type createMemberWithOAuthParams struct {
	payload      *user.CreateMemberDetailWithOAuth
	providerType entity.UserAuthProviderType
	redirectURI  string
}

func (s *service) createMemberWithOAuth(ctx context.Context, params *createMemberWithOAuthParams) (*entity.User, error) {
	event := &entity.UserAuthEvent{SessionID: params.payload.SessionID}
	if err := s.cache.Get(ctx, event); err != nil {
		return nil, internalError(err)
	}
	if event.Nonce != params.payload.Nonce {
		return nil, fmt.Errorf("service: invalid nonce: %w", exception.ErrFailedPrecondition)
	}
	if event.ProviderType != params.providerType {
		return nil, fmt.Errorf("service: invalid provider type: %w", exception.ErrFailedPrecondition)
	}

	// Cognitoユーザーの取得
	tokenParams := &cognito.GetAccessTokenParams{
		Code:        params.payload.Code,
		RedirectURI: params.redirectURI,
	}
	token, err := s.userAuth.GetAccessToken(ctx, tokenParams)
	if err != nil {
		return nil, internalError(err)
	}
	cuser, err := s.userAuth.GetUser(ctx, token.AccessToken)
	if err != nil {
		return nil, internalError(err)
	}
	s.logger.Debug("Creating User account", zap.Any("user", cuser))

	userParams := &entity.NewUserParams{
		Registered:    true,
		CognitoID:     cuser.Username,
		Email:         cuser.Email,
		ProviderType:  params.providerType,
		Username:      params.payload.Username,
		AccountID:     params.payload.AccountID,
		Lastname:      params.payload.Lastname,
		Firstname:     params.payload.Firstname,
		LastnameKana:  params.payload.LastnameKana,
		FirstnameKana: params.payload.FirstnameKana,
		PhoneNumber:   params.payload.PhoneNumber,
	}
	u := entity.NewUser(userParams)
	auth := func(_ context.Context) error {
		return nil // Cognitoへはすでに登録済みのため何もしない
	}
	if err := s.db.Member.Create(ctx, u, auth); err != nil {
		return nil, internalError(err)
	}
	return u, nil
}
