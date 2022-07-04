package messenger

import (
	"github.com/and-period/furumaru/api/internal/messenger/entity"
)

type NotifyRegisterAdminInput struct {
	AdminID  string `validate:"required"`
	Password string `validate:"required"`
}

type CreateNotificationInput struct {
	CreatedBy   string              `validate:"required"`
	CreatorName string              `validate:"required"`
	UpdatedBy   string              `validate:""`
	Title       string              `validate:"required,max=128"`
	Body        string              `validate:"required",max=2000`
	Targets     []entity.TargetType `validate:"min=1,max=2,dive,min=0,max=3"`
	Public      bool                `validate:""`
}
