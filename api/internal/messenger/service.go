//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package messenger

import (
	"context"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type Service interface {
	// 管理者登録通知
	NotifyRegisterAdmin(ctx context.Context, in *NotifyRegisterAdminInput) error
	// お知らせ作成
	CreateNotification(ctx context.Context, in *CreateNotificationInput) (*entity.Notification, error)
}
