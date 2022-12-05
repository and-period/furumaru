package entity

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/and-period/furumaru/api/internal/codes"
	"github.com/and-period/furumaru/api/pkg/uuid"
	"gorm.io/gorm"
)

// Address - 住所情報
type Address struct {
	ID             string         `gorm:"primaryKey;<-:create"` // 住所ID
	UserID         string         `gorm:""`                     // ユーザーID
	Hash           string         `gorm:""`                     // 住所識別ID(ハッシュ値)
	IsDefault      bool           `gorm:""`                     // デフォルト設定フラグ
	PostalCode     string         `gorm:""`                     // 郵便番号
	Prefecture     string         `gorm:""`                     // 都道府県
	PrefectureCode int64          `gorm:""`                     // 都道府県コード
	City           string         `gorm:""`                     // 市区町村
	AddressLine1   string         `gorm:""`                     // 町名・番地
	AddressLine2   string         `gorm:""`                     // ビル名・号室など
	CreatedAt      time.Time      `gorm:"<-:create"`            // 登録日時
	UpdatedAt      time.Time      `gorm:""`                     // 更新日時
	DeletedAt      gorm.DeletedAt `gorm:"default:null"`         // 退会日時
}

type NewAddressParams struct {
	UserID         string
	IsDefault      bool
	PostalCode     string
	Prefecture     string
	PrefectureCode string
	City           string
	AddressLine1   string
	AddressLine2   string
}

func NewAddress(params *NewAddressParams) (*Address, error) {
	prefecture, err := codes.ToPrefectureValue(params.PrefectureCode)
	if err != nil {
		return nil, err
	}
	return &Address{
		ID:             uuid.Base58Encode(uuid.New()),
		UserID:         params.UserID,
		Hash:           NewAddressHash(params.UserID, params.PostalCode, params.AddressLine1, params.AddressLine2),
		IsDefault:      params.IsDefault,
		PostalCode:     params.PostalCode,
		Prefecture:     params.Prefecture,
		PrefectureCode: prefecture,
		City:           params.City,
		AddressLine1:   params.AddressLine1,
		AddressLine2:   params.AddressLine2,
	}, nil
}

func NewAddressHash(userID, postalCode, addressLine1, addressLine2 string) string {
	fields := []string{userID, postalCode, addressLine1, addressLine2}
	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}
	hash := sha256.Sum256([]byte(strings.Join(fields, ":")))
	return hex.EncodeToString(hash[:])
}
