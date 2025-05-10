package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserNotification(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		userID string
		expect *UserNotification
	}{
		{
			name:   "success",
			userID: "user-id",
			expect: &UserNotification{
				UserID:   "user-id",
				Disabled: false,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewUserNotification(tt.userID)
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestUserNotification_Enabled(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		notification *UserNotification
		expect       bool
	}{
		{
			name: "success enabled",
			notification: &UserNotification{
				UserID:   "user-id",
				Disabled: false,
			},
			expect: true,
		},
		{
			name: "success disabled",
			notification: &UserNotification{
				UserID:   "user-id",
				Disabled: true,
			},
			expect: false,
		},
		{
			name:         "success nil",
			notification: nil,
			expect:       true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.notification.Enabled())
		})
	}
}

func TestUserNotifications_MapByUserID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		notifications UserNotifications
		expect        map[string]*UserNotification
	}{
		{
			name: "success",
			notifications: UserNotifications{
				{
					UserID:   "user-id",
					Disabled: false,
				},
			},
			expect: map[string]*UserNotification{
				"user-id": {
					UserID:   "user-id",
					Disabled: false,
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.notifications.MapByUserID()
			assert.Equal(t, tt.expect, actual)
		})
	}
}
