//nolint:lll
//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package messenger

import (
	"context"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type Service interface {
	// お知らせ一覧取得
	ListNotifications(ctx context.Context, in *ListNotificationsInput) (entity.Notifications, int64, error)
	// お知らせ取得
	GetNotification(ctx context.Context, in *GetNotificationInput) (*entity.Notification, error)
	// お知らせ作成
	CreateNotification(ctx context.Context, in *CreateNotificationInput) (*entity.Notification, error)
	// お知らせ編集
	UpdateNotification(ctx context.Context, in *UpdateNotificationInput) error
	// お知らせ削除
	DeleteNotification(ctx context.Context, in *DeleteNotificationInput) error
	// メッセージ一覧取得
	ListMessages(ctx context.Context, in *ListMessagesInput) (entity.Messages, int64, error)
	// メッセージ取得
	GetMessage(ctx context.Context, in *GetMessageInput) (*entity.Message, error)
	// 管理者登録通知
	NotifyRegisterAdmin(ctx context.Context, in *NotifyRegisterAdminInput) error
	// 管理者パスワードリセット通知
	NotifyResetAdminPassword(ctx context.Context, in *NotifyResetAdminPasswordInput) error
	// お知らせ通知
	NotifyNotification(ctx context.Context, in *NotifyNotificationInput) error
	// お問い合わせ一覧取得
	ListContacts(ctx context.Context, in *ListContactsInput) (entity.Contacts, int64, error)
	// お問い合わせ取得
	GetContact(ctx context.Context, in *GetContactInput) (*entity.Contact, error)
	// お問い合わせ作成
	CreateContact(ctx context.Context, in *CreateContactInput) (*entity.Contact, error)
	// お問い合わせ編集
	UpdateContact(ctx context.Context, in *UpdateContactInput) error
	// お問い合わせ削除
	DeleteContact(ctx context.Context, in *DeleteContactInput) error
	// お問い合わせ種別一覧取得
	ListContactCategories(ctx context.Context, in *ListContactCategoriesInput) (entity.ContactCategories, error)
	// お問い合わせ種別取得
	GetContactCategory(ctx context.Context, in *GetContactCategoryInput) (*entity.ContactCategory, error)
	// お問い合わせ会話履歴一覧取得(お問い合わせID指定)
	ListThreadsByContactID(ctx context.Context, in *ListThreadsByContactIDInput) (entity.Threads, int64, error)
	// お問い合わせ会話履歴取得
	GetThread(ctx context.Context, in *GetThreadInput) (*entity.Thread, error)
	// お問い合わせ会話履歴作成
	CreateThread(ctx context.Context, in *CreateThreadInput) (*entity.Thread, error)
	// お問い合わせ会話履歴編集
	UpdateThread(ctx context.Context, in *UpdateThreadInput) error
	// お問い合わせ会話履歴削除
	DeleteThread(ctx context.Context, in *DeleteThreadInput) error
}
