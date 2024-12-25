//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./../../mock/$GOPACKAGE/$GOFILE
package messenger

import (
	"context"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type Service interface {
	// Contact - お問い合わせ
	ListContacts(ctx context.Context, in *ListContactsInput) (entity.Contacts, int64, error) // 一覧取得
	GetContact(ctx context.Context, in *GetContactInput) (*entity.Contact, error)            // １件取得
	CreateContact(ctx context.Context, in *CreateContactInput) (*entity.Contact, error)      // 登録
	UpdateContact(ctx context.Context, in *UpdateContactInput) error                         // 更新
	DeleteContact(ctx context.Context, in *DeleteContactInput) error                         // 削除
	// ContactCategory - お問い合わせ種別
	ListContactCategories(ctx context.Context, in *ListContactCategoriesInput) (entity.ContactCategories, error)         // 一覧取得
	MultiGetContactCategories(ctx context.Context, in *MultiGetContactCategoriesInput) (entity.ContactCategories, error) // 一覧取得(ID指定)
	GetContactCategory(ctx context.Context, in *GetContactCategoryInput) (*entity.ContactCategory, error)                // １件取得
	// ContactRead - お問い合わせ既読管理
	CreateContactRead(ctx context.Context, in *CreateContactReadInput) (*entity.ContactRead, error) // 登録
	// Message - メッセージ
	ListMessages(ctx context.Context, in *ListMessagesInput) (entity.Messages, int64, error) // 一覧取得
	GetMessage(ctx context.Context, in *GetMessageInput) (*entity.Message, error)            // １件取得
	// Notification - お知らせ
	ListNotifications(ctx context.Context, in *ListNotificationsInput) (entity.Notifications, int64, error) // 一覧取得
	GetNotification(ctx context.Context, in *GetNotificationInput) (*entity.Notification, error)            // １件取得
	CreateNotification(ctx context.Context, in *CreateNotificationInput) (*entity.Notification, error)      // 登録
	UpdateNotification(ctx context.Context, in *UpdateNotificationInput) error                              // 更新
	DeleteNotification(ctx context.Context, in *DeleteNotificationInput) error                              // 削除
	// Notify - 通知関連(共通)
	NotifyNotification(ctx context.Context, in *NotifyNotificationInput) error // お知らせ通知
	// NotifyAdmin - 通知関連(管理者宛)
	NotifyRegisterAdmin(ctx context.Context, in *NotifyRegisterAdminInput) error           // 登録通知
	NotifyResetAdminPassword(ctx context.Context, in *NotifyResetAdminPasswordInput) error // パスワードリセット通知
	// NotifyUser - 通知関連(利用者宛)
	NotifyStartLive(ctx context.Context, in *NotifyStartLiveInput) error             // ライブ配信開始通知
	NotifyOrderAuthorized(ctx context.Context, in *NotifyOrderAuthorizedInput) error // 支払い完了通知
	NotifyOrderShipped(ctx context.Context, in *NotifyOrderShippedInput) error       // 発送完了通知
	// ReserveNotification - 通知予約関連
	ReserveNotification(ctx context.Context, in *ReserveNotificationInput) error // お知らせ通知予約
	ReserveStartLive(ctx context.Context, in *ReserveStartLiveInput) error       // ライブ配信開始通知予約
	// Threads - お問い合わせ会話履歴
	ListThreads(ctx context.Context, in *ListThreadsInput) (entity.Threads, int64, error) // 一覧取得
	GetThread(ctx context.Context, in *GetThreadInput) (*entity.Thread, error)            // １件取得
	CreateThread(ctx context.Context, in *CreateThreadInput) (*entity.Thread, error)      // 登録
	UpdateThread(ctx context.Context, in *UpdateThreadInput) error                        // 更新
	DeleteThread(ctx context.Context, in *DeleteThreadInput) error                        // 削除
}
