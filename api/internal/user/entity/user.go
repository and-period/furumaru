package entity

import (
	"time"

	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// ProviderType - 認証方法
type ProviderType int32

const (
	ProviderTypeUnknown ProviderType = 0
	ProviderTypeEmail   ProviderType = 1 // メールアドレス/SMS認証
	ProviderTypeOAuth   ProviderType = 2 // OAuth認証
)

// User - 購入者情報
type User struct {
	ID           string         `gorm:"primaryKey;<-:create"` // ユーザーID
	AccountID    string         `gorm:""`                     // ユーザーID (検索用)
	CognitoID    string         `gorm:""`                     // ユーザーID (Cognito用)
	ProviderType ProviderType   `gorm:""`                     // 認証方法
	Username     string         `gorm:""`                     // ユーザー名 (表示用)
	Email        string         `gorm:"default:null"`         // メールアドレス
	PhoneNumber  string         `gorm:"default:null"`         // 電話番号
	ThumbnailURL string         `gorm:""`                     // サムネイルURL
	CreatedAt    time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt    time.Time      `gorm:""`                     // 更新日時
	VerifiedAt   time.Time      `gorm:"default:null"`         // 確認日時
	DeletedAt    gorm.DeletedAt `gorm:"default:null"`         // 削除日時
}

func NewUser(cognitoID string, provider ProviderType, email, phoneNumber string) *User {
	return &User{
		ID:           uuid.Base58Encode(uuid.New()),
		CognitoID:    cognitoID,
		ProviderType: provider,
		Email:        email,
		PhoneNumber:  phoneNumber,
	}
}
