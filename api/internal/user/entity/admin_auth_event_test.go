package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAdminAuthEvent(t *testing.T) {
	t.Parallel()
	now := time.Now()
	params := &AdminAuthEventParams{
		AdminID: "admin-id",
		Now:     now,
		TTL:     3600,
	}
	event := NewAdminAuthEvent(params)
	t.Run("new", func(t *testing.T) {
		t.Parallel()
		expect := &AdminAuthEvent{
			AdminID:   "admin-id",
			Nonce:     event.Nonce, // ignore
			ExpiredAt: now.Add(3600),
			CreatedAt: now,
			UpdatedAt: now,
		}
		assert.Equal(t, expect, event)
	})
	t.Run("table name", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "admin-auth-events", event.TableName())
	})
	t.Run("primary key", func(t *testing.T) {
		t.Parallel()
		expect := map[string]interface{}{
			"admin_id": "admin-id",
		}
		assert.Equal(t, expect, event.PrimaryKey())
	})
}
