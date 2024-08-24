package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// ExperienceType - 体験種別情報
type ExperienceType struct {
	ID        string    `gorm:"primaryKey;<-:create"` // 体験種別ID
	Name      string    `gorm:""`                     // 体験種別名
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type ExperienceTypes []*ExperienceType

type NewExperienceTypeParams struct {
	Name string
}

func NewExperienceType(params *NewExperienceTypeParams) *ExperienceType {
	return &ExperienceType{
		ID:   uuid.Base58Encode(uuid.New()),
		Name: params.Name,
	}
}
