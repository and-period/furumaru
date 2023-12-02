package response

// User - 購入者情報
type User struct {
	ID         string `json:"id"`         // 購入者ID
	Registered bool   `json:"registered"` // 会員登録フラグ
	Email      string `json:"email"`      // メールアドレス
	Status     int32  `json:"status"`     // 購入者ステータス
	CreatedAt  int64  `json:"createdAt"`  // 登録日時
	UpdatedAt  int64  `json:"updateAt"`   // 更新日時
	*Address          // デフォルト設定の住所情報
}

// UserToList - 購入者一覧情報
type UserToList struct {
	ID             string `json:"id"`             // 購入者ID
	Lastname       string `json:"lastname"`       // 姓
	Firstname      string `json:"firstname"`      // 名
	Email          string `json:"email"`          // メールアドレス
	Registered     bool   `json:"registered"`     // 会員登録フラグ
	PrefectureCode int32  `json:"prefectureCode"` // 都道府県
	City           string `json:"city"`           // 市区町村
	TotalOrder     int64  `json:"totalOrder"`     // 購入回数
	TotalAmount    int64  `json:"totalAmount"`    // 購入金額
}

// UserOrder - 購入者注文情報
type UserOrder struct {
	OrderID   string `json:"orderId"`   // 注文情報ID
	Status    int32  `json:"status"`    // 支払い状況
	Total     int64  `json:"total"`     // 支払い合計金額
	OrderedAt int64  `json:"orderedAt"` // 注文日時
	PaidAt    int64  `json:"paidAt"`    // 支払日時
}

type UserResponse struct {
	User *User `json:"user"` // 購入者情報
}

type UsersResponse struct {
	Users []*UserToList `json:"users"` // 購入者一覧
	Total int64         `json:"total"` // 購入者合計数
}

type UserOrdersResponse struct {
	Orders      []*UserOrder `json:"orders"`      // 注文履歴一覧
	Total       int64        `json:"total"`       // 注文履歴合計数
	TotalAmount int64        `json:"totalAmount"` // 購入金額
}
