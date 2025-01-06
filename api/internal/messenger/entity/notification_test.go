package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotification(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name   string
		params *NewNotificationParams
		expect *Notification
	}{
		{
			name: "success system",
			params: &NewNotificationParams{
				Type: NotificationTypeSystem,
				Targets: []NotificationTarget{
					NotificationTargetUsers,
					NotificationTargetProducers,
					NotificationTargetCoordinators,
					NotificationTargetAdministrators,
				},
				Title:       "タイトル",
				Body:        "<html>本文<html>",
				Note:        "<html>備考<html>",
				PublishedAt: now,
				PromotionID: "",
				CreatedBy:   "admin-id",
			},
			expect: &Notification{
				ID:    "", // ignore
				Type:  NotificationTypeSystem,
				Title: "タイトル",
				Body:  "<html>本文<html>",
				Note:  "<html>備考<html>",
				Targets: []NotificationTarget{
					NotificationTargetUsers,
					NotificationTargetProducers,
					NotificationTargetCoordinators,
					NotificationTargetAdministrators,
				},
				PublishedAt: now,
				PromotionID: "",
				CreatedBy:   "admin-id",
				UpdatedBy:   "admin-id",
			},
		},
		{
			name: "success promotion",
			params: &NewNotificationParams{
				Type: NotificationTypeSystem,
				Targets: []NotificationTarget{
					NotificationTargetUsers,
				},
				Title:       "",
				Body:        "<html>本文<html>",
				Note:        "<html>備考<html>",
				PublishedAt: now,
				PromotionID: "promotion-id",
				CreatedBy:   "admin-id",
			},
			expect: &Notification{
				ID:    "", // ignore
				Type:  NotificationTypeSystem,
				Title: "",
				Body:  "<html>本文<html>",
				Note:  "<html>備考<html>",
				Targets: []NotificationTarget{
					NotificationTargetUsers,
				},
				PublishedAt: now,
				PromotionID: "promotion-id",
				CreatedBy:   "admin-id",
				UpdatedBy:   "admin-id",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := NewNotification(tt.params)
			actual.ID = "" // ignore
			assert.Equal(t, tt.expect, actual)
		})
	}
}

func TestNotification_Validate(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name         string
		notification *Notification
		expect       error
	}{
		{
			name: "success system",
			notification: &Notification{
				Type:        NotificationTypeSystem,
				Title:       "タイトル",
				Body:        "本文",
				Note:        "備考",
				Targets:     []NotificationTarget{NotificationTargetUsers},
				PublishedAt: now.AddDate(0, 0, 1),
				PromotionID: "",
				CreatedBy:   "admin-id",
				UpdatedBy:   "admin-id",
			},
			expect: nil,
		},
		{
			name: "success promotion",
			notification: &Notification{
				Type:        NotificationTypePromotion,
				Title:       "",
				Body:        "本文",
				Note:        "備考",
				Targets:     []NotificationTarget{NotificationTargetUsers},
				PublishedAt: now.AddDate(0, 0, 1),
				PromotionID: "promotion-id",
				CreatedBy:   "admin-id",
				UpdatedBy:   "admin-id",
			},
			expect: nil,
		},
		{
			name: "already published",
			notification: &Notification{
				Type:        NotificationTypeSystem,
				Title:       "タイトル",
				Body:        "本文",
				Note:        "備考",
				Targets:     []NotificationTarget{NotificationTargetUsers},
				PublishedAt: now.AddDate(0, 0, -1),
				PromotionID: "",
				CreatedBy:   "admin-id",
				UpdatedBy:   "admin-id",
			},
			expect: ErrNotificationAlreadyPublished,
		},
		{
			name: "incorrect targets",
			notification: &Notification{
				Type:        NotificationTypeSystem,
				Title:       "タイトル",
				Body:        "本文",
				Note:        "備考",
				Targets:     []NotificationTarget{},
				PublishedAt: now.AddDate(0, 0, 1),
				PromotionID: "",
				CreatedBy:   "admin-id",
				UpdatedBy:   "admin-id",
			},
			expect: ErrNotificationIncorrectTargets,
		},
		{
			name: "duplicated targets",
			notification: &Notification{
				Type:  NotificationTypeSystem,
				Title: "タイトル",
				Body:  "本文",
				Note:  "備考",
				Targets: []NotificationTarget{
					NotificationTargetUsers,
					NotificationTargetUsers,
				},
				PublishedAt: now.AddDate(0, 0, 1),
				PromotionID: "",
				CreatedBy:   "admin-id",
				UpdatedBy:   "admin-id",
			},
			expect: ErrNotificationDuplicatedTargets,
		},
		{
			name: "required title",
			notification: &Notification{
				Type:        NotificationTypeSystem,
				Title:       "",
				Body:        "本文",
				Note:        "備考",
				Targets:     []NotificationTarget{NotificationTargetUsers},
				PublishedAt: now.AddDate(0, 0, 1),
				PromotionID: "",
				CreatedBy:   "admin-id",
				UpdatedBy:   "admin-id",
			},
			expect: ErrNotificationRequiredTitle,
		},
		{
			name: "required promotion id",
			notification: &Notification{
				Type:        NotificationTypePromotion,
				Title:       "",
				Body:        "本文",
				Note:        "備考",
				Targets:     []NotificationTarget{NotificationTargetUsers},
				PublishedAt: now.AddDate(0, 0, 1),
				PromotionID: "",
				CreatedBy:   "admin-id",
				UpdatedBy:   "admin-id",
			},
			expect: ErrNotificationRequiredPromotionID,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.notification.Validate(now)
			assert.ErrorIs(t, err, tt.expect)
		})
	}
}

func TestNotification_Fill(t *testing.T) {
	t.Parallel()
	now := time.Now()
	tests := []struct {
		name         string
		notification *Notification
		expect       *Notification
	}{
		{
			name: "success waiting",
			notification: &Notification{
				ID:          "notification-id",
				Title:       "title",
				Body:        "<html>本文<html>",
				PublishedAt: now.AddDate(0, 0, 1),
			},
			expect: &Notification{
				ID:          "notification-id",
				Title:       "title",
				Body:        "<html>本文<html>",
				Status:      NotificationStatusWaiting,
				Targets:     []NotificationTarget{},
				PublishedAt: now.AddDate(0, 0, 1),
			},
		},
		{
			name: "success notified",
			notification: &Notification{
				ID:          "notification-id",
				Title:       "title",
				Body:        "<html>本文<html>",
				PublishedAt: now.AddDate(0, 0, -1),
			},
			expect: &Notification{
				ID:          "notification-id",
				Title:       "title",
				Body:        "<html>本文<html>",
				Status:      NotificationStatusNotified,
				PublishedAt: now.AddDate(0, 0, -1),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.notification.Fill(now)
			assert.Equal(t, tt.expect, tt.notification)
		})
	}
}

func TestNotification_TemplateID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		notification *Notification
		expect       MessageTemplateID
	}{
		{
			name: "system",
			notification: &Notification{
				Type: NotificationTypeSystem,
			},
			expect: MessageTemplateIDNotificationSystem,
		},
		{
			name: "live",
			notification: &Notification{
				Type: NotificationTypeLive,
			},
			expect: MessageTemplateIDNotificationLive,
		},
		{
			name: "promotion",
			notification: &Notification{
				Type: NotificationTypePromotion,
			},
			expect: MessageTemplateIDNotificationPromotion,
		},
		{
			name: "other",
			notification: &Notification{
				Type: NotificationTypeOther,
			},
			expect: MessageTemplateIDNotificationOther,
		},
		{
			name: "unknown",
			notification: &Notification{
				Type: NotificationTypeUnknown,
			},
			expect: MessageTemplateIDNotificationOther,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.notification.TemplateID())
		})
	}
}

func TestNotification_HasTarget(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                   string
		notification           *Notification
		hasUserTarget          bool
		hasAdminTarget         bool
		hasAdministratorTarget bool
		hasCoordinatorTarget   bool
		hasProducerTarget      bool
	}{
		{
			name: "contain user target",
			notification: &Notification{
				Targets: []NotificationTarget{NotificationTargetUsers},
			},
			hasUserTarget:          true,
			hasAdminTarget:         false,
			hasAdministratorTarget: false,
			hasCoordinatorTarget:   false,
			hasProducerTarget:      false,
		},
		{
			name: "contain administrator target",
			notification: &Notification{
				Targets: []NotificationTarget{NotificationTargetAdministrators},
			},
			hasUserTarget:          false,
			hasAdminTarget:         true,
			hasAdministratorTarget: true,
			hasCoordinatorTarget:   false,
			hasProducerTarget:      false,
		},
		{
			name: "contain coordinator target",
			notification: &Notification{
				Targets: []NotificationTarget{NotificationTargetCoordinators},
			},
			hasUserTarget:          false,
			hasAdminTarget:         true,
			hasAdministratorTarget: false,
			hasCoordinatorTarget:   true,
			hasProducerTarget:      false,
		},
		{
			name: "contain producer target",
			notification: &Notification{
				Targets: []NotificationTarget{NotificationTargetProducers},
			},
			hasUserTarget:          false,
			hasAdminTarget:         true,
			hasAdministratorTarget: false,
			hasCoordinatorTarget:   false,
			hasProducerTarget:      true,
		},
		{
			name: "contain unknown target",
			notification: &Notification{
				Targets: []NotificationTarget{NotificationTargetUnknown},
			},
			hasUserTarget:          false,
			hasAdminTarget:         false,
			hasAdministratorTarget: false,
			hasCoordinatorTarget:   false,
			hasProducerTarget:      false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.hasUserTarget, tt.notification.HasUserTarget())
			assert.Equal(t, tt.hasAdminTarget, tt.notification.HasAdminTarget())
			assert.Equal(t, tt.hasAdministratorTarget, tt.notification.HasAdministratorTarget())
			assert.Equal(t, tt.hasCoordinatorTarget, tt.notification.HasCoordinatorTarget())
			assert.Equal(t, tt.hasProducerTarget, tt.notification.HasProducerTarget())
		})
	}
}

func TestNotifications_AdminIDs(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tests := []struct {
		name          string
		notifications Notifications
		expect        []string
		hasErr        bool
	}{
		{
			name: "success",
			notifications: Notifications{
				{
					ID:          "notification-id01",
					Title:       "title",
					Body:        "<html>本文<html>",
					Type:        NotificationTypeSystem,
					PromotionID: "invalid-id",
					PublishedAt: now.AddDate(0, 0, -1),
					CreatedBy:   "admin-id01",
					UpdatedBy:   "admin-id02",
				},
				{
					ID:          "notification-id02",
					Title:       "title",
					Body:        "<html>本文<html>",
					Type:        NotificationTypePromotion,
					PromotionID: "promotion-id",
					PublishedAt: now.AddDate(0, 0, -1),
					CreatedBy:   "admin-id02",
					UpdatedBy:   "admin-id02",
				},
				{
					ID:          "notification-id03",
					Title:       "title",
					Body:        "<html>本文<html>",
					Type:        NotificationTypePromotion,
					PromotionID: "promotion-id",
					PublishedAt: now.AddDate(0, 0, -1),
					CreatedBy:   "admin-id03",
					UpdatedBy:   "admin-id03",
				},
			},
			expect: []string{"admin-id01", "admin-id02", "admin-id03"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.notifications.AdminIDs()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestNotifications_PromotionIDs(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tests := []struct {
		name          string
		notifications Notifications
		expect        []string
		hasErr        bool
	}{
		{
			name: "success",
			notifications: Notifications{
				{
					ID:          "notification-id01",
					Title:       "title",
					Body:        "<html>本文<html>",
					Type:        NotificationTypeSystem,
					PromotionID: "invalid-id",
					PublishedAt: now.AddDate(0, 0, -1),
				},
				{
					ID:          "notification-id02",
					Title:       "title",
					Body:        "<html>本文<html>",
					Type:        NotificationTypePromotion,
					PromotionID: "promotion-id",
					PublishedAt: now.AddDate(0, 0, -1),
				},
				{
					ID:          "notification-id03",
					Title:       "title",
					Body:        "<html>本文<html>",
					Type:        NotificationTypePromotion,
					PromotionID: "promotion-id",
					PublishedAt: now.AddDate(0, 0, -1),
				},
			},
			expect: []string{"promotion-id"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.notifications.PromotionIDs()
			assert.ElementsMatch(t, tt.expect, actual)
		})
	}
}

func TestNotifications_Fill(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tests := []struct {
		name          string
		notifications Notifications
		expect        Notifications
	}{
		{
			name: "success",
			notifications: Notifications{
				{
					ID:          "notification-id",
					Title:       "title",
					Body:        "<html>本文<html>",
					PublishedAt: now.AddDate(0, 0, -1),
				},
			},
			expect: Notifications{
				{
					ID:          "notification-id",
					Title:       "title",
					Body:        "<html>本文<html>",
					Status:      NotificationStatusNotified,
					PublishedAt: now.AddDate(0, 0, -1),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.notifications.Fill(now)
			assert.Equal(t, tt.expect, tt.notifications)
		})
	}
}
