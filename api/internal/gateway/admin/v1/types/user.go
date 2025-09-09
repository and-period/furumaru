package types

// User - 購入者情報
type User struct {
	ID            string `json:"id"`            // 購入者ID
	Status        int32  `json:"status"`        // 購入者ステータス
	Registered    bool   `json:"registered"`    // 会員登録フラグ
	Username      string `json:"username"`      // ユーザー名（表示名）
	AccountID     string `json:"accountId"`     // ユーザー名（検索用）
	Lastname      string `json:"lastname"`      // 氏名（姓）
	Firstname     string `json:"firstname"`     // 氏名（名）
	LastnameKana  string `json:"lastnameKana"`  // 氏名(姓:かな)
	FirstnameKana string `json:"firstnameKana"` // 氏名(名:かな)
	Email         string `json:"email"`         // メールアドレス
	PhoneNumber   string `json:"phoneNumber"`   // 電話番号
	ThumbnailURL  string `json:"thumbnailUrl"`  // サムネイルURL
	CreatedAt     int64  `json:"createdAt"`     // 登録日時
	UpdatedAt     int64  `json:"updateAt"`      // 更新日時
}

// UserToList - 購入者一覧情報
type UserToList struct {
	ID                string `json:"id"`                // 購入者ID
	Lastname          string `json:"lastname"`          // 姓
	Firstname         string `json:"firstname"`         // 名
	Email             string `json:"email"`             // メールアドレス
	Status            int32  `json:"status"`            // 購入者ステータス
	Registered        bool   `json:"registered"`        // 会員登録フラグ
	PrefectureCode    int32  `json:"prefectureCode"`    // 都道府県
	City              string `json:"city"`              // 市区町村
	PaymentTotalCount int64  `json:"paymentTotalCount"` // 支払い回数
}

// UserOrder - 購入者注文情報
type UserOrder struct {
	OrderID   string `json:"orderId"`   // 注文情報ID
	Status    int32  `json:"status"`    // 支払い状況
	SubTotal  int64  `json:"subtotal"`  // 商品合計金額
	Total     int64  `json:"total"`     // 支払い合計金額
	OrderedAt int64  `json:"orderedAt"` // 注文日時
	PaidAt    int64  `json:"paidAt"`    // 支払日時
}

type UserResponse struct {
	User    *User    `json:"user"`    // 購入者情報
	Address *Address `json:"address"` // デフォルト設定の住所情報
}

type UsersResponse struct {
	Users []*UserToList `json:"users"` // 購入者一覧
	Total int64         `json:"total"` // 購入者合計数
}

type UserOrdersResponse struct {
	Orders             []*UserOrder `json:"orders"`             // 注文履歴一覧
	OrderTotalCount    int64        `json:"orderTotalCount"`    // 注文合計回数
	PaymentTotalCount  int64        `json:"paymentTotalCount"`  // 支払い合計回数
	ProductTotalAmount int64        `json:"productTotalAmount"` // 購入商品合計金額
	PaymentTotalAmount int64        `json:"paymentTotalAmount"` // 支払い合計金額
}
