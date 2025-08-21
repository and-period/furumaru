package response

// AuthUser - 認証済みユーザー情報
type AuthUser struct {
	ID            string `json:"id"`            // ユーザーID
	Lastname      string `json:"lastname"`      // 姓
	Firstname     string `json:"firstname"`     // 名
	LastnameKana  string `json:"lastnameKana"`  // 姓（かな）
	FirstnameKana string `json:"firstnameKana"` // 名（かな）
	Email         string `json:"email"`         // メールアドレス
	PhoneNumber   string `json:"phoneNumber"`   // 電話番号
	LastCheckInAt int64  `json:"lastCheckInAt"` // 最新のチェックイン日時
}

type AuthUserResponse struct {
	*AuthUser
}
