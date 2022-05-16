package request

type CreateUserRequest struct {
	Email                string `json:"email,omitempty"`                // メールアドレス
	PhoneNumber          string `json:"phoneNumber,omitempty"`          // 電話番号
	Password             string `json:"password,omitempty"`             // パスワード
	PasswordConfirmation string `json:"passwordConfirmation,omitempty"` // パスワード (確認用)
}

type VerifyUserRequest struct {
	ID         string `json:"id,omitempty"`         // ユーザーID
	VerifyCode string `json:"verifyCode,omitempty"` // 検証コード
}

type InitializeUserRequest struct {
	ID        string `json:"id,omitempty"`        // ユーザーID(PK)
	Username  string `json:"username,omitempty"`  // ユーザー名
	AccountID string `json:"accountId,omitempty"` // ユーザーID(表示名)
}

type UpdateUserEmailRequest struct {
	Email string `json:"email,omitempty"` // メールアドレス
}

type VerifyUserEmailRequest struct {
	VerifyCode string `json:"verifyCode,omitempty"` // 検証コード
}

type UpdateUserPasswordRequest struct {
	OldPassword          string `json:"oldPassword,omitempty"`          // 現在のパスワード
	NewPassword          string `json:"newPassword,omitempty"`          // 新しいパスワード
	PasswordConfirmation string `json:"passwordConfirmation,omitempty"` // パスワード (確認用)
}

type ForgotUserPasswordRequest struct {
	Email string `json:"email,omitempty"` // メールアドレス
}

type ResetUserPasswordRequest struct {
	Email                string `json:"email,omitempty"`                // メールアドレス
	VerifyCode           string `json:"verifyCode,omitempty"`           // 検証コード
	Password             string `json:"password,omitempty"`             // パスワード
	PasswordConfirmation string `json:"passwordConfirmation,omitempty"` // パスワード (確認用)
}
