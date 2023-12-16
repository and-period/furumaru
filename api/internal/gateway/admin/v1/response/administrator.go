package response

// Administrator - システム管理者情報
type Administrator struct {
	ID            string `json:"id"`            // 管理者ID
	Status        int32  `json:"status"`        // 管理者の状態
	Lastname      string `json:"lastname"`      // 姓
	Firstname     string `json:"firstname"`     // 名
	LastnameKana  string `json:"lastnameKana"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana"` // 名(かな)
	Email         string `json:"email"`         // メールアドレス
	PhoneNumber   string `json:"phoneNumber"`   // 電話番号
	CreatedAt     int64  `json:"createdAt"`     // 登録日時
	UpdatedAt     int64  `json:"updatedAt"`     // 更新日時
}

type AdministratorResponse struct {
	Administrator *Administrator `json:"administrator"` // システム管理者情報
}

type AdministratorsResponse struct {
	Administrators []*Administrator `json:"administrators"` // システム管理者一覧
	Total          int64            `json:"total"`          // 合計数
}
