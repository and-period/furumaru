package request

type SignInRequest struct {
	Username string `json:"username,omitempty"` // ユーザー名 (メールアドレス,電話番号)
	Password string `json:"password,omitempty"` // パスワード
}

type RefreshAuthTokenRequest struct {
	RefreshToken string `json:"refreshToken,omitempty"` // 更新トークン
}

type CreateAuthRequest struct {
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

type VerifyAuthRequest struct {
	ID         string `json:"id,omitempty"`         // ユーザーID
	VerifyCode string `json:"verifyCode,omitempty"` // 検証コード
}

type CreateAuthWithOAuthRequest struct {
	Username      string `json:"username,omitempty"`      // ユーザー名
	AccountID     string `json:"accountId,omitempty"`     // ユーザーID(表示名)
	Lastname      string `json:"lastname,omitempty"`      // 姓
	Firstname     string `json:"firstname,omitempty"`     // 名
	LastnameKana  string `json:"lastnameKana,omitempty"`  // 姓（かな）
	FirstnameKana string `json:"firstnameKana,omitempty"` // 名（かな）
	PhoneNumber   string `json:"phoneNumber,omitempty"`   // 電話番号
}

type UpdateAuthEmailRequest struct {
	Email string `json:"email,omitempty"` // メールアドレス
}

type VerifyAuthEmailRequest struct {
	VerifyCode string `json:"verifyCode,omitempty"` // 検証コード
}

type UpdateAuthPasswordRequest struct {
	OldPassword          string `json:"oldPassword,omitempty"`          // 現在のパスワード
	NewPassword          string `json:"newPassword,omitempty"`          // 新しいパスワード
	PasswordConfirmation string `json:"passwordConfirmation,omitempty"` // パスワード (確認用)
}

type ForgotAuthPasswordRequest struct {
	Email string `json:"email,omitempty"` // メールアドレス
}

type ResetAuthPasswordRequest struct {
	Email                string `json:"email,omitempty"`                // メールアドレス
	VerifyCode           string `json:"verifyCode,omitempty"`           // 検証コード
	Password             string `json:"password,omitempty"`             // パスワード
	PasswordConfirmation string `json:"passwordConfirmation,omitempty"` // パスワード (確認用)
}
