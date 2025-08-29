package request

type SignInRequest struct {
	Username string `json:"username"` // ユーザー名 (メールアドレス,電話番号)
	Password string `json:"password"` // パスワード
}

type RefreshAuthTokenRequest struct {
	RefreshToken string `json:"refreshToken"` // 更新トークン
}

type UpdateAuthPasswordRequest struct {
	OldPassword          string `json:"oldPassword"`          // 現在のパスワード
	NewPassword          string `json:"newPassword"`          // 新しいパスワード
	PasswordConfirmation string `json:"passwordConfirmation"` // パスワード (確認用)
}

type ForgotAuthPasswordRequest struct {
	Email string `json:"email"` // メールアドレス
}

type ResetAuthPasswordRequest struct {
	Email                string `json:"email"`                // メールアドレス
	VerifyCode           string `json:"verifyCode"`           // 検証コード
	Password             string `json:"password"`             // パスワード
	PasswordConfirmation string `json:"passwordConfirmation"` // パスワード (確認用)
}
