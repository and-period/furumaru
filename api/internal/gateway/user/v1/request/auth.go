package request

type SignInRequest struct {
	Username string `json:"username,omitempty"` // ユーザー名 (メールアドレス,電話番号)
	Password string `json:"password,omitempty"` // パスワード
}

type RefreshAuthTokenRequest struct {
	RefreshToken string `json:"refreshToken,omitempty"` // 更新トークン
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
