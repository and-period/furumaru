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

func TestNotification_HasTarget(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                 string
		notification         *Notification
		hasUserTarget        bool
		hasAdminTarget       bool
		hasCoordinatorTarget bool
		hasProducerTarget    bool
	}{
		{
			name: "contain user target",
			notification: &Notification{
				Targets: []TargetType{PostTargetUsers},
			},
			hasUserTarget:        true,
			hasAdminTarget:       false,
			hasCoordinatorTarget: false,
			hasProducerTarget:    false,
		},
		{
			name: "contain coordinator target",
			notification: &Notification{
				Targets: []TargetType{PostTargetCoordinators},
			},
			hasUserTarget:        false,
			hasAdminTarget:       true,
			hasCoordinatorTarget: true,
			hasProducerTarget:    false,
		},
		{
			name: "contain producer target",
			notification: &Notification{
				Targets: []TargetType{PostTargetProducers},
			},
			hasUserTarget:        false,
			hasAdminTarget:       true,
			hasCoordinatorTarget: false,
			hasProducerTarget:    true,
		},
		{
			name: "contain unknown target",
			notification: &Notification{
				Targets: []TargetType{PostTargetUnknown},
			},
			hasUserTarget:        false,
			hasAdminTarget:       false,
			hasCoordinatorTarget: false,
			hasProducerTarget:    false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.hasUserTarget, tt.notification.HasUserTarget())
			assert.Equal(t, tt.hasAdminTarget, tt.notification.HasAdminTarget())
			assert.Equal(t, tt.hasCoordinatorTarget, tt.notification.HasCoordinatorTarget())
			assert.Equal(t, tt.hasProducerTarget, tt.notification.HasProducerTarget())
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

func TestNotification_Marshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		targets []TargetType
		expect  []byte
		hasErr  bool
	}{
		{
			name: "success",
			targets: []TargetType{
				PostTargetProducers,
				PostTargetCoordinators,
			},
			expect: []byte(`[2,3]`),
			hasErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual, err := Marshal(tt.targets)
			assert.Equal(t, tt.hasErr, err != nil, err)
			assert.Equal(t, tt.expect, actual)
		})
	}
}
