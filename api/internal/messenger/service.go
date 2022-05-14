//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package messenger

import (
	"context"
)

type MessengerService interface {
	NotifyRegisterAdmin(ctx context.Context, in *NotifyRegisterAdminInput) error
}
