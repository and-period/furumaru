package entity

import "time"

// FacilityUser - 施設利用者情報
type FacilityUser struct {
	UserID        string               `gorm:"primaryKey;<-:create"` // ユーザーID
	ExternalID    string               `gorm:""`                     // ユーザーID（外部用）
	ProducerID    string               `gorm:""`                     // 生産者ID（施設ID）
	Lastname      string               `gorm:""`                     // 姓
	Firstname     string               `gorm:""`                     // 名
	LastnameKana  string               `gorm:""`                     // 姓（かな）
	FirstnameKana string               `gorm:""`                     // 名（かな）
	ProviderType  UserAuthProviderType `gorm:""`                     // 認証方法
	Email         string               `gorm:""`                     // メールアドレス
	PhoneNumber   string               `gorm:""`                     // 電話番号
	CreatedAt     time.Time            `gorm:"<-:create"`            // 登録日時
	UpdatedAt     time.Time            `gorm:""`                     // 更新日時
}

type FacilityUsers []*FacilityUser

func (fs FacilityUsers) Map() map[string]*FacilityUser {
	res := make(map[string]*FacilityUser, len(fs))
	for _, f := range fs {
		res[f.UserID] = f
	}
	return res
}
