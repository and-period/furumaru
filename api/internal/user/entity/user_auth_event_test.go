package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserAuthEvent(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &UserAuthEventParams{
		SessionID:    "session-id",
		ProviderType: UserAuthProviderTypeGoogle,
		Now:          now,
		TTL:          3600,
	}
	event := NewUserAuthEvent(params)
	t.Run("new", func(t *testing.T) {
		t.Parallel()
		expect := &UserAuthEvent{
			SessionID:    "session-id",
			ProviderType: UserAuthProviderTypeGoogle,
			Nonce:        event.Nonce, // ignore
			ExpiredAt:    now.Add(3600),
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		assert.Equal(t, expect, event)
	})
	t.Run("table name", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "user-auth-events", event.TableName())
	})
	t.Run("primary key", func(t *testing.T) {
		t.Parallel()
		expect := map[string]interface{}{
			"session_id": "session-id",
		}
		assert.Equal(t, expect, event.PrimaryKey())
	})
}
