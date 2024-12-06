package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
)

// SpotType - スポット種別情報
type SpotType struct {
	ID        string    `gorm:"primaryKey;<-:create"` // スポット種別ID
	Name      string    `gorm:""`                     // スポット種別名
	CreatedAt time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt time.Time `gorm:""`                     // 更新日時
}

type SpotTypes []*SpotType

type NewSpotTypeParams struct {
	Name string
}

func NewSpotType(params *NewSpotTypeParams) *SpotType {
	return &SpotType{
		ID:   uuid.Base58Encode(uuid.New()),
		Name: params.Name,
	}
}
