package service

import (
	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/types"
	"github.com/and-period/furumaru/api/internal/gateway/util"
	uentity "github.com/and-period/furumaru/api/internal/user/entity"
)

type AuthProviderType int32

const (
	AuthProviderTypeUnknown AuthProviderType = 0
	AuthProviderTypeGoogle  AuthProviderType = 1
	AuthProviderTypeLINE    AuthProviderType = 2
)

func NewAuthProviderType(t uentity.AdminAuthProviderType) AuthProviderType {
	switch t {
	case uentity.AdminAuthProviderTypeGoogle:
		return AuthProviderTypeGoogle
	case uentity.AdminAuthProviderTypeLINE:
		return AuthProviderTypeLINE
	default:
		return AuthProviderTypeUnknown
	}
}

func (t AuthProviderType) Response() int32 {
	return int32(t)
}

type Auth struct {
	types.Auth
	GroupIDs []string
}

func NewAuth(auth *uentity.AdminAuth) *Auth {
	return &Auth{
		Auth: types.Auth{
			AdminID:      auth.AdminID,
			Type:         NewAdminType(auth.Type).Response(),
			AccessToken:  auth.AccessToken,
			RefreshToken: auth.RefreshToken,
			ExpiresIn:    auth.ExpiresIn,
			TokenType:    util.AuthTokenType,
		},
		GroupIDs: auth.GroupIDs,
	}
}

func (a *Auth) Response() *types.Auth {
	return &a.Auth
}

type AuthUser struct {
	types.AuthUser
}

func (a *AuthUser) Response() *types.AuthUser {
	return &a.AuthUser
}

type AuthProvider struct {
	types.AuthProvider
}

type AuthProviders []*AuthProvider

func NewAuthProvider(provider *uentity.AdminAuthProvider) *AuthProvider {
	return &AuthProvider{
		AuthProvider: types.AuthProvider{
			Type:        NewAuthProviderType(provider.ProviderType).Response(),
			ConnectedAt: provider.UpdatedAt.Unix(),
		},
	}
}

func (p *AuthProvider) Response() *types.AuthProvider {
	return &p.AuthProvider
}

func NewAuthProviders(providers uentity.AdminAuthProviders) AuthProviders {
	res := make(AuthProviders, len(providers))
	for i := range providers {
		res[i] = NewAuthProvider(providers[i])
	}
	return res
}

func (ps AuthProviders) Response() []*types.AuthProvider {
	res := make([]*types.AuthProvider, len(ps))
	for i := range ps {
		res[i] = ps[i].Response()
	}
	return res
}
