package entity

import (
	"testing"
	"time"

	"github.com/and-period/furumaru/api/pkg/cognito"
	"github.com/stretchr/testify/assert"
)

func TestAdminAuth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		admin  *Admin
		result *cognito.AuthResult
		expect *AdminAuth
	}{
		{
			name: "success",
			admin: &Admin{
				ID:            "admin-id",
				CognitoID:     "cognito-id",
				Type:          AdminTypeAdministrator,
				GroupIDs:      []string{"group-id"},
				Lastname:      "&.",
				Firstname:     "スタッフ",
				LastnameKana:  "あんどどっと",
				FirstnameKana: "すたっふ",
			},
			result: &cognito.AuthResult{
				IDToken:      "id-token",
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
			expect: &AdminAuth{
				AdminID:      "admin-id",
				Type:         AdminTypeAdministrator,
				GroupIDs:     []string{"group-id"},
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    3600,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewAdminAuth(tt.admin, tt.result)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestAdminAuthEvent(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &AdminAuthEventParams{
		AdminID:      "admin-id",
		ProviderType: "provider-type",
		Now:          now,
		TTL:          3600,
	}
	event := NewAdminAuthEvent(params)
	t.Run("new", func(t *testing.T) {
		t.Parallel()
		expect := &AdminAuthEvent{
			AdminID:      "admin-id",
			ProviderType: "provider-type",
			Nonce:        event.Nonce, // ignore
			ExpiredAt:    now.Add(3600),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		assert.Equal(t, expect, event)
	})
	t.Run("table name", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "admin_oauth_events", event.TableName())
	})
	t.Run("primary key", func(t *testing.T) {
		t.Parallel()
		expect := map[string]interface{}{
			"admin_id":      "admin-id",
			"provider_type": "provider-type",
		}
		assert.Equal(t, expect, event.PrimaryKey())
	})
}
