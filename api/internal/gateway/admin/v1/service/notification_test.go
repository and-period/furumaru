package service

import (
	"testing"

	"github.com/and-period/furumaru/api/internal/gateway/admin/v1/response"
	"github.com/and-period/furumaru/api/internal/messenger/entity"
	"github.com/and-period/furumaru/api/pkg/jst"
	"github.com/stretchr/testify/assert"
)

func TestNotification(t *testing.T) {
	t.Parallel()
	var date int64 = 1640962800

	tests := []struct {
		name         string
		notification *entity.Notification
		expect       *Notification
	}{
		{
			name: "success",
			notification: &entity.Notification{
				ID:        "notification-id",
				CreatedBy: "admin-id",
				UpdatedBy: "admin-id",
				Title:     "キャベツ祭り開催",
				Body:      "旬のキャベツを大安売り",
				Targets: []entity.NotificationTarget{
					entity.NotificationTargetUsers,
					entity.NotificationTargetProducers,
				},
				PublishedAt: jst.ParseFromUnix(date),
				CreatedAt:   jst.ParseFromUnix(date),
				UpdatedAt:   jst.ParseFromUnix(date),
			},
			expect: &Notification{
				Notification: response.Notification{
					ID:          "notification-id",
					CreatedBy:   "admin-id",
					UpdatedBy:   "admin-id",
					Title:       "キャベツ祭り開催",
					Body:        "旬のキャベツを大安売り",
					Targets:     []int32{1, 2},
					PublishedAt: 1640962800,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewNotification(tt.notification))
		})
	}
}

func TestNotification_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		notification *Notification
		promotion    *Promotion
		expect       *Notification
	}{
		{
			name: "success",
			notification: &Notification{
				Notification: response.Notification{
					ID:          "notification-id",
					CreatedBy:   "admin-id",
					UpdatedBy:   "admin-id",
					Title:       "キャベツ祭り開催",
					Body:        "旬のキャベツを大安売り",
					Targets:     []int32{3, 4},
					PublishedAt: 1640962800,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
			},
			promotion: &Promotion{
				Promotion: response.Promotion{
					ID:           "promotion-id",
					Title:        "セール情報",
					Description:  "セール詳細",
					Public:       true,
					DiscountType: DiscountTypeAmount.Response(),
					DiscountRate: 3980,
					Code:         "code",
					StartAt:      1640962800,
					EndAt:        1640962800,
					UsedCount:    0,
					UsedAmount:   0,
					CreatedAt:    1640962800,
					UpdatedAt:    1640962800,
				},
			},
			expect: &Notification{
				Notification: response.Notification{
					ID:          "notification-id",
					CreatedBy:   "admin-id",
					UpdatedBy:   "admin-id",
					Title:       "キャベツ祭り開催",
					Body:        "旬のキャベツを大安売り",
					Targets:     []int32{3, 4},
					PublishedAt: 1640962800,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.notification.Fill(tt.promotion)
			assert.Equal(t, tt.expect, tt.notification)
		})
	}
}

func TestNotification_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		notification *Notification
		expect       *response.Notification
	}{
		{
			name: "success",
			notification: &Notification{
				Notification: response.Notification{
					ID:          "notification-id",
					CreatedBy:   "admin-id",
					UpdatedBy:   "admin-id",
					Title:       "キャベツ祭り開催",
					Body:        "旬のキャベツを大安売り",
					Targets:     []int32{3, 4},
					PublishedAt: 1640962800,

					CreatedAt: 1640962800,
					UpdatedAt: 1640962800,
				},
			},
			expect: &response.Notification{
				ID:          "notification-id",
				CreatedBy:   "admin-id",
				UpdatedBy:   "admin-id",
				Title:       "キャベツ祭り開催",
				Body:        "旬のキャベツを大安売り",
				Targets:     []int32{3, 4},
				PublishedAt: 1640962800,
				CreatedAt:   1640962800,
				UpdatedAt:   1640962800,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.notification.Response())
		})
	}
}

func TestNotifications(t *testing.T) {
	t.Parallel()

	var date int64 = 1640962800
	tests := []struct {
		name          string
		notifications entity.Notifications
		expect        Notifications
	}{
		{
			name: "success",
			notifications: entity.Notifications{
				{
					ID:        "notification-id",
					CreatedBy: "admin-id",
					UpdatedBy: "admin-id",
					Title:     "キャベツ祭り開催",
					Body:      "旬のキャベツを大安売り",
					Targets: []entity.NotificationTarget{
						entity.NotificationTargetUsers,
						entity.NotificationTargetProducers,
					},
					PublishedAt: jst.ParseFromUnix(date),
					CreatedAt:   jst.ParseFromUnix(date),
					UpdatedAt:   jst.ParseFromUnix(date),
				},
			},
			expect: Notifications{
				{
					Notification: response.Notification{
						ID:          "notification-id",
						CreatedBy:   "admin-id",
						UpdatedBy:   "admin-id",
						Title:       "キャベツ祭り開催",
						Body:        "旬のキャベツを大安売り",
						Targets:     []int32{1, 2},
						PublishedAt: 1640962800,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewNotifications(tt.notifications))
		})
	}
}

func TestNotifications_AdminIDs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		notifications Notifications
		expect        []string
	}{
		{
			name: "success",
			notifications: Notifications{
				{
					Notification: response.Notification{
						ID:          "notification-id",
						CreatedBy:   "admin-id",
						UpdatedBy:   "admin-id",
						Title:       "キャベツ祭り開催",
						Body:        "旬のキャベツを大安売り",
						Targets:     []int32{3, 4},
						PublishedAt: 1640962800,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
			},
			expect: []string{"admin-id"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.ElementsMatch(t, tt.expect, tt.notifications.AdminIDs())
		})
	}
}

func TestNotifications_Fill(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		notifications Notifications
		promotions    map[string]*Promotion
		expect        Notifications
	}{
		{
			name: "success",
			notifications: Notifications{
				{
					Notification: response.Notification{
						ID:          "notification-id",
						Type:        NotificationTypeSystem.Response(),
						CreatedBy:   "admin-id",
						UpdatedBy:   "admin-id",
						Title:       "キャベツ祭り開催",
						Body:        "旬のキャベツを大安売り",
						Targets:     []int32{3, 4},
						PublishedAt: 1640962800,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
				{
					Notification: response.Notification{
						ID:          "notification-id",
						Type:        NotificationTypePromotion.Response(),
						CreatedBy:   "admin-id",
						UpdatedBy:   "admin-id",
						Title:       "",
						Body:        "旬のキャベツを大安売り",
						Targets:     []int32{3, 4},
						PromotionID: "promotion-id",
						PublishedAt: 1640962800,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
			},
			promotions: map[string]*Promotion{
				"promotion-id": {
					Promotion: response.Promotion{
						ID:           "promotion-id",
						Title:        "セール情報",
						Description:  "セール詳細",
						Public:       true,
						DiscountType: DiscountTypeAmount.Response(),
						DiscountRate: 3980,
						Code:         "code",
						StartAt:      1640962800,
						EndAt:        1640962800,
						UsedCount:    0,
						UsedAmount:   0,
						CreatedAt:    1640962800,
						UpdatedAt:    1640962800,
					},
				},
			},
			expect: Notifications{
				{
					Notification: response.Notification{
						ID:          "notification-id",
						Type:        NotificationTypeSystem.Response(),
						CreatedBy:   "admin-id",
						UpdatedBy:   "admin-id",
						Title:       "キャベツ祭り開催",
						Body:        "旬のキャベツを大安売り",
						Targets:     []int32{3, 4},
						PublishedAt: 1640962800,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
				{
					Notification: response.Notification{
						ID:          "notification-id",
						Type:        NotificationTypePromotion.Response(),
						CreatedBy:   "admin-id",
						UpdatedBy:   "admin-id",
						Title:       "セール情報",
						Body:        "旬のキャベツを大安売り",
						Targets:     []int32{3, 4},
						PromotionID: "promotion-id",
						PublishedAt: 1640962800,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.notifications.Fill(tt.promotions)
			assert.Equal(t, tt.expect, tt.notifications)
		})
	}
}

func TestNotifications_Response(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		notifications Notifications
		expect        []*response.Notification
	}{
		{
			name: "success",
			notifications: Notifications{
				{
					Notification: response.Notification{
						ID:          "notification-id",
						CreatedBy:   "admin-id",
						UpdatedBy:   "admin-id",
						Title:       "キャベツ祭り開催",
						Body:        "旬のキャベツを大安売り",
						Targets:     []int32{3, 4},
						PublishedAt: 1640962800,
						CreatedAt:   1640962800,
						UpdatedAt:   1640962800,
					},
				},
			},
			expect: []*response.Notification{
				{
					ID:          "notification-id",
					CreatedBy:   "admin-id",
					UpdatedBy:   "admin-id",
					Title:       "キャベツ祭り開催",
					Body:        "旬のキャベツを大安売り",
					Targets:     []int32{3, 4},
					PublishedAt: 1640962800,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, tt.notifications.Response())
		})
	}
}
