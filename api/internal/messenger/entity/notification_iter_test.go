package entity

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotifications_All(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name          string
		notifications Notifications
	}{
		{
			name: "success",
			notifications: Notifications{
				{
					ID:          "notification-id01",
					Type:        NotificationTypeOther,
					Title:       "お知らせ1",
					CreatedBy:   "admin-id01",
					UpdatedBy:   "admin-id01",
					PublishedAt: now,
				},
				{
					ID:          "notification-id02",
					Type:        NotificationTypePromotion,
					Title:       "お知らせ2",
					PromotionID: "promotion-id01",
					CreatedBy:   "admin-id02",
					UpdatedBy:   "admin-id02",
					PublishedAt: now,
				},
			},
		},
		{
			name:          "empty",
			notifications: Notifications{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var indices []int
			var ids []string
			for i, n := range tt.notifications.All() {
				indices = append(indices, i)
				ids = append(ids, n.ID)
			}
			for i, n := range tt.notifications {
				if i < len(indices) {
					assert.Equal(t, i, indices[i])
					assert.Equal(t, n.ID, ids[i])
				}
			}
			assert.Len(t, indices, len(tt.notifications))
		})
	}
}

func TestNotifications_All_EarlyBreak(t *testing.T) {
	t.Parallel()
	notifications := Notifications{
		{ID: "notification-id01"},
		{ID: "notification-id02"},
		{ID: "notification-id03"},
	}
	var count int
	for range notifications.All() {
		count++
		if count == 2 {
			break
		}
	}
	assert.Equal(t, 2, count)
}

func TestNotifications_IterPromotionIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		notifications Notifications
		expect        []string
	}{
		{
			name: "success with mixed types",
			notifications: Notifications{
				{ID: "notification-id01", Type: NotificationTypeOther, Title: "お知らせ"},
				{ID: "notification-id02", Type: NotificationTypePromotion, PromotionID: "promotion-id01"},
				{ID: "notification-id03", Type: NotificationTypeSystem, Title: "システム"},
				{ID: "notification-id04", Type: NotificationTypePromotion, PromotionID: "promotion-id02"},
			},
			expect: []string{"promotion-id01", "promotion-id02"},
		},
		{
			name: "no promotion type",
			notifications: Notifications{
				{ID: "notification-id01", Type: NotificationTypeOther, Title: "お知らせ"},
			},
			expect: nil,
		},
		{
			name:          "empty",
			notifications: Notifications{},
			expect:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := slices.Collect(tt.notifications.IterPromotionIDs())
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestNotifications_IterAdminIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		notifications Notifications
		expect        []string
	}{
		{
			name: "success",
			notifications: Notifications{
				{ID: "notification-id01", CreatedBy: "admin-id01", UpdatedBy: "admin-id01"},
				{ID: "notification-id02", CreatedBy: "admin-id02", UpdatedBy: "admin-id03"},
			},
			expect: []string{"admin-id01", "admin-id01", "admin-id02", "admin-id03"},
		},
		{
			name:          "empty",
			notifications: Notifications{},
			expect:        nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := slices.Collect(tt.notifications.IterAdminIDs())
			assert.Equal(t, tt.expect, actual)
		})
	}
}
