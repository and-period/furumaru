package response

// Admin - 管理者情報
//
//nolint:staticcheck
type Admin struct {
	ID            string `json:"id"`            // 管理者ID
	Type          int32  `json:"role"`          // 管理者種別
	Lastname      string `json:"lastname"`      // 姓
	Firstname     string `json:"firstname"`     // 名
	LastnameKana  string `json:"lastnameKana"`  // 姓(かな)
	FirstnameKana string `json:"firstnameKana"` // 名(かな)
	Email         string `json:"email"`         // メールアドレス
	CreatedAt     int64  `json:"createdAt"`     // 登録日時
	UpdatedAt     int64  `json:"updateAt"`      // 更新日時
}
