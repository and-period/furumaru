package service

import (
	"context"
	"fmt"

	"github.com/and-period/furumaru/api/internal/exception"
	"github.com/and-period/furumaru/api/internal/user"
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

func (s *service) InitialGoogleAdminAuth(ctx context.Context, in *user.InitialGoogleAdminAuthInput) (string, error) {
	if err := s.validator.Struct(in); err != nil {
		return "", internalError(err)
	}
	// TODO: すでに連携済みかの検証
	eventParams := &entity.AdminAuthEventParams{
		AdminID: in.AdminID,
		Now:     s.now(),
		TTL:     s.adminAuthTTL,
	}
	event := entity.NewAdminAuthEvent(eventParams)
	if err := s.cache.Insert(ctx, event); err != nil {
		return "", internalError(err)
	}
	redirectURL := s.adminAuthGoogleRedirectURL
	if in.RedirectURI != "" {
		redirectURL = in.RedirectURI
	}
	params := &cognito.GenerateAuthURLParams{
		State:       in.State,
		Nonce:       event.Nonce,
		RedirectURI: redirectURL,
	}
	authURL, err := s.adminAuth.GenerateAuthURL(ctx, params)
	return authURL, internalError(err)
}

func (s *service) ConnectGoogleAdminAuth(ctx context.Context, in *user.ConnectGoogleAdminAuthInput) error {
	if err := s.validator.Struct(in); err != nil {
		return internalError(err)
	}
	event := &entity.AdminAuthEvent{AdminID: in.AdminID}
	if err := s.cache.Get(ctx, event); err != nil {
		return internalError(err)
	}
	if event.Nonce != in.Nonce {
		return fmt.Errorf("service: invalid nonce for google auth: %w", exception.ErrFailedPrecondition)
	}
	admin, err := s.db.Admin.Get(ctx, in.AdminID)
	if err != nil {
		return internalError(err)
	}
	redirectURI := s.adminAuthGoogleRedirectURL
	if in.RedirectURI != "" {
		redirectURI = in.RedirectURI
	}
	tokenParams := &cognito.GetAccessTokenParams{
		Code:        in.Code,
		RedirectURI: redirectURI,
	}
	token, err := s.adminAuth.GetAccessToken(ctx, tokenParams)
	if err != nil {
		return internalError(err)
	}
	username, err := s.adminAuth.GetUsername(ctx, token.AccessToken)
	if err != nil {
		return internalError(err)
	}
	// Cognitoの仕様で、すでにサインイン済みの場合は連携できないため
	if err := s.adminAuth.DeleteUser(ctx, username); err != nil {
		return internalError(err)
	}
	linkParams := &cognito.LinkProviderParams{
		Username:     admin.CognitoID,
		ProviderType: cognito.ProviderTypeGoogle,
		AccountID:    username,
	}
	if err := s.adminAuth.LinkProvider(ctx, linkParams); err != nil {
		return internalError(err)
	}
	// TODO: 連携情報をDBに登録
	return nil
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
