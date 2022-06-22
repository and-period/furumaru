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
				TargetsJSON: datatypes.JSON([]byte(`[{"postTarget":1}, {"postTarget":2}]`)),
			},
			expect: &Notification{
				ID:    "notification-id",
				Title: "title",
				Body:  "<html>本文<html>",
				Targets: PostTargetList{
					{
						PostTarget: 1,
					},
					{
						PostTarget: 2,
					},
				},
				TargetsJSON: datatypes.JSON([]byte(`[{"postTarget":1}, {"postTarget":2}]`)),
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
				Targets: PostTargetList{
					{
						PostTarget: 1,
					},
					{
						PostTarget: 2,
					},
				},
			},
			expect: &Notification{
				ID:    "notification-id",
				Title: "title",
				Body:  "<html>本文<html>",
				Targets: PostTargetList{
					{
						PostTarget: 1,
					},
					{
						PostTarget: 2,
					},
				},
				TargetsJSON: datatypes.JSON([]byte(`[{"postTarget":1},{"postTarget":2}]`)),
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
