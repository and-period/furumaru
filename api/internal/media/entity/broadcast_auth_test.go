package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBroadcastAuth(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name   string
		target func(params *BroadcastAuthParams) *BroadcastAuth
		params *BroadcastAuthParams
		expect *BroadcastAuth
	}{
		{
			name:   "success youtube",
			target: NewYouTubeBroadcastAuth,
			params: &BroadcastAuthParams{
				SessionID:  "session-id",
				Account:    "account",
				ScheduleID: "schedule-id",
				Now:        now,
				TTL:        1 * time.Hour,
			},
			expect: &BroadcastAuth{
				SessionID:  "session-id",
				Type:       BroadcastAuthTypeYouTube,
				Account:    "account",
				ScheduleID: "schedule-id",
				ExpiredAt:  now.Add(1 * time.Hour),
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			primaryKey := map[string]interface{}{
				"session_id": tt.expect.SessionID,
			}
			actual := tt.target(tt.params)
			assert.Equal(t, tt.expect, actual)
			assert.Equal(t, "broadcast-auth", actual.TableName())
			assert.Equal(t, primaryKey, actual.PrimaryKey())
		})
	}
}

func TestBroadcastAuth_ValidYouTubeAuth(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name   string
		auth   *BroadcastAuth
		email  string
		expect bool
	}{
		{
			name: "success",
			auth: &BroadcastAuth{
				SessionID:  "session-id",
				Type:       BroadcastAuthTypeYouTube,
				Account:    "test@example.com",
				ScheduleID: "schedule-id",
				ExpiredAt:  now,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			email:  "test@example.com",
			expect: true,
		},
		{
			name:   "empty",
			auth:   nil,
			email:  "test@example.com",
			expect: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.auth.ValidYouTubeAuth(tt.email)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
