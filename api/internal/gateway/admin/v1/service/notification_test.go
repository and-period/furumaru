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
				ID:          "notification-id",
				CreatedBy:   "admin-id",
				CreatorName: "登録者",
				UpdatedBy:   "admin-id",
				Title:       "キャベツ祭り開催",
				Body:        "旬のキャベツを大安売り",
				Targets: []entity.TargetType{
					entity.PostTargetUsers,
					entity.PostTargetProducers,
				},
				Public:      true,
				PublishedAt: jst.ParseFromUnix(date),
				CreatedAt:   jst.ParseFromUnix(date),
				UpdatedAt:   jst.ParseFromUnix(date),
			},
			expect: &Notification{
				Notification: response.Notification{
					ID:          "notification-id",
					CreatedBy:   "admin-id",
					CreatorName: "登録者",
					UpdatedBy:   "admin-id",
					Title:       "キャベツ祭り開催",
					Body:        "旬のキャベツを大安売り",
					Targets: []response.TargetType{
						response.PostTargetUsers,
						response.PostTargetProducers,
					},
					PublishedAt: 1640962800,
					Public:      true,
					CreatedAt:   1640962800,
					UpdatedAt:   1640962800,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expect, NewNotification(tt.notification))
		})
	}
}
