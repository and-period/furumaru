package entity

import "time"

// ShippingCarrier - 配送会社
type ShippingCarrier int32

const (
	ShippingCarrierUnknown ShippingCarrier = 0
	ShippingCarrierYamato  ShippingCarrier = 1 // ヤマト運輸
	ShippingCarrierSagawa  ShippingCarrier = 2 // 佐川急便
)

// ShippingSize - 配送時の箱の大きさ
type ShippingSize int32

const (
	ShippingSizeUnknown ShippingSize = 0
	ShippingSize60      ShippingSize = 1 // 箱のサイズ:60
	ShippingSize80      ShippingSize = 2 // 箱のサイズ:80
	ShippingSize100     ShippingSize = 3 // 箱のサイズ:100
)

// OrderFulfillment - 配送情報
type OrderFulfillment struct {
	ID              string          `gorm:"primaryKey;<-:create"` // 配送情報
	OrderID         string          `gorm:""`                     // 注文履歴ID
	ShippingID      string          `gorm:""`                     // 配送設定ID
	TrackingNumber  string          `gorm:"default:null"`         // 伝票番号
	ShippingCarrier ShippingCarrier `gorm:""`                     // 配送会社
	ShippingMethod  DeliveryType    `gorm:""`                     // 配送方法
	BoxSize         ShippingSize    `gorm:""`                     // 箱の大きさ
	BoxCount        int64           `gorm:""`                     // 箱の個数
	WeightTotal     int64           `gorm:""`                     // 総重量(g)
	Lastname        string          `gorm:""`                     // 配送先情報 姓
	Firstname       string          `gorm:""`                     // 配送先情報 名
	PostalCode      string          `gorm:""`                     // 配送先情報 郵便番号
	Prefecture      string          `gorm:""`                     // 配送先情報 都道府県
	City            string          `gorm:""`                     // 配送先情報 市区町村
	AddressLine1    string          `gorm:""`                     // 配送先情報 町名・番地
	AddressLine2    string          `gorm:""`                     // 配送先情報 ビル名・号室など
	PhoneNumber     string          `gorm:""`                     // 配送先情報 電話番号
	CreatedAt       time.Time       `gorm:"<-:create"`            // 登録日時
	UpdatedAt       time.Time       `gorm:""`                     // 更新日時
}

type OrderFulfillments []*OrderFulfillment

func (fs OrderFulfillments) MapByOrderID() map[string]*OrderFulfillment {
	res := make(map[string]*OrderFulfillment, len(fs))
	for _, f := range fs {
		res[f.OrderID] = f
	}
	return res
}
