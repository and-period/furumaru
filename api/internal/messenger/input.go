package messenger

import (
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type ListContactsInput struct {
	Limit  int64 `validate:"required,max=200"`
	Offset int64 `validate:"min=0"`
}

type GetContactInput struct {
	ContactID string `validate:"required"`
}

type CreateContactInput struct {
	Title       string `validate:"required,max=64"`
	Content     string `validate:"required,max=2000"`
	Username    string `validate:"required,max=64"`
	Email       string `validate:"required,max=256,email"`
	PhoneNumber string `validate:"min=12,max=18,phone_number"`
}

type UpdateContactInput struct {
	ContactID string                 `validate:"required"`
	Status    entity.ContactStatus   `validate:"required,oneof=1 2 3 4"`
	Priority  entity.ContactPriority `validate:"required,oneof=1 2 3"`
	Note      string                 `validate:"max=2000"`
}

type CreateNotificationInput struct {
	CreatedBy   string              `validate:"required"`
	Title       string              `validate:"required,max=128"`
	Body        string              `validate:"required,max=2000"`
	Targets     []entity.TargetType `validate:"min=1,max=3,dive,min=1,max=3"`
	Public      bool                `validate:""`
	PublishedAt time.Time           `validate:"required"`
}

type NotifyRegisterAdminInput struct {
	AdminID  string `validate:"required"`
	Password string `validate:"required"`
}

type NotifyReceivedContactInput struct {
	ContactID string `validate:"required"`
	Username  string `validate:"required"`
	Email     string `validate:"required"`
}
