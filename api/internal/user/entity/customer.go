package entity

import "time"

// Customer - 購入者情報
type Customer struct {
	UserID        string    `gorm:"primaryKey;<-:create"` // ユーザーID
	Lastname      string    `gorm:""`                     // 姓
	Firstname     string    `gorm:""`                     // 名
	LastnameKana  string    `gorm:""`                     // 姓(かな)
	FirstnameKana string    `gorm:""`                     // 名(かな)
	PostalCode    string    `gorm:""`                     // 郵便番号
	Prefecture    int64     `gorm:""`                     // 都道府県
	City          string    `gorm:""`                     // 市区町村
	AddressLine1  string    `gorm:""`                     // 町名・番地
	AddressLine2  string    `gorm:""`                     // ビル名・号室など
	CreatedAt     time.Time `gorm:"<-:create"`            // 登録日時
	UpdatedAt     time.Time `gorm:""`                     // 更新日時
}

type Customers []*Customer

func (cs Customers) Map() map[string]*Customer {
	res := make(map[string]*Customer, len(cs))
	for _, c := range cs {
		res[c.UserID] = c
	}
	return res
}
