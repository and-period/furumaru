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
)

func (s *service) SignInAdmin(ctx context.Context, in *user.SignInAdminInput) (*entity.AdminAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	rs, err := s.adminAuth.SignIn(ctx, in.Key, in.Password)
	if err != nil {
		return nil, internalError(err)
	}
	auth, err := s.getAdminAuth(ctx, rs)
	if err != nil {
		return nil, internalError(err)
	}
	if err := s.db.Admin.UpdateSignInAt(ctx, auth.AdminID); err != nil {
		return nil, internalError(err)
	}
	return auth, nil
}

func (s *service) SignOutAdmin(ctx context.Context, in *user.SignOutAdminInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.adminAuth.SignOut(ctx, in.AccessToken)
	return internalError(err)
}

func (s *service) GetAdminAuth(ctx context.Context, in *user.GetAdminAuthInput) (*entity.AdminAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	rs := &cognito.AuthResult{AccessToken: in.AccessToken}
	auth, err := s.getAdminAuth(ctx, rs)
	return auth, internalError(err)
}

func (s *service) RefreshAdminToken(
	ctx context.Context, in *user.RefreshAdminTokenInput,
) (*entity.AdminAuth, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	rs, err := s.adminAuth.RefreshToken(ctx, in.RefreshToken)
	if err != nil {
		return nil, internalError(err)
	}
	auth, err := s.getAdminAuth(ctx, rs)
	if err != nil {
		return nil, internalError(err)
	}
	if err := s.db.Admin.UpdateSignInAt(ctx, auth.AdminID); err != nil {
		return nil, internalError(err)
	}
	return auth, internalError(err)
}

func (s *service) RegisterAdminDevice(ctx context.Context, in *user.RegisterAdminDeviceInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	err := s.db.Admin.UpdateDevice(ctx, in.AdminID, in.Device)
	return internalError(err)
}

func (s *service) UpdateAdminEmail(ctx context.Context, in *user.UpdateAdminEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	username, err := s.adminAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return internalError(err)
	}
	admin, err := s.db.Admin.GetByCognitoID(ctx, username, "email")
	if err != nil {
		return internalError(err)
	}
	if admin.Email == in.Email {
		return fmt.Errorf("this admin does not need to be changed email: %w", exception.ErrFailedPrecondition)
	}
	params := &cognito.ChangeEmailParams{
		AccessToken: in.AccessToken,
		Username:    username,
		OldEmail:    admin.Email,
		NewEmail:    in.Email,
	}
	err = s.adminAuth.ChangeEmail(ctx, params)
	return internalError(err)
}

func (s *service) VerifyAdminEmail(ctx context.Context, in *user.VerifyAdminEmailInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	username, err := s.adminAuth.GetUsername(ctx, in.AccessToken)
	if err != nil {
		return internalError(err)
	}
	admin, err := s.db.Admin.GetByCognitoID(ctx, username, "id", "role")
	if err != nil {
		return internalError(err)
	}
	params := &cognito.ConfirmChangeEmailParams{
		AccessToken: in.AccessToken,
		Username:    username,
		VerifyCode:  in.VerifyCode,
	}
	email, err := s.adminAuth.ConfirmChangeEmail(ctx, params)
	if err != nil {
		return internalError(err)
	}
	err = s.db.Admin.UpdateEmail(ctx, admin.ID, email)
	return internalError(err)
}

func (s *service) UpdateAdminPassword(ctx context.Context, in *user.UpdateAdminPasswordInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &cognito.ChangePasswordParams{
		AccessToken: in.AccessToken,
		OldPassword: in.OldPassword,
		NewPassword: in.NewPassword,
	}
	err := s.adminAuth.ChangePassword(ctx, params)
	return internalError(err)
}

func (s *service) getAdminAuth(ctx context.Context, rs *cognito.AuthResult) (*entity.AdminAuth, error) {
	username, err := s.adminAuth.GetUsername(ctx, rs.AccessToken)
	if err != nil {
		return nil, err
	}
	admin, err := s.db.Admin.GetByCognitoID(ctx, username)
	if err != nil {
		return nil, err
	}
	auth := entity.NewAdminAuth(admin, rs)
	return auth, nil
}

func (s *service) ListAdminAuthProviders(ctx context.Context, in *user.ListAdminAuthProvidersInput) (entity.AdminAuthProviders, error) {
	if err := s.validator.Struct(in); err != nil {
		return nil, internalError(err)
	}
	params := &database.ListAdminAuthProvidersParams{
		AdminID: in.AdminID,
	}
	providers, err := s.db.AdminAuthProvider.List(ctx, params)
	return providers, internalError(err)
}

func (s *service) InitialGoogleAdminAuth(ctx context.Context, in *user.InitialGoogleAdminAuthInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	params := &initialAdminAuthParams{
		adminID:      in.AdminID,
		state:        in.State,
		providerType: entity.AdminAuthProviderTypeGoogle,
		redirectURI:  s.adminAuthGoogleRedirectURL,
	}
	if in.RedirectURI != "" {
		params.redirectURI = in.RedirectURI
	}
	return s.initialAdminAuth(ctx, params)
}

func (s *service) ConnectGoogleAdminAuth(ctx context.Context, in *user.ConnectGoogleAdminAuthInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	params := &connectAdminAuthParams{
		adminID:     in.AdminID,
		code:        in.Code,
		nonce:       in.Nonce,
		redirectURI: s.adminAuthGoogleRedirectURL,
	}
	if in.RedirectURI != "" {
		params.redirectURI = in.RedirectURI
	}
	return s.connectAdminAuth(ctx, params)
}

type initialAdminAuthParams struct {
	adminID      string
	state        string
	providerType entity.AdminAuthProviderType
	redirectURI  string
}

func (s *service) initialAdminAuth(ctx context.Context, params *initialAdminAuthParams) (string, error) {
	provider, err := s.db.AdminAuthProvider.Get(ctx, params.adminID, params.providerType)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return "", internalError(err)
	}
	if provider != nil {
		return "", fmt.Errorf("this admin has already connected: %w", exception.ErrFailedPrecondition)
	}
	eventParams := &entity.AdminAuthEventParams{
		AdminID:      params.adminID,
		ProviderType: params.providerType,
		Now:          s.now(),
		TTL:          s.adminAuthTTL,
	}
	event := entity.NewAdminAuthEvent(eventParams)
	if err := s.cache.Insert(ctx, event); err != nil {
		return "", internalError(err)
	}
	authParams := &cognito.GenerateAuthURLParams{
		State:        params.state,
		Nonce:        event.Nonce,
		ProviderType: params.providerType.ToCognito(),
		RedirectURI:  params.redirectURI,
	}
	authURL, err := s.adminAuth.GenerateAuthURL(ctx, authParams)
	return authURL, internalError(err)
}

type connectAdminAuthParams struct {
	adminID     string
	code        string
	nonce       string
	redirectURI string
}

func (s *service) connectAdminAuth(ctx context.Context, params *connectAdminAuthParams) error {
	event := &entity.AdminAuthEvent{AdminID: params.adminID}
	if err := s.cache.Get(ctx, event); err != nil {
		return internalError(err)
	}
	if event.Nonce != params.nonce {
		return fmt.Errorf("service: invalid nonce: %w", exception.ErrForbidden)
	}

	// Cognitoユーザーの取得
	admin, err := s.db.Admin.Get(ctx, params.adminID)
	if err != nil {
		return internalError(err)
	}

	// 外部アカウントの取得
	tokenParams := &cognito.GetAccessTokenParams{
		Code:        params.code,
		RedirectURI: params.redirectURI,
	}
	token, err := s.adminAuth.GetAccessToken(ctx, tokenParams)
	if err != nil {
		return internalError(err)
	}

	// Cognitoの仕様で「すでにサインイン済みの場合は連携できない」ため、登録済みのGoogleアカウントを削除
	user, err := s.adminAuth.GetUser(ctx, token.AccessToken)
	if err != nil {
		return internalError(err)
	}
	providerParams := &entity.AdminAuthProviderParams{
		AdminID:      admin.ID,
		ProviderType: event.ProviderType,
		Auth:         user,
	}
	provider, err := entity.NewAdminAuthProvider(providerParams)
	if err != nil {
		return internalError(err)
	}
	if err := s.adminAuth.DeleteUser(ctx, user.Username); err != nil {
		return internalError(err)
	}

	// 外部アカウントとCognitoアカウントを連携
	linkParams := &cognito.LinkProviderParams{
		Username:     admin.CognitoID,
		ProviderType: provider.ProviderType.ToCognito(),
		AccountID:    provider.AccountID,
	}
	if err := s.adminAuth.LinkProvider(ctx, linkParams); err != nil {
		return internalError(err)
	}

	err = s.db.AdminAuthProvider.Upsert(ctx, provider)
	return internalError(err)
}
