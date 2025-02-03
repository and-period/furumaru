package entity

import (
	"testing"

	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/stretchr/testify/assert"
)

func TestAdminAuthProviderType_ToCognito(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		auth   AdminAuthProviderType
		expect cognito.ProviderType
	}{
		{
			name:   "google",
			auth:   AdminAuthProviderTypeGoogle,
			expect: cognito.ProviderTypeGoogle,
		},
		{
			name:   "unknown",
			auth:   AdminAuthProviderTypeUnknown,
			expect: cognito.ProviderTypeUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.auth.ToCognito())
		})
	}
}

func TestAdminAuthProvider(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *AdminAuthProviderParams
		expect *AdminAuthProvider
		err    error
	}{
		{
			name: "success",
			params: &AdminAuthProviderParams{
				AdminID:      "admin-id",
				ProviderType: AdminAuthProviderTypeGoogle,
				Auth: &cognito.AuthUser{
					Username: "google_123",
					Email:    "test@example.com",
				},
			},
			expect: &AdminAuthProvider{
				AdminID:      "admin-id",
				ProviderType: AdminAuthProviderTypeGoogle,
				AccountID:    "123",
				Email:        "test@example.com",
			},
			err: nil,
		},
		{
			name: "invalid username",
			params: &AdminAuthProviderParams{
				AdminID:      "admin-id",
				ProviderType: AdminAuthProviderTypeGoogle,
				Auth: &cognito.AuthUser{
					Username: "google",
					Email:    "test@example.com",
				},
			},
			expect: nil,
			err:    errInvalidAuthUsername,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			auth, err := NewAdminAuthProvider(tt.params)
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.expect, auth)
		})
	}
}
