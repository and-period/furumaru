//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package messenger

import (
	"context"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type Service interface {
	// お問い合わせ一覧取得
	ListContacts(ctx context.Context, in *ListContactsInput) (entity.Contacts, int64, error)
	// お問い合わせ取得
	GetContact(ctx context.Context, in *GetContactInput) (*entity.Contact, error)
	// お問い合わせ登録
	CreateContact(ctx context.Context, in *CreateContactInput) (*entity.Contact, error)
	// お問い合わせ更新
	UpdateContact(ctx context.Context, in *UpdateContactInput) error
	// お知らせ作成
	CreateNotification(ctx context.Context, in *CreateNotificationInput) (*entity.Notification, error)
	// メッセージ一覧取得
	ListMessages(ctx context.Context, in *ListMessagesInput) (entity.Messages, int64, error)
	// メッセージ取得
	GetMessage(ctx context.Context, in *GetMessageInput) (*entity.Message, error)
	// 管理者登録通知
	NotifyRegisterAdmin(ctx context.Context, in *NotifyRegisterAdminInput) error
	// 管理者パスワードリセット通知
	NotifyResetAdminPassword(ctx context.Context, in *NotifyResetAdminPasswordInput) error
	// お問い合わせ受領通知
	NotifyReceivedContact(ctx context.Context, in *NotifyReceivedContactInput) error
	// お知らせ通知
	NotifyNotification(ctx context.Context, in *NotifyNotificationInput) error
}
