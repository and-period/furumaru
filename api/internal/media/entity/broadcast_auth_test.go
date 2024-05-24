package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
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
			target: NewYoutubeBroadcastAuth,
			params: &BroadcastAuthParams{
				SessionID:  "session-id",
				Account:    "account",
				ScheduleID: "schedule-id",
				Now:        now,
				TTL:        1 * time.Hour,
			},
			expect: &BroadcastAuth{
				SessionID:  "session-id",
				Type:       BroadcastAuthTypeYoutube,
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

func TestBroadcastAuth_SetToken(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name   string
		auth   *BroadcastAuth
		token  *oauth2.Token
		hasErr bool
	}{
		{
			name: "success",
			auth: &BroadcastAuth{
				SessionID:  "session-id",
				Type:       BroadcastAuthTypeYoutube,
				Account:    "account",
				ScheduleID: "schedule-id",
				ExpiredAt:  now,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			token: &oauth2.Token{
				AccessToken:  "access-token",
				TokenType:    "Bearer",
				RefreshToken: "refresh-token",
				Expiry:       now.AddDate(0, 0, 1),
			},
			hasErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.auth.SetToken(tt.token)
			assert.Equal(t, tt.hasErr, err != nil, err)
		})
	}
}

func TestBroadcastAuth_GetToken(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name   string
		auth   *BroadcastAuth
		token  *oauth2.Token
		hasErr bool
	}{
		{
			name: "success",
			auth: &BroadcastAuth{
				SessionID:  "session-id",
				Type:       BroadcastAuthTypeYoutube,
				Account:    "account",
				ScheduleID: "schedule-id",
				ExpiredAt:  now,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			token: &oauth2.Token{
				AccessToken:  "access-token",
				TokenType:    "Bearer",
				RefreshToken: "refresh-token",
				Expiry:       now.AddDate(0, 0, 1),
			},
			hasErr: false,
		},
		{
			name: "empty",
			auth: &BroadcastAuth{
				SessionID:  "session-id",
				Type:       BroadcastAuthTypeYoutube,
				Account:    "account",
				ScheduleID: "schedule-id",
				ExpiredAt:  now,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			token:  nil,
			hasErr: true,
		},
		{
			name: "invalid token",
			auth: &BroadcastAuth{
				SessionID:  "session-id",
				Type:       BroadcastAuthTypeYoutube,
				Account:    "account",
				ScheduleID: "schedule-id",
				ExpiredAt:  now,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			token: &oauth2.Token{
				AccessToken:  "access-token",
				TokenType:    "Bearer",
				RefreshToken: "refresh-token",
				Expiry:       now.AddDate(0, 0, -1),
			},
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.auth.SetToken(tt.token)
			require.NoError(t, err)
			tt.auth.SetToken(tt.token)
			token, err := tt.auth.GetToken()
			if tt.hasErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.token, token)
		})
	}
}

func TestBroadcastAuth_ValidYoutubeAuth(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name     string
		auth     *BroadcastAuth
		channels []*youtube.Channel
		expect   bool
	}{
		{
			name: "success",
			auth: &BroadcastAuth{
				SessionID:  "session-id",
				Type:       BroadcastAuthTypeYoutube,
				Account:    "@handle",
				ScheduleID: "schedule-id",
				ExpiredAt:  now,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			channels: []*youtube.Channel{
				{
					Id:      "channel-id",
					Snippet: &youtube.ChannelSnippet{CustomUrl: "@handle"},
				},
			},
			expect: true,
		},
		{
			name: "invalid",
			auth: &BroadcastAuth{
				SessionID:  "session-id",
				Type:       BroadcastAuthTypeYoutube,
				Account:    "test@example.com",
				ScheduleID: "schedule-id",
				ExpiredAt:  now,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			channels: []*youtube.Channel{},
			expect:   false,
		},
		{
			name:     "empty",
			auth:     nil,
			channels: []*youtube.Channel{},
			expect:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.auth.ValidYoutubeAuth(tt.channels)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
