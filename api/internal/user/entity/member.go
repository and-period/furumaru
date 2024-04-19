package entity

import (
	"strings"
	"time"
)

// ProviderType - 認証方法
type ProviderType int32

const (
	ProviderTypeUnknown ProviderType = 0
	ProviderTypeEmail   ProviderType = 1 // メールアドレス/SMS認証
	ProviderTypeOAuth   ProviderType = 2 // OAuth認証
)

// Member - 会員情報
type Member struct {
	UserID        string       `gorm:"primaryKey;<-:create"` // ユーザーID
	CognitoID     string       `gorm:""`                     // ユーザーID (Cognito用)
	AccountID     string       `gorm:"default:null"`         // ユーザーID (検索用)
	Username      string       `gorm:""`                     // 表示名
	Lastname      string       `gorm:""`                     // 姓
	Firstname     string       `gorm:""`                     // 名
	LastnameKana  string       `gorm:""`                     // 姓（かな）
	FirstnameKana string       `gorm:""`                     // 名（かな）
	ProviderType  ProviderType `gorm:""`                     // 認証方法
	Email         string       `gorm:"default:null"`         // メールアドレス
	PhoneNumber   string       `gorm:"default:null"`         // 電話番号
	ThumbnailURL  string       `gorm:""`                     // サムネイルURL
	CreatedAt     time.Time    `gorm:"<-:create"`            // 登録日時
	UpdatedAt     time.Time    `gorm:""`                     // 更新日時
	VerifiedAt    time.Time    `gorm:"default:null"`         // 確認日時
}

type Members []*Member

func (m *Member) Name() string {
	str := strings.Join([]string{m.Lastname, m.Firstname}, " ")
	return strings.TrimSpace(str)
}

func (ms Members) Map() map[string]*Member {
	res := make(map[string]*Member, len(ms))
	for _, m := range ms {
		res[m.UserID] = m
	}
	return res
}
