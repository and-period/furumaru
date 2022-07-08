package messenger

import (
	"time"

	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type NotifyRegisterAdminInput struct {
	AdminID  string `validate:"required"`
	Password string `validate:"required"`
}

type CreateNotificationInput struct {
	CreatedBy   string              `validate:"required"`
	Title       string              `validate:"required,max=128"`
	Body        string              `validate:"required,max=2000"`
	Targets     []entity.TargetType `validate:"min=1,max=3,dive,min=1,max=3"`
	Public      bool                `validate:""`
	PublishedAt time.Time           `validate:"required"`
}
