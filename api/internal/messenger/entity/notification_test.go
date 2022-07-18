package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestNotification_Fill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		notification *Notification
		expect       *Notification
		hasErr       bool
	}{
		{
			name: "success",
			notification: &Notification{
				ID:          "notification-id",
				Title:       "title",
				Body:        "<html>本文<html>",
				TargetsJSON: datatypes.JSON([]byte(`[1,2,3]`)),
			},
			expect: &Notification{
				ID:    "notification-id",
				Title: "title",
				Body:  "<html>本文<html>",
				Targets: []TargetType{
					PostTargetUsers,
					PostTargetProducers,
					PostTargetCoordinators,
				},
				TargetsJSON: datatypes.JSON([]byte(`[1,2,3]`)),
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.notification.Fill()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.notification)
		})
	}
}

func TestNotification_FillJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		notification *Notification
		expect       *Notification
		hasErr       bool
	}{
		{
			name: "success",
			notification: &Notification{
				ID:    "notification-id",
				Title: "title",
				Body:  "<html>本文<html>",
				Targets: []TargetType{
					PostTargetUsers,
					PostTargetProducers,
					PostTargetCoordinators,
				},
			},
			expect: &Notification{
				ID:    "notification-id",
				Title: "title",
				Body:  "<html>本文<html>",
				Targets: []TargetType{
					PostTargetUsers,
					PostTargetProducers,
					PostTargetCoordinators,
				},
				TargetsJSON: datatypes.JSON([]byte(`[1,2,3]`)),
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.notification.FillJSON()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.notification)
		})
	}
}

func TestNotifications_Fill(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		notifications Notifications
		expect        Notifications
		hasErr        bool
	}{
		{
			name: "success",
			notifications: Notifications{
				{
					ID:          "notification-id",
					Title:       "title",
					Body:        "<html>本文<html>",
					TargetsJSON: datatypes.JSON([]byte(`[1,2,3]`)),
				},
			},
			expect: Notifications{
				{
					ID:    "notification-id",
					Title: "title",
					Body:  "<html>本文<html>",
					Targets: []TargetType{
						PostTargetUsers,
						PostTargetProducers,
						PostTargetCoordinators,
					},
					TargetsJSON: datatypes.JSON([]byte(`[1,2,3]`)),
				},
			},
			hasErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.notifications.Fill()
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, tt.notifications)
		})
	}
}
