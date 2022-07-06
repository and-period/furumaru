//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package messenger

import (
	"context"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type Service interface {
	// お問い合わせ一覧取得
	ListContacts(ctx context.Context, in *ListContactsInput) (entity.Contacts, error)
	// お問い合わせ取得
	GetContact(ctx context.Context, in *GetContactInput) (*entity.Contact, error)
	// お問い合わせ登録
	CreateContact(ctx context.Context, in *CreateContactInput) (*entity.Contact, error)
	// お問い合わせ更新
	UpdateContact(ctx context.Context, in *UpdateContactInput) error
	// 管理者登録通知
	NotifyRegisterAdmin(ctx context.Context, in *NotifyRegisterAdminInput) error
}
