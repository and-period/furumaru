package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	"github.com/and-period/furumaru/api/internal/user/entity"
)

type AuthProviderType int32

const (
	AuthProviderTypeUnknown AuthProviderType = 0
	AuthProviderTypeGoogle  AuthProviderType = 1
)

func NewAuthProviderType(t entity.AdminAuthProviderType) AuthProviderType {
	switch t {
	case entity.AdminAuthProviderTypeGoogle:
		return AuthProviderTypeGoogle
	default:
		return AuthProviderTypeUnknown
	}
}

func (t AuthProviderType) Response() int32 {
	return int32(t)
}

type Auth struct {
	response.Auth
}

type AuthUser struct {
	response.AuthUser
}

func NewAuth(auth *entity.AdminAuth) *Auth {
	return &Auth{
		Auth: response.Auth{
			AdminID:      auth.AdminID,
			Type:         NewAdminType(auth.Type).Response(),
			AccessToken:  auth.AccessToken,
			RefreshToken: auth.RefreshToken,
			ExpiresIn:    auth.ExpiresIn,
			TokenType:    util.AuthTokenType,
		},
	}
}

func (a *Auth) Response() *response.Auth {
	return &a.Auth
}

func (a *AuthUser) Response() *response.AuthUser {
	return &a.AuthUser
}

type AuthProvider struct {
	response.AuthProvider
}

type AuthProviders []*AuthProvider

func NewAuthProvider(provider *entity.AdminAuthProvider) *AuthProvider {
	return &AuthProvider{
		AuthProvider: response.AuthProvider{
			Type:        NewAuthProviderType(provider.ProviderType).Response(),
			ConnectedAt: provider.UpdatedAt.Unix(),
		},
	}
}

func (p *AuthProvider) Response() *response.AuthProvider {
	return &p.AuthProvider
}

func NewAuthProviders(providers entity.AdminAuthProviders) AuthProviders {
	res := make(AuthProviders, len(providers))
	for i := range providers {
		res[i] = NewAuthProvider(providers[i])
	}
	return res
}

func (ps AuthProviders) Response() []*response.AuthProvider {
	res := make([]*response.AuthProvider, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}
