package request

type CreateAuthUserRequest struct {
	Username             string `json:"username,omitempty"`             // ユーザー名
	AccountID            string `json:"accountId,omitempty"`            // ユーザーID(表示名)
	Lastname             string `json:"lastname,omitempty"`             // 姓
	Firstname            string `json:"firstname,omitempty"`            // 名
	LastnameKana         string `json:"lastnameKana,omitempty"`         // 姓（かな）
	FirstnameKana        string `json:"firstnameKana,omitempty"`        // 名（かな）
	Email                string `json:"email,omitempty"`                // メールアドレス
	PhoneNumber          string `json:"phoneNumber,omitempty"`          // 電話番号
	Password             string `json:"password,omitempty"`             // パスワード
	PasswordConfirmation string `json:"passwordConfirmation,omitempty"` // パスワード (確認用)
}

type VerifyAuthUserRequest struct {
	ID         string `json:"id,omitempty"`         // ユーザーID
	VerifyCode string `json:"verifyCode,omitempty"` // 検証コード
}

type CreateAuthUserWithOAuthRequest struct {
	Username      string `json:"username,omitempty"`      // ユーザー名
	AccountID     string `json:"accountId,omitempty"`     // ユーザーID(表示名)
	Lastname      string `json:"lastname,omitempty"`      // 姓
	Firstname     string `json:"firstname,omitempty"`     // 名
	LastnameKana  string `json:"lastnameKana,omitempty"`  // 姓（かな）
	FirstnameKana string `json:"firstnameKana,omitempty"` // 名（かな）
	PhoneNumber   string `json:"phoneNumber,omitempty"`   // 電話番号
}

type UpdateAuthUserEmailRequest struct {
	Email string `json:"email,omitempty"` // メールアドレス
}

type VerifyAuthUserEmailRequest struct {
	VerifyCode string `json:"verifyCode,omitempty"` // 検証コード
}

type UpdateAuthUserThumbnailRequest struct {
	ThumbnailURL string `json:"thumbnailUrl,omitempty"` // サムネイルURL
}
