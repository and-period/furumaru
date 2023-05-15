package messenger

import (
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type ListNotificationsInput struct {
	Limit         int64                     `validate:"required,max=200"`
	Offset        int64                     `validate:"min=0"`
	Since         time.Time                 `validate:""`
	Until         time.Time                 `validate:""`
	OnlyPublished bool                      `validate:""`
	Orders        []*ListNotificationsOrder `validate:"omitempty,dive,required"`
}

type ListNotificationsOrder struct {
	Key        entity.NotificationOrderBy `validate:"required"`
	OrderByASC bool                       `validate:""`
}

type GetNotificationInput struct {
	NotificationID string `validate:"required"`
}

type CreateNotificationInput struct {
	CreatedBy   string              `validate:"required"`
	Title       string              `validate:"required,max=128"`
	Body        string              `validate:"required,max=2000"`
	Targets     []entity.TargetType `validate:"min=1,max=3,dive,min=1,max=3"`
	Public      bool                `validate:""`
	PublishedAt time.Time           `validate:"required"`
}

type UpdateNotificationInput struct {
	NotificationID string              `validate:"required"`
	Title          string              `validate:"required"`
	Body           string              `validate:"required"`
	Targets        []entity.TargetType `validate:"min=1,max=3,dive,min=1,max=3"`
	Public         bool                `validate:""`
	PublishedAt    time.Time           `validate:"required"`
	UpdatedBy      string              `validete:"required"`
}

type DeleteNotificationInput struct {
	NotificationID string `validate:"required"`
}

type ListMessagesInput struct {
	UserType entity.UserType      `validate:"required,oneof=1 2"`
	UserID   string               `validate:"required"`
	Limit    int64                `validate:"required,max=200"`
	Offset   int64                `validate:"min=0"`
	Orders   []*ListMessagesOrder `validate:"omitempty,dive,required"`
}

type ListMessagesOrder struct {
	Key        entity.MessageOrderBy `validate:"required"`
	OrderByASC bool                  `validate:""`
}

type GetMessageInput struct {
	MessageID string          `validate:"required"`
	UserType  entity.UserType `validate:"omitempty,oneof=1 2"`
	UserID    string          `validate:"omitempty"`
}

type NotifyRegisterAdminInput struct {
	AdminID  string `validate:"required"`
	Password string `validate:"required"`
}

type NotifyResetAdminPasswordInput struct {
	AdminID  string `validate:"required"`
	Password string `validate:"required"`
}

type NotifyReceivedContactInput struct {
	ContactID string `validate:"required"`
}

type NotifyNotificationInput struct {
	NotificationID string `validate:"required"`
}
