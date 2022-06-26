//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package messenger

import (
	"context"
)

type Service interface {
	// 管理者登録通知
	NotifyRegisterAdmin(ctx context.Context, in *NotifyRegisterAdminInput) error
}
